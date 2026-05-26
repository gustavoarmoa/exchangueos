package domain

import (
	"time"

	"github.com/google/uuid"
)

type DomainEvent interface {
	EventName() string
	When() time.Time
}

type EventCycleOpened struct {
	CycleID        uuid.UUID
	TenantID       uuid.UUID
	CycleDate      time.Time
	OpenedAt       time.Time
	ScheduledClose time.Time
}

func (e EventCycleOpened) EventName() string { return "cls_cycle.opened.v1" }
func (e EventCycleOpened) When() time.Time   { return e.OpenedAt }

type EventCycleTradeAttached struct {
	CycleID    uuid.UUID
	TradeID    uuid.UUID
	OccurredAt time.Time
}

func (e EventCycleTradeAttached) EventName() string { return "cls_cycle.trade_attached.v1" }
func (e EventCycleTradeAttached) When() time.Time   { return e.OccurredAt }

type EventCyclePayInOpened struct {
	CycleID uuid.UUID
	At      time.Time
}

func (e EventCyclePayInOpened) EventName() string { return "cls_cycle.payin_opened.v1" }
func (e EventCyclePayInOpened) When() time.Time   { return e.At }

type EventCycleSettling struct {
	CycleID uuid.UUID
	At      time.Time
}

func (e EventCycleSettling) EventName() string { return "cls_cycle.settling.v1" }
func (e EventCycleSettling) When() time.Time   { return e.At }

type EventCycleClosed struct {
	CycleID  uuid.UUID
	ClosedAt time.Time
}

func (e EventCycleClosed) EventName() string { return "cls_cycle.closed.v1" }
func (e EventCycleClosed) When() time.Time   { return e.ClosedAt }

type EventCycleFailed struct {
	CycleID uuid.UUID
	Reason  string
	At      time.Time
}

func (e EventCycleFailed) EventName() string { return "cls_cycle.failed.v1" }
func (e EventCycleFailed) When() time.Time   { return e.At }
