# Cross-Repo PR: `cockroachdb/modules/exchangeos/` Hub TLS Registration

> Target repo: `revenu-tech/cockroachdb` (shared CRDB hub)
> Target path: `cockroachdb/modules/exchangeos/`
> Required before: production cluster bootstrap
>
> **Status: PREPARED 2026-06-19** — local branch `feat/exchangeos-module`
> commit `b3d0689` in `revenu-platform/cockroachdb/` carries the
> scaffold (`modules/exchangeos/docker-compose.yml` + `INTEGRATION.md`).
> PR opening is blocked on remote configuration for the cockroachdb
> repo; once the remote exists, `git push -u origin feat/exchangeos-module`
> + `gh pr create` finishes the loop. Track via MS-024l.

## Note on this spec vs the delivered shape

The "What to add" block below describes an early-pass vision (database.sql
+ users.sql + Taskfile.yml + env/{dev,staging,production}.env). The
**delivered shape** in cockroachdb commit `b3d0689` follows the actual
convention used by `modules/authorityos/` and `modules/onboardos/`:

- `modules/exchangeos/docker-compose.yml` — verbatim copy of
  `_template/docker-compose.yml` (reads MODULE from per-module `.env`
  which is gitignored + scaffolded by `make new-module NAME=exchangeos`).
- `modules/exchangeos/INTEGRATION.md` — operator guide covering host /
  Docker integration, DSN format, migrations, the read-only
  `exchangeos_auditor` role, and Vault path for the prod cert.

Per-module `database.sql` and `users.sql` files are not the hub repo
convention — the database name is inferred from `MODULE`, and role
grants are scripted into INTEGRATION.md so the shared-infra team can
review them without diving into SQL files. If the shared-infra team
prefers explicit SQL files, the spec below can be revived during PR
review.

## Why

ExchangeOS uses the **shared CRDB hub TLS** pattern: every Revenu Platform module
registers under `cockroachdb/modules/<module-name>/` with its own database +
TLS user/cert + role grants. This avoids per-module CRDB clusters AND keeps
every connection mTLS-authenticated (NEVER `--insecure`).

## What to add in `cockroachdb/`

```
cockroachdb/modules/exchangeos/
├── README.md
├── database.sql           # CREATE DATABASE exchangeos + grants
├── users.sql              # CREATE USER exchangeos_app + grants + cert subject mapping
├── tls/
│   ├── exchangeos.crt     # client cert (issued by hub CA)
│   ├── exchangeos.key     # client key (managed via SOPS / age-encrypted)
│   └── README.md          # rotation cadence + revocation procedure
├── Taskfile.yml           # task up / task migrate / task psql / task rotate-cert
└── env/
    ├── dev.env            # local compose DSN (uses --insecure, single-node)
    ├── staging.env        # staging cluster DSN
    └── production.env     # prod cluster DSN — verify-full + sslrootcert
```

## `database.sql`

```sql
-- Tenant DB for exchangeos.
CREATE DATABASE IF NOT EXISTS exchangeos;

-- Optional: per-tenant logical schemas would live here if multi-tenancy
-- migrates from row-level to schema-level (FX-CP-* future expansion).
```

## `users.sql`

```sql
CREATE USER IF NOT EXISTS exchangeos_app
  WITH PASSWORD NULL;        -- NULL because we authenticate via TLS client cert

-- Cert subject mapping: the client cert's CN must equal the user name.
-- (CRDB enforces this for cert-based auth.)

GRANT CONNECT ON DATABASE exchangeos TO exchangeos_app;
GRANT ALL ON ALL TABLES IN SCHEMA exchangeos.public TO exchangeos_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA exchangeos.public
  GRANT ALL ON TABLES TO exchangeos_app;

-- Read-only auditor role (used by SISBACEN report generator).
CREATE USER IF NOT EXISTS exchangeos_auditor WITH PASSWORD NULL;
GRANT CONNECT ON DATABASE exchangeos TO exchangeos_auditor;
GRANT SELECT ON ALL TABLES IN SCHEMA exchangeos.public TO exchangeos_auditor;
```

## Production DSN template

```
postgres://exchangeos_app@crdb-prod-cluster.internal:26257/exchangeos
  ?sslmode=verify-full
  &sslrootcert=/etc/ssl/crdb-hub-ca.crt
  &sslcert=/etc/ssl/exchangeos.crt
  &sslkey=/etc/ssl/exchangeos.key
```

Secrets sourced via External Secrets Operator from Vault (`secret/data/exchangeos/db`).

## Acceptance criteria for the PR

- [ ] `database.sql` + `users.sql` reviewed by DBA team
- [ ] Cert issued by hub CA with CN=`exchangeos_app` + SAN=`exchangeos_app`
- [ ] `task up` in hub repo produces a working dev cluster
- [ ] `task migrate` against staging cluster applies migrations 000001-000009 successfully
- [ ] DSN template added to `deploy/helm/exchangeos/values-{staging,production}.yaml` via External Secrets reference
- [ ] CRDB metrics scraped by Prometheus + dashboard imported

## Rollback

The new tenant DB + user are additive — rollback = `DROP DATABASE exchangeos
CASCADE` + `DROP USER exchangeos_app, exchangeos_auditor`. The cluster itself
is shared infrastructure; the rollback does not affect sibling modules.
