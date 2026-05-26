# ExchangeOS Migrations

CockroachDB v24.3.32 migrations executed by `cmd/migrator` via shared hub TLS
(`cockroachdb/modules/exchangeos/`).

## Conventions

- Filename: `NNNNNN_description.{up,down}.sql` (golang-migrate format)
- Range 000001-000020 reserved for MS-023a foundation
- All migrations wrapped in `BEGIN`/`COMMIT` for atomicity
- All `CREATE TABLE` use `IF NOT EXISTS` for idempotence
- Money/Rate columns: `DECIMAL(36,18)` — NEVER `FLOAT`/`DOUBLE`
- Currency codes: `STRING(3)` (ISO 4217 alpha-3)
- Country codes: `STRING(2)` (ISO 3166 alpha-2)
- BIC: `STRING(11)` (ISO 9362)
- LEI: `STRING(20)` (ISO 17442)
- IDs: `UUID` with `gen_random_uuid()` default
- Timestamps: `TIMESTAMPTZ NOT NULL DEFAULT current_timestamp()`

## Roadmap (000001-000020)

| # | File | Scope |
|---|------|-------|
| 000001 | `create_tenants` | tenants, actors, audit_events, schema_migrations |
| 000002 | `create_fx_trades` | counterparties, fx_trades, trade_amendments |
| 000003 | `create_quotes` | quotes, quote_streams |
| 000004 | `create_settlement` | cls_cycles, payin_instructions, net_reports |
| 000005 | `create_refdata` | currencies, calendars, bic_records, ssis |
| 000006 | `create_risk` | limits, exposures, breach_log |
| 000007 | `create_positions` | positions, position_snapshots |
| 000008 | `create_compliance` | classifications, iof_computations, bacen_reports, screening_results |
| 000009 | `create_admin` | system_events, eod_jobs |
| 000010 | `create_cfets` | cfets_captures, cfets_confirmations |
| 000011-000020 | reserved | indexes, materialized views, RBAC tables |

## Running

```bash
task db:migrate          # apply pending up migrations
task db:reset CONFIRM=yes # roll back all + reapply (destructive)
```
