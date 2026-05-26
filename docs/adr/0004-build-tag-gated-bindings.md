# ADR 0004 — Build-tag-gated optional bindings (grpcgen, kafka)

- Status: Accepted
- Date: 2026-05-24

## Context

Two integrations carry heavy dependencies that we want to keep out of the default build:

- **gRPC adapter binding** — depends on `proto/gen/` materialised by `buf generate`. Bringing the generated package in unconditionally would force every contributor to run buf locally.
- **Kafka client** — `franz-go/kgo` is ~150kLOC + transitive deps. Tests + local-only flows don't need it.

## Decision

**Build tags + paired no-op + real files.**

Two tag conventions:

- `//go:build grpcgen` — gRPC adapters in `modules/<bc>/api/grpc_server.go` + service registration in `cmd/api/grpc_register_proto.go`. Paired no-op: `cmd/api/grpc_register_default.go` with no tag.
- `//go:build kafka` — `pkg/outbox/kafka/publisher.go` + `cmd/worker/publisher_kafka.go`. Paired no-op: `cmd/worker/publisher_default.go` with no tag (logs only).

Default `go build` works without flags; enable via:

```bash
task proto:gen          # produce proto/gen
go build -tags grpcgen ./...
go build -tags kafka ./...
go build -tags "grpcgen kafka" ./...
```

## Consequences

### Positive

- **Fast default build** — no buf or franz-go in the default code path
- **Clean dep tree** for `go vet`, `golangci-lint`, IDE indexing
- **Explicit opt-in** — operators choose what to enable for their deployment
- **No conditional imports** — paired-file pattern keeps each file unambiguous

### Negative

- **Two files per integration** — slight boilerplate
- **CI must run both build tags** — addressed in `.github/workflows/ci.yml` matrix
- **First-time contributors must understand the pattern** — addressed in onboarding doc + `FX-GP-005` pattern catalog

## Alternatives considered

- **Always-on imports** — fast to start but bloats binary + slows lint
- **Plugins (Go shared objects)** — runtime cost + dynamic-linking pain on multi-arch
- **Separate modules** — would force a multi-module repo + go.work file; complexity > benefit at our size
