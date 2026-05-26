# ADR 0001 — Shared CRDB hub with TLS over per-module clusters

- Status: Accepted
- Date: 2026-05-24
- Decider(s): Platform team, DBA team

## Context

Every Revenu Platform module needs a transactional database. Two patterns are common:

- **A. Per-module CRDB cluster** — each module owns its own cluster; complete isolation.
- **B. Shared CRDB hub** — single multi-tenant cluster registers each module under `cockroachdb/modules/<name>/`; each module gets its own database + TLS user/cert.

13+ modules ship on the platform. Each cluster costs ~3 nodes minimum, plus operator overhead.

## Decision

**Option B — Shared CRDB hub with mTLS per module.**

Each module registers under `cockroachdb/modules/<name>/` (cross-repo) with:

- A dedicated database (`CREATE DATABASE exchangeos`)
- A TLS user (`exchangeos_app` + `exchangeos_auditor`) authenticated by client cert (CN = username)
- Client certs issued by the shared hub CA + rotated quarterly

**NEVER `--insecure` in any production DSN.** Local docker-compose may use `--insecure` for the single-node dev cluster only.

## Consequences

### Positive

- **Operational cost** — 3 nodes for 13+ modules vs 39+ nodes
- **Centralised observability** — single Prometheus + Grafana dashboard surface
- **Multi-tenancy at module granularity** — already enforced by CRDB's database-level isolation
- **Cross-module read patterns possible** — read-only auditor role can JOIN across modules for compliance reports

### Negative

- **Blast radius** — a CRDB cluster issue takes down all 13+ modules at once
- **Schema-change coordination** — large schema changes across multiple modules need careful sequencing
- **Cross-repo PR friction** — every new module must land a PR in the `cockroachdb` repo before its production deploy (see `docs/operations/crdb-hub-tls-pr.md`)

### Mitigations

- **Blast radius** addressed by multi-region replication + 5-min RPO target (see `docs/security/dr-runbook.md`)
- **Schema changes** queued through DBA team + per-module migration runners (ExchangeOS = `cmd/migrator`)
- **PR friction** addressed by the spec doc making the cross-repo PR routine

## Alternatives considered

- **Cloud-managed DB (AlloyDB / Spanner)** — vendor-lock-in; cost > shared CRDB; SQL dialect drift
- **Per-module managed CRDB Cloud** — same cost problem as Option A
- **PostgreSQL + Citus** — losing CRDB's geo-distribution + ZeroETL changefeed
