# ADR 0005 — shopspring/decimal mandatory; NEVER float for money

- Status: Accepted
- Date: 2026-05-24

## Context

FX trading involves arithmetic on money and rates with regulatory-grade precision (BACEN reports + IOF calculation + CLS net amounts). IEEE-754 floating point silently loses precision (`0.1 + 0.2 ≠ 0.3`). Even `float64` runs out of precision at amounts > 2^53 cents.

## Decision

**`decimal.Decimal` from `github.com/shopspring/decimal` is the only acceptable type for money or rate values throughout the codebase.**

Enforced at three layers:

1. **Lint** — `.golangci.yml` `forbidigo` rule:
   ```yaml
   forbidigo:
     forbid:
       - pattern: '^float64\s+'
         msg: "NEVER float64 for money/rate — use shopspring/decimal.Decimal"
   ```
2. **Storage** — every `migrations/*.sql` uses `DECIMAL(36,18)` for money/rate columns
3. **Review** — explicit RN_FX_026 cited in `modules/trade/domain/fxtrade.go` validation

Display rounding (4 decimals for majors, 2 for JPY) happens at the boundary via `decimal.RoundBank(n)` — banker's rounding (half-even) to avoid bias.

## Consequences

### Positive

- **Zero precision loss** in money arithmetic
- **Regulatory-grade audit trail** — every BACEN/IOF computation reproducible exactly
- **Decimal serialisation** in proto via string Amount fields — preserves precision across wire
- **Forbidigo lint catches violations** at PR-time, not in production

### Negative

- **~3× slower** than float64 arithmetic — benchmarked at < 5µs per CIP forward (acceptable for our throughput)
- **More verbose** — `decimal.NewFromInt(1_000_000)` vs `1000000.0`
- **Test setup** uses `decimal.RequireFromString` helpers; small ergonomic cost

### Mitigations

- Pricing engine benchmarked + tracked via `.github/workflows/benchmarks.yml`
- Test ergonomics improved via package-local `dec(s string)` helper in test files

## Alternatives considered

- **`int64` cents** — works for USD but breaks for JPY (0 decimals) and BHD (3 decimals); also loses precision in rate arithmetic
- **`big.Float`** — arbitrary precision but no native banker's rounding; community uses shopspring instead
- **Custom Money type** — reinventing the wheel; shopspring is battle-tested + maintained
