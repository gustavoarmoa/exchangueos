# FX-DDD-* — Domain-Driven Design Patterns (35 patterns)

ExchangeOS DDD adaptations for the FX bounded contexts.

## Catalog (representative)

| # | Title | Status |
|---|-------|--------|
| FX-DDD-001 | Aggregate Root as single entry point | ✅ |
| FX-DDD-002 | Reference other aggregates by ID only | ✅ |
| FX-DDD-003 | Optimistic concurrency via `version` field | ✅ |
| FX-DDD-004 | Domain events via RecordEvent + outbox flush | ✅ |
| FX-DDD-005..035 | (extend on demand) | ⏳ |

---

## FX-DDD-001 — Aggregate Root as single entry point

**Context:** Domain layer organisation.

**Problem:** Multiple entry points to the same aggregate let invariants slip.

**Solution:** Each aggregate exposes ONE root struct + private fields. All state changes go through methods on the root. Other entities inside the aggregate (value objects, child entities) are not exported.

**Example:** `modules/trade/domain/fxtrade.go:FXTrade` — root with private fields + Confirm/Cancel/MarkSettling/MarkSettled methods. `TradeAmendment` lives as a separate aggregate referenced by trade_id.

**Anti-pattern:** Exposing `*FXTrade.eventStore` for callers to inject events.

**Related:** FX-DDD-002, FX-DDD-004, FX-GP-001.

---

## FX-DDD-002 — Reference other aggregates by ID only

**Context:** Cross-aggregate associations.

**Problem:** Holding pointers to other aggregates creates implicit transactions across boundaries and breaks consistency rules.

**Solution:** Reference other aggregates by `uuid.UUID`, NEVER by pointer. Loading is a separate concern in the application layer.

**Example:** `modules/cls_settlement/domain/cycle.go:CLSCycle.tradeIDs []uuid.UUID` — never `[]*FXTrade`. Application service loads trade details when needed.

**Anti-pattern:** `type Cycle struct { trades []*FXTrade }`.

**Related:** FX-DDD-001, FX-EDA-001 (eventual consistency).

---

## FX-DDD-003 — Optimistic concurrency via `version`

**Context:** Concurrent writes to the same aggregate.

**Problem:** Last-write-wins corrupts state.

**Solution:** Every aggregate has a `version int` field; incremented on each mutation. Postgres UPSERT includes `version = $N` so subsequent writes detect stale state and fail (or retry).

**Example:** `modules/trade/domain/fxtrade.go:FXTrade.version` + `mutate` pipeline increments after each domain method.

**Anti-pattern:** Timestamps as concurrency tokens — they collide under high write load.

**Related:** FX-DDD-001, FX-CP-009 (ON CONFLICT DO UPDATE).

---

## FX-DDD-004 — Domain events via RecordEvent + outbox flush

**Context:** Cross-context reactions (Quote accepted → Trade booked).

**Problem:** Direct calls couple bounded contexts; transactional 2PC is unavailable.

**Solution:** Aggregate records DomainEvent on every state transition via private `recordEvent`. Application Save flushes events to outbox (same DB tx); worker dispatches asynchronously.

**Example:** `modules/quote/domain/quote.go:Quote.Accept` records `EventQuoteAccepted` → outbox → `tradeapp.QuoteAcceptedHandler.Handle` → `BookTrade`.

**Anti-pattern:** Aggregate calls another bounded context's service directly.

**Related:** FX-EDA-001 (Outbox), FX-DDD-001.
