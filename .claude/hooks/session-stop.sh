#!/usr/bin/env bash
# .claude/hooks/session-stop.sh — Session ended (graceful)
set -euo pipefail

# Append session summary para memory
TIMESTAMP=$(date -u +%FT%TZ)
echo "$TIMESTAMP | session-stop | $(pwd)" >> .claude/memory/sessions.log 2>/dev/null || true

# Cache prune se > 100MB
CACHE_SIZE=$(du -sm .claude/cache 2>/dev/null | awk '{print $1}' || echo 0)
if [[ "$CACHE_SIZE" -gt 100 ]]; then
    echo "🗑 Cache > 100MB, considering cleanup" >&2
fi
