# Cost Review — <YYYY-Qn> — ExchangeOS

> Copy to `.audit-bundles/cost-reviews/YYYY-Qn.md`. Owned by Engineering Lead +
> Finance partner. Cadence: quarterly during normal ops, monthly if any budget
> threshold breached the previous month.

## Metadata

| Field | Value |
|-------|-------|
| Quarter | YYYY-Qn |
| Reviewer | Engineering Lead |
| Finance partner | Name |
| Date | YYYY-MM-DD |
| Budget version | terraform module exchangeos-budget v? |

## Executive summary

> 3 sentences. Spend vs budget. Any threshold breached. Action items.

## Spend (USD)

| Bucket | Budget | Actual | Variance | Trend vs prev Q |
|--------|--------|--------|----------|-----------------|
| Total | | | | |
| Compute (GKE) | | | | |
| Storage (GCS + KMS) | | | | |
| Networking | | | | |
| Observability (Mimir/Tempo/Loki) | | | | |
| Other | | | | |

## Per-BC allocation

| Bounded Context | Spend | % of total | Notes |
|----------------|-------|-----------|-------|
| trade | | | |
| quote | | | |
| cls_settlement | | | |
| compliance | | | |
| shared (worker, migrator) | | | |
| ... | | | |

Source: `sum by (bc) (gcp_billing_cost_usd{module="exchangeos"})` (Grafana board "Cost by BC").

## Threshold events this quarter

| Date | Budget | Threshold % | Reason | Resolution |
|------|--------|-------------|--------|------------|
| | | | | |

## Right-sizing opportunities identified

| Resource | Current | Recommended | Estimated saving / mo |
|----------|---------|-------------|----------------------|
| | | | |

## Architectural decisions impacting cost (next quarter)

- ...

## Action items

| Action | Owner | Target |
|--------|-------|--------|
| | | |

## Sign-off

| Role | Name | Date |
|------|------|------|
| Engineering Lead | | |
| Finance Partner | | |
| CTO (if any threshold > 100% breached) | | |
