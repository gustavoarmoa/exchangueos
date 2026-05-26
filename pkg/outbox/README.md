# pkg/outbox — Transactional Outbox Pattern

At-least-once event delivery without distributed transactions.

## Architecture

```
┌────────────────────┐  same DB tx   ┌──────────────┐
│ Repository.Save    │ ─────────────▶│ outbox_events│
│ aggregate + outbox │               └──────┬───────┘
└────────────────────┘                       │
                                             │ poll (cmd/worker)
                                             ▼
                                  ┌──────────────────────┐
                                  │ outbox.Dispatch loop │
                                  │  store.Pending(...)  │
                                  │  pub.Publish(...)    │
                                  │  store.MarkDispatched│
                                  └──────────┬───────────┘
                                             │
                                             ▼
                                          Kafka
```

## Integration

1. Aggregate `Save` writes both state + outbox row in one transaction:

   ```go
   func (r *TradeRepo) Save(ctx context.Context, t *domain.FXTrade) error {
       tx, _ := r.pool.Begin(ctx); defer tx.Rollback(ctx)
       // INSERT INTO fx_trades ...
       // INSERT INTO outbox_events ... (one row per t.PendingEvents())
       return tx.Commit(ctx)
   }
   ```

2. `cmd/worker` runs a loop:

   ```go
   for {
       n, err := outbox.Dispatch(ctx, store, kafkaPub, 100)
       if err != nil { logger.Warn(...) }
       if n == 0 { time.Sleep(500*time.Millisecond) }
   }
   ```

3. Real Kafka client lives in `pkg/outbox/kafka` (next iteration; choose between
   Sarama / franz-go / kgo — `kgo` recommended for performance).

## Schema

See `migrations/000009_create_outbox.up.sql`:

- `outbox_events` — pending rows + per-row attempt counter + last_error.
- `outbox_dispatched_archive` — optional historical view; worker may archive after retention.

Indexes:

- `idx_outbox_pending` (partial WHERE dispatched_at IS NULL) — drives the worker loop.
- `idx_outbox_aggregate` — debug + aggregate-level audit.
- `idx_outbox_failed` — observability for stuck rows.

## Why?

- **No 2PC** — same DB tx covers aggregate + outbox; broker is decoupled.
- **At-least-once** — workers retry until MarkDispatched succeeds. Consumers
  must be idempotent (use event_name + outbox_id as dedupe key).
- **Replayable** — historical events live in the archive table.
