# First Canary Runbook — ExchangeOS API

Run this once for the first production deploy. Subsequent rollouts are
automated via ArgoCD + Argo Rollouts AnalysisTemplate gates.

## Pre-flight (T-24h)

- [ ] All MS-023a..x acceptance criteria green (see `delivered/` milestones)
- [ ] `cockroachdb/modules/exchangeos/` PR merged + production cluster bootstrapped
- [ ] Vault seeded (`scripts/vault-seed.sh` executed in prod Vault)
- [ ] Helm chart installed in `exchangeos` namespace with `replicas: 0` (dry run)
- [ ] First image tag (`v0.1.0`) built + signed via `.github/workflows/slsa-attestation.yml`
- [ ] Cosign verify chain passes for image digest
- [ ] cert-manager Certificate `exchangeos-api-tls` ready (`kubectl get cert -n exchangeos`)
- [ ] Migrations applied: `kubectl exec -it exchangeos-migrator -- exchangeos-migrator status`
- [ ] Seeds loaded: currencies/calendars/counterparties/SSIs all present

## T-1h — Smoke environment

- [ ] Scale `exchangeos-api` to `replicas: 1` via Argo Rollouts
- [ ] Smoke endpoints:
  ```bash
  curl https://api.exchangeos.revenu.tech/healthz
  curl https://api.exchangeos.revenu.tech/readyz
  curl https://api.exchangeos.revenu.tech/version
  curl https://api.exchangeos.revenu.tech/v1/refdata/currencies?active_only=true
  ```
- [ ] Confirm OTel traces flow to Tempo: `grep service.name=exchangeos-api`
- [ ] Confirm Prometheus scrape: `up{job="exchangeos-api"} == 1`
- [ ] grpcurl: `grpcurl -import-path proto/exchangeos/v1 -d '{}' grpc.exchangeos.revenu.tech:443 grpc.health.v1.Health/Check`

## T-0 — Initiate canary

1. **Promote via Argo Rollouts** (manual gate for first cut):
   ```bash
   kubectl argo rollouts promote exchangeos-api -n exchangeos
   ```
2. Watch dashboard: `dashboard.observability/d/exchangeos-delivery`
3. Argo Rollouts UI: `kubectl argo rollouts dashboard`

## During canary windows

| Step | Weight | Pause | What to watch |
|------|--------|-------|---------------|
| 1 | 10% | 5min | 5xx rate < 1%, p99 < 500ms (AnalysisTemplate gates auto-rollback) |
| 2 | 30% | 10min | Same + CPU/memory trends + DB pool exhaustion |
| 3 | 60% | 10min | Same + per-route traffic ratio |
| 4 | 100% | —   | Full traffic; observe for 24h |

## Manual rollback trigger

If observability dashboards show a regression NOT caught by AnalysisTemplate:

```bash
kubectl argo rollouts abort exchangeos-api -n exchangeos
kubectl argo rollouts undo exchangeos-api -n exchangeos
```

The previous ReplicaSet stays warm for instant traffic swap (RollingUpdate
strategy default).

## Post-canary (T+24h)

- [ ] Run Grafana SLO report — error budget consumption
- [ ] Run cost report: `task dash-update` + check Grafana cost panel
- [ ] Update `roadmap/delivery-dashboard.md` "Last update" + push commit
- [ ] Schedule retro for T+72h

## Incident escalation

- Page: PagerDuty service `exchangeos-api-prod`
- Slack: `#exchangeos-incidents`
- DRI (rotating): see `docs/operations/oncall-rotation.md`

## Related

- `deploy/k8s/argo-rollouts/api-rollout.yaml` — canary spec
- `docs/security/incident-response.md` — IR playbook
- `docs/security/dr-runbook.md` — regional failover
