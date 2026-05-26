package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ReconstitutePayIn hydrates a PayInInstruction from persisted state.
//
// USAGE: persistence boundary only. Bypasses NewPayInInstruction validation
// because the row is presumed to have been validated when first persisted.
// Never call from application logic — use NewPayInInstruction instead.
func ReconstitutePayIn(
	id, tenantID, cycleID uuid.UUID,
	currency string,
	amount decimal.Decimal,
	band DeadlineBand,
	deadline time.Time,
	status PayInStatus,
	submittedAt, confirmedAt time.Time,
	failureReason string,
	version int,
) *PayInInstruction {
	return &PayInInstruction{
		id:            id,
		tenantID:      tenantID,
		cycleID:       cycleID,
		currency:      currency,
		amount:        amount,
		band:          band,
		deadline:      deadline.UTC(),
		status:        status,
		submittedAt:   submittedAt,
		confirmedAt:   confirmedAt,
		failureReason: failureReason,
		version:       version,
	}
}
