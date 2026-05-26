# ExchangeOS — SLI / SLO Definitions

> Source of truth for what "reliable" means. Tied to the AnalysisTemplate gates
> in `deploy/k8s/argo-rollouts/api-rollout.yaml` + Grafana alert rules.
> Owner: Platform team + Product owner

## Core principles

- **SLI** = a measurement we can compute from raw telemetry
- **SLO** = a target value for an SLI over a window (e.g. 99.9% over 30d)
- **SLA** = customer-facing commitment ≤ SLO (with a safety margin)
- **Error budget** = (100% - SLO) — burn rate triggers alerts

We size SLOs so that meeting them = customer happy. We don't over-engineer.

## SLI/SLO catalogue

### 1. API availability

| Field | Value |
|-------|-------|
| SLI | `1 - (sum(rate(http_requests_total{job="exchangeos-api",code=~"5.."}[5m])) / sum(rate(http_requests_total{job="exchangeos-api"}[5m])))` |
| SLO | **99.9% over rolling 30 days** (43.2 min downtime budget/month) |
| SLA | 99.5% (offered to customers; 21 min/month buffer over SLO) |
| Alert | Burn-rate > 14.4× over 1h (2% budget consumed) → page |

### 2. API latency — p99 quote latency

| Field | Value |
|-------|-------|
| SLI | `histogram_quantile(0.99, sum(rate(http_request_duration_ms_bucket{job="exchangeos-api",endpoint="GetQuote"}[5m])) by (le))` |
| SLO | **< 100ms p99 over rolling 7 days** |
| Rationale | Quote latency = trader's experienced delay; > 100ms feels sluggish |
| Alert | p99 > 100ms for 5 min → warn; > 200ms for 5 min → page |

### 3. API latency — p99 trade booking

| Field | Value |
|-------|-------|
| SLI | `histogram_quantile(0.99, sum(rate(http_request_duration_ms_bucket{job="exchangeos-api",endpoint="BookTrade"}[5m])) by (le))` |
| SLO | **< 200ms p99 over rolling 7 days** |
| Rationale | Trade is heavier than quote (DB write + outbox insert + risk check); 2× quote budget reasonable |
| Alert | p99 > 200ms for 5 min → warn; > 500ms for 5 min → page |

### 4. Outbox dispatch lag

| Field | Value |
|-------|-------|
| SLI | `max(time() - max(outbox_pending_oldest_occurred_at_seconds))` |
| SLO | **< 5 min p99 over rolling 24 hours** |
| Rationale | Downstream consumers (LedgerOS) tolerate some lag; > 5 min = ops attention |
| Alert | Lag > 5 min for 2 min → warn; > 15 min for 2 min → page |

### 5. CLS cycle on-time close

| Field | Value |
|-------|-------|
| SLI | `count(cls_cycle{status="CLOSED",closed_after_deadline="false"}) / count(cls_cycle{status="CLOSED"})` |
| SLO | **100% over rolling 30 days** (no missed cycles) |
| Rationale | Missed CLS close = regulatory + customer impact; zero tolerance |
| Alert | Any cycle close > scheduled_close + 5 min → page |

### 6. BACEN report submission success

| Field | Value |
|-------|-------|
| SLI | `count(bacen_reports{status="ACCEPTED"}) / count(bacen_reports{status=~"ACCEPTED|REJECTED"})` |
| SLO | **≥ 99% over rolling 30 days** |
| Alert | Any rejection → Compliance Slack channel; > 3 rejections in 24h → page Compliance Officer |

### 7. Worker (outbox dispatcher) availability

| Field | Value |
|-------|-------|
| SLI | `up{job="exchangeos-worker"}` |
| SLO | **99.9% over 30d** |
| Alert | Down for > 2 min → page |

### 8. CRDB query latency (p99)

| Field | Value |
|-------|-------|
| SLI | `histogram_quantile(0.99, sum(rate(pgxpool_query_duration_ms_bucket{job=~"exchangeos.*"}[5m])) by (le))` |
| SLO | **< 50ms p99 over rolling 7 days** |
| Rationale | DB latency is upstream of API latency; tighter budget |
| Alert | p99 > 50ms for 5 min → warn; > 100ms for 5 min → page |

## Error budget policy

Each SLO has an implicit error budget = (1 - SLO).

### Burn-rate alerts (2 windows)

We use the [Google SRE Workbook two-window burn-rate alerting pattern](https://sre.google/workbook/alerting-on-slos/):

- **Page:** burning 2% of monthly budget in 1h (14.4× normal rate)
- **Ticket:** burning 5% of monthly budget in 6h (6× normal rate)
- **Slack-only:** burning 10% in 3d (1× normal rate)

### Policy on budget exhaustion

If we exhaust > 50% of the monthly budget mid-month:

1. **Stop all non-emergency deploys** until budget recovers
2. Post-mortem reviews focused on RELIABILITY (not features)
3. Engineering capacity reallocated to remediation

If we exhaust > 100% of the monthly budget (i.e. SLO missed):

1. Customer-facing apology + RCA published
2. SLA credits per contract
3. Quarterly review of SLO targets — were they unrealistic, or did we regress?

## Alert routing

| Severity | Routing |
|----------|---------|
| Page | PagerDuty `exchangeos-api-prod` → on-call rotation |
| Ticket | Slack `#exchangeos-incidents` + create JIRA |
| Slack | Slack `#exchangeos` thread |

## Dashboards

- **Grafana / exchangeos-api-slo** — live SLO + burn-rate panels
- **Grafana / exchangeos-delivery** — milestone + release health (auto-updated hourly)

## Review cadence

- **Monthly** — Platform Lead reviews previous month's SLO performance + tightens/loosens as needed
- **Quarterly** — Product Owner reviews SLA vs customer feedback; signs off any SLO change
- **Annually** — Full re-baseline with fresh data
