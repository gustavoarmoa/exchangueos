# MS-023e — risk-position-ledger

| Field | Value |
|-------|-------|
| **Code** | MS-023e |
| **Name** | risk-position-ledger |
| **Phase** | F8 + F10 |
| **Sprint** | 6 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Started** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023d (delivered), MS-023d2 (in-progress, not blocking) |

## Delivery Notes

**Acceptance criteria met:**
- ✅ Limit aggregate (v4.10.0) — 5 types COUNTERPARTY/CURRENCY/TENOR/DV01/VAR, scope required for first 3, Reserve returns ErrBreached without partial commit, Release clamped at zero, SetUtilised for reconciliation + 9 tests; RN_FX_015 cited
- ✅ Risk application service with CheckLimit returning CheckResult{Allowed, BreachedLimits, Explanation} + Reserve/Release + memory repo (v4.10.0) + 5 application tests
- ✅ Position aggregate (v4.10.0) — Long/Short totals + Net signed + IsLong/IsShort/IsFlat helpers + TradeLeg{BUY|SELL} + ApplyTradeLeg + 7 tests
- ✅ Position application service with Get/List/ApplyTradeLeg upsert-on-miss + memory repo (v4.10.0) + 4 application tests
- ✅ Migration 000007_create_risk_position (UNIQUE constraints + partial index for breached + check cap > 0)
- ✅ Postgres LimitRepo + PositionRepo with Reconstitute helpers (v4.11.0); container wires postgres backend
- ✅ RiskServiceServer + PositionServiceServer gRPC adapters under grpcgen tag (v4.12.0)

**Deferred:**
- ⏳ Real-time NOP monitoring via Flink CEP — MS-023g (EDA scope).
- ⏳ DV01 / VaR calculation engine — separate quant track.

## Description

Risk Management + Position Keeping + Dual-Ledger multi-currency: NOP realtime, limites enforced pre-trade, MTM EOD diario, postings atomicos multi-CCY via LedgerOS.

## Acceptance Criteria

- [ ] TradingLimit + NOPSnapshot + VaRSnapshot funcionais
- [ ] Position aggregate com real-time keeping
- [ ] cmd/eod/main.go batch MTM funcionando
- [ ] LedgerGateway multi-CCY com 3 adapters (gRPC/REST/stub)
- [ ] PostMultiLegTransaction atomico PvP
- [ ] 80+ tests

## Deliverables

- modules/risk/, modules/position/
- pkg/ledger/ com LedgerGateway interface
- cmd/eod/main.go
- Adapters: GRPCPostingAdapter, TemenosFXAdapter, StubAdapter

## Cross-References

- Plano monolitico: Fase F8 + F10
- Workstream: 02-core-domain + 05-integrations
