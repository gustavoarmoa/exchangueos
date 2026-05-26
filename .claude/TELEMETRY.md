# Telemetry — ExchangeOS Claude Code

## What we track (locally only)

| Metric | Storage | Purpose |
|--------|---------|---------|
| Sessions per day | `.claude/memory/sessions.log` | Usage patterns |
| Subagent invocations | `.claude/memory/agent-runs.log` | Effectiveness |
| Bash command audit | `.claude/cache/bash-audit.log` | Security audit |
| Bash failures | `.claude/cache/bash-failures.log` | Reliability |
| Hook bypass | `.git/audit-bypass.log` | Compliance |
| MCP calls | `.claude/cache/mcp-audit.log` | MCP usage |
| Notifications | `.claude/cache/notifications.log` | Alert history |

## What we DON'T track

- Personal info (CPF, CNPJ, account numbers — sanitized by hook)
- Secrets (blocked by deny list in settings.json)
- Prompts content (privacy)

## Reports

- `/cost-savings-report` — weekly cost reporting via lefthook telemetry
- `.claude/memory/sessions.log` — session frequency analysis

## Retention

`cleanupPeriodDays: 30` em settings.json — cache auto-pruned.
Memory logs sao persistent (knowledge graph) sem auto-cleanup.
