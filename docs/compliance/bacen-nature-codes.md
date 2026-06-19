# BACEN Nature Codes — Circular 3.690

> Operational guide for the BACEN nature-code catalogue used by ExchangeOS
> for DEC submission, IOF lookup, and free-text classification.
> Source-of-truth: `data/bacen/nature-codes-circ-3690-v20260101.csv`.
> Owner: Compliance team.

## What this is

Every FX operation booked through ExchangeOS must carry a BACEN **codigo de
natureza** — the regulatory identifier from Circular 3.690 that tells the
SISBACEN downstream report what kind of flow happened (import, royalty,
inheritance, etc.). There are 95 codes in the canonical Circular; ExchangeOS
currently ships a curated subset of **46 codes** covering the 10 categories
most commonly seen in cross-border BRL flows.

The codes drive:

- **DEC submission** (`Declaração Eletrônica de Câmbio`) — the code goes into
  the regulatory payload. An invalid code = rejected DEC.
- **IOF calculation** — each code maps to an `iof_op_type` consumed by
  `pkg/bacen/IOFCalculator`. Cite `RN_FX_037`.
- **Free-text classification** — when an operator books a trade with only a
  description, `Classifier.Classify` resolves a code via the keyword index.
  Audited against a golden corpus at every commit.

## Source-of-truth

| Path | Purpose | Owner |
|------|---------|-------|
| `data/bacen/nature-codes-circ-3690-v20260101.csv` | Machine-readable catalogue (1 row per code) | Compliance |
| `pkg/bacen/codes_full.go` | **Generated** from the CSV via `pkg/bacen/codegen` | Build pipeline |
| `pkg/bacen/classifier.go` | Boot-time keyword rules + legacy seed | Platform |
| `data/bacen/golden-classifications.csv` | 144 PT/EN phrases mapped to expected codes — accuracy harness | Compliance + Platform |

### CSV schema

```
code,category,subcategory,description_pt,description_en,direction,doc_requirement,iof_op_type,active
```

- **code** — 5-digit BACEN identifier (string, no leading-zero loss).
- **category / subcategory** — taxonomy (COMERCIAL / SERVICOS / CAPITAL /
  TRANSFERENCIAS / TURISMO / CARTAO / RENDA / DERIVATIVOS / OUTROS / VASP).
- **description_pt / description_en** — bilingual canonical description.
  PT is the regulatory wording; EN is for cross-team comprehension.
- **direction** — `INGRESSO` (inbound BRL) | `REMESSA` (outbound BRL) |
  `BIDIRECTIONAL` (conversion / derivative).
- **doc_requirement** — `+`-separated list of accepted document types
  (`INVOICE+DUE`, `SCE-IED+CONTRACT`, etc.). Consumed by the DEC submitter
  to flag missing evidence pre-submit.
- **iof_op_type** — string consumed by `pkg/bacen/IOFCalculator` (`EXPORT`,
  `IMPORT`, `DEFAULT`, `TRAVEL_CARD`, `TRAVEL_CASH`, `CREDIT_CARD`,
  `INSURANCE_FOREIGN`, `LOAN_SHORT`, `INVESTMENT`).
- **active** — `true` / `false`. Deprecated codes stay in the CSV but are
  filtered out at lookup time so historic rows can still be hydrated.

## Workflows

### Adding / updating a code

1. Edit `data/bacen/nature-codes-circ-3690-v20260101.csv`.
2. Run `go generate ./pkg/bacen/...` (regenerates `codes_full.go`
   deterministically — output is byte-stable, safe to commit).
3. If the description introduces new common phrasing, add a representative
   phrase to `data/bacen/golden-classifications.csv` and (if needed) a
   keyword rule to `pkg/bacen/classifier.go::defaultKeywordRules()`.
4. `go test ./pkg/bacen/...` — both the existence tests and the accuracy
   harness must stay green (≥ 95% threshold).
5. Commit. The CI drift guard (`scripts/check-bacen-codegen.sh`, wired into
   GitHub Actions `lint` job + lefthook `pre-push`) re-runs codegen and
   diffs against the committed file — drift fails the build.

### Bumping the CSV version

The filename embeds a version stamp (`-v20260101`). BACEN republishes the
Circular periodically; bump policy:

- **Additive change** (new code added): keep the filename, append the row,
  regenerate. Tag the commit `chore(bacen): add code NNNNN per Circ 3.690 X/Y`.
- **Breaking change** (code retired or semantics flipped): bump the filename
  (`-v<YYYYMMDD>`), update the `//go:generate` directive in
  `pkg/bacen/classifier.go`, and add an ADR documenting the migration plan
  for in-flight rows that referenced the retired code.

## Classifier resolution order

```
Classifier.ByCode(code)
  1. Look up in AllNatureCodes (generated catalogue).
  2. Fall back to legacy builtin map (~20 seed codes from the bootstrap era).
  3. Return false on miss.

Classifier.Classify(hint)
  1. Lowercase + trim the hint.
  2. Walk defaultKeywordRules() in order; first substring match wins.
  3. Resolve the rule's target code via ByCode (same path as above).
  4. Return ErrUnknown on no match.
```

Keyword rules are ordered **most specific first**: direction-specific phrases
(`recebimento em criptomoeda`, `dividendos para investidor`) precede generic
ones (`cripto`, `investidor`). Adding a new general rule near the top of the
slice will flip previously-correct classifications — re-run the accuracy
harness before merging.

## Accuracy harness

`pkg/bacen/accuracy_test.go::TestClassifier_AccuracyAgainstGoldenCorpus`
runs every phrase in `data/bacen/golden-classifications.csv` through
`Classify` and asserts the hit ratio meets the threshold (currently 95%).
The current measured accuracy against the 144-phrase corpus is **99.3%**.

When a phrase newly misses (after a rule change or a new code), the test
output shows the first 15 misses with the actual rule fire — debuggable
without spelunking. Fix by either adding a more-specific rule above the
generic one OR by expanding the corpus if the phrase was ambiguous.

## Known gap — 46 of 95 codes shipped

The Circular 3.690 publishes 95 codes; the ExchangeOS CSV currently covers
46 across the 10 most-trafficked categories. The remaining 49 require an
authoritative source from the Compliance team because BACEN does not publish
the codes in a machine-readable format — they live inside the regulatory PDF
plus its addenda, and inventing identifiers carries hard downstream risk
(DEC rejection at SISBACEN).

**Unblock path:**

1. Compliance extracts the remaining 49 rows (or provides the PDF + a
   per-row mapping spreadsheet).
2. Append rows to `data/bacen/nature-codes-circ-3690-vYYYYMMDD.csv`
   (bumping the version stamp).
3. `go generate ./pkg/bacen/...` + commit.
4. Expand `golden-classifications.csv` with ~2-3 phrases per new code.
5. Re-run accuracy harness; tune keyword rules if needed.

Until step 1 lands, `Classifier.Classify` returns `ErrUnknown` for any
phrase whose target code is among the missing 49 — by design (refuse to
guess) rather than fabricate a code that the regulator will reject.

## Cross-references

- Circular BACEN 3.690 — [BACEN website](https://www.bcb.gov.br/) (search "Circular 3690")
- `RN_FX_028` — "Codigo natureza valido" business rule (`.base/aasc/ontology/compliance/bacen-cambio-shapes.ttl`)
- `RN_FX_037` — IOF rate per `iof_op_type` (`pkg/bacen/iof.go`)
- DEC submission flow — `modules/compliance/application/service.go::SubmitDEC`
- Master plan reference — `.base/plans/09-compliance/index.md`
