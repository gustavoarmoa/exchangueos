# Integration — ExchangeOS ↔ PaymentOS

> Owner: Platform team
> Compatible since: ExchangeOS v4.10.0 (Settlement) + PaymentOS TBD
> Status: 🟡 Spec — PaymentOS side pending

## Purpose

CLS settlement requires PvP (Payment-versus-Payment): both legs of an FX trade
must settle atomically. ExchangeOS orchestrates the CLS cycle (open/payin/settle/
close); PaymentOS executes the actual wire transfers to CLS settlement member
accounts at the deadline windows.

## Direction

```
ExchangeOS ── settlement.payin_requested.v1 ─▶ PaymentOS
                                                    │
                                                    ▼
                                          executes wire transfer
                                                    │
PaymentOS  ── payment.settled.v1 / .failed.v1 ─▶ ExchangeOS
                                                    │
                                                    ▼
                                      PayInInstruction.Confirm() / .Fail()
```

## Events ExchangeOS produces

### `settlement.payin_requested.v1`

Emitted when a PayInInstruction is created + submitted (within the CLS deadline).

```json
{
  "event_id": "uuid",
  "occurred_at": "RFC3339",
  "tenant_id": "uuid",
  "instruction_id": "uuid",
  "cycle_id": "uuid",
  "currency": "USD",
  "amount": "1080000.00",
  "band": "PIN3",
  "deadline": "2026-05-26T10:00:00+02:00",
  "beneficiary_bic": "CLSBUS33",
  "ssi_account": "..."
}
```

Topic: `exchangeos.payin.events`
Key: `cycle_id`

## Events ExchangeOS consumes

### `payment.settled.v1`

PaymentOS confirms the wire cleared.

Effect: `PayInService.Confirm(instruction_id, at)`. If after deadline, the
deadline was already missed and instruction is in FAILED — Confirm is rejected
with `ErrInvalidTransition`.

### `payment.failed.v1`

PaymentOS reports failure (insufficient funds, beneficiary rejected, network).

Effect: `PayInService.Fail(instruction_id, at, reason)`. Triggers CLS cycle
investigation; may cascade to `cls_cycle.failed.v1`.

## Sync RPCs

### ExchangeOS → PaymentOS

```protobuf
service PaymentExecution {
  // Hot-path PvP commit — called by SettlementService when CLS cycle enters SETTLING.
  rpc CommitPvP(CommitPvPRequest) returns (CommitPvPResponse);
}
```

10s timeout (PvP must be atomic and fast). On error: cycle moves to FAILED.

### PaymentOS → ExchangeOS

None. PaymentOS publishes events; ExchangeOS subscribes.

## Failure semantics

- **CommitPvP timeout:** cycle FAILED + emergency Slack alert. Manual remediation by Treasury.
- **payment.settled.v1 lost:** outbox retry on PaymentOS side. ExchangeOS reconciles via daily SQL query.
- **Concurrent CLS cycles:** PvP is per-cycle; no concurrent execution within the same `cycle_id`.

## Open questions

- [ ] PvP atomicity contract: who is the "authority of last resort" — CLS, PaymentOS, or ExchangeOS?
- [ ] Partial settle: if 6 of 8 legs clear, does the whole cycle abort or proceed?
- [ ] Settlement currency precedence: USD as common pivot, or pair-direct?
- [ ] PaymentOS rate limit on CommitPvP — coordinate with CLS PIN deadlines
