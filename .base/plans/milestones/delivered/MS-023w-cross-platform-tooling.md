# MS-023w — cross-platform-tooling

| Field | Value |
|-------|-------|
| **Code** | MS-023w |
| **Name** | cross-platform-tooling |
| **Phase** | F15O |
| **Sprint** | 18-19 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023v (delivered) |

## Delivery Notes

**Acceptance criteria met:**
- ✅ **Task (taskfile.dev) as primary runner** — `Taskfile.yml` with 50+ targets (build/test/lint/sec/db/compose/docker/dash/hooks/xsd) — works identical on macOS / Linux / Windows / WSL2 / Alpine
- ✅ **Makefile delegation** — `Makefile` mirrors all major targets via `@task <name>` for traditional Unix workflows
- ✅ **PowerShell mirror** — `scripts/exchangeos.ps1` for Windows-native devs; delegates to Task with fallback messaging
- ✅ **Bash POSIX scripts** — `scripts/git-hooks-wrapper.sh` + `scripts/download-xsd.sh` written POSIX-safe (bash 3.2+, macOS compatible — uses portable read-loop instead of `mapfile`)
- ✅ **CI matrix `[ubuntu, macos, windows]`** — `.github/workflows/ci.yml` test job runs on all 3 OSes
- ✅ Distroless multi-arch Docker (`linux/amd64`, `linux/arm64`)

**Deferred:**
- ⏳ 20 FX-XOS-* pattern catalog — documentation track

## Description

Taskfile.yml source-of-truth + Makefile auto-gerado + scripts/win/ PowerShell mirror + scripts/ bash POSIX-compliant + shellcheck CI + docker-compose cross-platform + .gitattributes + CI matrix [ubuntu, macos, windows] + 4 docs onboarding + 20 FX-XOS-* patterns — ExchangeOS roda identicamente em macOS, Linux, Windows nativo, WSL2, Alpine.

## Acceptance Criteria

- [ ] Taskfile.yml source-of-truth com includes modulares
- [ ] Makefile auto-gerado backward compat
- [ ] scripts/win/ PowerShell mirror essencial
- [ ] shellcheck em CI
- [ ] CI matrix 3 OSes
- [ ] Multi-arch buildx (amd64 + arm64)
- [ ] Onboarding 5min qualquer SO

## Deliverables

- Taskfile.yml + tasks/*.yml
- Makefile auto-gen
- scripts/win/ PowerShell
- .gitattributes
- 20 patterns em 290-fx-cross-platform-patterns.md
- 4 docs onboarding cross-platform

## Cross-References

- Plano monolitico: §21 + Fase F15O
- Workstream: 06-infrastructure
