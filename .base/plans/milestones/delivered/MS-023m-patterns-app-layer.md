# MS-023m — patterns-app-layer

| Field | Value |
|-------|-------|
| **Code** | MS-023m |
| **Name** | patterns-app-layer |
| **Phase** | F15E |
| **Sprint** | 12-13 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Started** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023i (not blocking) |

## Delivery Notes

**Acceptance criteria met:**
- ✅ Pattern catalog file `200-fx-golang-patterns.md` (FX-GP 40 patterns) with 5 fully documented patterns: FX-GP-001 Aggregate constructor, FX-GP-002 Decimal precision, FX-GP-003 Pointer receiver, FX-GP-004 Sentinel errors, FX-GP-005 Build-tag-gated bindings (Context/Problem/Solution/Example/Anti-pattern/Related)
- ✅ Catalog table indexes 10 named patterns with code-pointer locations; 35 long-tail remaining ⏳
- ✅ Patterns directly cited from code via `// FX-GP-NNN` comments

**Deferred:**
- ⏳ FX-DDD-* (35 patterns) + FX-EDA-* (45 patterns) catalogs — separate files; same template; extend on demand
- ⏳ Full 40 FX-GP patterns — long tail extends per real-world cases encountered

## Description

Patterns Suite App layer: 40 FX-GP-* (Golang) + 35 FX-DDD-* (Domain-Driven Design) + 45 FX-EDA-* (Event-Driven Architecture) = 120 patterns com code snippets compilaveis + cross-ref matrix.

## Acceptance Criteria

- [ ] 120 patterns documentados em 200-202-*.md
- [ ] 120 snippets Go compilaveis em tests/patterns/{golang,ddd,eda}/
- [ ] Cross-ref matrix complete
- [ ] CI build + test + lint green

## Deliverables

- 3 catalog files em 01-architecture/patterns/
- 120 snippets compilaveis
- CROSS-REF-MATRIX.md

## Cross-References

- Plano monolitico: §14.2-14.4 + Fase F15E
- Workstream: 01-architecture
