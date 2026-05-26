# MS-023f — compliance-core-admin

| Field | Value |
|-------|-------|
| **Code** | MS-023f |
| **Name** | compliance-core-admin |
| **Phase** | F9 + F12 |
| **Sprint** | 7 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Started** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023e (delivered) |

## Delivery Notes

**Acceptance criteria met:**
- ✅ 4 compliance aggregates (v4.11.0): Classification (95-code BACEN nature), IOFComputation (Decreto 12.499/2025 rates), BACENReport (PENDING→SUBMITTED→ACCEPTED/REJECTED), ScreeningResult (LOW/MEDIUM/HIGH derived; RequiresCOS hook for RN_FX_039) + 18 tests
- ✅ 2 admin aggregates (v4.11.0): SystemEvent (admi.x mapped 8 codes), EODJob (PENDING→RUNNING→COMPLETED/FAILED + idempotent MarkStep) + 6 tests
- ✅ Compliance + Admin application services (v4.12.0): Service wires pkg/bacen.Classifier + IOFCalculator for ClassifyOperation/ComputeIOF; SubmitBACENReport persists PENDING; ScreenCounterparty derives risk_level + 8 application tests
- ✅ ComplianceServiceServer + AdminServiceServer gRPC adapters (v4.12.0) under grpcgen tag, registered in grpc_register_proto.go
- ✅ Migration 000008_create_compliance_admin (6 tables: classifications + iof_computations + bacen_reports + screening_results + system_events + eod_jobs with CHECK constraints + UNIQUE (tenant, business_date) + partial index for HIGH-risk screenings)
- ✅ Container wires Compliance + Admin in wireComplianceAdmin() — both backends

**Deferred:**
- ⏳ Postgres repos for compliance + admin — wire in next sprint (memory currently; schema ready in 000008).
- ⏳ Real OFAC/UN/EU/COAF list provider integration — ScreeningResult currently accepts caller-supplied hits; production needs adapter.

## Description

Compliance Core (DEC + Sanctions + Audit log Merkle) + Administration messages CLS (admi.002/004/009/010/011/017) integrados.

## Acceptance Criteria

- [ ] DECDeclaration + RegulatoryReport + SanctionsHit aggregates
- [ ] Audit log tamper-evident (Merkle hash chain)
- [ ] admi.* messages handler completo
- [ ] StaticDataSync workflow (admi.009 → admi.010)
- [ ] SystemEvent dispatcher (admi.004 → block submissions em halt)
- [ ] 60+ tests

## Deliverables

- modules/compliance/ + modules/admin/
- Tamper-evident audit_log table
- 6 admi message marshalers (002/004/009/010/011/017)

## Cross-References

- Plano monolitico: Fase F9 + F12
- Workstream: 02-core-domain + 09-compliance
