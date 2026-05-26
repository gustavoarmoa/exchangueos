# MS-024h — Complete Postgres Repository Layer

| Field | Value |
|-------|-------|
| **Code** | MS-024h |
| **Name** | complete-postgres-repos |
| **Phase** | F-OPS-PROD |
| **Sprint** | 1 of MS-024 cycle |
| **Status** | ACTIVE |
| **Owner** | Platform |
| **Dependencies** | None |

## Why this milestone

Six bounded contexts ship with only in-memory repositories:
- `modules/payin/` (memory only)
- `modules/netreport/` (memory only)
- `modules/compliance/` (memory only — Classification, IOFComputation, BACENReport, Screening)
- `modules/admin/` (memory only — SystemEvent, EODJob)
- `modules/cfets_capture/` (memory only)
- `modules/cfets_confirmation/` (memory only)

Running on memory means restart = data loss + no durability + no real concurrency story. Production needs all 14 BCs persisted.

## Description

Add `postgres.*Repository` implementations for the 6 BCs above. Each follows the established pattern (see `modules/refdata/infrastructure/postgres/repos.go` + `modules/cls_settlement/infrastructure/postgres/repos.go`): pgxpool + Reconstitute helper + UPSERT on Save + tx scoping for compound writes. Migrations already exist (000006 settlement, 000007 risk/position, 000008 compliance/admin, plus need 000011 cos_cases, 000012 cfets if not covered).

## Acceptance Criteria

- [ ] `modules/payin/infrastructure/postgres/repos.go` — `PayInRepo` against `payin_instructions` table
- [ ] `modules/netreport/infrastructure/postgres/repos.go` — `NetReportRepo` against `net_reports`
- [ ] `modules/compliance/infrastructure/postgres/repos.go` — 4 repos (Classification, IOF, BACENReport, Screening) against migration 000008 tables
- [ ] `modules/admin/infrastructure/postgres/repos.go` — `SystemEventRepo` + `EODJobRepo`
- [ ] `modules/cfets_capture/infrastructure/postgres/repos.go` — `CFETSCaptureRepo`
- [ ] `modules/cfets_confirmation/infrastructure/postgres/repos.go` — `CFETSConfirmationRepo`
- [ ] If schemas missing: migration 000012 covering cfets_captures + cfets_confirmations
- [ ] `internal/container/container.go::wirePostgres` extended to wire all 6 with backend switch
- [ ] Integration tests per repo against testcontainers CRDB (8 tests per repo: insert + get + update + list + concurrent-version-bump + filter + paging + tx-rollback)
- [ ] Pre-existing memory repo tests pass against postgres backend via dual-execution test harness
- [ ] Performance: each repo Save < 10ms p99 against local CRDB (benchmark added)

## Deliverables

- 6 × `repos.go` files
- 6 × `repos_test.go` integration tests
- `migrations/000012_create_cfets.up.sql` + `.down.sql` (if needed)
- Reconstitute helpers per aggregate in domain packages
- Updated container wiring
- Updated `.env.example` documenting full postgres backend coverage

## Cross-References

- Existing reference impls: `modules/refdata/infrastructure/postgres/`, `modules/cls_settlement/infrastructure/postgres/`
- `.claude/rules/migrations.md` — migration conventions
- ISO 27001 control 8.13 (information backup) — durability prerequisite
- MS-024a + MS-024b consume some of these repos
