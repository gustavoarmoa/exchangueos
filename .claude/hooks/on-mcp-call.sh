#!/usr/bin/env bash
# .claude/hooks/on-mcp-call.sh — Reage quando uma ferramenta MCP e invocada
# Triggered antes/apos chamadas MCP (logging + validation)
set -euo pipefail

MCP_TOOL="${1:-unknown}"
TIMESTAMP="$(date -u +%FT%TZ)"

# Audit log local
echo "{\"timestamp\":\"$TIMESTAMP\",\"tool\":\"$MCP_TOOL\",\"event\":\"mcp_call\"}" >> .claude/mcp-audit.log

# Se tool e sensitive (ex: filesystem write), valida
case "$MCP_TOOL" in
    *write*|*delete*|*exec*)
        echo "⚠ MCP tool sensitive: $MCP_TOOL — audit logged"
        ;;
esac
