# ExchangeOS Disaster Recovery Runbook

> ISO 27001 A.5.29 + A.8.13 + A.8.14 compliance evidence
> Owner: Platform team
> Last drill: TBD (quarterly cadence)

## RTO + RPO targets

| Target | Value | Notes |
|--------|-------|-------|
| RTO (Recovery Time Objective) | 4 hours | Time to restore full service after regional outage |
| RPO (Recovery Point Objective) | 5 minutes | Max data loss tolerated |

The CRDB hub multi-region replication + Kafka MirrorMaker 2 + GitOps-managed
infrastructure together enable these targets. Application state (in CRDB) is
the only data; the rest is reconstructible from declarative manifests.

## Primary region

- **us-east1** (Moncks Corner, SC)
- GKE Autopilot cluster `exchangeos`
- Shared CRDB hub cluster (cross-platform; not exchangeos-managed)
- Kafka cluster `kafka-prod-us-east1`
- Vault HA cluster (separate VPC)

## Secondary / failover region

- **us-central1** (Council Bluffs, IA)
- GKE Autopilot cluster `exchangeos-dr` (standby)
- CRDB hub replica zone (read-only secondary)
- Kafka cluster `kafka-prod-us-central1` (MirrorMaker 2 from primary)
- Vault replica (DR replication mode)

## Failover scenarios

### Scenario A: Single AZ failure in primary region

**Action:** None — GKE Autopilot + CRDB multi-zone + Kafka RF=3 handle automatically.

### Scenario B: Full primary-region outage (GCP us-east1 down)

#### Pre-flight (T-0)

1. Confirm outage via GCP status page + internal observability
2. Open incident channel + IC takes lead
3. Notify customers via status page

#### Failover (T+0 to T+1h)

1. **Verify secondary region health:**
   ```bash
   gcloud container clusters get-credentials exchangeos-dr --region us-central1
   kubectl get nodes
   kubectl get applications -n argocd | grep exchangeos
   ```
2. **Promote CRDB secondary to primary:**
   - Engage CRDB cluster owner (cross-platform)
   - Update DSN in Vault: `secret/data/exchangeos/db` → secondary cluster endpoint
   - External Secrets Operator auto-refreshes K8s Secret within reconciliation interval (~1 min)
3. **Promote Kafka MirrorMaker target → source-of-truth:**
   - Disable MirrorMaker from primary
   - Update `EXCHANGEOS_KAFKA_BROKERS` in Vault to secondary cluster
   - Restart worker pods to pick up new brokers
4. **Scale up secondary GKE:**
   - ArgoCD already syncs; verify replicas match production values
   - Update DNS: api.exchangeos.revenu.tech CNAME → us-central1 LB
   - DNS TTL is 60s; expect customer recovery within 2 min after DNS flip

#### Validation (T+1h to T+2h)

- [ ] /healthz green on all binaries
- [ ] Smoke tests: book test trade in dev tenant + observe settlement
- [ ] Trace span verifies end-to-end flow
- [ ] No 5xx spikes for 30 min
- [ ] Customer comms: "Service restored, monitoring closely"

#### Post-failover (T+2h to T+4h)

- [ ] Full E2E suite (`task test:e2e`) green
- [ ] Outbox dispatch caught up (`outbox_events WHERE dispatched_at IS NULL` ~0)
- [ ] Argo Rollouts canary smoke pass
- [ ] Status page: "All systems operational"

### Scenario C: Data corruption / ransomware

1. **DO NOT failover** — corruption may have replicated
2. Stop all writes: scale api + worker to 0 replicas
3. Identify last known good backup timestamp (CRDB BACKUP audit log)
4. Restore CRDB to point-in-time (cluster owner)
5. Replay outbox events from last-known-good timestamp via worker
6. Audit + reconcile against external sources (CLS NetReports, BACEN reports)
7. Customer comms + LGPD notification if PII involved

## Backup verification

- CRDB cluster: backups every 15 min (incremental) + 24h (full); retention 90 days
- Verified weekly by automated restore-to-staging test (cross-platform script)
- Outbox archive table retains 30 days; full audit_events retains 7 years (regulatory)

## Drill schedule

- **Quarterly DR drill:** full failover simulation in staging using production-like data
- **Monthly mini-drill:** Vault secret rotation + ESO refresh validation
- **Annual:** external red-team exercise

## Drill log

| Date | Scenario | Outcome | Action items |
|------|----------|---------|--------------|
| TBD  | First drill scheduled within 30 days of production go-live | — | — |

## Decision tree

```
Outage detected
    ├── Single AZ → no action (auto)
    ├── Full region → Scenario B failover
    ├── Data corruption → Scenario C (point-in-time restore)
    └── Security breach → engage Security IR (incident-response.md) first
```

## Communications template

```
Subject: [INCIDENT] ExchangeOS regional failover initiated

We are failing over the ExchangeOS service from us-east1 to us-central1 due to
[a GCP regional outage / data integrity event]. Expected recovery time: 2 hours.

Status updates every 15 minutes at status.exchangeos.revenu.tech.

— Platform Team
```
