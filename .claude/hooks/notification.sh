#!/usr/bin/env bash
# .claude/hooks/notification.sh — Reage a notification events
set -euo pipefail

INPUT=$(cat)
TYPE=$(echo "$INPUT" | jq -r '.notification.type // "info"' 2>/dev/null || echo info)
MESSAGE=$(echo "$INPUT" | jq -r '.notification.message // ""' 2>/dev/null || echo "")

# Log local
echo "$(date -u +%FT%TZ) | $TYPE | $MESSAGE" >> .claude/cache/notifications.log 2>/dev/null || true

# Critical notifications → desktop notification (macOS)
if [[ "$TYPE" == "error" || "$TYPE" == "critical" ]]; then
    if command -v osascript >/dev/null 2>&1; then
        osascript -e "display notification \"$MESSAGE\" with title \"ExchangeOS Claude\""
    fi
fi

exit 0
