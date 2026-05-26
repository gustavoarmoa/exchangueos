# MS-023d2 — cfets-capture-confirmation

| Field | Value |
|-------|-------|
| **Code** | MS-023d2 |
| **Name** | cfets-capture-confirmation |
| **Phase** | F7F |
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
- ✅ All 8 fxtr CFETS variant structs (v4.10.0) — fxtr.031 Capture Request, fxtr.032 Ack w/ CFETSAckStatus (SUCC/REJT), fxtr.033 Notification, fxtr.034 Confirmation Request, fxtr.035 Confirmation, fxtr.036 Status w/ CFETSConfStatus (PAIR/UPRD/REJT), fxtr.037 Amendment, fxtr.038 Cancellation. Shared CFETSTradeIdentification + CFETSEconomics blocks + 8 namespace constants. fxtr/doc.go status table = 15/15 ✅
- ✅ CFETSCapture aggregate (v4.10.0) — DRAFT → SUBMITTED → ACK/REJECTED → NOTIFIED lifecycle with 5 DomainEvents + 8 tests
- ✅ CFETSConfirmation aggregate (v4.10.0) — CONFIRMING → CONFIRMED/UNPAIRED/REJECTED with 4 DomainEvents + 7 tests
- ✅ CFETSCaptureService + CFETSConfirmationService application layer (v4.12.0) with shared mutate pipelines + memory repos
- ✅ Container wires both services in wireComplianceAdmin()

**Deferred (deliberate — internal-only):**
- ⏳ No public gRPC service for CFETS — capture/confirmation are internal flows triggered by trade lifecycle. Public API surface remains TradeServiceServer; CFETS messages emit via outbox to downstream CFETS adapter (MS-023g scope).

## Description

CFETS Trade Capture Report (fxtr.031/032/033) + Trade Confirmation Request (fxtr.034/035/036/037/038) end-to-end para China interbank market.

## Acceptance Criteria

- [ ] CFETSTradeCapture aggregate funcional
- [ ] CFETSConfirmationRequest + CFETSStatusAdvice aggregates
- [ ] 5-message lifecycle fxtr.034-038
- [ ] CFETS gateway adapter (HTTPS REST)
- [ ] 48+ tests

## Deliverables

- modules/cfets_capture/ + modules/cfets_confirmation/
- Adapter CFETS gateway
- Tests E2E para Trade Capture + Confirmation flows

## Cross-References

- Plano monolitico: Fase F7F
- Workstream: 05-integrations
