# LGPD Right-to-Erasure Workflow

> Operational procedure for handling Art. 18 IV (LGPD) erasure requests.
> Target SLA: respond within 15 days (Art. 19). Owner: DPO.

## Intake

Requests arrive via:
- DPO email mailbox (LGPD-only; segregated)
- Customer portal "exercise data rights" button (planned)
- Postal mail (rare; redirect to DPO email)

Per request, create a tracking ticket `LGPD-YYYY-NNNN` in the DPO's secure
issue tracker.

## Stage 1 — Identity verification (within 3 days)

Required:
- CPF/CNPJ of the data subject
- Government-issued ID copy
- Second factor: SMS code via Identos OR notarised request

If unable to verify within 7 days: reject with cause (Art. 18 § 1º — entity must
demonstrate identity to exercise the right).

Tracker fields:
```yaml
ticket: LGPD-2026-0001
received_at: 2026-05-24T10:30:00-03:00
subject_cpf: hashed-sha256-of-cpf
verification_method: identos-sms
verified_at: 2026-05-26T14:00:00-03:00
```

## Stage 2 — Eligibility check (within 7 days)

Run `scripts/lgpd-eligibility.sh <subject-id>` to query CRDB across all tables:

```bash
bash scripts/lgpd-eligibility.sh <subject-id>
```

Output for each table the subject appears in:
- table name
- row count
- earliest occurred_at
- regulatory hold status (calculated from `data-lifecycle/README.md` retention)
- erasure eligibility decision (ELIGIBLE / DEFERRED_UNTIL_<date> / FROZEN_REGULATORY)

If subject has rows in `fx_trades`, `bacen_reports`, `audit_events` etc. → those
rows are FROZEN (regulatory hold). PII fields in `actors` may still be redactable
without affecting the hold (the trade record stays; the name on it is replaced
with `[REDACTED PER LGPD ART 18 IV LGPD-2026-0001]`).

## Stage 3 — Execution (within 12 days)

Generate execution plan as YAML:

```yaml
ticket: LGPD-2026-0001
operations:
  - table: actors
    where: id IN ('uuid-1', 'uuid-2')
    op: redact
    fields: [name, email, tax_id, phone, address]
  - table: counterparties
    where: id IN ('uuid-3')
    op: redact
    fields: [name, beneficial_owner_name, beneficial_owner_doc]
  - table: quote_streams
    where: requester_id IN ('uuid-1')
    op: hard_delete  # no regulatory hold
  - table: screening_results
    where: actor_id IN ('uuid-1')
    op: redact
    fields: [hit_details, raw_evidence]
```

DPO + Compliance Officer co-sign the plan. Then:

```bash
go run cmd/erasure-worker --ticket LGPD-2026-0001 --plan plan.yaml --dry-run
# review output, then:
go run cmd/erasure-worker --ticket LGPD-2026-0001 --plan plan.yaml --execute
```

Every operation:
- Runs inside a single CRDB transaction per table
- Emits `audit_event(type='LGPD_ERASURE', ticket=LGPD-2026-0001, table=X, row_count=N, before_hash, after_hash)`
- Updates the outbox with `lgpd.erasure_completed.v1` for downstream (LedgerOS,
  ComplOS) to apply equivalent erasure

`cmd/erasure-worker/` is planned — tracked in `.base/plans/00-governance/lgpd-backlog.md`.
Until implemented, execute via reviewed SQL scripts under 4-eyes (DPO + Platform Lead).

## Stage 4 — Response to subject (within 15 days)

Send signed confirmation:
- What was erased / redacted
- What was retained + LEGAL BASIS (cite the regulation)
- Subject's right to escalate to ANPD if dissatisfied

Save the signed confirmation under `.audit-bundles/lgpd-requests/YYYY/LGPD-YYYY-NNNN/`.

## Stage 5 — Post-completion audit

Within 30 days, DPO verifies:
- No queries against erased fields succeed (sample-query test)
- Backups older than the request date contain the data (expected — retention not
  retroactive to backups; backups age out per schedule)
- After full backup-retention period (90 days), erasure is "complete" — record
  closure in tracker

## Edge cases

### Subject appears only in audit logs

Audit logs cannot be edited (would break integrity). Respond explaining the
regulatory basis. Optionally, encrypt the row with a tenant-key + destroy the
key on next legal-basis expiry → "crypto-shredding" (planned, ISO 27001 8.10).

### Subject is also an employee actor

Employee records have separate retention (10 years per CLT). DPO consults HR
before any erasure of `actors` rows with role='EMPLOYEE'.

### Request volume spike

If > 5 simultaneous requests, DPO may invoke Art. 19 § 4º (extension by up to
2 months with notification) — must notify both the subject AND ANPD.

## Reporting

Quarterly to leadership:
- Total requests received
- Within-SLA % (target ≥ 95%)
- Deferred / rejected breakdown
- Average time-to-completion
- Escalations to ANPD (target: zero)

Annual to ANPD per regulator request (Art. 38 LGPD).

## Cross-references

- `README.md` (data classification + retention)
- `docs/security/iso27001-controls-mapping.md` controls 5.34, 8.10
- `docs/security/incident-response.md` § S-3 PII breach
- `.base/plans/00-governance/lgpd-backlog.md` for `cmd/erasure-worker` implementation tracking
