# Integration — ExchangeOS ↔ RiskOS

> Owner: Risk team
> Compatible since: future (RiskOS not in initial delivery)
> Status: ⏳ Out of scope for initial release

## Purpose

RiskOS is a platform-wide risk aggregation layer that consumes per-module
exposure events and computes cross-module risk metrics (concentration limits
across multiple business lines, group-level VaR, regulatory capital).

ExchangeOS owns FX-specific risk (counterparty + currency + tenor + DV01 + VaR
per-trade) — see `modules/risk/`. RiskOS would consume our position +
breach events for cross-module aggregation.

## Direction (future)

```
ExchangeOS ── risk.breach.v1 ─────────▶ RiskOS
            ── position.snapshot.v1 ──▶
                                            │
                                            ▼
                                aggregates across modules
                                            │
RiskOS     ── group_risk.limit_pressure.v1 ─▶ ExchangeOS  (warning, not gate)
```

## Events ExchangeOS would produce

### `risk.breach.v1`

Already emitted to `exchangeos.risk.events` topic (see `deploy/kafka/topics.yaml`):

```json
{
  "limit_id": "uuid",
  "limit_type": "COUNTERPARTY|CURRENCY|TENOR|DV01|VAR",
  "scope": "DEUTDEFF",
  "tenant_id": "uuid",
  "proposed_exposure": "1000000.00",
  "cap": "10000000.00",
  "utilised": "9500000.00",
  "blocked_trade_id": "uuid"
}
```

### `position.snapshot.v1` (planned)

Emitted by `cmd/eod` after the daily snapshot — one event per (tenant, currency).

## Events ExchangeOS would consume

### `group_risk.limit_pressure.v1` (planned)

RiskOS warns that ExchangeOS's contribution to a group-level limit is approaching
the cap. ExchangeOS surfaces this as an admin.system_event but does NOT
auto-block trades — group risk is a soft warning.

## Sync RPCs

None initially. RiskOS reads via topics + has its own read replica for ad-hoc
queries.

## Failure semantics

- **RiskOS down:** ExchangeOS unaffected — local risk continues to gate trades
- **Group limit breach signal lost:** detected at EOD reconciliation; minor (not regulatory)

## Open questions

- [ ] When will RiskOS exist?
- [ ] Does RiskOS expose a hard-block API (would couple ExchangeOS hot path) or stays advisory?
- [ ] Cross-module limit semantics: sum of per-module utilisation, max, or weighted?
