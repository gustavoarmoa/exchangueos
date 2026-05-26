package domain

import (
	"time"

	"github.com/google/uuid"
)

// ReconstituteCycle rebuilds a CLSCycle from persisted state. Persistence-boundary only.
func ReconstituteCycle(
	id, tenantID uuid.UUID,
	cycleDate time.Time,
	status CycleStatus,
	openedAt, pin1, pin2, pin3, scheduledClose, closedAt time.Time,
	failureReason string,
	tradeIDs []uuid.UUID,
	version int,
) *CLSCycle {
	return &CLSCycle{
		id:             id,
		tenantID:       tenantID,
		cycleDate:      cycleDate.UTC(),
		status:         status,
		openedAt:       openedAt.UTC(),
		pin1Deadline:   pin1.UTC(),
		pin2Deadline:   pin2.UTC(),
		pin3Deadline:   pin3.UTC(),
		scheduledClose: scheduledClose.UTC(),
		closedAt:       closedAt.UTC(),
		failureReason:  failureReason,
		tradeIDs:       append([]uuid.UUID(nil), tradeIDs...),
		version:        version,
	}
}
