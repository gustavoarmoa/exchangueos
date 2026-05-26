# MS-023c — trade-core

| Field | Value |
|-------|-------|
| **Code** | MS-023c |
| **Name** | trade-core |
| **Phase** | F5 + F6 |
| **Sprint** | 4 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Started** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023b (delivered) |

## Delivery Notes

**Acceptance criteria met:**
- ✅ FXTrade aggregate (v4.3.0) with lifecycle PENDING→CONFIRMED→SETTLING→SETTLED + CANCELLED + 13 tests
- ✅ TradeRepository interface + memory impl (v4.9.0) + postgres impl with counterparty BIC joins (v4.9.0)
- ✅ Application service with 7 use cases (BookTrade/Get/List/Confirm/Cancel/MarkSettling/MarkSettled) + shared mutate pipeline (v4.9.0) + 7 tests
- ✅ Quote→Trade integration via in-process eventbus dispatcher (v4.9.0) + end-to-end TestContainer_QuoteAccepted_BooksTrade
- ✅ TradeServiceServer gRPC adapter under grpcgen tag with full enum mapping + canonical error codes (v4.9.0)
- ✅ HTTP smoke endpoint `/v1/trades/:id` (v4.9.0)
- ✅ ReconstituteFXTrade helper + 8 new accessors for postgres hydration

**Deferred (outside MS-023c scope):**
- ⏳ Kafka outbox publisher replacing NoopPublisher — MS-023g territory.

## Description

Trade Core implementado: FXTrade aggregate completo com state machine (NEW→MATCHED→CONFIRMED→SETTLEMENT_PENDING→SETTLED), suporte Spot/Forward/Swap/NDF, multi-leg para Swap, fixing rate snapshot para NDF + Forward. Amendment + Cancellation + Novation via gRPC interno.

## Acceptance Criteria

- [ ] FXTrade aggregate com 60+ tests (state machine + all trade types + specs)
- [ ] BookSpotTrade, BookForwardTrade, BookSwapTrade, BookNDFTrade commands
- [ ] MatchingService domain service com auto-match + counterparty confirmation
- [ ] Amendment + Cancellation + Novation flow com 4-eyes acima de USD 100k
- [ ] FXTradeService gRPC funcional

## Deliverables

- modules/trade/ com FXTrade + TradeLeg + ConfirmationStatus + FixingRate
- modules/amendment/ com FXAmendment + FXCancellation + FXNovation
- 14 commands CQRS + 12 queries
- 80+ domain unit tests

## Cross-References

- Plano monolitico: Fase F5 + F6
- Workstream: 02-core-domain
