# 00 — Governance

> **Workstream:** Governance
> **Versao:** 1.0.0
> **Status:** DRAFT

## Conteudo

| Documento | Status | Descricao |
|-----------|--------|-----------|
| `program-charter.md` | TODO | Visao + escopo + principios + governance model |
| `glossary.md` | TODO | Ubiquitous Language (FX domain + FIBO + ISO 20022) |
| `quality-gates.md` | TODO | 6 gates GS-1..GS-6 + Allenty Activate pipeline |
| `risk-register.md` | TODO | 90+ riscos identificados + mitigacoes |
| `open-questions.md` | TODO | 108 open questions categorizadas |
| `cross-references.md` | TODO | 200+ fontes oficiais consultadas |
| `integration-audit.md` | TODO | Matrix 4 vetores × 13 modulos (Kafka/DB/gRPC/Sync) — §20 monolitico |
| `enterprise-benchmark.md` | TODO | Allenty Triad (Oracle x SAP x Apple) — 12 principios |
| `document-framework.md` | TODO | Padrao de estruturacao de docs do workstream |
| `decision-records/` | TODO | ADRs (ADR-001..050+) |
| `audits/` | TODO | Audit reports |
| `legal/` | TODO | Regulatory + contracts |
| `product/` | TODO | Product roadmap context |
| `standards/` | TODO | Internal standards |

## Decision Records (ADRs)

| ADR | Title | Status |
|-----|-------|--------|
| ADR-001 | Microservice Architecture | ACCEPTED |
| ADR-009 | AccountOS Architecture | ACCEPTED |
| ADR-010 | OnboardOS Architecture | ACCEPTED |
| ADR-011 | ExchangeOS Architecture (standalone vs embedded, multi-CCY dual-ledger, CLS strategy, SWIFT FIN vs MX) | DRAFT |
| ADR-012 | PVP Strategy (CLS first vs gross with Herstatt threshold) | DRAFT |
| ADR-013 | ISO 20022 FX Coverage (fxtr + CLS + CFETS + decomposicao trea → fxti/fxmt obsolete) | DRAFT |
| ADR-014 | Dual-Ledger Strategy (internal LedgerOS / production Temenos / stub) | ACCEPTED |
| ADR-015 | Shared CRDB Hub Adoption (NUNCA inline insecure) | DRAFT |

## Sources

- §1 (Visao Geral) + §5 (Milestones) + §6 (Cross-References) + §7 (Riscos) + §9 (Open Questions) + §20 (Integration Audit) do [snapshot monolitico v3.11.7](../_archive/allenty-v3.11.7-monolithic-plan.md)
- Pattern de referencia: [LedgerOS 00-governance](../../../../ledgeros/.base/plans/00-governance/)
