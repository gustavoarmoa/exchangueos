# ExchangeOS Proto Contracts

9 services under `exchangeos.v1`:

| Service | Bounded Context | ISO 20022 Mapping |
|---------|-----------------|-------------------|
| `TradeService` | trade | fxtr.014/015/016 |
| `QuoteService` | quote | internal (no ISO mapping) |
| `AmendmentService` | amendment | internal (no ISO mapping) |
| `SettlementService` | cls_settlement + payin + netreport | camt.061/062/063/088 |
| `RefDataService` | refdata | reda.060/061 |
| `AdminService` | admin | admi.002/004/009/010/011/017 |
| `RiskService` | risk | internal |
| `PositionService` | position | internal |
| `ComplianceService` | compliance | BACEN Lei 14.286 + Circulares |

## Generation

```bash
task proto:lint      # buf lint
task proto:gen       # buf generate → proto/gen/...
task proto:breaking  # buf breaking (vs main)
```

## Conventions

- **Money/Rate:** `string` decimal — NEVER `float`/`double`. Server unmarshals to `shopspring/decimal`.
- **IDs:** UUIDv7 strings.
- **Tenant context** required on every request (FX-API-001).
- **Audit envelope** present on every persisted message.
- **Pagination:** cursor-based via `PageRequest`/`PageResponse`.
- **Errors:** canonical `ErrorCode` enum + `google.rpc.Status` details.

## Versioning

- Major versions = package suffix (`v1`, `v2`).
- Minor changes: backward-compatible additions only.
- Breaking changes: gated by `buf breaking` in CI.
