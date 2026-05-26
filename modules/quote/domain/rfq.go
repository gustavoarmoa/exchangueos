package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RFQStatus is the RFQ lifecycle state.
type RFQStatus string

const (
	RFQRequested RFQStatus = "REQUESTED"
	RFQQuoted    RFQStatus = "QUOTED"
	RFQAccepted  RFQStatus = "ACCEPTED"
	RFQRejected  RFQStatus = "REJECTED"
	RFQExpired   RFQStatus = "EXPIRED"
)

// RFQ — Request-For-Quote aggregate.
type RFQ struct {
	id        uuid.UUID
	tenantID  uuid.UUID
	requester string
	baseCCY   string
	quoteCCY  string
	status    RFQStatus
	quoteIDs  []uuid.UUID
	createdAt time.Time
	version   int
	events    []DomainEvent
}

// NewRFQInput parameterises construction.
type NewRFQInput struct {
	TenantID  uuid.UUID
	Requester string
	BaseCCY   string
	QuoteCCY  string
}

// NewRFQ constructs a fresh RFQ in REQUESTED state.
func NewRFQ(in NewRFQInput) (*RFQ, error) {
	if in.TenantID == uuid.Nil {
		return nil, fmt.Errorf("%w: tenant_id required", ErrInvalidInput)
	}
	if in.Requester == "" {
		return nil, fmt.Errorf("%w: requester required", ErrInvalidInput)
	}
	if err := validatePair(in.BaseCCY, in.QuoteCCY); err != nil {
		return nil, err
	}
	id := uuid.New()
	r := &RFQ{
		id:        id,
		tenantID:  in.TenantID,
		requester: in.Requester,
		baseCCY:   in.BaseCCY,
		quoteCCY:  in.QuoteCCY,
		status:    RFQRequested,
		createdAt: time.Now().UTC(),
		version:   1,
	}
	r.recordEvent(EventRFQRequested{
		RFQID:      id,
		TenantID:   in.TenantID,
		BaseCCY:    in.BaseCCY,
		QuoteCCY:   in.QuoteCCY,
		Requester:  in.Requester,
		OccurredAt: r.createdAt,
	})
	return r, nil
}

func (r *RFQ) ID() uuid.UUID         { return r.id }
func (r *RFQ) TenantID() uuid.UUID   { return r.tenantID }
func (r *RFQ) Requester() string     { return r.requester }
func (r *RFQ) BaseCCY() string       { return r.baseCCY }
func (r *RFQ) QuoteCCY() string      { return r.quoteCCY }
func (r *RFQ) Status() RFQStatus     { return r.status }
func (r *RFQ) Version() int          { return r.version }
func (r *RFQ) CreatedAt() time.Time  { return r.createdAt }
func (r *RFQ) QuoteIDs() []uuid.UUID { return append([]uuid.UUID(nil), r.quoteIDs...) }
func (r *RFQ) PendingEvents() []DomainEvent {
	return append([]DomainEvent(nil), r.events...)
}
func (r *RFQ) MarkEventsCommitted() { r.events = nil }

// AttachQuote transitions REQUESTED → QUOTED on first quote, or appends additional quotes if already QUOTED.
func (r *RFQ) AttachQuote(quoteID uuid.UUID) error {
	if r.status != RFQRequested && r.status != RFQQuoted {
		return fmt.Errorf("%w: cannot attach quote from %s", ErrInvalidTransition, r.status)
	}
	r.quoteIDs = append(r.quoteIDs, quoteID)
	if r.status == RFQRequested {
		r.status = RFQQuoted
	}
	r.version++
	r.recordEvent(EventRFQQuoted{RFQID: r.id, QuoteID: quoteID, OccurredAt: time.Now().UTC()})
	return nil
}

// Accept transitions QUOTED → ACCEPTED, capturing which quoteID was chosen.
func (r *RFQ) Accept(quoteID uuid.UUID, actor string) error {
	if r.status != RFQQuoted {
		return fmt.Errorf("%w: accept requires QUOTED, got %s", ErrInvalidTransition, r.status)
	}
	found := false
	for _, q := range r.quoteIDs {
		if q == quoteID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("%w: quote_id not part of this RFQ", ErrInvalidInput)
	}
	r.status = RFQAccepted
	r.version++
	r.recordEvent(EventRFQAccepted{
		RFQID:      r.id,
		QuoteID:    quoteID,
		Actor:      actor,
		OccurredAt: time.Now().UTC(),
	})
	return nil
}

// Reject transitions REQUESTED|QUOTED → REJECTED.
func (r *RFQ) Reject(reason string) error {
	if r.status != RFQRequested && r.status != RFQQuoted {
		return fmt.Errorf("%w: reject from %s not allowed", ErrInvalidTransition, r.status)
	}
	if reason == "" {
		return fmt.Errorf("%w: reject reason required", ErrInvalidInput)
	}
	r.status = RFQRejected
	r.version++
	r.recordEvent(EventRFQRejected{RFQID: r.id, Reason: reason, OccurredAt: time.Now().UTC()})
	return nil
}

// Expire transitions REQUESTED|QUOTED → EXPIRED when its TTL passes.
func (r *RFQ) Expire() error {
	if r.status != RFQRequested && r.status != RFQQuoted {
		return fmt.Errorf("%w: expire from %s not allowed", ErrInvalidTransition, r.status)
	}
	r.status = RFQExpired
	r.version++
	r.recordEvent(EventRFQExpired{RFQID: r.id, OccurredAt: time.Now().UTC()})
	return nil
}

func (r *RFQ) recordEvent(e DomainEvent) { r.events = append(r.events, e) }
