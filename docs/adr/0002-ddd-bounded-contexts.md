# ADR 0002 — 14 bounded contexts with DDD aggregates

- Status: Accepted
- Date: 2026-05-24

## Context

FX trading + settlement spans many concerns: pricing, quote/RFQ, trade booking, CLS settlement, CFETS capture, compliance (BACEN), risk limits, position keeping, admin events. A single monolithic domain model would entangle these concerns.

## Decision

Split into **14 bounded contexts** under `modules/<bc>/`:

trade · quote · amendment · cls_settlement · payin · netreport · cfets_capture · cfets_confirmation · settlement · refdata · admin · risk · position · compliance

Each BC has the same internal structure:

```
modules/<bc>/
├── domain/         — aggregates + value objects + sentinel errors + events
├── application/    — Service + Repository interfaces (pure Go, no infra deps)
├── infrastructure/
│   ├── memory/     — in-memory repo (tests + bootstrap)
│   └── postgres/   — pgx-backed repo (production)
└── api/            — gRPC adapter (//go:build grpcgen)
```

Aggregates reference each other by ID only (FX-DDD-002). Cross-BC reactions go through event bus (in-process now, Kafka outbox at scale).

## Consequences

### Positive

- **Clear ownership** — each BC has a single team responsible for its domain layer
- **Independent deployment future-proofing** — could split into microservices later without rewriting domain
- **Testability** — repositories are interfaces; in-memory impls for fast unit tests + postgres impls for integration
- **Discoverability** — a developer asking "where does trade state live?" answers themselves in 5 seconds

### Negative

- **Boilerplate** — each BC repeats the same scaffolding (sentinel errors, ReconstituteX helpers, NoopPublisher)
- **Cross-BC operations need explicit eventing** — direct calls forbidden by convention; learning curve

### Mitigations

- **Boilerplate** — accepted as the cost of long-term clarity; partially addressed by code templates in `modules/<canonical>/` (trade is the reference template)
- **Cross-BC discipline** — enforced via `.claude/rules/modules-domain.md` + code review

## Alternatives considered

- **Single domain package** — fast initial, would calcify around the first developer's mental model
- **Microservices from day 1** — premature; we'd be debugging gRPC plumbing instead of building the FX engine
