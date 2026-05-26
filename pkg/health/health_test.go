package health

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRegistry_EmptyReturnsServing(t *testing.T) {
	r := NewRegistry()
	res := r.Check(context.Background(), 100*time.Millisecond)
	if res.Status != StatusServing {
		t.Fatalf("want SERVING, got %s", res.Status)
	}
}

func TestRegistry_OneFailingMakesNotServing(t *testing.T) {
	r := NewRegistry()
	r.Register("ok", func(ctx context.Context) error { return nil })
	r.Register("bad", func(ctx context.Context) error { return errors.New("boom") })
	res := r.Check(context.Background(), 100*time.Millisecond)
	if res.Status != StatusNotServing {
		t.Fatalf("want NOT_SERVING, got %s", res.Status)
	}
	if res.Details["bad"] != "boom" {
		t.Fatalf("expected boom detail, got %q", res.Details["bad"])
	}
}

func TestRegistry_TimeoutDegrades(t *testing.T) {
	r := NewRegistry()
	r.Register("slow", func(ctx context.Context) error {
		select {
		case <-time.After(500 * time.Millisecond):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	res := r.Check(context.Background(), 20*time.Millisecond)
	if res.Status != StatusDegraded {
		t.Fatalf("want DEGRADED, got %s", res.Status)
	}
}
