# FX-CP-* — CockroachDB Patterns (50 patterns)

ExchangeOS CRDB usage patterns.

## Catalog (representative)

| # | Title | Status | Where |
|---|-------|--------|-------|
| FX-CP-001 | DECIMAL(36,18) for money/rate | ✅ | migrations/000002 (fx_trades), 000007 (positions) |
| FX-CP-002 | `gen_random_uuid()` PK (avoid hot-spotting) | ✅ | every CREATE TABLE |
| FX-CP-003 | `TIMESTAMPTZ NOT NULL DEFAULT current_timestamp()` | ✅ | every created_at/updated_at |
| FX-CP-004 | Partial index for hot subset | ✅ | `idx_outbox_pending WHERE dispatched_at IS NULL` |
| FX-CP-005 | Composite tenant-scoped UNIQUE | ✅ | `UNIQUE (tenant_id, business_date)` on eod_jobs |
| FX-CP-006 | CHECK constraints for enums | ✅ | `status IN (...)` everywhere |
| FX-CP-007 | Shared CRDB hub TLS (no inline --insecure) | ⏳ | cockroachdb/modules/exchangeos (cross-repo) |
| FX-CP-008 | pgx.Pool with MaxConnLifetime 30min | ✅ | internal/db/pool.go |
| FX-CP-009 | `ON CONFLICT DO UPDATE` upsert pattern | ✅ | repos.Save in postgres impls |
| FX-CP-010 | Transactional outbox via Begin/Commit | ✅ | cls_settlement/infrastructure/postgres/repos.go |
| FX-CP-011..050 | (extend on demand) | ⏳ | — |

---

## FX-CP-001 — DECIMAL(36,18) for money/rate

**Context:** Money + rate columns.

**Problem:** Floating-point types (REAL/DOUBLE) silently lose precision. INT can't represent fractional amounts cleanly.

**Solution:** `DECIMAL(36,18)` — 36 total digits with 18 after the point. Sufficient for IDR (no decimals, large notionals) through JPY (0 decimals) and BHD (3 decimals) edge cases.

**Example:**

```sql
fx_trades.bought_amount  DECIMAL(36,18) NOT NULL CHECK (bought_amount > 0)
fx_trades.deal_rate      DECIMAL(36,18) NOT NULL CHECK (deal_rate > 0)
```

Go side maps to `shopspring/decimal.Decimal` via pgx scan.

**Anti-pattern:** `NUMERIC(20,4)` — too narrow for cumulative gross/net in 8-digit-precision FX rates.

**Related:** FX-GP-002 (Go decimal), FX-CP-009 (Save pattern).

---

## FX-CP-004 — Partial index for hot subset

**Context:** Most queries hit a small subset of a large table.

**Problem:** Full B-tree index wastes RAM + write amplification.

**Solution:** `WHERE <predicate>` clause on `CREATE INDEX` — index only the rows that match.

**Example:**

```sql
CREATE INDEX idx_outbox_pending ON outbox_events (occurred_at)
    WHERE dispatched_at IS NULL;

CREATE INDEX idx_limits_breaching ON risk_limits (tenant_id)
    WHERE utilised >= cap;
```

**Anti-pattern:** Full index on `dispatched_at` with most rows non-null (dispatched).

**Related:** FX-CP-010 (transactional outbox), FX-CP-005 (composite uniqueness).

---

## FX-CP-002 — `gen_random_uuid()` PK avoids hot-spotting

**Context:** CRDB ranges split by primary-key range.

**Problem:** Monotonic PKs (SERIAL, time-prefixed) concentrate writes on a single hot range.

**Solution:** UUIDv4 via `gen_random_uuid()` distributes inserts across the keyspace.

**Example:** Every CREATE TABLE in `migrations/` uses `<id> UUID PRIMARY KEY DEFAULT gen_random_uuid()`. Time-ordered access uses a separate indexed `created_at` column.

**Anti-pattern:** `id SERIAL PRIMARY KEY`.

---

## FX-CP-008 — pgxpool with bounded lifetime

**Context:** Long-running services + DB-side connection limits.

**Problem:** Stale connections accumulate behind NAT timeouts + rolling updates; unbounded pools breach `max_connections`.

**Solution:** `pgxpool.Config` with `MaxConnLifetime: 30m` + `MaxConnIdleTime: 5m`. Pool size bounded by `Max/MinConn` from config.

**Example:** `internal/db/pool.go:New(ctx, cfg)` builds the pool + pings on connect.

**Anti-pattern:** Default `pgxpool.New` without lifetime caps.

---

## FX-CP-009 — `ON CONFLICT DO UPDATE` upsert idiom

**Context:** Aggregate Save where repository can't be sure whether the row exists.

**Problem:** SELECT-then-INSERT-or-UPDATE is two round-trips + a race window.

**Solution:** UPSERT — `ON CONFLICT (pk) DO UPDATE SET col = EXCLUDED.col, version = EXCLUDED.version, updated_at = current_timestamp()`. CRDB executes atomically.

**Example:** `modules/risk/infrastructure/postgres/repos.go:LimitRepo.Save` bumps utilised + version in one statement.

**Anti-pattern:** `BEGIN; SELECT; INSERT-or-UPDATE; COMMIT;` — slower + race-prone.

**Related:** FX-CP-010 (transactional outbox), FX-CP-005 (composite uniqueness).
