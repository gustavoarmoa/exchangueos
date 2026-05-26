package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ReportStatus is the lifecycle state of a BACENReport.
type ReportStatus string

const (
	StatusPending   ReportStatus = "PENDING"
	StatusSubmitted ReportStatus = "SUBMITTED"
	StatusAccepted  ReportStatus = "ACCEPTED"
	StatusRejected  ReportStatus = "REJECTED"
)

// ReportType identifies which BACEN endpoint the report targets.
type ReportType string

const (
	ReportSISBACEN ReportType = "SISBACEN"
	ReportCCS      ReportType = "BCB-CCS"
	ReportCambio   ReportType = "BCB-CAMBIO"
)

// BACENReport — tracks the submission of a regulatory report.
type BACENReport struct {
	id              uuid.UUID
	tenantID        uuid.UUID
	reportType      ReportType
	referenceDate   time.Time
	payloadHash     string // SHA-256 of the serialised payload (audit trail)
	status          ReportStatus
	submittedAt     time.Time
	respondedAt     time.Time
	rejectionReason string
	version         int
}

// NewBACENReportInput parameterises construction.
type NewBACENReportInput struct {
	TenantID      uuid.UUID
	ReportType    ReportType
	ReferenceDate time.Time
	PayloadHash   string
}

// NewBACENReport constructs a report in PENDING state.
func NewBACENReport(in NewBACENReportInput) (*BACENReport, error) {
	if in.TenantID == uuid.Nil {
		return nil, fmt.Errorf("%w: tenant_id required", ErrInvalidInput)
	}
	switch in.ReportType {
	case ReportSISBACEN, ReportCCS, ReportCambio:
	default:
		return nil, fmt.Errorf("%w: report_type %q unsupported", ErrInvalidInput, in.ReportType)
	}
	if in.ReferenceDate.IsZero() {
		return nil, fmt.Errorf("%w: reference_date required", ErrInvalidInput)
	}
	if strings.TrimSpace(in.PayloadHash) == "" {
		return nil, fmt.Errorf("%w: payload_hash required (sha256 hex)", ErrInvalidInput)
	}
	return &BACENReport{
		id:            uuid.New(),
		tenantID:      in.TenantID,
		reportType:    in.ReportType,
		referenceDate: in.ReferenceDate.UTC(),
		payloadHash:   strings.ToLower(in.PayloadHash),
		status:        StatusPending,
		version:       1,
	}, nil
}

// Accessors
func (r *BACENReport) ID() uuid.UUID            { return r.id }
func (r *BACENReport) TenantID() uuid.UUID      { return r.tenantID }
func (r *BACENReport) Type() ReportType         { return r.reportType }
func (r *BACENReport) ReferenceDate() time.Time { return r.referenceDate }
func (r *BACENReport) PayloadHash() string      { return r.payloadHash }
func (r *BACENReport) Status() ReportStatus     { return r.status }
func (r *BACENReport) SubmittedAt() time.Time   { return r.submittedAt }
func (r *BACENReport) RespondedAt() time.Time   { return r.respondedAt }
func (r *BACENReport) RejectionReason() string  { return r.rejectionReason }
func (r *BACENReport) Version() int             { return r.version }

// MarkSubmitted transitions PENDING → SUBMITTED.
func (r *BACENReport) MarkSubmitted(at time.Time) error {
	if r.status != StatusPending {
		return fmt.Errorf("%w: submit requires PENDING, got %s", ErrInvalidTransition, r.status)
	}
	r.status = StatusSubmitted
	r.submittedAt = at.UTC()
	r.version++
	return nil
}

// MarkAccepted transitions SUBMITTED → ACCEPTED.
func (r *BACENReport) MarkAccepted(at time.Time) error {
	if r.status != StatusSubmitted {
		return fmt.Errorf("%w: accept requires SUBMITTED, got %s", ErrInvalidTransition, r.status)
	}
	r.status = StatusAccepted
	r.respondedAt = at.UTC()
	r.version++
	return nil
}

// MarkRejected transitions SUBMITTED → REJECTED with a reason.
func (r *BACENReport) MarkRejected(at time.Time, reason string) error {
	if r.status != StatusSubmitted {
		return fmt.Errorf("%w: reject requires SUBMITTED, got %s", ErrInvalidTransition, r.status)
	}
	if reason == "" {
		return fmt.Errorf("%w: rejection reason required", ErrInvalidInput)
	}
	r.status = StatusRejected
	r.rejectionReason = reason
	r.respondedAt = at.UTC()
	r.version++
	return nil
}
