# ExchangeOS Production Go-Live Checklist

Owner: Platform team
Target date: TBD (post-MS-023a..x delivery confirmation in production environment)

## ✅ Code + tests

- [x] All 26 milestones delivered (plan 100% — see `.base/plans/milestones/delivered/`)
- [x] ~271 unit tests + 4 E2E tests green
- [x] `task hooks:pre-merge` passes (golangci-lint + integration + trivy)
- [x] Security workflow green: gitleaks + govulncheck + trivy + CodeQL + SBOM

## ⬜ Infrastructure prerequisites

- [ ] Cross-repo PR merged: `cockroachdb/modules/exchangeos/` (see `docs/operations/crdb-hub-tls-pr.md`)
- [ ] CRDB production cluster bootstrapped + cert issued
- [ ] Kafka cluster reachable from GKE Autopilot (cross-VPC peering or PSC)
- [ ] GCP project provisioned: `revenu-platform-prod`
- [ ] Terraform state bucket `gs://revenu-platform-tfstate` exists + versioning ON
- [ ] WIF pool created: `projects/<NUM>/locations/global/workloadIdentityPools/github`
- [ ] DNS zone `revenu-tech` exists with apex `exchangeos.revenu.tech` records

## ⬜ Terraform provisioning

- [ ] `cd deploy/terraform/environments/production && terraform init`
- [ ] `terraform plan -out=plan.out` reviewed by 2 platform engineers
- [ ] `terraform apply plan.out` — creates VPC + IAM + GKE
- [ ] GKE Autopilot cluster `exchangeos` reachable: `gcloud container clusters get-credentials exchangeos --region us-east1`
- [ ] Binary Authorization policy applied (PROJECT_SINGLETON_POLICY_ENFORCE)
- [ ] KMS HSM keyring + etcd encryption key active

## ⬜ Vault + Secrets

- [ ] Vault HA cluster provisioned (separate from this module)
- [ ] `scripts/vault-seed.sh` executed against production Vault
- [ ] External Secrets Operator installed in `exchangeos` namespace
- [ ] ClusterSecretStore `vault-backend` reachable
- [ ] ExternalSecret CRs materialise Secrets:
  - [ ] `exchangeos-db`
  - [ ] `exchangeos-oidc`
  - [ ] `exchangeos-kafka` (if Kafka in use)

## ⬜ Observability

- [ ] OTel Collector deployed in `observability` namespace
- [ ] Tempo + Mimir + Loki + Grafana reachable from `observability/otel-collector:4317`
- [ ] Prometheus scrapes `exchangeos-api` (verify `up == 1`)
- [ ] Grafana dashboard `exchangeos-delivery.json` provisioned
- [ ] PagerDuty service `exchangeos-api-prod` wired + on-call schedule active

## ⬜ Security + Compliance

- [ ] cert-manager `ClusterIssuer/letsencrypt-prod` deployed
- [ ] `Certificate/exchangeos-api-tls` reaches Ready state
- [ ] All container images SLSA L3 attested + cosign-signed
- [ ] `cosign verify` chain validated in deploy pipeline
- [ ] ISO 27001 Annex A controls mapping reviewed (`docs/security/iso27001-controls-mapping.md`)
- [ ] Threat model reviewed (`docs/security/threat-model-stride.md`)
- [ ] Segregation of Duties matrix approved (`docs/security/sod-matrix.md`)
- [ ] Incident Response playbook tabletop-tested (`docs/security/incident-response.md`)

## ⬜ GitOps

- [ ] ArgoCD installed in `argocd` namespace
- [ ] AppProject `revenu-platform` applied
- [ ] Application `exchangeos` applied — initial sync paused
- [ ] Sync window approved for go-live time

## ⬜ Canary deploy

- [ ] Follow `docs/operations/canary-runbook.md` step-by-step
- [ ] Sign off after T+24h observation period

## ⬜ Day 2 operations

- [ ] On-call rotation populated
- [ ] Runbooks linked from PagerDuty service
- [ ] Cost monitoring dashboard (Grafana) reviewed
- [ ] Backup + DR drill scheduled within 30 days (`docs/security/dr-runbook.md`)
- [ ] First post-mortem template ready for unexpected incidents

## Sign-off

| Role | Name | Signature | Date |
|------|------|-----------|------|
| Platform Lead | | | |
| Security Officer | | | |
| Compliance Officer | | | |
| Product Owner | | | |
