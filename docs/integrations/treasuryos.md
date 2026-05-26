# Integration — ExchangeOS ↔ TreasuryOS

> Owner: Treasury team
> Compatible since: future (TreasuryOS not in initial delivery)
> Status: ⏳ Out of scope for initial release

## Purpose

TreasuryOS is the platform's liquidity + funding manager — it tracks nostro
balances across correspondent banks, projects funding needs, manages collateral
pools, and orchestrates intraday liquidity. ExchangeOS is one of its biggest
liquidity consumers (every settled FX trade moves money on nostros) and one of
its biggest sources of FX exposure to hedge.

## Direction (future)

```
ExchangeOS ── settlement.payment_required.v1 ──▶ TreasuryOS  (funding need)
            ── position.snapshot.v1 ───────────▶              (hedge book)
            ── cls.payin_required.v1 ─────────▶              (CLS USD funding)
                                                       │
                                                       ▼
                                       reserves liquidity / arranges hedge
                                                       │
TreasuryOS ── liquidity.unavailable.v1 ─────────▶ ExchangeOS  (block large trade)
            ── nostro.balance_snapshot.v1 ────▶              (refdata refresh)
            ── hedge.proposal.v1 ─────────────▶              (auto-cover NDF)
```

## Events ExchangeOS would produce

### `settlement.payment_required.v1`

For each leg of a non-CLS gross settlement, ExchangeOS signals the future
payment amount + value date so TreasuryOS can plan nostro funding.

```json
{
  "settlement_id": "uuid",
  "currency": "USD",
  "amount": "1000000.00",
  "value_date": "2026-05-28",
  "nostro_hint": "JPMUS33",
  "tenant_id": "uuid"
}
```

### `cls.payin_required.v1`

Per-CCY PayIn requirement after CLS NetReport. TreasuryOS must arrange the
funding to land in the CLS account before the deadline (RTGS cut-off).

### `position.snapshot.v1`

Daily EOD position per (tenant, currency) — TreasuryOS uses this to compute
the hedge book.

## Events ExchangeOS would consume

### `liquidity.unavailable.v1`

TreasuryOS warns that a proposed trade would breach a nostro intraday limit.
ExchangeOS surfaces as warning (advisory) for spot but may HARD-BLOCK
for forwards > USD 10M.

### `nostro.balance_snapshot.v1`

Hourly snapshot of nostro balances per (currency, correspondent). Used to
populate the `settlement.nostro_hint` field on new trades.

### `hedge.proposal.v1`

Treasury proposes an offsetting trade to neutralise FX risk. ExchangeOS may
auto-book if signed by Treasury role + within pre-approved limits.

## Sync RPCs

### ExchangeOS → TreasuryOS

```protobuf
service LiquidityQuery {
  rpc CheckFunding(CheckFundingRequest) returns (CheckFundingResponse);
  rpc ReserveLiquidity(ReserveLiquidityRequest) returns (ReserveLiquidityResponse);
}
```

Pre-trade check for large notionals. 200ms p99 timeout — if Treasury times out,
fall back to local nostro-hint cache (24h staleness).

### TreasuryOS → ExchangeOS

```protobuf
service PositionQuery {
  rpc GetCurrentPosition(...) returns (...);
  rpc GetFXBookForHedging(...) returns (...);
}
```

## Failure semantics

- **TreasuryOS down:** ExchangeOS continues for spot < USD 1M; rejects new
  large trades (USD > 10M) with HTTP 503 + Retry-After
- **Stale nostro hints:** mitigated by 24h fallback + EOD reconciliation
- **Missed PayIn signal:** detected by CLS PayIn deadline alert (3 deadlines/cycle); CRITICAL incident

## Open questions

- [ ] When will TreasuryOS exist?
- [ ] Does Treasury auto-book hedges or always require human sign-off?
- [ ] Should `liquidity.unavailable` be a hard block or advisory? (operational policy decision)
- [ ] CLS PayIn ownership: Treasury arranges funding but ExchangeOS owns the camt.061 submission — need clear hand-off
