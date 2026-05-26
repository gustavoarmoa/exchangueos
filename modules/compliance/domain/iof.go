package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// IOFComputation — Imposto sobre Operações Financeiras applied to an FX trade.
// Cite RN_FX_037; rates per Decreto 12.499/2025.
type IOFComputation struct {
	id            uuid.UUID
	tradeID       uuid.UUID
	tenantID      uuid.UUID
	operationType string // e.g. "EXPORT", "IMPORT", "TRAVEL", "CREDIT_CARD", "LOAN", "INVESTMENT"
	notional      decimal.Decimal
	notionalCCY   string
	rate          decimal.Decimal // as fraction (0.0038 for 0.38%)
	iofAmount     decimal.Decimal // notional * rate, in notionalCCY
	computedAt    time.Time
}

// NewIOFInput parameterises construction.
type NewIOFInput struct {
	TenantID      uuid.UUID
	TradeID       uuid.UUID
	OperationType string
	Notional      decimal.Decimal
	NotionalCCY   string
	Rate          decimal.Decimal
}

// NewIOFComputation validates + computes the tax amount.
func NewIOFComputation(in NewIOFInput) (*IOFComputation, error) {
	if in.TenantID == uuid.Nil || in.TradeID == uuid.Nil {
		return nil, fmt.Errorf("%w: tenant_id + trade_id required", ErrInvalidInput)
	}
	if !in.Notional.IsPositive() {
		return nil, fmt.Errorf("%w: notional must be > 0", ErrInvalidInput)
	}
	if in.Rate.IsNegative() {
		return nil, fmt.Errorf("%w: rate cannot be negative", ErrInvalidInput)
	}
	if in.Rate.GreaterThan(decimal.RequireFromString("1")) {
		return nil, fmt.Errorf("%w: rate must be a fraction (<= 1), got %s", ErrInvalidInput, in.Rate)
	}
	if len(in.NotionalCCY) != 3 {
		return nil, fmt.Errorf("%w: notional_ccy must be ISO 4217 alpha-3", ErrInvalidInput)
	}
	if in.OperationType == "" {
		return nil, fmt.Errorf("%w: operation_type required", ErrInvalidInput)
	}
	amount := in.Notional.Mul(in.Rate).RoundBank(2) // 2 decimals for tax money
	return &IOFComputation{
		id:            uuid.New(),
		tenantID:      in.TenantID,
		tradeID:       in.TradeID,
		operationType: in.OperationType,
		notional:      in.Notional,
		notionalCCY:   in.NotionalCCY,
		rate:          in.Rate,
		iofAmount:     amount,
		computedAt:    time.Now().UTC(),
	}, nil
}

// Accessors
func (i *IOFComputation) ID() uuid.UUID              { return i.id }
func (i *IOFComputation) TenantID() uuid.UUID        { return i.tenantID }
func (i *IOFComputation) TradeID() uuid.UUID         { return i.tradeID }
func (i *IOFComputation) OperationType() string      { return i.operationType }
func (i *IOFComputation) Notional() decimal.Decimal  { return i.notional }
func (i *IOFComputation) NotionalCCY() string        { return i.notionalCCY }
func (i *IOFComputation) Rate() decimal.Decimal      { return i.rate }
func (i *IOFComputation) IOFAmount() decimal.Decimal { return i.iofAmount }
func (i *IOFComputation) ComputedAt() time.Time      { return i.computedAt }
