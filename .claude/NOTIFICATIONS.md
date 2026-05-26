# Notifications — ExchangeOS Claude Code

## Routing

| Event | Destination | Severity threshold |
|-------|-------------|-------------------|
| Bash failure (CI gate) | Slack #exchangeos-platform | warn+ |
| Pre-commit bypass | Slack + audit log | always |
| MCP error | Local log + desktop | error+ |
| Subagent timeout | Local log | warn+ |
| Session end summary | (silent) | info |

## Configuration

Set Slack webhook in `CLAUDE.local.md` or env:
```bash
export SLACK_WEBHOOK_URL='https://hooks.slack.com/...'
```

## Notification Hook

`.claude/hooks/notification.sh` handles routing:
- Critical → desktop notification (osascript on macOS)
- All → local log `.claude/cache/notifications.log`
