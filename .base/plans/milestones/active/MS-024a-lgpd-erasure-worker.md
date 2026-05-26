# MS-024a — LGPD Erasure Worker

| Field | Value |
|-------|-------|
| **Code** | MS-024a |
| **Name** | lgpd-erasure-worker |
| **Phase** | F-OPS-PROD (production hardening cycle) |
| **Sprint** | 1 of MS-024 cycle |
| **Status** | ACTIVE |
| **Owner** | Compliance + Platform |
| **Created** | 2026-05-24 |
| **Dependencies** | MS-024h (postgres repos) for full-table coverage |

## Why this milestone

The LGPD right-to-erasure workflow is **specified** in `docs/security/data-lifecycle/erasure-workflow.md` (5-stage process within 15-day SLA) but the executor `cmd/erasure-worker/` does not exist. Current state: erasures executed via reviewed SQL scripts under 4-eyes. That's brittle, slow, and audit-fragile.

## Description

Implement `cmd/erasure-worker/` — a one-shot Go binary that takes a signed
execution plan (YAML) + executes the redactions/deletions inside a transactional
boundary per table, emits audit events, and updates the outbox so downstream
modules (LedgerOS, ComplOS) apply equivalent erasure.

## Acceptance Criteria

- [ ] `cmd/erasure-worker/main.go` accepting `--ticket LGPD-YYYY-NNNN`, `--plan plan.yaml`, `--dry-run|--execute`
- [ ] Plan validation against JSON schema (`schemas/erasure-plan-v1.json`)
- [ ] Per-table operation: `redact` (UPDATE fields to `[REDACTED PER LGPD ART 18 IV <ticket>]`) and `hard_delete` (DELETE within tx)
- [ ] Every operation wrapped in single CRDB transaction with `SET TRANSACTION PRIORITY HIGH`
- [ ] Emits `audit_event(type='LGPD_ERASURE', ticket, table, row_count, before_hash, after_hash)` per op
- [ ] Emits outbox `lgpd.erasure_completed.v1` for downstream
- [ ] Dry-run prints the SQL + affected row count without mutating
- [ ] 4-eyes enforcement: refuses execute unless `EXCHANGEOS_ERASURE_APPROVERS=dpo,compliance_officer` env contains both roles
- [ ] Integration test against ephemeral CRDB (testcontainers) covering 3 redact + 1 hard-delete + 1 regulatory-frozen rejection
- [ ] Helm CronJob template `templates/erasure-worker-cronjob.yaml` (manual-trigger only — not scheduled)
- [ ] Runbook entry in `docs/operations/runbook-index.md`
- [ ] Cross-link from `docs/security/data-lifecycle/erasure-workflow.md` → "Stage 3 executor" updated to mention this binary

## Deliverables

- `cmd/erasure-worker/main.go` + handlers
- `internal/erasure/plan.go` (parser + validator)
- `internal/erasure/executor.go` (transactional executor)
- `internal/erasure/audit.go` (audit event emitter)
- `schemas/erasure-plan-v1.json` JSON schema
- `tests/integration/erasure_test.go` (`//go:build integration`)
- `deploy/helm/exchangeos/templates/erasure-worker-cronjob.yaml`
- `.audit-bundles/lgpd-requests/` template directory + README

## Cross-References

- `docs/security/data-lifecycle/erasure-workflow.md` — workflow spec
- `docs/security/data-lifecycle/README.md` — retention schedule + classification
- `scripts/lgpd-eligibility.sh` — discovery feeds into plan generation
- ISO 27001 controls 5.34, 8.10
- LGPD Art. 18 IV, Art. 19, Art. 16
