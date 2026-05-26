#!/usr/bin/env bash
# .claude/hooks/subagent-stop.sh — Subagent completed
set -euo pipefail

INPUT=$(cat)
AGENT=$(echo "$INPUT" | jq -r '.subagent.type // "unknown"' 2>/dev/null || echo unknown)

echo "$(date -u +%FT%TZ) | subagent-stop | $AGENT" >> .claude/memory/agent-runs.log 2>/dev/null || true

exit 0
