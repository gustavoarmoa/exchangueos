# 07 — CI/CD

> **Workstream:** CI/CD
> **Versao:** 1.0.0
> **Status:** DRAFT

## Conteudo

| Documento | Status | Descricao |
|-----------|--------|-----------|
| `cicd.md` | TODO | CI/CD strategy + GitHub Actions pipelines |
| `git-flow.md` | TODO | GitHub Flow modificado + branch protection + Conventional Commits |
| `pre-commit-enforcement.md` | TODO | Pre-Commit HARD Enforcement Pipeline (§22 monolitico) — 3 tiers SLO + git wrapper |
| `release-management.md` | TODO | Release management + SemVer + tags + changelog |
| `supply-chain-security.md` | TODO | SLSA L3 + Cosign keyless + SBOM CycloneDX + Binary Authorization |
| `slsa-l3-roadmap.md` | TODO | SLSA L3 adoption roadmap |
| `github-actions-workflows.md` | TODO | Catalogo dos 8 workflows |

## 8 GitHub Actions Workflows

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `ci.yml` | PR + push | Build + lint + test + cross-platform matrix |
| `security.yml` | PR + push | Security gates (SAST + SCA + secrets + container + IaC) |
| `deploy-dev.yml` | push main | Auto-deploy dev environment |
| `deploy-staging.yml` | push main (apos dev) | Auto-deploy staging com smoke tests |
| `deploy-prod.yml` | manual + 2 approvers | Production deploy com canary Argo Rollouts |
| `nightly-scan.yml` | cron daily | Daily security scan independente do PR |
| `sbom-publish.yml` | release tag | SBOM CycloneDX publication |
| `preview-env.yml` | PR open/close | Ephemeral GKE namespace per PR |
| `integration-audit.yml` | cron quarterly | Integration audit 4 vetores × 13 modulos |
| `quality-gates.yml` | PR + push | Espelha local gates Tier 1+2+3 (zero discrepancia) |

## Pre-Commit HARD Enforcement (§22 monolitico)

| Tier | Trigger | SLO | Coverage |
|------|---------|-----|----------|
| **Tier 1** | `git commit` | < 30s | TDD impactados + lint + secrets + SAST + Trivy fs |
| **Tier 2** | `git push` | < 3min | Tier 1 + govulncheck + integration impactados + Supply Chain |
| **Tier 3** | `make premerge` | < 15min | Todos + full integration + E2E + IaC + SLSA L3 |

**Hard enforcement:** `--no-verify` BLOQUEADO via `scripts/git-hooks-wrapper.sh`; emergency override apenas via `EMERGENCY_BYPASS=true` + reason + audit + Slack alerta.

## Sources

- §22 (Pre-Commit HARD Enforcement) do [snapshot monolitico v3.11.7](../_archive/allenty-v3.11.7-monolithic-plan.md)
- Pattern de referencia: [LedgerOS 07-cicd](../../../../ledgeros/.base/plans/07-cicd/)
