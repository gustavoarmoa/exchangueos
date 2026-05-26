package pricing

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// Pair represents one FX quote: `Rate` units of `QuoteCCY` per 1 unit of `BaseCCY`.
//
// Example: Pair{BaseCCY:"EUR", QuoteCCY:"USD", Rate:1.08} means 1 EUR = 1.08 USD.
type Pair struct {
	BaseCCY  string
	QuoteCCY string
	Rate     decimal.Decimal
}

// Invert returns the same Pair viewed from the opposite direction:
//
//	Pair{EUR, USD, 1.08}.Invert() == Pair{USD, EUR, 1/1.08 == 0.925925…}
//
// Returns ErrZeroDenominator if Rate is zero.
func (p Pair) Invert() (Pair, error) {
	if p.Rate.IsZero() {
		return Pair{}, ErrZeroDenominator
	}
	return Pair{
		BaseCCY:  p.QuoteCCY,
		QuoteCCY: p.BaseCCY,
		Rate:     decimal.NewFromInt(1).Div(p.Rate).RoundBank(internalScale),
	}, nil
}

// validate enforces structural invariants on a Pair.
func (p Pair) validate() error {
	base := strings.ToUpper(strings.TrimSpace(p.BaseCCY))
	quote := strings.ToUpper(strings.TrimSpace(p.QuoteCCY))
	if len(base) != 3 || len(quote) != 3 {
		return fmt.Errorf("%w: pair ccy must be ISO 4217 alpha-3", ErrInvalidInput)
	}
	if base == quote {
		return fmt.Errorf("%w: pair base and quote ccy must differ", ErrInvalidInput)
	}
	if !p.Rate.IsPositive() {
		return fmt.Errorf("%w: pair rate must be > 0", ErrInvalidInput)
	}
	return nil
}

// Cross triangulates `a` and `b` (which must share exactly one currency — the pivot)
// and returns a new Pair representing the non-shared sides.
//
// Examples:
//
//	Cross(EURUSD=1.08, USDBRL=5.20)   → EURBRL = 1.08 × 5.20  = 5.6160
//	Cross(EURUSD=1.08, BRLUSD=0.1923) → EURBRL = 1.08 / 0.1923 ≈ 5.6160
//	Cross(USDEUR=0.925, USDBRL=5.20)  → EURBRL = 5.20 / 0.925 ≈ 5.6216
//
// Same-pair input (BaseCCY/QuoteCCY identical on both sides) is rejected.
// Zero pivot intersection is rejected.
func Cross(a, b Pair) (Pair, error) {
	if err := a.validate(); err != nil {
		return Pair{}, fmt.Errorf("a: %w", err)
	}
	if err := b.validate(); err != nil {
		return Pair{}, fmt.Errorf("b: %w", err)
	}

	a.BaseCCY = strings.ToUpper(a.BaseCCY)
	a.QuoteCCY = strings.ToUpper(a.QuoteCCY)
	b.BaseCCY = strings.ToUpper(b.BaseCCY)
	b.QuoteCCY = strings.ToUpper(b.QuoteCCY)

	// Identical pairs would compute as ratio of themselves — disallow.
	if a.BaseCCY == b.BaseCCY && a.QuoteCCY == b.QuoteCCY {
		return Pair{}, fmt.Errorf("%w: cross of identical pair is undefined", ErrInvalidInput)
	}

	// Normalise so the shared currency sits on a.QuoteCCY and b.BaseCCY.
	// Then result = a.Rate × b.Rate, BaseCCY = a.BaseCCY, QuoteCCY = b.QuoteCCY.
	pivot, ok := sharedCurrency(a, b)
	if !ok {
		return Pair{}, fmt.Errorf("%w: no shared currency between %s/%s and %s/%s",
			ErrInvalidInput, a.BaseCCY, a.QuoteCCY, b.BaseCCY, b.QuoteCCY)
	}

	if a.BaseCCY == pivot {
		inv, err := a.Invert()
		if err != nil {
			return Pair{}, err
		}
		a = inv
	}
	if b.QuoteCCY == pivot {
		inv, err := b.Invert()
		if err != nil {
			return Pair{}, err
		}
		b = inv
	}
	// Sanity (should be guaranteed by sharedCurrency logic).
	if a.QuoteCCY != pivot || b.BaseCCY != pivot {
		return Pair{}, fmt.Errorf("%w: failed to normalise via pivot %s", ErrInvalidInput, pivot)
	}

	return Pair{
		BaseCCY:  a.BaseCCY,
		QuoteCCY: b.QuoteCCY,
		Rate:     a.Rate.Mul(b.Rate).RoundBank(internalScale),
	}, nil
}

// sharedCurrency returns the currency present in both pairs (or false if none / more than one).
func sharedCurrency(a, b Pair) (string, bool) {
	matches := map[string]struct{}{}
	for _, ac := range []string{a.BaseCCY, a.QuoteCCY} {
		for _, bc := range []string{b.BaseCCY, b.QuoteCCY} {
			if ac == bc {
				matches[ac] = struct{}{}
			}
		}
	}
	if len(matches) != 1 {
		return "", false
	}
	for k := range matches {
		return k, true
	}
	return "", false
}
