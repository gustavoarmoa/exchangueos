---
glob: "pkg/pricing/**/*.go"
---

# Rule: pkg/pricing/*.go

## Pricing Precision
- `shopspring/decimal` obrigatorio
- 8 decimais internos / 4-5 displayed
- Banker's rounding (half-even)
- NUNCA float64 ou math.* funcoes que retornam float

## Required Tests
- Golden tests com casos reais (BIS, CME, BACEN PTAX) em `tests/integration/golden/pricing/`
- Property-based tests (gopter): inverse, round-trip, cross-rate consistency
- Benchmark: `ForwardSimple` < 5µs target

## Cite Source
Reference formula iotafinance.com para Forward FX no comment.
