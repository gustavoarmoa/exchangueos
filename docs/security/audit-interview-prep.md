# ISO 27001 Audit Interview Prep

> What the auditor is likely to ask + the canned answer + the evidence file pointer.
> Refresh before each audit + use during dry-run interviews.

## Who gets interviewed

| Role | Likely questions | Prep doc section |
|------|------------------|------------------|
| CISO / Security Officer | Risk acceptance, top residual risks, governance | A + B + C |
| Platform Lead | Operational controls, change management, incident response | C + D |
| DBA / Cluster Owner | Backup, encryption, access logs | D + E |
| Developer (any IC) | Secure coding, code review, peer testing | E + F |
| Compliance Officer | Regulatory reporting, screening, SoD | F + G |

---

## A. Governance + organisational

**Q1: How does ExchangeOS demonstrate top-management commitment to information security?**
**A:** Security objectives are tracked in `.base/plans/milestones/` (8 of 26 explicitly include security acceptance criteria — MS-023q, MS-023t, MS-023x in particular). Per-release sign-off in `docs/operations/go-live-checklist.md` requires Security Officer + Compliance Officer + Platform Lead + Product Owner. CISO holds the residual-risk register (`docs/security/iso27001-gap-tracker.md`).
**Evidence:** `.base/plans/milestones/delivered/MS-023x-precommit-hard-enforcement.md`, `docs/security/iso27001-gap-tracker.md`.

**Q2: How are risks identified and managed?**
**A:** Two artefacts: the STRIDE+DREAD threat model (`docs/security/threat-model-stride.md`) lists 15 threats with mitigations; the gap tracker (`docs/security/iso27001-gap-tracker.md`) closes Annex A gaps. Threat model is re-run on every major release + post-incident.
**Evidence:** `docs/security/threat-model-stride.md`, `docs/security/iso27001-gap-tracker.md`.

---

## B. Access control + identity

**Q3: How are M2M credentials rotated?**
**A:** `cmd/cred-rotator` is a CronJob (monthly per Helm values) that pulls 14 M2M client_secrets, generates new ones via KeycloakOS API, and pushes to Vault SPI. Every rotation emits an OTel span + admin SystemEvent.
**Evidence:** `cmd/cred-rotator/main.go`, `deploy/helm/exchangeos/values.yaml` (`cred-rotator.cronJob.schedule: "0 3 1 * *"`), `modules/admin/domain/event.go`.

**Q4: How is privileged access controlled?**
**A:** Workload Identity Federation (zero JSON keys); 6 least-privilege GCP roles bound to the GKE ServiceAccount via Terraform. Vault tokens TTL ≤ 1h. SoD matrix (`docs/security/sod-matrix.md`) lists 23 critical actions with Δ-4-eyes annotation; enforced via GitHub branch protection + EMERGENCY_BYPASS audit log.
**Evidence:** `deploy/terraform/modules/exchangeos-iam/main.tf`, `docs/security/sod-matrix.md`, `scripts/git-hooks-wrapper.sh`.

---

## C. Cryptography

**Q5: What encryption is in transit / at rest?**
**A:** **In transit:** TLS 1.3 minimum (cited in CLAUDE.md); cert-manager + Let's Encrypt for public; mTLS inter-service via shared CA. **At rest:** CRDB cluster encryption + GKE etcd via KMS HSM CMEK (`exchangeos-gke` Terraform module, 90-day rotation).
**Evidence:** `deploy/k8s/cert-manager/cluster-issuer.yaml`, `deploy/terraform/modules/exchangeos-gke/main.tf` (KMS keyring), `CLAUDE.md` (TLS rule).

**Q6: Where do you derive cryptographic material from?**
**A:** Vault HA cluster (sealed by Shamir secret sharing — operated by Security Ops team). Application secrets fetched via External Secrets Operator at pod startup; never written to disk.
**Evidence:** `scripts/vault-seed.sh`, `deploy/helm/exchangeos/values.yaml` (externalSecrets section).

---

## D. Operations + change management

**Q7: How do production deploys happen?**
**A:** GitOps via ArgoCD pulls `deploy/helm/exchangeos`. Argo Rollouts canary 10→30→60→100 with Prometheus AnalysisTemplate gates (5xx rate < 1%, p99 < 500ms, failureLimit=3 auto-rollback). All container images Cosign-signed (keyless OIDC) + SLSA L3 provenance attested.
**Evidence:** `deploy/argocd/application.yaml`, `deploy/k8s/argo-rollouts/api-rollout.yaml`, `.github/workflows/slsa-attestation.yml`.

