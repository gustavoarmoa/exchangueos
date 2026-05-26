#!/usr/bin/env bash
# .claude/hooks/pre-bash.sh — Pre-flight check antes de Bash exec
set -euo pipefail

INPUT=$(cat)
CMD=$(echo "$INPUT" | jq -r '.tool_input.command // ""' 2>/dev/null || echo "")

# Audit log
echo "$(date -u +%FT%TZ) | $(git config user.email 2>/dev/null || echo unknown) | $CMD" \
    >> .claude/cache/bash-audit.log 2>/dev/null || true

# Block destructive patterns que escaparam do deny list
DANGER_PATTERNS=("rm -rf /" "rm -rf ~" "dd if=" "mkfs" ":(){ :|:& };:")
for pattern in "${DANGER_PATTERNS[@]}"; do
    if [[ "$CMD" == *"$pattern"* ]]; then
        echo "❌ BLOCKED: dangerous pattern detected: $pattern" >&2
        exit 2
    fi
done

# Warn on slow commands sem timeout
if [[ "$CMD" == *"go test ./..."* && "$CMD" != *"-timeout"* ]]; then
    echo "⚠ tip: add -timeout 5m to avoid hung tests" >&2
fi

exit 0
