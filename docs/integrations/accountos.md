# Integration — ExchangeOS ↔ AccountOS

> Owner: Platform team
> Compatible since: ExchangeOS v4.2.0 (tenants table created in 000001)
> Status: ✅ Conceptually wired (tenant SoT in AccountOS via CDC materialised view)

## Purpose

AccountOS is the **single source of truth (SoT)** for tenants + actors across
the Revenu Platform. ExchangeOS holds a read-only mirror via CDC + materialised
view so that every fx_trades row can FK to a stable tenant_id without coupling
to AccountOS at write-path latency.

## Direction

```
AccountOS ──── tenant.created.v1 ────▶ exchangeos.refdata.tenants (Kafka topic)
            ── tenant.updated.v1 ───▶
            ── actor.created.v1 ──▶
            ── actor.disabled.v1 ─▶
                                            │
                                            ▼
                               (cmd/worker materialises into
                                exchangeos.tenants table —
                                CDC pattern, eventually consistent)
```

## Events ExchangeOS consumes

| Event | Action |
|-------|--------|
| `tenant.created.v1`  | `INSERT INTO tenants ... ON CONFLICT (code) DO NOTHING` |
| `tenant.updated.v1`  | `UPDATE tenants SET name/country/status WHERE tenant_id = ...` |
| `tenant.suspended.v1`| `UPDATE tenants SET status = 'SUSPENDED'` — application service refuses new trades |
| `actor.created.v1`   | INSERT into actors with external_sub (OIDC subject) |
| `actor.disabled.v1`  | UPDATE actors SET status = 'DISABLED' — JWT validation still passes but no actions allowed |

Schema authority lives in AccountOS; ExchangeOS uses Reconstitute helpers for
backward-compat when fields are added.

## Events ExchangeOS produces

None — ExchangeOS is read-only WRT tenants + actors.

## Sync RPCs (ExchangeOS → AccountOS)

Fall-back only: when a tenant_id is referenced in a trade booking but is not
in the local mirror (CDC lag or first activation), exchangeos-api may call:

```protobuf
service AccountSync {
  rpc ResolveTenant(ResolveTenantRequest) returns (Tenant);
}
```

5s timeout + 3 retries + circuit breaker. Result cached for 5 min in-process.

## Failure semantics

- **CDC lag > 5 min:** alert ops. Trades for unmirrored tenants fall through
  to the sync RPC, which may slow request latency.
- **AccountOS down:** previously-mirrored tenants continue to work. New tenant
  onboarding stalls until AccountOS recovers.
- **Tenant suspended:** application services check `tenants.status = 'ACTIVE'`
  on every BookTrade — suspended tenants get `403 PERMISSION_DENIED`.

## Schema mapping

| AccountOS field | ExchangeOS column | Notes |
|-----------------|-------------------|-------|
| `account_id`    | `tenants.tenant_id` (UUID) | Same UUID across platform |
| `account_code`  | `tenants.code` (UNIQUE)    | Human-friendly identifier |
| `legal_name`    | `tenants.name`             | |
| `jurisdiction`  | `tenants.country` (ISO 3166 alpha-2) | Drives BACEN classification rules |
| `lifecycle_status` | `tenants.status`        | ACTIVE \| SUSPENDED \| ARCHIVED |

## Open questions

- [ ] Concrete Kafka topic name in AccountOS: is it `accountos.tenant.events` or `revenu.tenant.events`?
- [ ] Are deleted tenants soft-delete (status=ARCHIVED) or hard-delete? (we assume soft)
- [ ] Multi-region: which region owns the canonical SoT?
- [ ] LGPD right-to-erasure: how do we handle `actor.deleted.v1` while preserving 7-year audit retention?
