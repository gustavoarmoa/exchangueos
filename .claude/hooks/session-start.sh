#!/usr/bin/env bash
# .claude/hooks/session-start.sh — Loads context at session start
set -euo pipefail

echo "📋 ExchangeOS session start"
echo "  Module: ExchangeOS — Standalone FX Module"
echo "  Plan version: $(cat .base/plans/version.md | grep '^**Current:**' | awk '{print $2}' | tr -d '\`')"
echo "  Active milestone: (check .base/plans/milestones/active/)"

# Check git status briefly
git status --short --branch | head -5

# Check pending pre-commit hooks
if [[ -f lefthook.yml && ! -d .git/hooks ]]; then
    echo "⚠ lefthook hooks nao instalados. Rodar: make install-hooks"
fi
