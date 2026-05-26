# Archival Cron — `cmd/archiver`

> Nightly job (planned) that ages CRDB rows out to GCS Coldline per
> `README.md` retention schedule. Owned by Platform team. Currently tracked
> as future work — see `.base/plans/00-governance/lgpd-backlog.md`.

## Schedule

`02:00 America/Sao_Paulo` daily (chosen to avoid overlap with EOD batch
which runs at SP close + 30min, typically ~17:30 SP).

Driven by Kubernetes CronJob via Helm:

```yaml
# deploy/helm/exchangeos/values.yaml (planned addition)
archiver:
  enabled: true
  schedule: "0 5 * * *"  # 05:00 UTC = 02:00 SP year-round (handles DST)
  resources:
    requests: { cpu: 200m, memory: 512Mi }
    limits:   { cpu: 1,    memory: 2Gi }
  serviceAccount: exchangeos-archiver
```

## Per-table policy (data-driven from a config table)

```sql
CREATE TABLE archive_policy (
    table_name       TEXT PRIMARY KEY,
    age_threshold    INTERVAL NOT NULL,    -- when to archive
    retention_after  INTERVAL,             -- when to hard-delete from CRDB (NULL = keep)
    gcs_prefix       TEXT NOT NULL,
    legal_hold       BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed (mirrors README.md retention schedule)
INSERT INTO archive_policy VALUES
('quotes',           interval '90 days', interval '5 years',  'quotes/',       false, now()),
('quote_streams',    interval '30 days', interval '1 year',   'quote_streams/',false, now()),
('eod_jobs',         interval '90 days', interval '1 year',   'eod_jobs/',     false, now()),
('outbox_dispatched_archive', interval '90 days', interval '7 years', 'outbox/', false, now()),
('audit_events',     interval '1 year',  NULL,                'audit_events/', true,  now()),  -- never delete
('fx_trades',        interval '1 year',  NULL,                'fx_trades/',    true,  now()),
('bacen_reports',    interval '1 year',  NULL,                'bacen_reports/',true,  now()),
('cls_cycles',       interval '1 year',  NULL,                'cls_cycles/',   true,  now()),
('payin_instructions', interval '1 year',NULL,                'payin/',        true,  now()),
('net_reports',      interval '1 year',  NULL,                'net_reports/',  true,  now())
ON CONFLICT (table_name) DO UPDATE SET updated_at = now();
```

## Algorithm (per table)

```
for each policy in archive_policy where updated_at > 0:
    1. SELECT * FROM <table> WHERE archived_at IS NULL
                                AND occurred_at < now() - age_threshold
        ORDER BY occurred_at ASC LIMIT 10_000
    2. if no rows: continue
    3. Bulk write to GCS:
        gs://exchangeos-archive-<env>/<gcs_prefix>YYYY/MM/DD/batch-<uuid>.parquet.gz
    4. Compute SHA256 of bytes uploaded
    5. UPDATE <table> SET archived_at = now() WHERE id IN (...) RETURNING id
    6. INSERT INTO archive_manifest (table_name, batch_uuid, row_count, sha256, gcs_uri, completed_at)
    7. If policy.retention_after IS NOT NULL AND NOT legal_hold:
        DELETE FROM <table> WHERE archived_at < now() - retention_after
    8. Emit outbox event archival.batch_completed.v1
```

## Manifest table

```sql
CREATE TABLE archive_manifest (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_name      TEXT NOT NULL,
    batch_uuid      UUID NOT NULL UNIQUE,
    row_count       INTEGER NOT NULL,
    sha256          TEXT NOT NULL,
    gcs_uri         TEXT NOT NULL,
    completed_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    legal_hold      BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX ON archive_manifest(table_name, completed_at DESC);
```

The manifest is itself C3-classified + never deleted (audit chain).

## GCS bucket lifecycle (Terraform)

```hcl
resource "google_storage_bucket" "exchangeos_archive" {
  name     = "exchangeos-archive-${var.env}"
  location = "SOUTHAMERICA-EAST1"
  storage_class = "COLDLINE"
  versioning { enabled = true }

  lifecycle_rule {
    condition { age = 365 }
    action {
      type          = "SetStorageClass"
      storage_class = "ARCHIVE"
    }
  }

  encryption {
    default_kms_key_name = google_kms_crypto_key.archive.id
  }

  retention_policy {
    retention_period = 10 * 365 * 24 * 3600  # 10 years for legal-hold safety
  }
}
```

## Restore-from-archive procedure

Quarterly drill (DR runbook § restore). Process:

1. Pick a random batch from `archive_manifest` (≥ 30 days old)
2. `gsutil cp gs://...batch-<uuid>.parquet.gz /tmp/`
3. Verify SHA256 matches manifest
4. Load into staging CRDB via `gcloud sql import` or `cockroach import`
5. Spot-check 10 rows match expected schema + content
6. Tear down staging table

Drill writeup: `.audit-bundles/restore-drills/YYYY-Qn.md`.

## SLOs for archival

| SLI | SLO |
|-----|-----|
| Daily archival job success | 100% / 30d (page on any miss) |
| Archive completeness (rows pending archival vs total) | < 1% of eligible rows behind |
| Restore time-to-verify | < 30 min for a single batch |

## Cross-references

- `README.md` (data classification + retention schedule)
- `erasure-workflow.md` (LGPD interaction with archival)
- `docs/security/dr-runbook.md` § restore-from-archive
- ISO 27001 control 8.10 (information deletion) + 8.13 (information backup)
