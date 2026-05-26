#!/usr/bin/env bash
# .claude/hooks/session-end.sh — Session totally ended
set -euo pipefail

# Compute session metrics + save
SESSION_LOG=".claude/memory/sessions.log"
if [[ -f "$SESSION_LOG" ]]; then
    SESSIONS_TODAY=$(grep "$(date -u +%F)" "$SESSION_LOG" 2>/dev/null | wc -l || echo 0)
    echo "📊 Sessions today: $SESSIONS_TODAY" >&2
fi

# Mark session end
echo "$(date -u +%FT%TZ) | session-end" >> .claude/memory/sessions.log 2>/dev/null || true
