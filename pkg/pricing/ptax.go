package pricing

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// PTAX models the BACEN daily PTAX fixing computed from 4 windows
// surveyed at 10:00 / 11:00 / 12:00 / 13:00 São Paulo local time.
// Reference: Resolução BCB 277/2022 (Capítulo I, art. 2-3).
//
// Convention: rates are quoted as BRL per 1 USD (the public PTAX standard).
type PTAX struct {
	Date    time.Time     // business date (date-only, UTC normalised)
	Windows [4]PTAXWindow // ordered: 0=10h, 1=11h, 2=12h, 3=13h
}

// PTAXWindow holds the survey result for a single window.
// Each window aggregates dealer quotes — the BACEN methodology computes an
// arithmetic mean per window AFTER trimming outliers. For ExchangeOS purposes
// the window already arrives summarised (Bid + Ask per window).
type PTAXWindow struct {
	Hour int             // 10, 11, 12, or 13 (SP)
	Bid  decimal.Decimal // BRL per USD (bid side of dealer survey)
	Ask  decimal.Decimal // BRL per USD (ask side)
}

// Mid returns the per-window mid rate.
func (w PTAXWindow) Mid() decimal.Decimal {
	if w.Bid.IsZero() && w.Ask.IsZero() {
		return decimal.Zero
	}
	return w.Bid.Add(w.Ask).Div(decimal.NewFromInt(2))
}

// WeightedFixing returns the official PTAX rate per Resolução BCB 277:
// arithmetic mean of the 4 window mids, rounded with banker's rounding to 4 decimals
// (BACEN's published display precision).
//
// Errors:
//   - ErrInvalidInput if any window has zero bid+ask (incomplete survey)
//   - ErrInvalidInput if any window hour is outside {10,11,12,13}
//   - ErrInvalidInput if PTAX.Date is zero
func (p PTAX) WeightedFixing() (decimal.Decimal, error) {
	if err := p.validate(); err != nil {
		return decimal.Zero, err
	}
	sum := decimal.Zero
	for _, w := range p.Windows {
		sum = sum.Add(w.Mid())
	}
	mean := sum.Div(decimal.NewFromInt(4))
	return mean.RoundBank(4), nil // BACEN displays 4 decimals
}

// BidFixing returns the mean of window bids (used for "PTAX compra").
func (p PTAX) BidFixing() (decimal.Decimal, error) {
	if err := p.validate(); err != nil {
		return decimal.Zero, err
	}
	sum := decimal.Zero
	for _, w := range p.Windows {
		sum = sum.Add(w.Bid)
	}
	return sum.Div(decimal.NewFromInt(4)).RoundBank(4), nil
}

// AskFixing returns the mean of window asks (used for "PTAX venda").
func (p PTAX) AskFixing() (decimal.Decimal, error) {
	if err := p.validate(); err != nil {
		return decimal.Zero, err
	}
	sum := decimal.Zero
	for _, w := range p.Windows {
		sum = sum.Add(w.Ask)
	}
	return sum.Div(decimal.NewFromInt(4)).RoundBank(4), nil
}

func (p PTAX) validate() error {
	if p.Date.IsZero() {
		return fmt.Errorf("%w: ptax date required", ErrInvalidInput)
	}
	wantHours := [4]int{10, 11, 12, 13}
	for i, w := range p.Windows {
		if w.Hour != wantHours[i] {
			return fmt.Errorf("%w: window[%d].Hour=%d, want %d", ErrInvalidInput, i, w.Hour, wantHours[i])
		}
		if !w.Bid.IsPositive() || !w.Ask.IsPositive() {
			return fmt.Errorf("%w: window[%d] bid/ask must be > 0", ErrInvalidInput, i)
		}
		if w.Bid.GreaterThan(w.Ask) {
			return fmt.Errorf("%w: window[%d] bid (%s) > ask (%s)", ErrInvalidInput, i, w.Bid, w.Ask)
		}
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Fetcher abstraction
//
// PTAXFetcher decouples pricing logic from the HTTP/network layer. Production
// uses an OLINDA-backed implementation under modules/refdata/infrastructure/;
// tests use a fake.
// ─────────────────────────────────────────────────────────────────────────────

// PTAXFetcher retrieves a published PTAX for the given business date.
type PTAXFetcher interface {
	FetchPTAX(ctx context.Context, businessDate time.Time) (PTAX, error)
}

// PTAXFetcherFunc adapts a function to the PTAXFetcher interface.
type PTAXFetcherFunc func(ctx context.Context, businessDate time.Time) (PTAX, error)

// FetchPTAX implements PTAXFetcher.
func (f PTAXFetcherFunc) FetchPTAX(ctx context.Context, businessDate time.Time) (PTAX, error) {
	return f(ctx, businessDate)
}
