# FX-DS-* — DevSecOps + CI/CD Patterns (50 patterns)

ExchangeOS supply-chain + delivery patterns.

## Catalog (representative)

| # | Title | Status | Where |
|---|-------|--------|-------|
| FX-DS-001 | Lefthook 3-tier HARD enforcement | ✅ | lefthook.yml + scripts/git-hooks-wrapper.sh |
| FX-DS-002 | gitleaks + govulncheck + trivy in pre-push | ✅ | .github/workflows/security.yml |
| FX-DS-003 | golangci-lint with `forbidigo` float ban | ✅ | .golangci.yml |
| FX-DS-004 | Distroless multi-stage Docker (nonroot 65532, readonly fs) | ✅ | docker/api/Dockerfile + helm templates |
| FX-DS-005 | SLSA L3 build provenance (Cosign keyless) | ⏳ | .github/workflows/ci.yml (sign step) |
| FX-DS-006 | SBOM CycloneDX per release artifact | ✅ | .github/workflows/security.yml (sbom job) |
| FX-DS-007 | Argo Rollouts canary with Prometheus AnalysisTemplate | ✅ | deploy/k8s/argo-rollouts/api-rollout.yaml |
| FX-DS-008 | Workload Identity Federation (zero JSON keys) | ✅ | deploy/terraform/modules/exchangeos-iam |
| FX-DS-009 | CMEK HSM key for GKE etcd encryption | ✅ | deploy/terraform/modules/exchangeos-gke |
| FX-DS-010 | Binary Authorization enforcement on cluster | ✅ | exchangeos-gke module |
| FX-DS-011..050 | (extend on demand) | ⏳ | — |

---

## FX-DS-001 — Lefthook 3-tier HARD enforcement

**Context:** Pre-commit/pre-push gates must stop bad code from leaving the dev machine.

**Problem:** Devs bypass softer hooks (`git commit --no-verify`); CI failures cost cycle time + run-cost.

**Solution:**

```
pre-commit  < 30s  →  fmt + vet + secrets + lint changed + proto-lint + yaml-lint
pre-push    < 3min →  unit tests + lint full + govulncheck
pre-merge   <15min →  full test + integration + trivy   (CI only)
```

`scripts/git-hooks-wrapper.sh` blocks `git --no-verify`. Emergency bypass requires
`EMERGENCY_BYPASS=true GIT_BYPASS_REASON="..."` — logged to `.git/audit-bypass.log` + Slack.

**Example:** `lefthook.yml` + `scripts/git-hooks-wrapper.sh`.

**Anti-pattern:** Optional hooks (default install). Hooks must be HARD gated.

**Related:** FX-DS-002 (security scans), FX-COMMIT-* catalog.

---

## FX-DS-007 — Argo Rollouts canary with Prometheus AnalysisTemplate

**Context:** Deploy a stateless API safely to thousands of req/s.

**Problem:** Big-bang rollouts risk silent regressions; rollbacks are slow when 100% traffic is hit.

**Solution:** `argoproj.io/v1alpha1/Rollout` with canary strategy 10→30→60→100 + AnalysisTemplate
gating on Prometheus metrics (5xx rate + p99 latency). Auto-rollback on
`failureLimit: 3`.

**Example:** `deploy/k8s/argo-rollouts/api-rollout.yaml` — 4 weight steps with 5/10/10m pauses
and two `AnalysisTemplate` metrics (http-5xx-rate < 0.01, http-p99-latency < 500).

**Anti-pattern:** `Deployment` with `RollingUpdate` strategy + manual rollback. The promotion
windows here are non-negotiable for FX-grade reliability.

**Related:** FX-DS-009 (CMEK), FX-DS-010 (Binary Authorization), FX-K8S-* (deployment patterns).

---

## FX-DS-002 — Multi-tool security scan in pre-push

**Context:** Secret leaks, known CVEs, and image vulnerabilities are common slip-throughs.

**Problem:** Single-tool scans miss orthogonal risks; CI-only checks let bad commits propagate locally.

**Solution:** Run three orthogonal scans on every pre-push: `gitleaks` (secret scan), `govulncheck` (Go CVE), `trivy` (filesystem + image). Block the push on HIGH/CRITICAL.

**Example:** `.github/workflows/security.yml` runs all three on push/PR + weekly cron; lefthook `hooks:pre-push` mirrors locally.

**Anti-pattern:** Picking only one — each catches what the others miss.

**Related:** FX-DS-001 (Lefthook), FX-DS-006 (SBOM).

---

## FX-DS-003 — `forbidigo` lint as policy gate

**Context:** Coding rules that aren't expressible as types (e.g. "NEVER float for money") need lint-time enforcement.

**Problem:** Code review catches most violations but not all; tests don't exercise dead code paths that violate.

**Solution:** golangci-lint `forbidigo` with regex rules + actionable error messages.

**Example:** `.golangci.yml` blocks `float64`/`float32` declarations outside tests + generated code:

```yaml
forbidigo:
  forbid:
    - pattern: '^float64\s+'
      msg: "NEVER float64 for money/rate — use shopspring/decimal.Decimal"
```

**Anti-pattern:** Verbal "we don't do floats here" — falls through reviewer attention.

**Related:** FX-GP-002 (decimal precision), FX-DS-002.

---

## FX-DS-008 — Workload Identity Federation (zero JSON keys)

**Context:** Pods need GCP API credentials.

**Problem:** Service-account JSON keys leak in repos/logs/Slack; they don't rotate.

**Solution:** WIF binds the K8s ServiceAccount to a GCP SA via `roles/iam.workloadIdentityUser`. Pods authenticate using their ambient K8s identity — no secrets to manage.

**Example:** `deploy/terraform/modules/exchangeos-iam/main.tf`:

```hcl
resource "google_service_account_iam_member" "wif" {
  service_account_id = google_service_account.exchangeos.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project_id}.svc.id.goog[${var.namespace}/${var.k8s_service_account}]"
}
```

Plus the ServiceAccount annotation `iam.gke.io/gcp-service-account: exchangeos@...`.

**Anti-pattern:** Mounting a Secret with a downloaded JSON key.

**Related:** FX-IAC-* (Terraform GCP), FX-IAM-* (M2M client_secret rotation).
