# Cost Allocation Policy — ExchangeOS

> Every GCP resource that runs ExchangeOS workloads MUST carry these labels.
> Budgets (`deploy/terraform/modules/exchangeos-budget/`) filter by them — an
> unlabeled resource is invisible to FinOps + can blow through limits silently.

## Mandatory labels

| Label | Allowed values | Source of truth |
|-------|---------------|-----------------|
| `module` | `exchangeos` | Hard-coded — never override |
| `env` | `dev`, `staging`, `production` | Terraform variable per environment |
| `bc` | `trade`, `quote`, `amendment`, `cls_settlement`, `payin`, `netreport`, `cfets_capture`, `cfets_confirmation`, `settlement`, `refdata`, `admin`, `risk`, `position`, `compliance`, `shared` | Per-BC binding; `shared` for cross-cutting (worker, migrator) |
| `tier` | `infra`, `data`, `runtime`, `observability` | Helps allocate to budget categories |
| `cost_center` | `engineering`, `compliance`, `security` | Internal allocation |

## Enforcement

### Terraform (infra-level)

Every module sets these via `default_labels` block on the provider:

```hcl
provider "google" {
  project = var.project_id
  default_labels = {
    module      = "exchangeos"
    env         = var.env
    tier        = "infra"
    cost_center = "engineering"
  }
}
```

Per-resource overrides add `bc` as appropriate. CI enforces presence via
`scripts/finops-lint-tf.sh` (planned).

### Helm (workload-level)

`deploy/helm/exchangeos/values.yaml` carries:

```yaml
commonLabels:
  module: exchangeos
  env: production
  tier: runtime
  cost_center: engineering
```

Per-deployment overrides add `bc: trade` etc. via `extraLabels` on the
Deployment template.

### Kafka topics

Topic ACLs already use `module=exchangeos` prefix; topic config labels added
via `deploy/kafka/topics.yaml` (`labels:` block per topic).

## Allocation rollups

The Grafana **Cost by BC** dashboard sums GCP billing export rows grouped by
`labels.bc` to attribute spend to a specific bounded context owner.

```promql
sum by (bc) (gcp_billing_cost_usd{module="exchangeos"})
```

Used in the quarterly cost review (`cost-review-template.md`).

## What if a resource can't carry a label?

Rare — but for managed services without label support (e.g. some Cloud Endpoints
configurations), bill-by-project is the fallback. Document the exception in
`.base/plans/00-governance/exceptions.md` with rationale + review date.

## Why this matters

Without strict allocation:
- A runaway pod in `trade` BC inflates the total budget without attribution
- We can't say "compliance costs $X/month" when justifying headcount or controls spend
- The 80% sub-budget alerts (compute / storage) become unreliable
- ISO 27001 control 5.10 (asset use rules) loses the cost dimension

Treat label drift as a P3 bug. Catch via `scripts/finops-lint-tf.sh` in pre-merge.
