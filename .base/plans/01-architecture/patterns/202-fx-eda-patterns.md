# FX-EDA-* — Event-Driven Architecture Patterns (45 patterns)

ExchangeOS asynchronous-messaging patterns.

## Catalog (representative)

| # | Title | Status |
|---|-------|--------|
| FX-EDA-001 | Transactional Outbox | ✅ |
| FX-EDA-002 | At-least-once delivery + consumer idempotency | ✅ |
| FX-EDA-003 | In-process bus for cross-context handoff | ✅ |
| FX-EDA-004 | Event naming `<context>.<action>.v<N>` | ✅ |

---

## FX-EDA-001 — Transactional Outbox

**Context:** Persist aggregate state AND publish events without 2PC.

**Problem:** Inline Kafka publish in the same Save can fail after the DB commits, losing the event. Or vice-versa, the event ships but the aggregate write rolls back.

**Solution:** Write aggregate + outbox row in ONE DB transaction. A separate worker (cmd/worker) polls outbox_events and publishes to Kafka, then marks dispatched. At-least-once delivery; consumers de-dupe.

**Example:** `migrations/000009_create_outbox.up.sql` + `pkg/outbox/postgres/store.go` + `pkg/outbox/kafka/publisher.go` (under -tags kafka) + `cmd/worker/main.go` dispatch loop.

**Anti-pattern:** Calling `kafkaClient.Publish` inline in a repo Save method.

**Related:** FX-EDA-002, FX-CP-010, FX-CP-004 (partial index on pending).

---

## FX-EDA-002 — At-least-once delivery + consumer idempotency

**Context:** Outbox can publish a message twice (worker crashes between Publish and MarkDispatched).

**Problem:** Consumers double-process events → duplicate trades, double IOF, etc.

**Solution:** Consumers MUST be idempotent. Use the `outbox_id` (UUID) as the de-dupe key persisted in a `processed_events` table per consumer; reject on duplicate.

**Example:** Compliance worker maintains `compliance_processed_events (outbox_id PK)` — `INSERT ... ON CONFLICT (outbox_id) DO NOTHING; if rowcount==0 then already processed; SKIP`.

**Anti-pattern:** Relying on Kafka exactly-once semantics across producer + consumer + sink — too brittle for FX.

**Related:** FX-EDA-001, FX-KP-005 (consumer group config).

---

## FX-EDA-003 — In-process bus for cross-context handoff

**Context:** Two bounded contexts in the same process need an event dispatch (e.g. Quote → Trade) before Kafka outbox is wired.

**Problem:** Direct service-to-service calls violate DDD context boundaries; full Kafka adds latency + ops burden for in-process flows.

**Solution:** `internal/eventbus.Bus` — synchronous Subscribe/Publish, thread-safe, handlers run inline. Application services use eventbus-backed publishers; cross-context handlers register via container `wireEventHandlers()`.

**Example:** `internal/eventbus/eventbus.go` + `eventbus.QuotePublisher` + `tradeapp.QuoteAcceptedHandler` wired by `container.wireEventHandlers()` for `quote.accepted.v1`.

**Anti-pattern:** Hard-call from QuoteService.AcceptQuote into TradeService.BookTrade.

**Related:** FX-EDA-001 (later swap to outbox), FX-DDD-004.

---

## FX-EDA-004 — Event naming `<context>.<action>.v<N>`

**Context:** Event names need to be stable across consumers.

**Problem:** Renaming breaks downstream; missing version makes upgrades painful.

**Solution:** Format `<context>.<action>.v<N>` (lowercase). e.g. `trade.created.v1`, `quote.accepted.v1`, `cls_cycle.opened.v1`. v2 ships alongside v1 during migration.

**Example:** `modules/trade/domain/events.go` — every Event impl exposes `EventName() string` returning the canonical name.

**Anti-pattern:** `TradeCreated` (PascalCase, no version) — every breaking change forces all consumers to redeploy in lockstep.

**Related:** FX-EDA-001, FX-KP-002 (topic naming).
