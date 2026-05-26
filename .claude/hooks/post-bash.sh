#!/usr/bin/env bash
# .claude/hooks/post-bash.sh — Post Bash exec hook
set -euo pipefail

INPUT=$(cat)
EXIT_CODE=$(echo "$INPUT" | jq -r '.tool_response.exit_code // 0' 2>/dev/null || echo 0)
CMD=$(echo "$INPUT" | jq -r '.tool_input.command // ""' 2>/dev/null || echo "")

# Log failures para review
if [[ "$EXIT_CODE" != "0" ]]; then
    echo "$(date -u +%FT%TZ) | EXIT=$EXIT_CODE | $CMD" \
        >> .claude/cache/bash-failures.log 2>/dev/null || true
fi

# Auto-suggest fixes for common failures
if [[ "$CMD" == *"go test"* && "$EXIT_CODE" != "0" ]]; then
    echo "💡 Tip: invoke testing-qa agent for test failure analysis" >&2
fi

if [[ "$CMD" == *"docker compose"* && "$EXIT_CODE" != "0" ]]; then
    echo "💡 Tip: check docker daemon + try 'make local-down && make local-up'" >&2
fi

exit 0
