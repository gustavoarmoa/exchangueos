#!/usr/bin/env bash
# .claude/scripts/generate-delivery-dashboard.sh
# Auto-generates delivery-dashboard.md from milestones state + git stats
# Usage: task dash-update  (ou direto: ./generate-delivery-dashboard.sh)

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
PLANS="$ROOT/.base/plans"
DASH="$PLANS/roadmap/delivery-dashboard.md"
SNAPSHOT_DIR="$PLANS/00-governance/audits"
TODAY="$(date -u +%F)"

# Count milestones per status
BACKLOG=$(find "$PLANS/milestones/backlog" -name "MS-*.md" 2>/dev/null | wc -l | tr -d ' ')
ACTIVE=$(find "$PLANS/milestones/active" -name "MS-*.md" 2>/dev/null | wc -l | tr -d ' ')
DELIVERED=$(find "$PLANS/milestones/delivered" -name "MS-*.md" 2>/dev/null | wc -l | tr -d ' ')
TOTAL=$((BACKLOG + ACTIVE + DELIVERED))

# Compute percentage delivered
if [[ $TOTAL -gt 0 ]]; then
    PCT=$(( DELIVERED * 100 / TOTAL ))
else
    PCT=0
fi

# Progress bar (20 chars)
FILLED=$(( PCT * 20 / 100 ))
EMPTY=$(( 20 - FILLED ))
BAR=$(printf '█%.0s' $(seq 1 $FILLED 2>/dev/null))$(printf '░%.0s' $(seq 1 $EMPTY 2>/dev/null))

# Determine current sprint (best-effort: parse from active milestone or fall back)
CURRENT_SPRINT=$(grep -r "Sprint" "$PLANS/milestones/active/" 2>/dev/null | head -1 | grep -oE "[0-9]+" | head -1 || echo "0")

# Git stats — defensive (em repos sem commits, retorna 0)
count_commits() {
    local since="$1"
    local n
    n=$(git -C "$ROOT" log --since="$since" --oneline 2>/dev/null | wc -l 2>/dev/null | tr -d '[:space:]') || n=0
    printf '%s' "${n:-0}"
}
COMMITS_30D=$(count_commits "30 days ago")
COMMITS_7D=$(count_commits "7 days ago")

# Bypass count
BYPASS_COUNT=0
if [[ -f "$ROOT/.git/audit-bypass.log" ]]; then
    BYPASS_COUNT=$(wc -l < "$ROOT/.git/audit-bypass.log" | tr -d ' ')
fi

# Plan version
PLAN_VERSION=$(grep -E '^\*\*Current:\*\*' "$PLANS/version.md" 2>/dev/null | awk '{print $2}' | tr -d '`' || echo "unknown")

# Generate snapshot HTML for archive
mkdir -p "$SNAPSHOT_DIR"
SNAPSHOT="$SNAPSHOT_DIR/dashboard-snapshot-$TODAY.html"
cat > "$SNAPSHOT" << HTML_EOF
<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>ExchangeOS Delivery Snapshot $TODAY</title>
<style>body{font-family:monospace;background:#0d1117;color:#c9d1d9;padding:20px}</style>
</head><body>
<h1>ExchangeOS Delivery Snapshot — $TODAY</h1>
<pre>
Overall:      $BAR  ${PCT}% delivered (${DELIVERED}/${TOTAL} MS)
Sprint:       ${CURRENT_SPRINT}
Commits 30d:  ${COMMITS_30D}
Commits 7d:   ${COMMITS_7D}
Bypass count: ${BYPASS_COUNT}
Plan version: ${PLAN_VERSION}

BACKLOG:   ${BACKLOG}
ACTIVE:    ${ACTIVE}
DELIVERED: ${DELIVERED}
</pre>
</body></html>
HTML_EOF

# Print summary to stdout
cat << EOF
═══════════════════════════════════════════════════════════════
  ExchangeOS Delivery Dashboard Update — $TODAY
═══════════════════════════════════════════════════════════════
  Overall:    $BAR ${PCT}% (${DELIVERED}/${TOTAL})
  BACKLOG:    ${BACKLOG}
  ACTIVE:     ${ACTIVE}
  DELIVERED:  ${DELIVERED}
  Sprint:     ${CURRENT_SPRINT}
  Commits 30d: ${COMMITS_30D}  |  7d: ${COMMITS_7D}
  Bypass:     ${BYPASS_COUNT}
  Plan:       v${PLAN_VERSION}

  Snapshot:   ${SNAPSHOT}
  Dashboard:  ${DASH}
═══════════════════════════════════════════════════════════════
EOF

# Update Snapshot Executivo block in delivery-dashboard.md (in-place patch)
if [[ -f "$DASH" ]]; then
    # Use python for safe in-place YAML-like block replacement
    python3 << PYEOF
import re, sys
from pathlib import Path

dash = Path("$DASH")
if not dash.exists():
    sys.exit(0)

content = dash.read_text()

# Replace the snapshot block
new_snapshot = f"""\`\`\`
┌──────────────────────────────────────────────────────────────────┐
│  ExchangeOS Delivery Status — Sprint $CURRENT_SPRINT of 19 (auto-updated)        │
├──────────────────────────────────────────────────────────────────┤
│  Overall:      {'█' * $FILLED}{'░' * $EMPTY}  $PCT% delivered ($DELIVERED/$TOTAL MS)     │
│  This sprint:  See active/ milestones                            │
│  Velocity:     Commits 7d: $COMMITS_7D │ 30d: $COMMITS_30D                          │
│  Health:       {'🟢 ON TRACK' if $DELIVERED > 0 else '🟡 DRAFT — pending approval'}                          │
│  Last update:  $TODAY                                          │
└──────────────────────────────────────────────────────────────────┘
\`\`\`
"""

pattern = re.compile(r'\`\`\`\n┌─+┐\n│\s*ExchangeOS Delivery Status.*?└─+┘\n\`\`\`\n', re.DOTALL)
content = pattern.sub(new_snapshot, content, count=1)

dash.write_text(content)
print(f"✅ Updated snapshot block in {dash}")
PYEOF
fi

# Slack notification (if webhook configured)
if [[ -n "${SLACK_WEBHOOK_URL:-}" ]]; then
    curl -X POST -H 'Content-Type: application/json' \
        -d "{\"text\":\"📊 ExchangeOS Delivery: ${PCT}% (${DELIVERED}/${TOTAL}) | Sprint ${CURRENT_SPRINT} | Commits 7d: ${COMMITS_7D}\"}" \
        "$SLACK_WEBHOOK_URL" 2>/dev/null || true
fi

exit 0
