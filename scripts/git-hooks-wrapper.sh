#!/usr/bin/env bash
# ExchangeOS — git-hooks-wrapper (FX-COMMIT-002)
# Blocks `git commit --no-verify` and `git push --no-verify`.
# Emergency bypass: EMERGENCY_BYPASS=true GIT_BYPASS_REASON="..." git ...
#   → appends to .git/audit-bypass.log + emits SLACK alert (if SLACK_WEBHOOK_URL set).

set -euo pipefail

ARGS=("$@")
SUBCMD="${1:-}"

ROOT="$(git rev-parse --show-toplevel 2>/dev/null || echo ".")"
LOG="$ROOT/.git/audit-bypass.log"

# Detect --no-verify / -n with subcommand commit|push|merge
for a in "${ARGS[@]}"; do
    case "$a" in
        --no-verify|-n)
            if [[ "${EMERGENCY_BYPASS:-false}" != "true" ]]; then
                echo "❌ ExchangeOS: --no-verify is BLOCKED."
                echo "   Pre-commit/pre-push hooks are HARD-enforced per FX-COMMIT-* patterns."
                echo "   For genuine emergencies set: EMERGENCY_BYPASS=true GIT_BYPASS_REASON=\"...\" git $SUBCMD ..."
                exit 1
            fi
            if [[ -z "${GIT_BYPASS_REASON:-}" ]]; then
                echo "❌ EMERGENCY_BYPASS=true requires GIT_BYPASS_REASON=\"<reason>\""
                exit 1
            fi
            mkdir -p "$(dirname "$LOG")"
            printf '%s\t%s\t%s\t%s\n' \
                "$(date -u +%FT%TZ)" "$(whoami)" "$SUBCMD" "$GIT_BYPASS_REASON" >> "$LOG"
            if [[ -n "${SLACK_WEBHOOK_URL:-}" ]]; then
                curl -fsS -X POST -H 'Content-Type: application/json' \
                    -d "{\"text\":\"🚨 ExchangeOS bypass: $(whoami) — $SUBCMD — $GIT_BYPASS_REASON\"}" \
                    "$SLACK_WEBHOOK_URL" >/dev/null 2>&1 || true
            fi
            echo "⚠️  EMERGENCY BYPASS logged: $LOG"
            ;;
    esac
done

exec git "${ARGS[@]}"
