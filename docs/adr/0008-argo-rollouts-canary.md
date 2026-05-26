# ADR 0008 — Argo Rollouts canary with Prometheus AnalysisTemplate

- Status: Accepted
- Date: 2026-05-24

## Context

Deploys to a financial-grade API need progressive traffic shift + automated rollback on regression. Plain `Deployment` with `RollingUpdate` does (1) but not (2) — by the time pods are healthy and traffic shifted, regressions are already customer-visible. Manual rollback wastes precious minutes.

## Decision

**Argo Rollouts `canary` strategy + Prometheus `AnalysisTemplate` gates + auto-rollback on `failureLimit: 3`.**

Spec in `deploy/k8s/argo-rollouts/api-rollout.yaml`:

```yaml
strategy:
  canary:
    analysis:
      templates: [{templateName: exchangeos-api-health}]
      startingStep: 2
    steps:
      - setWeight: 10
      - pause: { duration: 5m }
      - setWeight: 30
      - pause: { duration: 10m }
      - setWeight: 60
      - pause: { duration: 10m }
      - setWeight: 100
```

AnalysisTemplate checks at each pause:

- `http_5xx_rate < 0.01` (Prometheus query against `http_requests_total`)
- `http_p99_latency < 500ms` (histogram_quantile)

`failureLimit: 3` triggers automatic abort + revert to stable ReplicaSet.

## Consequences

### Positive

- **Automated rollback** within seconds of detected regression
- **Progressive traffic** limits blast radius — first 10% of users absorb risk
- **Objective gates** — Prometheus metrics, not human judgment, decide promotion
- **Manual override available** — `task canary:promote` / `task canary:abort` for IC discretion

### Negative

- **Total deploy duration ~30 min** vs ~2 min for RollingUpdate — acceptable for our cadence (1-2 deploys/day max)
- **Requires stable Prometheus metrics in place from day 1** — chicken-and-egg if metric ingestion broken
- **Tuning the gates is iterative** — 1% / 500ms thresholds are starting points

### Mitigations

- For zero-traffic regions (first prod deploy), the AnalysisTemplate starts at step 2 (`startingStep: 2`) — gating only kicks in after the 30% step where there's enough signal
- Operator runbook (`docs/operations/canary-runbook.md`) documents manual gates

## Alternatives considered

- **Blue/Green** — full duplicate at promotion is expensive at our pod count + offers no progressive risk reduction
- **Flagger** — equivalent feature set; chose Argo Rollouts for tighter ArgoCD integration
- **Plain RollingUpdate + manual rollback** — fails (2); already discounted in context
