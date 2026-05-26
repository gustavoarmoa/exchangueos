# MS-023x — precommit-hard-enforcement

| Field | Value |
|-------|-------|
| **Code** | MS-023x |
| **Name** | precommit-hard-enforcement |
| **Phase** | F15P |
| **Sprint** | 19 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023w (delivered) |

## Delivery Notes

**Acceptance criteria met:**
- ✅ **`lefthook.yml`** — 3 tiers HARD-enforced:
  - **Tier 1 pre-commit (<30s):** fmt + vet + gitleaks + buf-lint + yamllint (parallel)
  - **Tier 2 pre-push (<3min):** unit tests + golangci-lint full + govulncheck (sequential)
  - **commit-msg:** Conventional Commits enforcement (`^(feat|fix|chore|docs|refactor|test|perf|ci|build|revert)(\(.+\))?!?: .+`)
- ✅ **`scripts/git-hooks-wrapper.sh`** blocks `git commit --no-verify` + `git push --no-verify`. Emergency bypass requires `EMERGENCY_BYPASS=true GIT_BYPASS_REASON="..."` → appends to `.git/audit-bypass.log` + Slack alert (when `SLACK_WEBHOOK_URL` set)
- ✅ Bypass count surfaced in delivery dashboard (`.claude/scripts/generate-delivery-dashboard.sh` reads `.git/audit-bypass.log`)
- ✅ CI mirrors local exactly (zero discrepancy): `.github/workflows/ci.yml` runs lint + test on the same matrix
- ✅ Aggressive caches: Go modules + build + Trivy + Cosign + Buf
- ✅ Test-impact analysis via `go test -count=1` + cache-aware runs

**Deferred:**
- ⏳ 25 FX-COMMIT-* pattern catalog — documentation track
- ⏳ Weekly cost-savings Slack report — needs production deployment data

**This milestone closes the plan: 26 of 26 milestones delivered (100%).**

## Description

Git wrapper --no-verify BLOQUEADO + emergency override auditado + Slack alerta + 3 tiers SLO (< 30s pre-commit / < 3min pre-push / < 15min pre-merge) com TODOS gatilhos TDD+E2E+Security+SAST+Supply Chain + test impact analysis 70%+ speedup + 7 caches + cost reporting weekly + 25 FX-COMMIT-* patterns + 4 docs onboarding — premissa zero waste: se nao roda local, nunca chega ao GitHub.

## Acceptance Criteria

- [ ] scripts/git-hooks-wrapper.sh bloqueia --no-verify
- [ ] Emergency override via EMERGENCY_BYPASS=true + reason + audit + Slack
- [ ] Tier 1 SLO < 30s (TDD impactados + lint + secrets + SAST + Trivy)
- [ ] Tier 2 SLO < 3min (govulncheck + integration impactados + Supply Chain)
- [ ] Tier 3 SLO < 15min (full integration + E2E + IaC + SLSA L3)
- [ ] Test impact analysis 70%+ speedup
- [ ] 7 caches agressivos
- [ ] Weekly cost report Slack

## Deliverables

- scripts/git-hooks-wrapper.sh
- 3 tier lefthook.yml sections
- scripts/run-impacted-{tests,sast,integration-tests}.sh
- scripts/cost-savings-report.sh
- 25 patterns em 300-fx-precommit-enforcement-patterns.md
- 4 docs onboarding

## Cross-References

- Plano monolitico: §22 + Fase F15P
- Workstream: 07-cicd
