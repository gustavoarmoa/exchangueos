package domain

import (
	"time"

	"github.com/google/uuid"
)

// ReconstituteConfirmation hydrates a CFETSConfirmation from a persisted row.
//
// USAGE: persistence boundary only. Bypasses NewConfirmation validation because
// the row is presumed to have been validated when first persisted. Never call
// from application logic.
func ReconstituteConfirmation(
	id, tenantID, tradeID uuid.UUID,
	cfetsDealID string,
	status ConfirmationStatus,
	requestedAt, confirmedAt time.Time,
	rejectionReason string,
	version int,
) *CFETSConfirmation {
	return &CFETSConfirmation{
		id:              id,
		tenantID:        tenantID,
		tradeID:         tradeID,
		cfetsDealID:     cfetsDealID,
		status:          status,
		requestedAt:     requestedAt.UTC(),
		confirmedAt:     confirmedAt,
		rejectionReason: rejectionReason,
		version:         version,
	}
}
