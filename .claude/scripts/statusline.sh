#!/usr/bin/env bash
# .claude/scripts/statusline.sh — Custom status line for ExchangeOS
# Receives JSON via stdin with session state; outputs 1-line status

INPUT=$(cat)
MODEL=$(echo "$INPUT" | jq -r '.model.display_name // "claude"' 2>/dev/null || echo "claude")
CWD=$(echo "$INPUT" | jq -r '.workspace.current_dir // "."' 2>/dev/null || pwd)

# Git context
BRANCH=$(git -C "$CWD" branch --show-current 2>/dev/null || echo "no-git")
DIRTY=""
if [[ -n "$BRANCH" ]] && [[ -n "$(git -C "$CWD" status --porcelain 2>/dev/null)" ]]; then
    DIRTY="*"
fi

# Active milestone (read from .base/plans/milestones/active/)
ACTIVE_MS=""
if [[ -d "$CWD/.base/plans/milestones/active" ]]; then
    ACTIVE_MS=$(ls "$CWD/.base/plans/milestones/active/" 2>/dev/null | head -1 | sed 's/\.md$//')
fi

# Output: model | branch | active milestone
printf "🔷 %s │ ⎇ %s%s │ 🎯 %s" \
    "$MODEL" \
    "${BRANCH:-no-git}" \
    "$DIRTY" \
    "${ACTIVE_MS:-no-active-milestone}"
