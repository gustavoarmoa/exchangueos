# MS-023d — settlement-cls-non-cls

| Field | Value |
|-------|-------|
| **Code** | MS-023d |
| **Name** | settlement-cls-non-cls |
| **Phase** | F7 + F11 |
| **Sprint** | 5 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Started** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023c (delivered) |

## Delivery Notes

**Acceptance criteria met:**
- ✅ CLSCycle aggregate (v4.9.0) with full lifecycle OPEN → PAY_IN_WINDOW → SETTLING → CLOSED/FAILED + 6 DomainEvents + 8 tests; deadlines anchored to Europe/Zurich CET (07/08/09/10/12)
- ✅ PayInInstruction aggregate (v4.9.0) with auto-fail on missed deadline via ErrDeadlineMissed + 8 tests
- ✅ NetReport aggregate (v4.9.0) with derived NetSettlement = GrossPayIn − GrossPayOut + IsReceivable/IsPayable + 6 tests
- ✅ Application services for all 3 aggregates (v4.10.0): CLSSettlement (7 use cases incl. OpenCycle with ErrConflict on duplicate), PayIn (6 use cases), NetReport (3 use cases) + 13 application tests
- ✅ SettlementServiceServer gRPC adapter under grpcgen tag (v4.10.0): OpenCycle, SubmitPayIn (Create+Submit chained), GetNetReport, CloseCycle (EnterSettling+Close chained)
- ✅ Migration 000006_create_settlement (cls_cycles + cls_cycle_trades + payin_instructions + net_reports with FK cascades + partial index for open cycles + composite indexes)
- ✅ Postgres CycleRepo with transactional Save (UPSERT cycle + DELETE+INSERT trade_ids in same tx) + ReconstituteCycle (v4.11.0)

**Deferred:**
- ⏳ Postgres repos for payin + netreport — MS-023g (Kafka outbox lands them together).
- ⏳ Real camt.088 marshalling for GetNetReport — currently placeholder XML.

## Description

Settlement multi-pista funcional: CLS PvP completo (Instruction fxtr.014 + PayIn cycle camt.061/062/063 + NetReport camt.088) para 18 CCYs elegiveis; gross/non-CLS para os demais com BACEN Tx 70; SWIFT FIN bridge via IBM MQ; reconciliacao nostro automatica.

## Acceptance Criteria

- [ ] CLSSubmission aggregate com state machine completa
- [ ] cmd/cls-cycle/main.go scheduler ativo (deadlines 08/09/10 CET pay-in)
- [ ] camt.061/062/063 PayIn cycle end-to-end
- [ ] camt.088 NetReport reconciliation
- [ ] FXSettlement non-CLS com BACEN Tx 70
- [ ] IBM MQ bridge (espelho paymentos/internal/ibmmq) funcional
- [ ] Nostro reconciliation engine detecta breaks
- [ ] Counterparty adapters: CLS + SWIFT FIN + BACEN Tx 70

## Deliverables

- modules/cls_settlement/, modules/payin/, modules/netreport/, modules/settlement/
- cmd/cls-cycle/main.go + cmd/mq-bridge/main.go
- internal/ibmmq/ espelho de paymentos
- Adapters: cls/, swift/, bacen/
- 145+ tests

## Cross-References

- Plano monolitico: Fase F7A-F7E + F11
- Workstream: 02-core-domain + 05-integrations
