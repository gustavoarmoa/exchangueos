# Architecture Decision Records (ADRs)

Format: [MADR](https://adr.github.io/madr/) — Markdown Architecture Decision Records.

Each ADR captures a single decision: context, options considered, choice, and consequences. ADRs are **append-only** — when a decision changes, create a new ADR superseding the old one.

## Index

| # | Title | Status | Supersedes | Date |
|---|-------|--------|------------|------|
| [0001](0001-shared-crdb-hub-tls.md) | Shared CRDB hub with TLS over per-module clusters | Accepted | — | 2026-05-24 |
| [0002](0002-ddd-bounded-contexts.md) | 14 bounded contexts with DDD aggregates | Accepted | — | 2026-05-24 |
| [0003](0003-transactional-outbox.md) | Transactional outbox for async event publication | Accepted | — | 2026-05-24 |
| [0004](0004-build-tag-gated-bindings.md) | Build-tag-gated optional bindings (grpcgen, kafka) | Accepted | — | 2026-05-24 |
| [0005](0005-decimal-only-money.md) | shopspring/decimal mandatory; NEVER float for money | Accepted | — | 2026-05-24 |
| [0006](0006-wif-zero-keys.md) | Workload Identity Federation over JSON service-account keys | Accepted | — | 2026-05-24 |
| [0007](0007-gitops-argocd.md) | GitOps via ArgoCD over imperative kubectl | Accepted | — | 2026-05-24 |
| [0008](0008-argo-rollouts-canary.md) | Argo Rollouts canary with Prometheus AnalysisTemplate | Accepted | — | 2026-05-24 |

## Authoring a new ADR

1. Pick the next available 4-digit number.
2. Copy `template.md` (TBD) and fill in.
3. Open PR with the new file + add a row to this index.
4. Review by 2 platform engineers; once accepted, the file becomes append-only.
5. Status transitions: `Proposed → Accepted → Deprecated → Superseded by 00NN`.
