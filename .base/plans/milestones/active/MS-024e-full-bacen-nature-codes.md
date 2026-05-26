# MS-024e — Full BACEN Nature Code Catalogue (95 codes)

| Field | Value |
|-------|-------|
| **Code** | MS-024e |
| **Name** | full-bacen-nature-codes |
| **Phase** | F-OPS-PROD |
| **Sprint** | 2 of MS-024 cycle |
| **Status** | ACTIVE |
| **Owner** | Compliance |
| **Dependencies** | None |

## Why this milestone

`pkg/bacen/classifier.go` ships with 20 representative codes + a keyword fallback. Production-grade BACEN classification per Circular 3.690 requires the full **95-code catalogue** with proper categorisation, sub-codes, and authoritative descriptions. Without it, classifications will routinely fall through to imprecise keyword matching.

## Description

Source the complete BACEN Circular 3.690 nature-code table, encode as `pkg/bacen/codes_full.go` (generated from CSV), add lookup helpers, expand keyword rules so > 95% of typical free-text classifies confidently.

## Acceptance Criteria

- [ ] `data/bacen/nature-codes-circ-3690-v<latest>.csv` checked in (source of truth, machine-readable)
- [ ] `pkg/bacen/codes_full.go` generated via `go generate` from CSV (95 codes with code + category + description PT/EN + minimum_doc_set)
- [ ] `Classifier.ByCode(code)` returns full structured `NatureCode` (not just exists/not)
- [ ] `Classifier.FreeText` keyword rules cover all 95 categories — measured by golden corpus
- [ ] Golden test corpus `data/bacen/golden-classifications.csv` with ≥ 200 real-world phrases + expected codes (curated with Compliance team)
- [ ] Classifier accuracy ≥ 95% against golden corpus
- [ ] Code generation pipeline runs in CI to catch CSV drift
- [ ] Documentation: `docs/compliance/bacen-nature-codes.md` linking to Circ 3.690 + version + diff vs previous version
- [ ] Update business rule `RN_FX_028` documentation + SHACL shape
- [ ] Versioning policy: when BACEN publishes amendment, bump `nature-codes-circ-3690-vYYYYMMDD.csv` + run accuracy regression

## Deliverables

- `data/bacen/nature-codes-circ-3690-v<date>.csv`
- `data/bacen/golden-classifications.csv`
- `pkg/bacen/codes_full.go` (generated)
- `pkg/bacen/codegen/main.go` (generator)
- `pkg/bacen/classifier_full_test.go` (golden accuracy test)
- `docs/compliance/bacen-nature-codes.md`
- Updated `.base/aasc/ontology/compliance/bacen-cambio-shapes.ttl` for RN_FX_028

## Cross-References

- BACEN Circular 3.690/2013 + amendments
- `pkg/bacen/classifier.go` (current 20-code impl)
- `.base/plans/02-core-domain/business-rules.md` RN_FX_028
- MS-024f (DEC submission consumes nature code)
