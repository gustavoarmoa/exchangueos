package domain

import (
	"time"

	"github.com/google/uuid"
)

type DomainEvent interface {
	EventName() string
	When() time.Time
}

type EventConfirmationRequested struct {
	ConfirmationID uuid.UUID
	TenantID       uuid.UUID
	TradeID        uuid.UUID
	CFETSDealID    string
	OccurredAt     time.Time
}

func (e EventConfirmationRequested) EventName() string { return "cfets_confirmation.requested.v1" }
func (e EventConfirmationRequested) When() time.Time   { return e.OccurredAt }

type EventConfirmationPaired struct {
	ConfirmationID uuid.UUID
	At             time.Time
}

func (e EventConfirmationPaired) EventName() string { return "cfets_confirmation.paired.v1" }
func (e EventConfirmationPaired) When() time.Time   { return e.At }

type EventConfirmationUnpaired struct {
	ConfirmationID uuid.UUID
	At             time.Time
}

func (e EventConfirmationUnpaired) EventName() string { return "cfets_confirmation.unpaired.v1" }
func (e EventConfirmationUnpaired) When() time.Time   { return e.At }

type EventConfirmationRejected struct {
	ConfirmationID uuid.UUID
	At             time.Time
	Reason         string
}

func (e EventConfirmationRejected) EventName() string { return "cfets_confirmation.rejected.v1" }
func (e EventConfirmationRejected) When() time.Time   { return e.At }
