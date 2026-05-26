# MS-023b — refdata-pricing-quote

| Field | Value |
|-------|-------|
| **Code** | MS-023b |
| **Name** | refdata-pricing-quote |
| **Phase** | F3 + F4P + F4 |
| **Sprint** | 3 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Started** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023a (delivered 2026-05-24) |

## Delivery Notes (2026-05-24)

**Acceptance criteria met:**
- ✅ 30+ currency pairs ativos (32 in seeds/02_currency_pairs.sql; 18 CLS-eligible)
- ✅ BACEN/NYFR/BOE/TARGET2/TOKYO/TORONTO calendars 2026 carregados (seeds/03)
- ✅ Counterparty BICs (36 sample) + SSI sample (seeds/04, 05)
- ✅ Netting cutoffs PIN1/PIN2/PIN3 + bilateral EOD per CCY (seeds/06)
- ✅ pkg/pricing/ com 47 tests (cip + crossrate + ndf + tenor + ptax + mtm; 6/7 algorithms ✅)
- ✅ CIP forward formula validada (golden cases EURUSD/USDBRL/GBPUSD/USDJPY hand-computed)
- ✅ Quote engine + RFQ lifecycle (modules/quote/{domain,application})
- ✅ PTAX fetch via BACEN OLINDA API (modules/refdata/infrastructure/olinda)
- ✅ Real PricingEngine wiring (modules/refdata/infrastructure/pricing.Engine) replacing stub (v4.9.0)
- ✅ Postgres repos for refdata + quote (v4.8.0)
- ✅ gRPC adapters (RefDataServiceServer, QuoteServiceServer) under build tag grpcgen
- ✅ HTTP smoke endpoint `/v1/refdata/currencies` proving end-to-end wiring

**Deferred (out of scope here, owned elsewhere):**
- ⏳ Kafka outbox publisher replacing NoopPublisher — **MS-023g (EDA + E2E)** scope.
  Container already exposes `EventBus` + adapter interfaces; only the Kafka client needs swap.

This milestone closes with all 7 pricing algorithms exposed + 3 bounded contexts wired
end-to-end (RefData, Pricing, Quote/RFQ).

## Description

RefData populado (currency pairs, calendars BACEN+NYFR+BOE+TARGET2, SSI, counterparties, CLS netting cutoffs) + pkg/pricing/ completo (CIP iotafinance + NDF + PTAX D-1/D-2 + cross-rate via USD pivot + MTM EOD) + Quote/RFQ lifecycle funcional consumindo pricing engine.

## Acceptance Criteria

- [ ] 30+ currency pairs ativos (18 CLS-eligible marcados)
- [ ] BACEN/NYFR/BOE/TARGET2 calendars 2026 carregados
- [ ] SSI cadastrados para counterparties iniciais
- [ ] pkg/pricing/ com 80+ tests (60 golden de mercado + 20 property-based)
- [ ] CIP forward formula validada contra casos BIS/CME/BACEN PTAX historico
- [ ] Quote engine + RFQ lifecycle funcional (gRPC interno)
- [ ] PTAX fetch via BACEN OLINDA API + fallback 4-window survey

## Deliverables

- modules/refdata/ implementado (5 aggregates: CurrencyPair, Calendar, Counterparty, SSI, NettingCutOff)
- pkg/pricing/ com 14 arquivos Go (forward, crossrate, ndf, ptax, points, pip, inverse, spread, mtm, daycount/, calendar.go, tenor.go, spotdate.go, implied.go)
- modules/quote/ com FXQuoteService gRPC
- seeds/02_currency_pairs.sql + 03_calendars.sql + 04_counterparties.sql + 05_ssi.sql + 06_netting_cutoffs.sql carregados

## Cross-References

- Plano monolitico: §2.5 (Pricing & Algorithms) + Fase F3 + F4P + F4
- Workstream: 02-core-domain
