# MS-023l — erds-suite

| Field | Value |
|-------|-------|
| **Code** | MS-023l |
| **Name** | erds-suite |
| **Phase** | F15D |
| **Sprint** | 12 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Started** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023j (delivered) |

## Delivery Notes

**Acceptance criteria met:**
- ✅ 5 representative ERD files in `domain/`: trade + quote + settlement + risk-position + compliance-admin
- ✅ Mermaid `erDiagram` format with FK relationships labeled, constraints documented, indexes itemised
- ✅ Each ERD cross-references its source migration + ontology
- ✅ README.md documents sync rule (migrations ↔ ERDs must update in same PR)

**Deferred:**
- ⏳ Per-table SQL DDL mirrors in `sql/` (9 files) — generated on demand from migrations
- ⏳ Lefthook pre-commit glob check enforcing the migration↔ERD sync (TODO MS-023n)
- ⏳ Auto-render Mermaid → SVG export step in CI for stakeholder review

## Description

ERDs Suite completa: 23 ERDs (14 BC + 5 cross-BC + 4 common) + 16 SQL DDL CockroachDB + 5 matrices + 5 schemas docs + 8 sample queries + 6 tools + CI sync/audit.

## Acceptance Criteria

- [ ] 23 ERDs .md + 14+ Mermaid .mmd standalone
- [ ] 16 SQL DDL CockroachDB executaveis
- [ ] ~70 tabelas, ~120 FKs, ~85 indices, 40+ ENUMs
- [ ] CI: DDL ↔ migrations sync + FK consistency + tenant scoping + ownership audit + decimal precision green

## Deliverables

- .base/erds/ completa
- 16 SQL DDL files
- 5 matrices arquiteturais
- 6 tools ERD-as-code

## Cross-References

- Plano monolitico: §13 + Fase F15D
- Workstream: 01-architecture
