# 10 — Quality

> **Workstream:** Quality
> **Versao:** 1.0.0
> **Status:** DRAFT

## Conteudo

| Documento | Status | Descricao |
|-----------|--------|-----------|
| `local-quality-gates.md` | TODO | TDD + E2E + Security Local Gates (§18 monolitico) |
| `crud-tests.md` | TODO | ~290 CRUD integration tests (14 BCs × ~20) — §17 monolitico |
| `success-metrics.md` | TODO | Metricas de sucesso + KPIs (§8 monolitico) |
| `testing-strategy.md` | TODO | Testing strategy + 40 FX-TEST-* patterns |
| `tdd-workflow.md` | TODO | TDD Red-Green-Refactor playbook |
| `e2e-testing.md` | TODO | 10 cenarios canonicos E2E |
| `performance-baseline.md` | TODO | Performance baseline + benchmarks |
| `sli-slo-catalog.md` | TODO | 8+ SLIs com SLOs + error budget |
| `code-coverage.md` | TODO | Coverage gates (domain >= 80%, app >= 70%) |
| `chaos-engineering.md` | TODO | Chaos engineering por integration point |
| `certification/` | TODO | Certification docs |

## Testing Coverage Targets

| Tipo | Quantidade |
|------|------------|
| **Domain tests** | 470+ (per BC) |
| **Pricing tests** | 80+ (60 golden + 20 property-based) |
| **BACEN tests** | 196+ |
| **CRUD integration tests** | ~290 (14 BCs × ~20 tests) |
| **E2E tests** | 10 cenarios canonicos |
| **Contract tests** | gRPC + OpenAPI + AsyncAPI |
| **Load tests** | k6 + vegeta |
| **Compliance tests** | RN_FX coverage + SHACL + ISO 27001 evidence |
| **TOTAL** | **900+ tests** |

## SLI/SLO Targets

| SLI | SLO | Error Budget |
|-----|-----|--------------|
| RFQ latency p95 | < 50ms | 0.5% / 30d |
| Trade booking p95 | < 200ms | 0.5% / 30d |
| CLS submission p95 | < 500ms | 1.0% / 30d |
| PayIn deadline adherence | 99.99% | 0.01% / 30d (single miss = incident) |
| API availability | 99.95% | 0.05% / 30d |
| CLS daily cycle success | 100% | zero tolerance |
| SISCOAF filing SLA | 100% | zero tolerance |
| Audit log integrity | 100% | zero tolerance |

## Sources

- §8 (Metricas) + §17 (CRUD Tests Local) + §18 (TDD + E2E + Security Local Gates) do [snapshot monolitico v3.11.7](../_archive/allenty-v3.11.7-monolithic-plan.md)
- Pattern de referencia: [LedgerOS 10-quality](../../../../ledgeros/.base/plans/10-quality/)
