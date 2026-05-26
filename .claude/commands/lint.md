---
description: Run all linters (Go + Docker + YAML + Markdown)
allowed-tools: [Bash]
---

# /lint

Roda todos linters em paralelo:
- golangci-lint run
- hadolint Dockerfile*
- actionlint .github/workflows/*.yml
- yamllint .
- markdownlint .
