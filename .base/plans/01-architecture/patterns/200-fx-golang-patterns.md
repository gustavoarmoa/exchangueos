# FX-GP-* — Go Patterns (40 patterns)

ExchangeOS Go coding patterns. Cited via `// FX-GP-NNN` comments in code.

## Pattern template

Each pattern follows:

```
### FX-GP-NNN — <title>
Context: when does this apply?
Problem: what does it solve?
Solution: the canonical Go shape
Example: a snippet from the repo
Anti-pattern: what NOT to do
Related: FX-GP-MMM, FX-DDD-XXX
```

## Catalog (status legend: ✅ documented, ⏳ planned)

| # | Title | Status | Where in repo |
|---|-------|--------|---------------|
| FX-GP-001 | Aggregate constructor returns `(*T, error)` | ✅ | modules/trade/domain/fxtrade.go:NewFXTrade |
| FX-GP-002 | Decimal precision: `shopspring/decimal` NEVER float | ✅ | pkg/pricing/cip.go |
| FX-GP-003 | Pointer receiver for aggregate mutation | ✅ | modules/trade/domain/fxtrade.go |
| FX-GP-004 | Sentinel errors via `errors.New(...)` + `errors.Is` | ✅ | modules/*/domain/errors.go |
| FX-GP-005 | Build-tag-gated optional bindings | ✅ | modules/*/api/grpc_server.go (`//go:build grpcgen`) |
| FX-GP-006 | Repository as interface in application package | ✅ | modules/*/application/service.go |
| FX-GP-007 | `Reconstitute<T>` helpers for persistence boundary | ✅ | modules/*/domain/reconstitute.go |
| FX-GP-008 | Context propagation as first param | ✅ | every Service method |
| FX-GP-009 | `context.WithTimeout` per outbound call | ✅ | pkg/outbox/kafka/publisher.go |
| FX-GP-010 | `defer rows.Close()` immediately after Query | ✅ | modules/*/infrastructure/postgres/repos.go |
| FX-GP-011..040 | (full list — extend on demand) | ⏳ | — |

---

## FX-GP-001 — Aggregate constructor returns `(*T, error)`

**Context:** Every DDD aggregate needs validated construction.

**Problem:** Zero-value structs bypass invariants; tests must run on guaranteed-valid aggregates.

**Solution:**

```go
func NewFXTrade(in NewTradeInput) (*FXTrade, error) {
    if err := in.validate(); err != nil {
        return nil, err   // never return half-built aggregate
    }
    // construct + record creation event
    return t, nil
}
```

**Example:** `modules/trade/domain/fxtrade.go:NewFXTrade` validates RN_FX_001/026 before returning.

**Anti-pattern:** Exposing the struct literal — `FXTrade{}` would let callers skip validation. All fields are private; the constructor + `Reconstitute*` helper are the only entry points.

**Related:** FX-GP-006 (Repository iface), FX-GP-007 (Reconstitute for DB hydration), FX-DDD-001 (Aggregate Root).

---

## FX-GP-002 — Decimal precision

**Context:** Money + FX rate arithmetic.

**Problem:** float64 silently loses precision (e.g. 0.1 + 0.2 ≠ 0.3 in IEEE-754).

**Solution:** Use `decimal.Decimal` from `shopspring/decimal` for every money/rate value. golangci-lint `forbidigo` rule bans float in CI.

**Example:** `pkg/pricing/cip.go` uses `decimal.RequireFromString` + `RoundBank(8)` for the CIP formula.

**Anti-pattern:** `float64` anywhere in money path — blocked by lint.

**Related:** FX-GP-007 (Reconstitute preserves decimal), FX-CP-* (CockroachDB DECIMAL(36,18)).

---

## FX-GP-003 — Pointer receiver for aggregate mutation

**Context:** Aggregate methods that change state.

**Problem:** Value receivers return a modified copy; the caller easily loses the new state.

**Solution:** Pointer receiver on every state-changing method of the aggregate root. Value receivers are reserved for pure read accessors.

**Example:** `func (t *FXTrade) Confirm() error` + `func (t *FXTrade) MarkSettled(ref string) error` in modules/trade/domain/fxtrade.go.

**Anti-pattern:** Mixing — a value-receiver `MarkSettled` would persist nothing.

**Related:** FX-GP-001, FX-DDD-003 (Optimistic Concurrency `version`).

---

## FX-GP-004 — Sentinel errors + `errors.Is`

**Context:** Callers need to distinguish error categories without string matching.

**Problem:** Naive `==` breaks under wrapping; `strings.Contains` is fragile.

**Solution:** Each package declares typed sentinels in `errors.go`; callers use `errors.Is(err, pkg.ErrSentinel)`.

**Example:** `modules/*/application/service.go` declares `ErrInvalidInput`/`ErrNotFound`/`ErrConflict`. gRPC adapters map them via `mapErr` to canonical gRPC codes.

**Anti-pattern:** `strings.Contains(err.Error(), "not found")`.

**Related:** FX-GP-005, FX-GRPC-* (mapErr helpers).

---

## FX-GP-005 — Build-tag-gated optional bindings

**Context:** Optional integrations (gRPC stubs, Kafka client) carry heavy deps.

**Problem:** Always-on imports slow CI + bloat images.

**Solution:** Guard the integration package with `//go:build <tag>` and pair with a no-op default sibling so the default build stays slim.

**Example:**

```
modules/*/api/grpc_server.go          // +build grpcgen
cmd/api/grpc_register_default.go      // no tag (no-op)
cmd/api/grpc_register_proto.go        // +build grpcgen
pkg/outbox/kafka/publisher.go         // +build kafka
cmd/worker/publisher_default.go       // logs only
cmd/worker/publisher_kafka.go         // +build kafka
```

`go build` works without flags; `go build -tags "grpcgen kafka"` enables both.

**Anti-pattern:** Unconditional imports that drag in unused deps.

**Related:** FX-DOC-* (image size), FX-DS-* (CI matrix).
