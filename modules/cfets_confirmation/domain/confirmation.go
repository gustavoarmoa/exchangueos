package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ConfirmationStatus string

const (
	StatusConfirming ConfirmationStatus = "CONFIRMING"
	StatusConfirmed  ConfirmationStatus = "CONFIRMED"
	StatusUnpaired   ConfirmationStatus = "UNPAIRED"
	StatusRejected   ConfirmationStatus = "REJECTED"
)

// CFETSConfirmation is the aggregate root for one CFETS Trade Confirmation flow.
type CFETSConfirmation struct {
	id              uuid.UUID
	tenantID        uuid.UUID
	tradeID         uuid.UUID
	cfetsDealID     string
	status          ConfirmationStatus
	requestedAt     time.Time
	confirmedAt     time.Time
	rejectionReason string
	version         int
	events          []DomainEvent
}

type NewConfirmationInput struct {
	TenantID    uuid.UUID
	TradeID     uuid.UUID
	CFETSDealID string
}

func NewConfirmation(in NewConfirmationInput) (*CFETSConfirmation, error) {
	if in.TenantID == uuid.Nil || in.TradeID == uuid.Nil {
		return nil, fmt.Errorf("%w: tenant_id + trade_id required", ErrInvalidInput)
	}
	if in.CFETSDealID == "" {
		return nil, fmt.Errorf("%w: cfets_deal_id required", ErrInvalidInput)
	}
	id := uuid.New()
	c := &CFETSConfirmation{
		id:          id,
		tenantID:    in.TenantID,
		tradeID:     in.TradeID,
		cfetsDealID: in.CFETSDealID,
		status:      StatusConfirming,
		requestedAt: time.Now().UTC(),
		version:     1,
	}
	c.recordEvent(EventConfirmationRequested{
		ConfirmationID: id,
		TenantID:       in.TenantID,
		TradeID:        in.TradeID,
		CFETSDealID:    in.CFETSDealID,
		OccurredAt:     c.requestedAt,
	})
	return c, nil
}

// ── Accessors ──────────────────────────────────────────────────────────────
func (c *CFETSConfirmation) ID() uuid.UUID              { return c.id }
func (c *CFETSConfirmation) TenantID() uuid.UUID        { return c.tenantID }
func (c *CFETSConfirmation) TradeID() uuid.UUID         { return c.tradeID }
func (c *CFETSConfirmation) CFETSDealID() string        { return c.cfetsDealID }
func (c *CFETSConfirmation) Status() ConfirmationStatus { return c.status }
func (c *CFETSConfirmation) RequestedAt() time.Time     { return c.requestedAt }
func (c *CFETSConfirmation) ConfirmedAt() time.Time     { return c.confirmedAt }
func (c *CFETSConfirmation) RejectionReason() string    { return c.rejectionReason }
func (c *CFETSConfirmation) Version() int               { return c.version }
func (c *CFETSConfirmation) PendingEvents() []DomainEvent {
	return append([]DomainEvent(nil), c.events...)
}
func (c *CFETSConfirmation) MarkEventsCommitted() { c.events = nil }

// ── Lifecycle ──────────────────────────────────────────────────────────────

// MarkPaired transitions CONFIRMING → CONFIRMED.
func (c *CFETSConfirmation) MarkPaired(at time.Time) error {
	if c.status != StatusConfirming && c.status != StatusUnpaired {
		return fmt.Errorf("%w: pair requires CONFIRMING|UNPAIRED, got %s", ErrInvalidTransition, c.status)
	}
	c.status = StatusConfirmed
	c.confirmedAt = at.UTC()
	c.version++
	c.recordEvent(EventConfirmationPaired{ConfirmationID: c.id, At: c.confirmedAt})
	return nil
}

// MarkUnpaired transitions CONFIRMING → UNPAIRED (awaiting counterparty submission).
func (c *CFETSConfirmation) MarkUnpaired(at time.Time) error {
	if c.status != StatusConfirming {
		return fmt.Errorf("%w: unpair requires CONFIRMING, got %s", ErrInvalidTransition, c.status)
	}
	c.status = StatusUnpaired
	c.version++
	c.recordEvent(EventConfirmationUnpaired{ConfirmationID: c.id, At: at.UTC()})
	return nil
}

// MarkRejected transitions CONFIRMING|UNPAIRED → REJECTED.
func (c *CFETSConfirmation) MarkRejected(at time.Time, reason string) error {
	if c.status == StatusConfirmed || c.status == StatusRejected {
		return fmt.Errorf("%w: reject from %s not allowed", ErrInvalidTransition, c.status)
	}
	if reason == "" {
		return fmt.Errorf("%w: rejection reason required", ErrInvalidInput)
	}
	c.status = StatusRejected
	c.rejectionReason = reason
	c.version++
	c.recordEvent(EventConfirmationRejected{ConfirmationID: c.id, At: at.UTC(), Reason: reason})
	return nil
}

func (c *CFETSConfirmation) recordEvent(e DomainEvent) { c.events = append(c.events, e) }
