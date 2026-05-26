# 02 — Core Domain

> **Workstream:** Core Domain
> **Versao:** 1.0.0
> **Status:** DRAFT

## Conteudo

| Documento | Status | Descricao |
|-----------|--------|-----------|
| `exchangeos-engine.md` | TODO | ExchangeOS Engine — 14 BCs com aggregates + state machines + domain events (§2.1 monolitico) |
| `pricing-engine.md` | TODO | Pricing & Algorithms Module (CIP iotafinance + NDF + PTAX + MTM + cross-rate) — §2.5 monolitico |
| `cls-daily-cycle.md` | TODO | CLS Daily Settlement Day Timeline (07:00-12:00 CET) — §2.3 monolitico |
| `trade-aggregate.md` | TODO | FXTrade aggregate detail (Spot/Forward/Swap/NDF) |
| `cls-submission-aggregate.md` | TODO | CLSSubmission aggregate + fxtr lifecycle |
| `payin-aggregate.md` | TODO | PayInSchedule + PayInCall + PayInEvent aggregates |
| `cfets-capture-aggregate.md` | TODO | CFETSTradeCapture aggregate |
| `cfets-confirmation-aggregate.md` | TODO | CFETSConfirmationRequest aggregate |
| `position-aggregate.md` | TODO | Position aggregate (real-time NOP + MTM EOD) |
| `refdata-aggregates.md` | TODO | CurrencyPair, Calendar, SSI, Counterparty, NettingCutOff |
| `business-rules.md` | TODO | 50 RN_FX_001..050 detalhados |
| `domain-events.md` | TODO | ~24 domain events catalog |
| `use-cases.md` | TODO | Use cases canonicos por persona (Trader, Compliance Officer, Settlement Ops) |

## 14 Bounded Contexts

| BC | Aggregate Root | ISO 20022 Family |
|----|---------------|------------------|
| FX Trade | `FXTrade` | (interno; vira fxtr.014/031 na fronteira) |
| FX Quote / RFQ | `FXQuoteRequest` | (gRPC interno) |
| FX Amendment | `FXAmendment`, `FXCancellation`, `FXNovation` | (gRPC interno → fxtr.015/016/035/036) |
| FX CLS Settlement | `CLSSubmission`, `CLSStatusUpdate` | fxtr CLS (008/013/014/015/016/017/030) |
| FX PayIn Lifecycle | `PayInSchedule`, `PayInCall`, `PayInEvent` | camt.061/062/063 |
| FX Net Report | `NetReport` | camt.088 |
| FX CFETS Capture | `CFETSTradeCapture` | fxtr.031/032/033 |
| FX CFETS Confirmation | `CFETSConfirmationRequest`, `CFETSStatusAdvice` | fxtr.034/035/036/037/038 |
| FX Settlement (non-CLS) | `FXSettlement` | camt.052/053/054/056/060/087 |
| FX Reference Data | `CurrencyPair`, `Calendar`, `SSI`, `Counterparty`, `NettingCutOff` | reda.016/017/018/028/029/060/061 |
| FX Administration | `SystemEvent`, `StaticDataSync`, `MessageReject` | admi.002/004/009/010/011/017 |
| FX Risk & Limits | `TradingLimit` | (interno + RiskOS) |
| FX Position | `Position` | (interno + LedgerOS) |
| FX Compliance | `DECDeclaration`, `IEDRecord`, `CreditoRecord`, `CBEDeclaration`, `IOFEntry`, `COSReport`, `SanctionsHit` | (interno + AuthorityOS + BACEN) |

## Sources

- §2.1 (BCs + State Machines + Business Rules) + §2.3 (CLS Daily Cycle) + §2.5 (Pricing & Algorithms) do [snapshot monolitico v3.11.7](../_archive/allenty-v3.11.7-monolithic-plan.md)
- Pattern de referencia: [LedgerOS 02-core-domain](../../../../ledgeros/.base/plans/02-core-domain/)
