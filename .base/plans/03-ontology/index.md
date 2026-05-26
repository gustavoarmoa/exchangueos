# 03 — Ontology

> **Workstream:** Ontology
> **Versao:** 1.0.0
> **Status:** DRAFT

## Conteudo

| Documento | Status | Descricao |
|-----------|--------|-----------|
| `exchangeos-ontology-suite.md` | TODO | Master spec da Ontology Suite — `.base/aasc/ontology/` (35 TTL v1.2.0) — §10 monolitico |
| `coverage-metrics.md` | TODO | Triples + Cobertura (~4.300 triples + 100% SHACL + ≥80% FIBO + 100% ISO 20022) — §11 monolitico |
| `ttl-style-guide.md` | TODO | Convencoes TTL v1.2.0 (namespace + headers + imports) |
| `fibo-mapping.md` | TODO | FIBO mapping detalhado (50+ classes referenciadas) |
| `iso20022-mapping.md` | TODO | ISO 20022 mapping por message ID (32 schemas FX-specific) |
| `bacen-mapping.md` | TODO | BACEN regulatory mapping (RMCCI + Resolucoes BCB) |
| `shacl-validation-guide.md` | TODO | SHACL validation playbook (pyshacl + Apache Jena) |

## TTL Files (`.base/aasc/ontology/`)

| Categoria | Quantidade | Path |
|-----------|------------|------|
| Core (foundation + BCs + pricing) | 18 TTL | `core/` |
| Bridges (FIBO + ISO + BACEN + LedgerOS) | 9 TTL | `bridges/` |
| SHACL Shapes | 8 TTL | `shapes/` |
| Compliance Shapes (BACEN) | 5 TTL | `compliance/` |
| Domain-Specific (cls + cfets + bacen + pricing + swift-mt) | 16 TTL | `domains/` |
| FIBO Imports | 7 TTL | `imports/fibo/` |
| Fixtures | 6 TTL | `fixtures/` |
| **TOTAL** | **35 proprios + 7 imports + 6 fixtures** | |

## Coverage Targets

- **35 TTL files v1.2.0** OWL 2 DL profile compliant
- **~4.300 triples** semanticos
- **100% SHACL validation** das 50 RN_FX_001..050
- **≥ 80% FIBO coverage** (target medido em CI)
- **100% ISO 20022 coverage** dos 32 schemas FX-domain
- **OWL 2 DL + HermiT consistency check** passing

## Sources

- §10 (Ontology Suite) + §11 (Triples + Cobertura) do [snapshot monolitico v3.11.7](../_archive/allenty-v3.11.7-monolithic-plan.md)
- Pattern de referencia: [LedgerOS 03-ontology](../../../../ledgeros/.base/plans/03-ontology/)
