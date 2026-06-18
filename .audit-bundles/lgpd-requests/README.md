# LGPD Erasure Request Audit Bundles

> Evidence root for every LGPD Art. 18 IV right-to-erasure request handled
> through [`docs/security/data-lifecycle/erasure-workflow.md`](../../docs/security/data-lifecycle/erasure-workflow.md).
> 7-year retention per regulatory baseline (Art. 19 + ISO 27001 control 5.34).

## Layout per request

Each ticket lives in its own folder, named exactly after the ticket id:

```
.audit-bundles/lgpd-requests/
├── README.md                        # this file
├── template/                        # canonical layout — copy on intake
│   ├── intake.yaml
│   ├── eligibility-report.md
│   ├── plan.json
│   ├── executor-logs/
│   │   └── .gitkeep
│   ├── completion-proof.json
│   └── notes.md
└── LGPD-2026-0001/                  # one folder per real request
    ├── intake.yaml                  # Stage 1 — verified identity record
    ├── eligibility-report.md        # Stage 2 — scripts/lgpd-eligibility.sh output
    ├── plan.json                    # Stage 3 — co-signed plan that fed cmd/erasure-worker
    ├── executor-logs/               # Stage 3 — stdout/stderr from each erasure-worker run
    │   ├── dry-run-2026-05-25T14:00.json
    │   └── execute-2026-05-26T09:30.json
    ├── completion-proof.json        # Stage 5 — final audit_event + outbox event payloads
    └── notes.md                     # Free-form context (decisions, escalations)
```

## What each file holds

| File | Stage | Owner | Purpose |
|------|-------|-------|---------|
| `intake.yaml` | 1 | DPO | CPF/CNPJ hash, verification method, verified_at timestamp |
| `eligibility-report.md` | 2 | DPO + Compliance | Per-table eligibility decision (ELIGIBLE / DEFERRED / FROZEN) |
| `plan.json` | 3 | DPO + Compliance (co-signed) | Input to `cmd/erasure-worker --plan` — validates against `schemas/erasure-plan-v1.json` |
| `executor-logs/dry-run-*.json` | 3 | Platform on-call | SQL preview + row-count estimates |
| `executor-logs/execute-*.json` | 3 | Platform on-call | Live run stdout — JSON-formatted via slog |
| `completion-proof.json` | 5 | DPO | Captured rows from `audit_events` + `outbox_events` for the ticket |
| `notes.md` | any | DPO | Any human-readable context not captured above |

## Hard rules

1. **NEVER commit raw PII.** All identifying values must be hashed or
   tokenised before they land in this folder. The plan itself only carries
   table names + WHERE clauses on indexed UUIDs.
2. **Append-only.** Once a folder exists, edit only `notes.md`. Every other
   file is a frozen record; rerunning a stage creates a new file (e.g.
   `executor-logs/dry-run-<ts>.json`), it does NOT overwrite the previous one.
3. **DPO + Compliance Officer** signatures live inline in `plan.json`
   (`approvals: ["dpo","compliance_officer"]`) AND as detached signature files
   alongside (`plan.json.dpo.sig`, `plan.json.compliance.sig`) when the
   organisation uses GPG/age. Either form satisfies the audit.
4. **Retention.** Folders older than 7 years are eligible for archive to
   cold storage; that is itself a planned change and must be logged in an
   ADR. Until then: keep them here.

## Creating a new bundle

```bash
# At Stage 1, after identity verification succeeds:
cp -R .audit-bundles/lgpd-requests/template \
      .audit-bundles/lgpd-requests/LGPD-YYYY-NNNN
$EDITOR .audit-bundles/lgpd-requests/LGPD-YYYY-NNNN/intake.yaml
```

Subsequent stages fill in the other files. The bundle closes when Stage 5
completes and `completion-proof.json` is populated.

## Related

- Workflow: [`docs/security/data-lifecycle/erasure-workflow.md`](../../docs/security/data-lifecycle/erasure-workflow.md)
- Executor: [`cmd/erasure-worker/main.go`](../../cmd/erasure-worker/main.go)
- Plan schema: [`schemas/erasure-plan-v1.json`](../../schemas/erasure-plan-v1.json)
- Runbook: [`docs/operations/runbook-index.md`](../../docs/operations/runbook-index.md)
- ISO 27001 controls: 5.34 (privacy), 8.10 (information deletion), 8.15 (logging)
