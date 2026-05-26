package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// NewTradeInput is the immutable parameter object for NewFXTrade.
// All validation lives here so the aggregate constructor stays focused on state.
type NewTradeInput struct {
	TenantID       uuid.UUID
	ExternalRef    string
	TradeType      TradeType
	Venue          SettlementVenue
	BuyerBIC       string
	SellerBIC      string
	BoughtCurrency string
	BoughtAmount   decimal.Decimal
	SoldCurrency   string
	SoldAmount     decimal.Decimal
	DealRate       decimal.Decimal
	TradeDate      time.Time
	ValueDate      time.Time
}

// validate enforces structural invariants and the subset of RN_FX_* that the
// aggregate itself owns (currency-pair sanity, positive amounts, BIC presence).
// Cross-aggregate rules (refdata existence, sanctions, limits) belong upstream
// in the application layer.
func (i NewTradeInput) validate() error {
	if i.TenantID == uuid.Nil {
		return fmt.Errorf("%w: tenant_id required", ErrInvalidInput)
	}
	if !isValidTradeType(i.TradeType) {
		return fmt.Errorf("%w: trade_type %q", ErrInvalidInput, i.TradeType)
	}
	if !isValidVenue(i.Venue) {
		return fmt.Errorf("%w: venue %q", ErrInvalidInput, i.Venue)
	}

	if err := validateBIC(i.BuyerBIC); err != nil {
		return fmt.Errorf("buyer_bic: %w", err)
	}
	if err := validateBIC(i.SellerBIC); err != nil {
		return fmt.Errorf("seller_bic: %w", err)
	}
	if strings.EqualFold(i.BuyerBIC, i.SellerBIC) {
		return fmt.Errorf("%w: buyer and seller cannot be the same party", ErrInvalidInput)
	}

	// RN_FX_001 — currency pair valid (here: structural pair check; refdata-active check upstream).
	if err := validateCurrencyPair(i.BoughtCurrency, i.SoldCurrency); err != nil {
		return err
	}

	// RN_FX_026 — decimal precision; amounts and rate must be positive.
	if !i.BoughtAmount.IsPositive() {
		return fmt.Errorf("%w: bought_amount must be > 0", ErrInvalidInput)
	}
	if !i.SoldAmount.IsPositive() {
		return fmt.Errorf("%w: sold_amount must be > 0", ErrInvalidInput)
	}
	if !i.DealRate.IsPositive() {
		return fmt.Errorf("%w: deal_rate must be > 0", ErrInvalidInput)
	}

	if i.TradeDate.IsZero() {
		return fmt.Errorf("%w: trade_date required", ErrInvalidInput)
	}
	if i.ValueDate.IsZero() {
		return fmt.Errorf("%w: value_date required", ErrInvalidInput)
	}
	if i.ValueDate.Before(i.TradeDate) {
		return fmt.Errorf("%w: value_date must be >= trade_date", ErrInvalidInput)
	}
	return nil
}

func isValidTradeType(t TradeType) bool {
	switch t {
	case TradeTypeSpot, TradeTypeForward, TradeTypeNDF, TradeTypeSwap:
		return true
	}
	return false
}

func isValidVenue(v SettlementVenue) bool {
	switch v {
	case VenueCLS, VenueBilateral, VenueCFETS:
		return true
	}
	return false
}

func validateBIC(bic string) error {
	bic = strings.ToUpper(strings.TrimSpace(bic))
	switch len(bic) {
	case 8, 11:
		// ok — structural length check; refdata layer verifies the BIC actually exists.
	default:
		return fmt.Errorf("%w: BIC must be 8 or 11 chars, got %d", ErrInvalidInput, len(bic))
	}
	return nil
}

// validateCurrencyPair enforces RN_FX_001 at the structural level (ISO 4217 alpha-3, distinct).
// The application layer is responsible for refdata-active and CLS/CFETS eligibility checks.
func validateCurrencyPair(a, b string) error {
	if err := validateCCY(a); err != nil {
		return fmt.Errorf("bought_currency: %w", err)
	}
	if err := validateCCY(b); err != nil {
		return fmt.Errorf("sold_currency: %w", err)
	}
	if strings.EqualFold(a, b) {
		return fmt.Errorf("%w: bought and sold currency must differ (RN_FX_001)", ErrInvalidInput)
	}
	return nil
}

func validateCCY(c string) error {
	c = strings.ToUpper(strings.TrimSpace(c))
	if len(c) != 3 {
		return fmt.Errorf("%w: currency must be ISO 4217 alpha-3", ErrInvalidInput)
	}
	for _, r := range c {
		if r < 'A' || r > 'Z' {
			return fmt.Errorf("%w: currency must be alpha-only", ErrInvalidInput)
		}
	}
	return nil
}
