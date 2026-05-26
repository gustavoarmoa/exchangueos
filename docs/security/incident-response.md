# ExchangeOS Incident Response Playbook

> ISO 27001 A.5.24..A.5.30 compliance evidence
> On-call: see `docs/operations/oncall-rotation.md`
> Communication channels: PagerDuty `exchangeos-api-prod` + Slack `#exchangeos-incidents`

## Severity classification

| Sev | Definition | Response time | Owner |
|-----|------------|---------------|-------|
| Sev1 | Customer-facing outage OR data integrity loss OR security breach | < 5 min | On-call + Sec lead + Eng lead |
| Sev2 | Degraded service (slow, partial functionality) OR known security vuln in production | < 30 min | On-call |
| Sev3 | Minor degradation OR observability gap | < 4 h | On-call |
| Sev4 | Cosmetic / future risk | Next business day | Backlog grooming |

## Sev1 response — first 15 minutes

1. **Acknowledge page** — claim ownership via PagerDuty
2. **Open Incident Channel** — `#inc-<YYYYMMDD>-<short-id>` in Slack; pin the runbook link
3. **Engage IC (Incident Commander)** — first on-call defaults to IC unless rotated
4. **Initial assessment** — check Grafana dashboards:
   - `dashboard.observability/d/exchangeos-delivery` — health snapshot
   - `dashboard.observability/d/exchangeos-api-slo` — error budget burn
5. **Mitigation OR rollback decision** — within 5 min of assessment
   - **Mitigate:** scale up HPA, restart pods, adjust circuit breaker
   - **Rollback:** `kubectl argo rollouts abort + undo exchangeos-api` (see `docs/operations/canary-runbook.md`)
6. **Stabilise** — confirm 5xx rate < 1% + p99 < 500ms via Grafana
7. **Communicate** — status update every 15 min in incident channel; customer comms via marketing lead

## Sev1 response — first hour

- [ ] All affected customers identified + listed in incident channel
- [ ] Root-cause hypothesis documented (not yet root cause)
- [ ] Evidence preserved: container logs collected (`kubectl logs --previous`), OTel traces exported, DB snapshots taken if data corruption suspected
- [ ] Regulatory notifications considered:
  - BACEN: 24h notification for systemic events (FX-specific)
  - LGPD: 72h notification if PII breach
  - SISCOAF: COS if suspected fraud (per RN_FX_039)

## Post-incident (T+24h to T+72h)

1. **Blameless retro** — within 72h, scheduled by IC
2. **Root cause analysis** — Five Whys + Ishikawa diagram for complex incidents
3. **Action items** — assigned to owners with target dates
4. **Public post-mortem** — for Sev1 customer-impacting, within 7 days
5. **Update threat model** — `docs/security/threat-model-stride.md` if new threat surfaced
6. **Update runbook** — `docs/operations/canary-runbook.md` if mitigation pattern changed

## Common scenarios

### S-1 — API 5xx spike

1. Check `up{job="exchangeos-api"}` — pod loss?
2. Check DB latency `pgx.pool.acquired{job=...}` saturation
3. Check Kafka producer error rate (worker logs)
4. **Mitigation order:** restart api pods → rollback canary → switch traffic to last-known-good replica set

### S-2 — Outbox dispatch lag > 5 min

1. Check `outbox_events WHERE dispatched_at IS NULL` count via SQL console
2. Check worker logs for publish errors
3. If Kafka cluster healthy → scale worker replicas; if Kafka down → engage Kafka on-call (cross-platform)
4. Outbox tolerates 24h lag without data loss; users notified via WebSocket gracefully

### S-3 — Risk limit false breach (production)

1. Verify via SQL: `SELECT * FROM risk_limits WHERE limit_id = '<X>'`
2. If breach is correct → no action; trade was correctly rejected
3. If breach is incorrect → security-officer overrides via 4-eyes (see SoD matrix)
4. Post-incident: classify whether RN_FX_015 monitor threshold needs adjustment

### S-4 — BACEN report submission rejected

1. Check `bacen_reports.rejection_reason` for the failed report_id
2. Engage compliance-officer for content review
3. Resubmit via ComplianceService.SubmitBACENReport with corrected payload
4. If repeated rejections → escalate to BACEN technical liaison

### S-5 — Suspected PII exfiltration via Kafka topic

1. **DO NOT** delete any data — preserve for forensics
2. Isolate consumer groups via SASL/SCRAM credential rotation
3. Engage Security lead + LGPD officer
4. Begin LGPD 72h notification clock

## Tabletop drills

- Quarterly drill scheduled by Security team
- Scenarios rotated from S-1..S-5 + bespoke chaos scenarios
- Drill notes filed in `docs/security/drills/YYYY-Qn.md`

## Contacts

| Role | Primary | Secondary | Escalation |
|------|---------|-----------|------------|
| Platform on-call | PagerDuty rotation | PagerDuty rotation | Platform Lead |
| Security on-call | PagerDuty rotation | Security Officer | CISO |
| Compliance on-call | Compliance Officer | Backup CO | Compliance Lead |
| External vendors | CRDB Labs support | Kafka cluster owner | — |
