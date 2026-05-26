package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type DomainEvent interface {
	EventName() string
	When() time.Time
}

type EventPayInCreated struct {
	InstructionID uuid.UUID
	TenantID      uuid.UUID
	CycleID       uuid.UUID
	Currency      string
	Amount        decimal.Decimal
	Band          DeadlineBand
	Deadline      time.Time
	OccurredAt    time.Time
}

func (e EventPayInCreated) EventName() string { return "payin.created.v1" }
func (e EventPayInCreated) When() time.Time   { return e.OccurredAt }

type EventPayInSubmitted struct {
	InstructionID uuid.UUID
	At            time.Time
}

func (e EventPayInSubmitted) EventName() string { return "payin.submitted.v1" }
func (e EventPayInSubmitted) When() time.Time   { return e.At }

type EventPayInConfirmed struct {
	InstructionID uuid.UUID
	At            time.Time
}

func (e EventPayInConfirmed) EventName() string { return "payin.confirmed.v1" }
func (e EventPayInConfirmed) When() time.Time   { return e.At }

type EventPayInFailed struct {
	InstructionID uuid.UUID
	Reason        string
	At            time.Time
}

func (e EventPayInFailed) EventName() string { return "payin.failed.v1" }
func (e EventPayInFailed) When() time.Time   { return e.At }
