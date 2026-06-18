package domain

import (
	"time"

	"github.com/google/uuid"
)

// ReconstituteCapture hydrates a CFETSCapture from a persisted row.
//
// USAGE: persistence boundary only. Bypasses NewCapture validation because the
// row is presumed to have been validated when first persisted. Never call from
// application logic.
func ReconstituteCapture(
	id, tenantID, tradeID uuid.UUID,
	submitterRef, cfetsDealID string,
	status CaptureStatus,
	submittedAt, ackAt, notifiedAt time.Time,
	rejectionReason string,
	version int,
) *CFETSCapture {
	return &CFETSCapture{
		id:              id,
		tenantID:        tenantID,
		tradeID:         tradeID,
		submitterRef:    submitterRef,
		cfetsDealID:     cfetsDealID,
		status:          status,
		submittedAt:     submittedAt,
		ackAt:           ackAt,
		notifiedAt:      notifiedAt,
		rejectionReason: rejectionReason,
		version:         version,
	}
}
