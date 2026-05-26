// Package health provides liveness + readiness probe primitives composable from multiple checks.
package health

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Status string

const (
	StatusServing    Status = "SERVING"
	StatusNotServing Status = "NOT_SERVING"
	StatusDegraded   Status = "DEGRADED"
	StatusUnknown    Status = "UNKNOWN"
)

// CheckFunc runs a single readiness probe. Should respect ctx deadline.
type CheckFunc func(ctx context.Context) error

// Registry aggregates named checks for /readyz.
type Registry struct {
	mu     sync.RWMutex
	checks map[string]CheckFunc
}

func NewRegistry() *Registry {
	return &Registry{checks: make(map[string]CheckFunc)}
}

func (r *Registry) Register(name string, fn CheckFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.checks[name] = fn
}

// Result of evaluating the registry.
type Result struct {
	Status  Status
	Details map[string]string
	At      time.Time
}

// Check runs all probes in parallel with a per-probe timeout.
// If any probe fails, status = NOT_SERVING. If any times out (without others failing), DEGRADED.
func (r *Registry) Check(ctx context.Context, perProbeTimeout time.Duration) Result {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.checks) == 0 {
		return Result{Status: StatusServing, Details: map[string]string{}, At: time.Now().UTC()}
	}

	type pr struct {
		name string
		err  error
		dur  time.Duration
	}
	results := make(chan pr, len(r.checks))
	var wg sync.WaitGroup
	for name, fn := range r.checks {
		wg.Add(1)
		go func(name string, fn CheckFunc) {
			defer wg.Done()
			cctx, cancel := context.WithTimeout(ctx, perProbeTimeout)
			defer cancel()
			start := time.Now()
			err := fn(cctx)
			results <- pr{name: name, err: err, dur: time.Since(start)}
		}(name, fn)
	}
	wg.Wait()
	close(results)

	details := make(map[string]string, len(r.checks))
	status := StatusServing
	for p := range results {
		if p.err == nil {
			details[p.name] = "ok"
			continue
		}
		details[p.name] = p.err.Error()
		if errors.Is(p.err, context.DeadlineExceeded) {
			if status == StatusServing {
				status = StatusDegraded
			}
		} else {
			status = StatusNotServing
		}
	}
	return Result{Status: status, Details: details, At: time.Now().UTC()}
}