**Q8: Walk me through an incident response.**
**A:** PagerDuty page → IC claims via Slack `#exchangeos-incidents` → assess Grafana dashboards within 5 min → mitigate OR rollback (`task canary:abort` + `task canary:rollback`) → preserve evidence (logs + traces + DB snapshots) → regulatory notification clocks start (BACEN 24h, LGPD 72h, SISCOAF per RN_FX_039). Full runbook + 5 common scenarios in `incident-response.md`.
**Evidence:** `docs/security/incident-response.md`, `docs/operations/runbook-index.md`.

---

## E. Development security

**Q9: How does insecure code get blocked from reaching main?**
**A:** Lefthook 3-tier HARD enforcement:
- pre-commit (<30s): fmt + vet + gitleaks + buf-lint + yaml-lint
- pre-push (<3min): unit tests + golangci-lint full (incl. forbidigo banning `float64` outside tests) + govulncheck
- pre-merge (CI <15min): full test + integration + trivy

`scripts/git-hooks-wrapper.sh` blocks `git --no-verify`; emergency bypass requires `EMERGENCY_BYPASS=true GIT_BYPASS_REASON="..."` and is logged + Slack-alerted.
**Evidence:** `lefthook.yml`, `.golangci.yml`, `scripts/git-hooks-wrapper.sh`, `.git/audit-bypass.log` (review on demand).

**Q10: How are dependencies kept secure?**
**A:** `govulncheck` runs on every push + weekly cron. Trivy scans filesystem + container images. CodeQL on every PR. SBOM (CycloneDX) attached to every release artifact via the SLSA L3 workflow.
**Evidence:** `.github/workflows/security.yml`, `.github/workflows/slsa-attestation.yml`.

---

## F. Data + privacy

**Q11: How do you prevent cross-tenant data leakage?**
**A:** Every database table (except global refdata) has a mandatory `tenant_id` column; application services filter at query time. Repository interfaces accept `tenantID uuid.UUID` as first context parameter for every Find/Get/List. The threat model T-4 (tenant scoping) is scored 7.0 critical and explicitly mitigated.
**Evidence:** `migrations/000001..000009`, `modules/*/application/service.go`, `docs/security/threat-model-stride.md` (T-4).

**Q12: What PII does ExchangeOS handle?**
**A:** Counterparty BICs + LEIs (institutional identifiers, not personal). No customer PII per se — trades are between institutional counterparties. ScreeningResult.hits may contain individual names matched against sanctions lists; these are treated as Sensitive Personal Data under LGPD and stored only as opaque strings with 7-day Kafka topic retention.
**Evidence:** `deploy/kafka/topics.yaml` (compliance topic retention), `modules/compliance/domain/screening.go`.

---

## G. Regulatory + compliance

**Q13: How is BACEN compliance enforced?**
**A:** `pkg/bacen` provides the 95-code classifier (20 most-common pre-seeded) and the 6 IOF rates per Decreto 12.499/2025. ComplianceService classifies + computes on every booked trade. BACENReport aggregate tracks SISBACEN/BCB-CCS/BCB-CAMBIO submission status. SISCOAF COS emission flagged automatically when ScreeningResult.RiskLevel = HIGH (RN_FX_039).
**Evidence:** `pkg/bacen/classifier.go`, `pkg/bacen/iof.go`, `modules/compliance/domain/`, `modules/compliance/application/service.go`.

**Q14: How are 14 ISO 20022 message types handled correctly?**
**A:** `pkg/iso20022` registry catalogues 32 schemas; routes via OrganisationRouter (CLSBUS33 / CFETS / ISO fallback). Per-message Go structs in `pkg/iso20022/fxtr` (15 CLS+CFETS variants) + admi + camt + reda. Round-trip tested via marshaller + registry (see `pkg/iso20022/fxtr/fxtr_014_roundtrip_test.go`).
**Evidence:** `pkg/iso20022/`, fxtr round-trip test.

---

## Bring-along checklist

When the auditor visits, bring:

- [ ] This document (interview prep)
- [ ] `docs/security/iso27001-controls-mapping.md` printed + annotated
- [ ] `docs/security/threat-model-stride.md`
- [ ] `docs/security/sod-matrix.md`
- [ ] `docs/security/iso27001-gap-tracker.md` (shows what's deferred + why)
- [ ] Live access to GitHub (PR history) + Grafana (dashboards) + Vault audit log
- [ ] Last 30 days `.git/audit-bypass.log` (proves rare + justified)
- [ ] 90-day backup verification log (DBA team)
- [ ] Recent incident post-mortems (from `docs/security/drills/`)
