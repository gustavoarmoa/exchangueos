// Package eventbus — minimal in-process event dispatcher.
//
// This is a stop-gap until the Kafka outbox lands (MS-023g). It lets the
// container wire cross-context handlers (e.g. quote.accepted.v1 → trade.BookTrade)
// today without a real broker.
//
// Contract:
//   - Synchronous Publish (handlers run inline; errors logged but not returned).
//   - Subscribe(eventName, handler) registers by canonical event name.
//   - Thread-safe.
package eventbus

import (
	"context"
	"sync"
)

// Event is what handlers receive. Implementations live in module domain packages;
// the bus only needs the name + opaque payload.
type Event interface {
	EventName() string
}

// Handler is called once per matching event. Errors are returned for observability
// (the bus logs and continues; the publisher does not block on handler failure).
type Handler func(ctx context.Context, e Event) error

// Bus is a lightweight pub/sub. Construct via New.
type Bus struct {
	mu   sync.RWMutex
	subs map[string][]Handler
}

// New constructs an empty Bus.
func New() *Bus { return &Bus{subs: make(map[string][]Handler)} }

// Subscribe registers a handler for the given canonical event name.
func (b *Bus) Subscribe(eventName string, h Handler) {
	if h == nil || eventName == "" {
		return
	}
	b.mu.Lock()
	b.subs[eventName] = append(b.subs[eventName], h)
	b.mu.Unlock()
}

// Publish dispatches events to all matching handlers, in subscription order.
// Returns the slice of handler errors (or nil) so callers can log/observe.
func (b *Bus) Publish(ctx context.Context, events ...Event) []error {
	if len(events) == 0 {
		return nil
	}
	var errs []error
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, e := range events {
		for _, h := range b.subs[e.EventName()] {
			if err := h(ctx, e); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

// Counts returns the number of subscribers per event name — used by tests.
func (b *Bus) Counts() map[string]int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make(map[string]int, len(b.subs))
	for k, v := range b.subs {
		out[k] = len(v)
	}
	return out
}
