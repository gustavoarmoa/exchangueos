# MS-024l — CRDB Hub TLS Cross-repo PR

| Field | Value |
|-------|-------|
| **Code** | MS-024l |
| **Name** | crdb-hub-cross-repo-pr |
| **Phase** | F-OPS-PROD |
| **Sprint** | 1 of MS-024 cycle (blocking for prod deploy) |
| **Status** | DELIVERED-PENDING-REMOTE |
| **Owner** | Platform + Shared infra team |
| **Created** | 2026-05-24 |
| **Started** | 2026-06-19 |
| **Delivered** | 2026-06-19 |
| **Updated** | 2026-06-19 |
| **Dependencies** | None |
| **Cockroachdb branch** | `feat/exchangeos-module` |
| **Cockroachdb commit** | `b3d0689` |
| **Pending operator action** | Configure remote on `revenu-platform/cockroachdb/` (no remote present), then `git push -u origin feat/exchangeos-module` + `gh pr create`. Cert generation via `make new-module NAME=exchangeos` + Vault seed + staging deploy + go-live-checklist sign-off all remain operator-driven. |

## Why this milestone

`docs/operations/crdb-hub-tls-pr.md` specs the cross-repo PR for `cockroachdb/modules/exchangeos/` (database.sql + users.sql + cert template + DSN template + acceptance criteria + rollback). The PR has **not been opened**. Without it, ExchangeOS cannot connect to the shared CRDB hub in production — it remains stuck on docker-compose CRDB or per-pod sidecars (both unacceptable per CLAUDE.md "Shared CRDB hub TLS desde dia 1").

## Description

Open the cross-repo PR in `cockroachdb` repo per the spec, drive it to merge through the shared-infra review process, then verify end-to-end connectivity from a staging ExchangeOS pod.

## Acceptance Criteria

- [ ] PR opened in `cockroachdb` repo with title `exchangeos: add module registration with shared CA TLS` containing exactly the files from `docs/operations/crdb-hub-tls-pr.md`
- [ ] Reviewed + approved by 2 shared-infra maintainers
- [ ] Cert generated from shared CA + stored in Vault `secret/data/exchangeos/db`
- [ ] DSN updated in `.env.example` + Vault prod secret with `sslmode=verify-full` + path to shared CA
- [ ] Staging ExchangeOS api pod connects successfully (`SELECT version()` returns CRDB version + tenant scope works)
- [ ] Rollback procedure tested in staging (revert PR + confirm graceful degradation pattern documented)
- [ ] Update `docs/operations/go-live-checklist.md` § "CRDB connectivity" — change ⏳ to ✅
- [ ] Update `.base/plans/milestones/delivered/MS-023a-foundation-scaffolding.md` deferred-items section noting this milestone closed it

## Deliverables

- Merged cross-repo PR (link captured in delivery notes)
- Vault path `secret/data/exchangeos/db` populated in production
- `.env.example` documenting prod DSN format
- Staging proof-of-life log lines captured + linked in delivery notes
- Updated `crdb-hub-tls-pr.md` marked "executed YYYY-MM-DD, PR #NNN"

## Cross-References

- `docs/operations/crdb-hub-tls-pr.md` — PR spec
- `CLAUDE.md` — "Shared CRDB hub TLS desde dia 1" rule
- MS-023a foundation-scaffolding deferred items
- `scripts/vault-seed.sh` — credential seeding
