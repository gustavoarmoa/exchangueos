# Eligibility Report — LGPD-YYYY-NNNN

> Stage 2 output. Run `scripts/lgpd-eligibility.sh <subject-ref>` and paste
> the table below. Each row is a separate decision the DPO + Compliance
> Officer must explicitly approve before the plan can be drafted.

| Table | Row count | Earliest occurred_at | Retention rule | Decision | Notes |
|-------|-----------|----------------------|----------------|----------|-------|
| actors | N | YYYY-MM-DD | ISO27001 5.34 | ELIGIBLE | redact display_name + external_sub |
| fx_trades | N | YYYY-MM-DD | BACEN 5 yr hold | FROZEN_REGULATORY | retention until YYYY-MM-DD |
| ... | | | | | |

## Notes

- DEFERRED rows must be marked with the unfreeze date (e.g.
  `DEFERRED_UNTIL_2031-05-26`). At that point the bundle is reopened.
- FROZEN_REGULATORY rows are excluded from `plan.json`. Document the
  regulatory basis (Lei 14.286, BACEN Circ 3.690, etc.).
