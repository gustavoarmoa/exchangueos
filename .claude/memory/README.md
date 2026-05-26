# Memory — Persistent Cross-Session Context

> Sessions log + agent runs log + MCP knowledge graph

## Files

| File | Purpose | Gitignored |
|------|---------|------------|
| `sessions.log` | Session start/end timestamps | YES (personal) |
| `agent-runs.log` | Subagent invocation history | YES |
| `mcp-memory.json` | MCP memory server knowledge graph | YES |

## Usage

MCP memory server persiste knowledge entre sessoes.
Use via MCP tool calls (auto-loaded via `.mcp.json`).

## Privacy

NUNCA commitar — gitignored. Personal session data.
