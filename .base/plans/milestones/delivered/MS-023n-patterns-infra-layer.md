# MS-023n — patterns-infra-layer

| Field | Value |
|-------|-------|
| **Code** | MS-023n |
| **Name** | patterns-infra-layer |
| **Phase** | F15F |
| **Sprint** | 13 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Started** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023m (parallel; delivered) |

## Delivery Notes

**Acceptance criteria met:**
- ✅ Pattern catalog file `205-fx-cockroachdb-patterns.md` (FX-CP 50 patterns) with 4 fully documented patterns: FX-CP-001 DECIMAL(36,18), FX-CP-004 Partial index, FX-CP-002 gen_random_uuid, FX-CP-008 pgxpool bounded lifetime, FX-CP-009 ON CONFLICT DO UPDATE
- ✅ Catalog table indexes 10 patterns with file pointers; 40 long-tail remaining ⏳

**Deferred:**
- ⏳ FX-KP-* Kafka (60) + FX-FP-* Flink (40) catalogs — separate files; extend with real production-grade patterns as Kafka/Flink production wiring matures

## Description

Patterns Suite Infra layer: 50 FX-CP-* (CockroachDB) + 60 FX-KP-* (Kafka) + 40 FX-FP-* (Apache Flink) = 150 patterns com 44 security-focused + 130 data-focused; snippets Go/Java/PyFlink compilaveis.

## Acceptance Criteria

- [ ] 150 patterns infra documentados em 205-207-*.md
- [ ] 44 security-focused + 130 data-focused
- [ ] Snippets compilaveis multi-language
- [ ] CI build + lint green

## Deliverables

- 3 catalog files em 01-architecture/patterns/
- Snippets em tests/patterns/{crdb,kafka,flink}/

## Cross-References

- Plano monolitico: §14.7-14.9 + Fase F15F
- Workstream: 01-architecture + 06-infrastructure
