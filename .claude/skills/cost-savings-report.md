---
name: cost-savings-report
description: Generate weekly GitHub Actions cost savings report (lefthook local pre-commit telemetry)
allowed-tools: [Bash, Read]
---

# Skill: /cost-savings-report

## Trigger
`/cost-savings-report [--week <YYYY-WW>]` (default: current week)

## Workflow
1. Parse lefthook telemetry JSON: `lefthook --json telemetry --since <week_start> --until <week_end>`
2. Count local hooks runs (Tier 1 + Tier 2 + Tier 3)
3. Estimate GitHub Actions minutos saved (each local run ~3 min CI saved)
4. Calculate cost saved (ubuntu-latest: $0.008/min)
5. Generate HTML report em `.base/plans/00-governance/audits/cost-savings-<week>.html`
6. Update Grafana dashboard `exchangeos-cost-savings`
7. Post Slack #exchangeos-platform com summary

## Example Output
- Local hooks executed: 1,250
- Estimated CI minutes saved: 3,750 min
- Estimated cost saved: $30.00 USD
- Bypass count: 2 (within tolerance)
