# Integration — ExchangeOS ↔ LedgerOS

> Owner: Platform team
> Compatible since: ExchangeOS v4.10.0 (Trade core delivered) + LedgerOS TBD
> Status: 🟡 Spec — concrete LedgerOS adapter pending other-repo PR

## Purpose

ExchangeOS books FX trades + settles them through CLS / bilateral venues.
LedgerOS is the double-entry general ledger for the platform. Every settled
trade produces 2-4 journal entries representing the cash flows.

## Direction

```
ExchangeOS ──── trade.settled.v1 ────▶ LedgerOS (Kafka topic)
                                            │
                                            ▼
                                  posts JournalEntry
                                            │
                          ledger.posted.v1 ─┘  (optional ack — not required by ExchangeOS)
```

## Events ExchangeOS produces

### `trade.settled.v1`

Emitted when an FXTrade transitions to SETTLED. Payload:

```json
{
  "event_id": "uuid-v4",
  "occurred_at": "RFC3339",
  "tenant_id": "uuid",
  "trade_id": "uuid",
  "settlement_ref": "CLS-…",
  "venue": "CLS|BILATERAL|CFETS",
  "buyer_bic": "DEUTDEFF",
  "seller_bic": "CHASUS33",
  "bought_currency": "EUR",
  "bought_amount":   "1000000.00",
  "sold_currency":   "USD",
  "sold_amount":     "1080000.00",
  "deal_rate":       "1.08000000",
  "value_date":      "2026-05-26"
}
```

Topic: `exchangeos.trade.events` (key = `trade_id`)
Retention: 30 days (default per topic catalogue)
Schema source-of-truth: `modules/trade/domain/events.go:EventTradeSettled`

## Events ExchangeOS consumes from LedgerOS

None (read-mostly via sync RPC when reconciliation needed).

## Sync RPCs (LedgerOS → ExchangeOS)

The exchangeos-api gRPC service does NOT expose a ledger-specific endpoint;
LedgerOS reconciliation reads `audit_events` table via the auditor role
(`exchangeos_auditor` per `crdb-hub-tls-pr.md`).

## Sync RPCs (ExchangeOS → LedgerOS)

Optional post-settle reconciliation:

```protobuf
service LedgerReconciliation {
  rpc GetPostingsForTrade(GetPostingsRequest) returns (GetPostingsResponse);
}
```

ExchangeOS would call this only from the EOD batch (`cmd/eod`) when running
the position snapshot — to validate that the journal matches the trade.

## Failure semantics

- **trade.settled.v1 not delivered:** outbox retries indefinitely + alert
  ops if `attempt_count > 10` (means LedgerOS is consistently failing).
  LedgerOS reconciliation runs nightly to catch missed events.
- **LedgerOS rejects:** ExchangeOS does NOT roll back the trade — settlement
  is authoritative on the CLS side. Mismatches surface as reconciliation
  exceptions, handled manually by Treasury.

## Open questions

- [ ] LedgerOS-side: chart of accounts schema for FX cash/MTM/realised P&L?
- [ ] Multi-currency journal: one entry per CCY leg vs combined?
- [ ] Cross-tenant LedgerOS — shared single ledger or per-tenant? (affects topic partitioning)
- [ ] Reconciliation cadence: post-settle real-time, EOD batch, or both?
