# ADR 0003 — Transactional outbox for async event publication

- Status: Accepted
- Date: 2026-05-24

## Context

When an aggregate is saved (e.g. FXTrade.Confirm), we want downstream consumers (LedgerOS, AuthorityOS) to react. Two simple-looking options have well-known failure modes:

- **Inline publish in repo.Save** — if DB commits then Kafka publish fails, the event is lost. If Kafka publish succeeds then DB rolls back, consumers see a ghost event.
- **Distributed 2PC across DB + Kafka** — not supported by either side; would require an XA coordinator.

## Decision

**Transactional outbox pattern.**

1. `Repository.Save` writes aggregate state AND outbox rows in the **same DB transaction**.
2. A separate worker (`cmd/worker`) polls `outbox_events WHERE dispatched_at IS NULL`, publishes to Kafka, marks dispatched.
3. **At-least-once** delivery — consumers MUST be idempotent (use `outbox_id` UUID as dedupe key).

Schema in `migrations/000009_create_outbox.up.sql` + Go in `pkg/outbox/` + postgres impl + kgo-backed Kafka publisher (build tag `kafka`).

## Consequences

### Positive

- **No 2PC** — single DB transaction is atomic
- **Survives broker outages** — events queue in `outbox_events` until broker recovers
- **Replayable** — `outbox_dispatched_archive` table retains 30 days for audit + replay
- **Decoupled** — module code uses `outbox.Store` interface; broker choice deferred to `pkg/outbox/<adapter>`

### Negative

- **At-least-once requires consumer idempotency** — non-trivial discipline; mitigated by `outbox_id`-based dedupe pattern
- **Polling adds latency** — typical 100ms..1s between commit and dispatch (vs 0ms inline). Acceptable for our SLOs.
- **Storage overhead** — outbox table grows; mitigated by 30-day archive + nightly truncate of archive

## Alternatives considered

- **Inline publish** — rejected: lost-event risk fatal for financial events
- **Change Data Capture (CDC) on tables** — viable but tightly couples consumers to internal schema; chose explicit outbox for contract stability
- **Event sourcing** — too disruptive for this codebase + team's experience; outbox gives us most of the benefits
