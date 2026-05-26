#!/usr/bin/env bash
# .claude/hooks/on-file-save.sh — Reage quando arquivo e salvo
set -euo pipefail

FILE="${1:-}"

# Auto-format Go files
if [[ "$FILE" == *.go ]]; then
    gofumpt -l -w "$FILE" 2>/dev/null || true
    goimports -w "$FILE" 2>/dev/null || true
fi

# Auto-lint TTL ontology files
if [[ "$FILE" == *.ttl ]]; then
    echo "TTL saved: $FILE — run /ontology-validate to validate"
fi
