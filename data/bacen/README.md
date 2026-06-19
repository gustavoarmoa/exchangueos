# `data/bacen/` — BACEN regulatory data sources

> Machine-readable source-of-truth for BACEN Circular 3.690 nature codes + amendments.
> Consumed by `pkg/bacen/codegen` to produce `pkg/bacen/codes_full.go`.

## Files

| File | Purpose | Versioning |
|------|---------|------------|
| `nature-codes-circ-3690-vYYYYMMDD.csv` | Full nature-code catalogue | Date-versioned per BACEN amendment |
| `golden-classifications.csv` | Curated free-text → expected-code pairs for the classifier accuracy gate | Append-only |
| `CHANGELOG.md` | What changed between dated versions | Manual per amendment |

## Schema — `nature-codes-circ-3690-vYYYYMMDD.csv`

| Column | Type | Notes |
|--------|------|-------|
| `code` | string (5 digits) | BACEN nature code, e.g. `10000` |
| `category` | enum | `COMERCIAL`, `SERVICOS`, `CAPITAL`, `TRANSFERENCIAS`, `TURISMO`, `CARTAO`, `RENDA`, `DERIVATIVOS`, `OUTROS`, `VASP` |
| `subcategory` | snake-case string | finer grouping inside the category |
| `description_pt` | string | Portuguese description (canonical for BACEN comms) |
| `description_en` | string | English translation (for international stakeholders) |
| `direction` | enum | `INGRESSO`, `REMESSA`, `BIDIRECTIONAL` |
| `doc_requirement` | `+`-joined codes | minimum required document set, e.g. `INVOICE+DUE` |
| `iof_op_type` | string | Maps to `pkg/bacen.IOFCalculator` op type |
| `active` | bool | If false: not selectable by Classifier.ByCode |

## Versioning policy

1. BACEN publishes a Circular 3.690 amendment.
2. Compliance team produces a new `nature-codes-circ-3690-vYYYYMMDD.csv`.
3. PR includes:
   - CSV update
   - `CHANGELOG.md` entry summarising added / removed / modified codes
   - Re-run `go generate ./pkg/bacen/...` to regenerate `codes_full.go`
   - Accuracy regression run against `golden-classifications.csv` (≥ 95%)
4. Tag release `bacen-codes-vYYYYMMDD`.

## CSV → Go pipeline

```
nature-codes-circ-3690-vYYYYMMDD.csv
        │
        ▼  go run ./pkg/bacen/codegen --input <csv> --output ../codes_full.go
pkg/bacen/codes_full.go  (generated; checked in for reproducible builds)
        │
        ▼  imported by
pkg/bacen/classifier.go  (Classifier.ByCode + Classifier.FreeText)
        │
        ▼
modules/compliance/application/service.go  (ClassifyOperation use case)
```

## Source-of-truth chain

BACEN Circular 3.690 (authoritative) → Compliance team transcription → CSV in
this folder → generated Go code → classifier → audit trail per classification.
Any divergence between BACEN and CSV is a P1 bug.
