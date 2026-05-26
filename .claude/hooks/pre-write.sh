#!/usr/bin/env bash
# .claude/hooks/pre-write.sh — Pre-flight antes de file write
set -euo pipefail

INPUT=$(cat)
FILE=$(echo "$INPUT" | jq -r '.tool_input.file_path // ""' 2>/dev/null || echo "")

# Block writes em arquivos sensitive
if [[ "$FILE" == *.env || "$FILE" == *.key || "$FILE" == *.pem ]]; then
    echo "❌ BLOCKED: NUNCA write em $FILE (secrets/credentials)" >&2
    exit 2
fi

# Block writes em .git/
if [[ "$FILE" == *.git/* && "$FILE" != *.git/info/* ]]; then
    echo "❌ BLOCKED: NUNCA write direto em .git/" >&2
    exit 2
fi

# Warn on writes em _archive/
if [[ "$FILE" == *_archive/* ]]; then
    echo "⚠ Write em _archive/ (read-only snapshot) — confirme se intentional" >&2
fi

exit 0
