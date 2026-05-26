# Terraform Module — `exchangeos-budget`

Creates GCP Billing budgets + Pub/Sub alert routing for ExchangeOS.

## What it creates

| Resource | Purpose |
|----------|---------|
| `google_pubsub_topic.budget_alerts` | Alert sink (created if `pubsub_topic` not provided) |
| `google_billing_budget.exchangeos_total` | Total spend budget — 5 thresholds (50/80/100/120% current + 100% forecasted) |
| `google_billing_budget.exchangeos_compute` | GKE/Compute Engine sub-budget (60% of total) — early warning at 80% |
| `google_billing_budget.exchangeos_storage` | Cloud Storage + KMS sub-budget (15% of total) — early warning at 80% |

## Required labels (enforced upstream)

Budgets filter by `labels.module = "exchangeos"` + `labels.env = <env>`. Every
GCP resource ExchangeOS creates **must** carry these labels via Terraform module
defaults or Helm `commonLabels`. Without them, spend escapes the budget filter
and we get blind-sided.

See `docs/operations/cost-allocation.md` for the labeling policy.

## Usage

```hcl
module "exchangeos_budget_production" {
  source = "../../modules/exchangeos-budget"

  billing_account     = "XXXXXX-XXXXXX-XXXXXX"
  project_id          = "revenu-platform-prod"
  env                 = "production"
  monthly_budget_usd  = 50000  # adjust per capacity planning
}
```

## Alert routing

Alerts publish to Pub/Sub → Cloud Function → Slack `#exchangeos-finops`.
The Cloud Function is owned by the platform Finance/FinOps function (out of
scope for this repo). Pub/Sub topic name is exposed via `pubsub_topic` output.

## Operational policy

- **50%** — informational; tracked in monthly review
- **80%** — warning; engineering lead notified
- **100%** — page on-call FinOps + engineering lead
- **120%** — escalate to CTO; freeze non-critical infra growth
- **Forecast 100%** — early-warning: review scaling decisions + recent merges

## Quarterly review

See `docs/operations/cost-review-template.md` — one filled instance per quarter
under `.audit-bundles/cost-reviews/`.
