package domain

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// SpotRate — a single (base, quote) spot mid quoted at a point in time.
type SpotRate struct {
	BaseCCY  string
	QuoteCCY string
	Mid      decimal.Decimal
	AsOf     time.Time
}

// SpotRateBook is a thread-safe in-memory cache of spot mids keyed by (base, quote).
// Pricing engines look up rates here; market-data feeders write rates here on tick.
// Rates older than the configured MaxAge are reported as stale via LookupFresh.
type SpotRateBook struct {
	mu     sync.RWMutex
	byPair map[string]SpotRate
	maxAge time.Duration
}

// NewSpotRateBook constructs an empty book. `maxAge <= 0` disables freshness checks.
func NewSpotRateBook(maxAge time.Duration) *SpotRateBook {
	return &SpotRateBook{byPair: make(map[string]SpotRate), maxAge: maxAge}
}

// Put stores a spot rate. Validates input.
func (b *SpotRateBook) Put(r SpotRate) error {
	base := strings.ToUpper(strings.TrimSpace(r.BaseCCY))
	quote := strings.ToUpper(strings.TrimSpace(r.QuoteCCY))
	if len(base) != 3 || len(quote) != 3 || base == quote {
		return fmt.Errorf("%w: invalid pair %s/%s", ErrInvalidInput, base, quote)
	}
	if !r.Mid.IsPositive() {
		return fmt.Errorf("%w: mid must be > 0", ErrInvalidInput)
	}
	if r.AsOf.IsZero() {
		r.AsOf = time.Now().UTC()
	}
	r.BaseCCY, r.QuoteCCY = base, quote
	b.mu.Lock()
	b.byPair[pairKey(base, quote)] = r
	b.mu.Unlock()
	return nil
}

// Lookup returns the rate for a pair regardless of staleness.
func (b *SpotRateBook) Lookup(baseCCY, quoteCCY string) (SpotRate, bool) {
	base := strings.ToUpper(strings.TrimSpace(baseCCY))
	quote := strings.ToUpper(strings.TrimSpace(quoteCCY))
	b.mu.RLock()
	defer b.mu.RUnlock()
	r, ok := b.byPair[pairKey(base, quote)]
	return r, ok
}

// LookupFresh returns the rate only if it is within MaxAge of now. Returns
// (rate, true, nil) on success, (rate, false, ErrStale) when present but stale,
// and (zero, false, ErrNotFound) when the pair is unknown.
func (b *SpotRateBook) LookupFresh(baseCCY, quoteCCY string, now time.Time) (SpotRate, bool, error) {
	r, ok := b.Lookup(baseCCY, quoteCCY)
	if !ok {
		return SpotRate{}, false, ErrNotFound
	}
	if b.maxAge > 0 && !now.IsZero() && now.Sub(r.AsOf) > b.maxAge {
		return r, false, ErrStale
	}
	return r, true, nil
}

func pairKey(base, quote string) string { return base + "/" + quote }
