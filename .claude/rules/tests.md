---
glob: "tests/**/*_test.go"
---

# Rule: tests/**/*_test.go

## TDD Conventions
- Test name cita RN_FX_NNN quando aplicavel (`TestRN_FX_010_PvP_CLS_Eligible`)
- Table-driven tests preferido
- `testify/require` para fatal preconditions; `assert` para non-fatal
- `t.Parallel()` para concurrent safe tests
- `t.TempDir()` em vez de manual dirs
- `t.Cleanup()` em vez de defer

## Integration Tests
- Build tag `//go:build integration`
- testcontainers-go para CRDB + Kafka + Vault
- TLS verify-full pattern (shared CA)
- Per-test schema isolation (CREATE SCHEMA tenant_test_<id>)

## E2E Tests
- Build tag `//go:build e2e`
- `require.Eventually` (NUNCA `time.Sleep` arbitrario)
- docker-compose full stack via `make local-up`

## Coverage Gate
- Domain >= 80% (target 90%)
- Application >= 70%
