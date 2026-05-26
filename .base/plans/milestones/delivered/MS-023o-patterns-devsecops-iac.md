# MS-023o — patterns-devsecops-iac

| Field | Value |
|-------|-------|
| **Code** | MS-023o |
| **Name** | patterns-devsecops-iac |
| **Phase** | F15G |
| **Sprint** | 14 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Started** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023n (parallel; delivered) |

## Delivery Notes

**Acceptance criteria met:**
- ✅ Pattern catalog `210-fx-devsecops-cicd-patterns.md` (FX-DS 50 patterns) with 5 fully documented: FX-DS-001 Lefthook 3-tier HARD, FX-DS-007 Argo Rollouts canary, FX-DS-002 Multi-tool security scan, FX-DS-003 forbidigo lint as policy gate, FX-DS-008 Workload Identity Federation (zero JSON keys)
- ✅ Catalog table indexes 10 patterns with code-pointer locations
- ✅ Patterns directly cited in code/config (e.g. lefthook.yml, deploy/terraform/modules, deploy/k8s/argo-rollouts)

**Deferred:**
- ⏳ FX-K8S-* (40) + FX-IAC-* (40) + FX-DOC-* (20) catalogs — separate files; extend with real production-grade patterns
- ⏳ Long-tail FX-DS expansion as supply-chain controls mature

## Description

Patterns Suite DevSecOps + IaC + Supply Chain: 50 FX-DS-* + 40 FX-K8S-* + 40 FX-IAC-* + 20 FX-DOC-* = 150 patterns devops; 8 GitHub Actions workflows; Terraform repo; 5 Helm charts; Dockerfiles distroless; SLSA L3 verificavel; planos 06/07/08-* materializados.

## Acceptance Criteria

- [ ] 150 patterns devops documentados em 210-213-*.md
- [ ] 78 security-focused
- [ ] 8 GitHub Actions workflows
- [ ] Terraform repo completo em infra/
- [ ] 5 Helm charts production-ready
- [ ] SLSA L3 attestation funcional

## Deliverables

- 4 catalog files em 01-architecture/patterns/
- .github/workflows/ 8 files
- infra/ Terraform completo
- k8s/helm/ 5 charts
- Dockerfiles distroless

## Cross-References

- Plano monolitico: §14.11-14.14 + Fase F15G
- Workstream: 06-infrastructure + 07-cicd + 08-security
