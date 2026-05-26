package pricing

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Side represents the direction of a position.
type Side string

const (
	SideLong  Side = "LONG"  // bought base CCY at DealRate (positive P&L when market rate rises above DealRate)
	SideShort Side = "SHORT" // sold base CCY at DealRate (positive P&L when market rate falls below DealRate)
)

// Position models a single open FX position to revalue.
//
// Convention:
//   - NotionalBase = position size in BASE CCY (positive value; sign carried by Side)
//   - DealRate     = quoted CCY per 1 base CCY at trade time
//   - MarketRate   = quoted CCY per 1 base CCY at revaluation time
//   - Side         = LONG (long base CCY) or SHORT (short base CCY)
//
// P&L is reported in QUOTED currency.
type Position struct {
	NotionalBase decimal.Decimal
	BaseCCY      string
	QuoteCCY     string
	DealRate     decimal.Decimal
	MarketRate   decimal.Decimal
	Side         Side
}

// PositionMTM computes the unrealised P&L in quote-CCY units.
//
// Formula:
//
//	LONG:  P&L = notional_base × (market_rate − deal_rate)
//	SHORT: P&L = notional_base × (deal_rate   − market_rate)
//
// Examples:
//
//	LONG 1,000,000 EUR @ 1.0800 EURUSD, market 1.0850
//	  → P&L = 1,000,000 × (1.0850 − 1.0800) = +5,000 USD
//
//	SHORT 1,000,000 EUR @ 1.0800, market 1.0850
//	  → P&L = 1,000,000 × (1.0800 − 1.0850) = −5,000 USD
func PositionMTM(p Position) (decimal.Decimal, error) {
	if err := p.validate(); err != nil {
		return decimal.Zero, err
	}
	delta := p.MarketRate.Sub(p.DealRate)
	if p.Side == SideShort {
		delta = delta.Neg()
	}
	return p.NotionalBase.Mul(delta).RoundBank(internalScale), nil
}

func (p Position) validate() error {
	if !p.NotionalBase.IsPositive() {
		return fmt.Errorf("%w: notional_base must be > 0 (sign carried by Side)", ErrInvalidInput)
	}
	if p.BaseCCY == "" || p.QuoteCCY == "" {
		return fmt.Errorf("%w: base_ccy and quote_ccy required", ErrInvalidInput)
	}
	if p.BaseCCY == p.QuoteCCY {
		return fmt.Errorf("%w: base_ccy and quote_ccy must differ", ErrInvalidInput)
	}
	if !p.DealRate.IsPositive() {
		return fmt.Errorf("%w: deal_rate must be > 0", ErrInvalidInput)
	}
	if !p.MarketRate.IsPositive() {
		return fmt.Errorf("%w: market_rate must be > 0", ErrInvalidInput)
	}
	if p.Side != SideLong && p.Side != SideShort {
		return fmt.Errorf("%w: side must be LONG or SHORT, got %q", ErrInvalidInput, p.Side)
	}
	return nil
}

// PortfolioMTM aggregates per-position P&L. Returns a map keyed by quote CCY
// (positions whose quote currency differs sum into separate buckets).
//
// Callers wanting a single base-CCY total should convert via the cross-rate
// engine (pricing.Cross) at the rate of their choice (typically EOD PTAX or spot mid).
func PortfolioMTM(positions []Position) (map[string]decimal.Decimal, error) {
	out := make(map[string]decimal.Decimal, 4)
	for i, p := range positions {
		pnl, err := PositionMTM(p)
		if err != nil {
			return nil, fmt.Errorf("position[%d]: %w", i, err)
		}
		cur, ok := out[p.QuoteCCY]
		if !ok {
			cur = decimal.Zero
		}
		out[p.QuoteCCY] = cur.Add(pnl)
	}
	return out, nil
}
