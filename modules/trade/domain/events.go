package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// DomainEvent is the marker interface for trade-domain events.
// Events are appended to FXTrade.uncommittedEvts and flushed by an outbox.
type DomainEvent interface {
	EventName() string
	When() time.Time
}

type EventTradeCreated struct {
	TradeID        uuid.UUID
	TenantID       uuid.UUID
	TradeType      TradeType
	Venue          SettlementVenue
	BoughtCurrency string
	SoldCurrency   string
	DealRate       decimal.Decimal
	OccurredAt     time.Time
}

func (e EventTradeCreated) EventName() string { return "trade.created.v1" }
func (e EventTradeCreated) When() time.Time   { return e.OccurredAt }

type EventTradeConfirmed struct {
	TradeID    uuid.UUID
	OccurredAt time.Time
}

func (e EventTradeConfirmed) EventName() string { return "trade.confirmed.v1" }
func (e EventTradeConfirmed) When() time.Time   { return e.OccurredAt }

type EventTradeCancelled struct {
	TradeID    uuid.UUID
	Reason     string
	OccurredAt time.Time
}

func (e EventTradeCancelled) EventName() string { return "trade.cancelled.v1" }
func (e EventTradeCancelled) When() time.Time   { return e.OccurredAt }

type EventTradeSettling struct {
	TradeID    uuid.UUID
	OccurredAt time.Time
}

func (e EventTradeSettling) EventName() string { return "trade.settling.v1" }
func (e EventTradeSettling) When() time.Time   { return e.OccurredAt }

type EventTradeSettled struct {
	TradeID       uuid.UUID
	SettlementRef string
	OccurredAt    time.Time
}

func (e EventTradeSettled) EventName() string { return "trade.settled.v1" }
func (e EventTradeSettled) When() time.Time   { return e.OccurredAt }
