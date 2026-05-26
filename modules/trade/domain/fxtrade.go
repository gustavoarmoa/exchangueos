package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TradeType — RN_FX_002 / RN_FX_005 cite valid types.
type TradeType string

const (
	TradeTypeSpot    TradeType = "SPOT"
	TradeTypeForward TradeType = "FORWARD"
	TradeTypeNDF     TradeType = "NDF"
	TradeTypeSwap    TradeType = "SWAP"
)

// TradeStatus is the trade lifecycle state.
type TradeStatus string

const (
	StatusPending   TradeStatus = "PENDING"
	StatusConfirmed TradeStatus = "CONFIRMED"
	StatusSettling  TradeStatus = "SETTLING"
	StatusSettled   TradeStatus = "SETTLED"
	StatusCancelled TradeStatus = "CANCELLED"
	StatusRejected  TradeStatus = "REJECTED"
)

// SettlementVenue — RN_FX_010 (CLS PvP for 18 CLS-eligible CCYs).
type SettlementVenue string

const (
	VenueCLS       SettlementVenue = "CLS"
	VenueBilateral SettlementVenue = "BILATERAL"
	VenueCFETS     SettlementVenue = "CFETS"
)

// FXTrade is the aggregate root for the trade bounded context.
type FXTrade struct {
	id              uuid.UUID
	tenantID        uuid.UUID
	externalRef     string
	tradeType       TradeType
	status          TradeStatus
	venue           SettlementVenue
	buyerBIC        string
	sellerBIC       string
	boughtCurrency  string
	boughtAmount    decimal.Decimal
	soldCurrency    string
	soldAmount      decimal.Decimal
	dealRate        decimal.Decimal
	tradeDate       time.Time
	valueDate       time.Time
	version         int
	uncommittedEvts []DomainEvent
}

// NewFXTrade is the canonical constructor. Validates RN_FX_001 (currency pair valid + active),
// RN_FX_002 (spot T+2 default), and RN_FX_026 (decimal.Decimal mandatory).
//
// `currencyPair` callers verify pair existence in refdata before constructing — the
// aggregate enforces structural invariants only (no cross-aggregate refdata lookup).
func NewFXTrade(input NewTradeInput) (*FXTrade, error) {
	if err := input.validate(); err != nil {
		return nil, err
	}

	id := uuid.New()
	t := &FXTrade{
		id:             id,
		tenantID:       input.TenantID,
		externalRef:    input.ExternalRef,
		tradeType:      input.TradeType,
		status:         StatusPending,
		venue:          input.Venue,
		buyerBIC:       input.BuyerBIC,
		sellerBIC:      input.SellerBIC,
		boughtCurrency: input.BoughtCurrency,
		boughtAmount:   input.BoughtAmount,
		soldCurrency:   input.SoldCurrency,
		soldAmount:     input.SoldAmount,
		dealRate:       input.DealRate,
		tradeDate:      input.TradeDate.UTC(),
		valueDate:      input.ValueDate.UTC(),
		version:        1,
	}
	t.recordEvent(EventTradeCreated{
		TradeID:        id,
		TenantID:       input.TenantID,
		TradeType:      input.TradeType,
		Venue:          input.Venue,
		BoughtCurrency: input.BoughtCurrency,
		SoldCurrency:   input.SoldCurrency,
		DealRate:       input.DealRate,
		OccurredAt:     time.Now().UTC(),
	})
	return t, nil
}

// ── Accessors ──────────────────────────────────────────────────────────────
func (t *FXTrade) ID() uuid.UUID                 { return t.id }
func (t *FXTrade) TenantID() uuid.UUID           { return t.tenantID }
func (t *FXTrade) ExternalRef() string           { return t.externalRef }
func (t *FXTrade) Status() TradeStatus           { return t.status }
func (t *FXTrade) Venue() SettlementVenue        { return t.venue }
func (t *FXTrade) Type() TradeType               { return t.tradeType }
func (t *FXTrade) Version() int                  { return t.version }
func (t *FXTrade) BuyerBIC() string              { return t.buyerBIC }
func (t *FXTrade) SellerBIC() string             { return t.sellerBIC }
func (t *FXTrade) BoughtCurrency() string        { return t.boughtCurrency }
func (t *FXTrade) BoughtAmount() decimal.Decimal { return t.boughtAmount }
func (t *FXTrade) SoldCurrency() string          { return t.soldCurrency }
func (t *FXTrade) SoldAmount() decimal.Decimal   { return t.soldAmount }
func (t *FXTrade) DealRate() decimal.Decimal     { return t.dealRate }
func (t *FXTrade) TradeDate() time.Time          { return t.tradeDate }
func (t *FXTrade) ValueDate() time.Time          { return t.valueDate }

// PendingEvents returns events recorded since the last flush.
func (t *FXTrade) PendingEvents() []DomainEvent {
	return append([]DomainEvent(nil), t.uncommittedEvts...)
}

// MarkEventsCommitted is called by the outbox after persistence.
func (t *FXTrade) MarkEventsCommitted() { t.uncommittedEvts = nil }

// ── State transitions ──────────────────────────────────────────────────────

// Confirm transitions PENDING → CONFIRMED. Idempotent for repeat confirms.
func (t *FXTrade) Confirm() error {
	switch t.status {
	case StatusConfirmed:
		return nil
	case StatusPending:
		t.status = StatusConfirmed
		t.version++
		t.recordEvent(EventTradeConfirmed{TradeID: t.id, OccurredAt: time.Now().UTC()})
		return nil
	default:
		return fmt.Errorf("%w: cannot confirm from %s", ErrInvalidTransition, t.status)
	}
}

// Cancel transitions PENDING|CONFIRMED → CANCELLED. Forbidden after SETTLING.
func (t *FXTrade) Cancel(reason string) error {
	if reason == "" {
		return ErrCancelReasonRequired
	}
	switch t.status {
	case StatusPending, StatusConfirmed:
		t.status = StatusCancelled
		t.version++
		t.recordEvent(EventTradeCancelled{TradeID: t.id, Reason: reason, OccurredAt: time.Now().UTC()})
		return nil
	default:
		return fmt.Errorf("%w: cannot cancel from %s", ErrInvalidTransition, t.status)
	}
}

// MarkSettling transitions CONFIRMED → SETTLING (CLS opens the cycle).
func (t *FXTrade) MarkSettling() error {
	if t.status != StatusConfirmed {
		return fmt.Errorf("%w: settle requires CONFIRMED, got %s", ErrInvalidTransition, t.status)
	}
	t.status = StatusSettling
	t.version++
	t.recordEvent(EventTradeSettling{TradeID: t.id, OccurredAt: time.Now().UTC()})
	return nil
}

// MarkSettled transitions SETTLING → SETTLED with a settlement reference.
func (t *FXTrade) MarkSettled(settlementRef string) error {
	if t.status != StatusSettling {
		return fmt.Errorf("%w: settled requires SETTLING, got %s", ErrInvalidTransition, t.status)
	}
	if settlementRef == "" {
		return errors.New("settlement reference is required")
	}
	t.status = StatusSettled
	t.version++
	t.recordEvent(EventTradeSettled{TradeID: t.id, SettlementRef: settlementRef, OccurredAt: time.Now().UTC()})
	return nil
}

func (t *FXTrade) recordEvent(e DomainEvent) {
	t.uncommittedEvts = append(t.uncommittedEvts, e)
}
