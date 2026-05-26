package domain

import (
	"time"

	"github.com/google/uuid"
)

type DomainEvent interface {
	EventName() string
	When() time.Time
}

type EventCaptureDrafted struct {
	CaptureID    uuid.UUID
	TenantID     uuid.UUID
	TradeID      uuid.UUID
	SubmitterRef string
	OccurredAt   time.Time
}

func (e EventCaptureDrafted) EventName() string { return "cfets_capture.drafted.v1" }
func (e EventCaptureDrafted) When() time.Time   { return e.OccurredAt }

type EventCaptureSubmitted struct {
	CaptureID uuid.UUID
	At        time.Time
}

func (e EventCaptureSubmitted) EventName() string { return "cfets_capture.submitted.v1" }
func (e EventCaptureSubmitted) When() time.Time   { return e.At }

type EventCaptureAcked struct {
	CaptureID   uuid.UUID
	At          time.Time
	CFETSDealID string
}

func (e EventCaptureAcked) EventName() string { return "cfets_capture.acked.v1" }
func (e EventCaptureAcked) When() time.Time   { return e.At }

type EventCaptureRejected struct {
	CaptureID uuid.UUID
	At        time.Time
	Reason    string
}

func (e EventCaptureRejected) EventName() string { return "cfets_capture.rejected.v1" }
func (e EventCaptureRejected) When() time.Time   { return e.At }

type EventCaptureNotified struct {
	CaptureID uuid.UUID
	At        time.Time
}

func (e EventCaptureNotified) EventName() string { return "cfets_capture.notified.v1" }
func (e EventCaptureNotified) When() time.Time   { return e.At }
