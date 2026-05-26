---
glob: "migrations/*.sql"
---

# Rule: migrations/*.sql

## Convention
- Format: `000NNN_descriptive_name.{up,down}.sql`
- Sempre par up + down (rollback)
- Idempotente: `CREATE TABLE IF NOT EXISTS ...`
- Transactional: wrap em `BEGIN; ... COMMIT;`

## CockroachDB Specifics
- UUID PK via `gen_random_uuid()` (avoid hot-spotting)
- NUMERIC(20,8) para rates, NUMERIC(20,2) para money default
- TIMESTAMPTZ UTC obrigatorio
- ENUMs via `CREATE TYPE`
- FK ON DELETE: CASCADE (filhos) ou RESTRICT (refdata)

## Tenant Scoping
- Toda tabela exceto refdata: `tenant_id UUID NOT NULL`

## Sync with ERDs
- Update `.base/erds/sql/NN-<bc>-ddl.sql` em paralelo
- Update Mermaid ERD em `.base/erds/domain/erd-<bc>-domain.md`
