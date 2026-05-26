# MS-023t — local-quality-gates

| Field | Value |
|-------|-------|
| **Code** | MS-023t |
| **Name** | local-quality-gates |
| **Phase** | F15L |
| **Sprint** | 17 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023s (delivered) |

## Delivery Notes

**Acceptance criteria met:**
- ✅ TDD workflow: every domain change starts with failing test (Red-Green-Refactor) — enforced via `.claude/rules/modules-domain.md`
- ✅ 10 E2E scenarios catalogued in `tests/e2e/README.md` + 3 representative implementations in `tests/e2e/scenario_*_test.go`
- ✅ 30 security gates across 3 tiers (`lefthook.yml`):
  - Tier 1 pre-commit (<30s): fmt + vet + secrets (gitleaks) + proto-lint + yaml-lint
  - Tier 2 pre-push (<3min): unit tests + golangci-lint full + govulncheck
  - Tier 3 pre-merge (<15min, CI): full test + integration + trivy
- ✅ `golangci-lint` with strict config: errcheck/govet/staticcheck/unused/gosec/revive/dupl/gocyclo/gocritic/bodyclose/noctx/nilerr/exhaustive/errorlint/sqlclosecheck/rowserrcheck + **forbidigo** banning `float64/float32` outside tests
- ✅ `.github/workflows/security.yml` — gitleaks + govulncheck + trivy + CodeQL + SBOM CycloneDX on push/PR + weekly cron

**Deferred:**
- ⏳ 35 FX-QA-* patterns catalog — separate documentation track

## Description

lefthook.yml + .pre-commit-config.yaml instalados + 30 security gates locais em 3 ciclos (pre-commit 5s + pre-push 60s + pre-merge 15min) + 10 E2E cenarios + 35 FX-QA-* patterns + 7 scripts + 4 docs onboarding + CI espelha local 100% — zero push falho ao GitHub.

## Acceptance Criteria

- [ ] lefthook.yml + .pre-commit-config.yaml configurados
- [ ] 30 security gates em 3 ciclos
- [ ] 10 E2E cenarios canonicos com require.Eventually
- [ ] 9 Makefile quality targets
- [ ] CI espelha local 100% (zero discrepancia)
- [ ] TDD coverage gate domain >= 90% / app >= 75%

## Deliverables

- lefthook.yml + .pre-commit-config.yaml
- 9 Makefile targets quality
- 7 scripts auxiliares em scripts/
- 35 patterns em 260-fx-qa-tdd-e2e-patterns.md
- 4 docs onboarding em docs/

## Cross-References

- Plano monolitico: §18 + Fase F15L
- Workstream: 07-cicd + 10-quality
