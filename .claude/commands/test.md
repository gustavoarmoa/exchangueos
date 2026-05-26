---
description: Run tests (smart impact-based)
allowed-tools: [Bash]
argument-hint: [unit|integration|crud|e2e|all]
---

# /test

Run tests com smart impact analysis (apenas afetados se git diff disponivel):

- Sem args: roda unit tests impactados
- `unit`: go test -race -short ./...
- `integration`: go test -tags=integration ./tests/integration/...
- `crud`: go test -tags=integration -run '^TestCRUD' ./tests/integration/...
- `e2e`: make e2e-local (full docker-compose stack)
- `all`: tudo (Tier 1+2+3)
