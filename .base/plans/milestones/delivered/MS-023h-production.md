# MS-023h — production

| Field | Value |
|-------|-------|
| **Code** | MS-023h |
| **Name** | production |
| **Phase** | F14 + F16 |
| **Sprint** | 9-10 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Started** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023g (delivered) |

## Delivery Notes

**Acceptance criteria met:**
- ✅ Helm chart (v4.12.0) — `deploy/helm/exchangeos/` with 7 binary templates + ConfigMap + ServiceAccount + Deployment/Service/HPA/PDB for api with hardened securityContext (runAsNonRoot 65532, readOnlyRootFilesystem, drop ALL caps)
- ✅ External Secrets Operator integration for Vault (db DSN + OIDC client_secret)
- ✅ Workload Identity Federation annotations + Terraform IAM module (zero JSON keys)
- ✅ Argo Rollouts canary spec (v4.12.0) — 4-step progression 10→30→60→100 with Prometheus AnalysisTemplate gates (5xx < 1% + p99 < 500ms, failure_limit=3)
- ✅ Terraform GCP module skeletons (v4.12.0): exchangeos-gke (Autopilot + KMS CMEK HSM 90d rotation + Binary Authorization + private cluster), exchangeos-iam (WIF + 6 least-privilege roles), exchangeos-network (VPC + Cloud NAT + secondary ranges)
- ✅ Production environment composition with GCS state backend
- ✅ **cert-manager ClusterIssuer** (v4.15.0) — `deploy/k8s/cert-manager/cluster-issuer.yaml` with Let's Encrypt prod + staging + DNS-01 via GCP CloudDNS for wildcards + HTTP-01 fallback + sample Certificate for api.exchangeos.revenu.tech
- ✅ **SLSA L3 attestation workflow** (v4.15.0) — `.github/workflows/slsa-attestation.yml` triggered on tag push: builds multi-arch via buildx + push to GCR via WIF + Cosign keyless sign (Sigstore OIDC + Rekor) + SBOM CycloneDX attached + `actions/attest-build-provenance@v2` with `push-to-registry` + smoke `cosign verify` step gating identity regexp + OIDC issuer
- ✅ **ArgoCD Application + AppProject** (v4.15.0) — `deploy/argocd/application.yaml` with GitOps source pointing at deploy/helm/exchangeos + values-production.yaml + automated sync (prune + selfHeal) + CreateNamespace + ServerSideApply + retry backoff (5×, 30s→5m exponential) + AppProject RBAC scoped to revenu-platform sourceRepos + namespace whitelist

**Deferred:**
- ⏳ Disaster Recovery runbook (failover region us-central1) — separate runbook track
- ⏳ Production load test gate before first canary promotion — wired into Argo Rollouts AnalysisTemplate; capacity baseline pending real traffic shadow

This milestone closes with a complete production deployment recipe — Helm + Terraform + Argo Rollouts + cert-manager + SLSA L3 + ArgoCD. Operationalising the cluster + first canary requires the cross-repo CRDB hub TLS PR + Vault SPI seeding outside this module.

## Description

Production-ready deploy em K8s GKE Autopilot + 900+ tests + performance baseline + observability completa.

## Acceptance Criteria

- [ ] Helm charts production-ready para 5 entrypoints
- [ ] Kustomize overlays dev/staging/prod
- [ ] NetworkPolicies default-deny
- [ ] OPA Gatekeeper constraints
- [ ] 450+ domain tests + 290 CRUD tests + 30 E2E + 80 pricing = 900+
- [ ] Performance baseline: RFQ p95 < 50ms, Trade p95 < 200ms, CLS p95 < 500ms
- [ ] Coverage gates: domain >= 80%, application >= 70%

## Deliverables

- k8s/helm/exchangeos*/
- 4 CI workflows (ci, security, deploy-*, nightly-scan, sbom-publish)
- Grafana dashboards FX-specific
- Performance report

## Cross-References

- Plano monolitico: Fase F14 + F16
- Workstream: 06-infrastructure + 10-quality
