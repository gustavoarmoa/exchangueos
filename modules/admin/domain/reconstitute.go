package domain

import (
	"time"

	"github.com/google/uuid"
)

// ReconstituteSystemEvent hydrates a SystemEvent from a row.
//
// USAGE: persistence boundary only. Bypasses NewSystemEvent validation
// because the row is presumed to have been validated when first persisted.
func ReconstituteSystemEvent(
	id uuid.UUID,
	code EventCode,
	component, description string,
	at time.Time,
	iso20022Ref string,
) *SystemEvent {
	return &SystemEvent{
		id:          id,
		code:        code,
		component:   component,
		description: description,
		at:          at.UTC(),
		iso20022Ref: iso20022Ref,
	}
}

// ReconstituteEODJob hydrates an EODJob from a row.
func ReconstituteEODJob(
	id, tenantID uuid.UUID,
	businessDate time.Time,
	status EODStatus,
	startedAt, completedAt time.Time,
	failureReason string,
	stepsDone []string,
	version int,
) *EODJob {
	return &EODJob{
		id:            id,
		tenantID:      tenantID,
		businessDate:  time.Date(businessDate.Year(), businessDate.Month(), businessDate.Day(), 0, 0, 0, 0, time.UTC),
		status:        status,
		startedAt:     startedAt,
		completedAt:   completedAt,
		failureReason: failureReason,
		stepsDone:     append([]string(nil), stepsDone...),
		version:       version,
	}
}
