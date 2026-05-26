#!/usr/bin/env bash
# .claude/hooks/pre-push.sh — Type-check + tests antes do push
# Triggered automaticamente em git push
set -euo pipefail

echo "🔍 Pre-push hook (.claude/hooks/pre-push.sh)..."

# Roda Tier 2 via lefthook (60s budget)
if command -v lefthook >/dev/null 2>&1; then
    lefthook run pre-push
else
    echo "⚠ lefthook nao instalado — fallback minimo"
    go build ./... && go test -race -short ./...
fi

echo "✅ Pre-push checks passed"
