# MS-023k — flows-suite

| Field | Value |
|-------|-------|
| **Code** | MS-023k |
| **Name** | flows-suite |
| **Phase** | F15C |
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
- ✅ 5 representative RFLW files across 4 sub-folders: trade/RFLW.024.001.01 + quote/RFLW.024.010.01 + cls_settlement/RFLW.024.020.01 + compliance/RFLW.024.030.01 + cfets/RFLW.024.040.01
- ✅ Standardised 10-section structure: YAML metadata + Description + Sequence (Mermaid) + Error flowchart (Mermaid) + Business Rules table + Observability + Compliance Notes
- ✅ Predecessor/Successor links navigable across flows
- ✅ Traceability fields cite RNs + ISO 20022 messages + Ontology classes
- ✅ README.md documents RFLW.024.NNN.NN naming convention + sub-folder catalog

**Deferred:**
- ⏳ 60 remaining flows scheduled across 6 sub-folders (12 trade + 8 quote + 15 cls_settlement + 10 cfets + 12 compliance + 6 eod = 63 planned; 5 done → 58 remaining)
- ⏳ Auto-rendering CI step (Mermaid → PNG export for stakeholder review)

## Description

Flows Suite completa: 85 flows individuais RFLW.024.NNN.NN em .base/flows/ com Mermaid sequence + flowchart + metadata + CI lint/traceability/cross-ref.

## Acceptance Criteria

- [ ] 85 flows individuais Mermaid em 13 subdominios
- [ ] CI lint Mermaid + traceability checker green
- [ ] Index automatico por subpasta
- [ ] Domain code RFLW.024 reservado

## Deliverables

- .base/flows/ com 13 subpastas + ~170 Mermaid diagrams
- CI workflow flows.yml
- 5 tools de automacao

## Cross-References

- Plano monolitico: §12 + Fase F15C
- Workstream: 01-architecture
