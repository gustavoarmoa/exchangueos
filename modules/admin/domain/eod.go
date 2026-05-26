package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EODStatus is the lifecycle state of an EOD batch run.
type EODStatus string

const (
	EODPending   EODStatus = "PENDING"
	EODRunning   EODStatus = "RUNNING"
	EODCompleted EODStatus = "COMPLETED"
	EODFailed    EODStatus = "FAILED"
)

// EODJob tracks the end-of-day batch orchestration (PTAX fixing → MTM → position
// snapshot → BACEN report submission).
type EODJob struct {
	id            uuid.UUID
	tenantID      uuid.UUID
	businessDate  time.Time
	status        EODStatus
	startedAt     time.Time
	completedAt   time.Time
	failureReason string
	stepsDone     []string
	version       int
}

// NewEODJobInput parameterises construction.
type NewEODJobInput struct {
	TenantID     uuid.UUID
	BusinessDate time.Time
}

// NewEODJob constructs a PENDING job for (tenant, business_date).
func NewEODJob(in NewEODJobInput) (*EODJob, error) {
	if in.TenantID == uuid.Nil {
		return nil, fmt.Errorf("%w: tenant_id required", ErrInvalidInput)
	}
	if in.BusinessDate.IsZero() {
		return nil, fmt.Errorf("%w: business_date required", ErrInvalidInput)
	}
	return &EODJob{
		id:           uuid.New(),
		tenantID:     in.TenantID,
		businessDate: time.Date(in.BusinessDate.Year(), in.BusinessDate.Month(), in.BusinessDate.Day(), 0, 0, 0, 0, time.UTC),
		status:       EODPending,
		version:      1,
	}, nil
}

func (j *EODJob) ID() uuid.UUID           { return j.id }
func (j *EODJob) TenantID() uuid.UUID     { return j.tenantID }
func (j *EODJob) BusinessDate() time.Time { return j.businessDate }
func (j *EODJob) Status() EODStatus       { return j.status }
func (j *EODJob) StartedAt() time.Time    { return j.startedAt }
func (j *EODJob) CompletedAt() time.Time  { return j.completedAt }
func (j *EODJob) StepsDone() []string     { return append([]string(nil), j.stepsDone...) }
func (j *EODJob) FailureReason() string   { return j.failureReason }
func (j *EODJob) Version() int            { return j.version }

// Start transitions PENDING → RUNNING.
func (j *EODJob) Start(at time.Time) error {
	if j.status != EODPending {
		return fmt.Errorf("%w: start requires PENDING, got %s", ErrInvalidTransition, j.status)
	}
	j.status = EODRunning
	j.startedAt = at.UTC()
	j.version++
	return nil
}

// MarkStep records an idempotent step name (e.g. "PTAX", "MTM", "POSITION_SNAPSHOT").
func (j *EODJob) MarkStep(step string) error {
	if j.status != EODRunning {
		return fmt.Errorf("%w: mark-step requires RUNNING, got %s", ErrInvalidTransition, j.status)
	}
	if step == "" {
		return fmt.Errorf("%w: step name required", ErrInvalidInput)
	}
	for _, s := range j.stepsDone {
		if s == step {
			return nil
		}
	}
	j.stepsDone = append(j.stepsDone, step)
	j.version++
	return nil
}

// Complete transitions RUNNING → COMPLETED.
func (j *EODJob) Complete(at time.Time) error {
	if j.status != EODRunning {
		return fmt.Errorf("%w: complete requires RUNNING, got %s", ErrInvalidTransition, j.status)
	}
	j.status = EODCompleted
	j.completedAt = at.UTC()
	j.version++
	return nil
}

// Fail moves any non-terminal status to FAILED.
func (j *EODJob) Fail(at time.Time, reason string) error {
	switch j.status {
	case EODCompleted, EODFailed:
		return fmt.Errorf("%w: cannot fail from %s", ErrInvalidTransition, j.status)
	}
	if reason == "" {
		return fmt.Errorf("%w: reason required", ErrInvalidInput)
	}
	j.status = EODFailed
	j.failureReason = reason
	j.completedAt = at.UTC()
	j.version++
	return nil
}
