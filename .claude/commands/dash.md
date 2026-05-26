---
description: Update delivery dashboard (auto-gen from milestones state + git stats)
allowed-tools: [Bash, Read]
---

# /dash

Atualiza o **delivery dashboard** baseado em estado atual:

```bash
.claude/scripts/generate-delivery-dashboard.sh
```

Outputs:
- Atualiza `.base/plans/roadmap/delivery-dashboard.md` (snapshot block)
- Cria snapshot HTML em `.base/plans/00-governance/audits/dashboard-snapshot-<date>.html`
- (Se SLACK_WEBHOOK_URL set) posta Slack notification

Cron equivalent: `.github/workflows/dashboard-update.yml` (hourly + on milestone changes)

## Output Example
```
═══════════════════════════════════════════════════════════════
  ExchangeOS Delivery Dashboard Update — 2026-05-24
═══════════════════════════════════════════════════════════════
  Overall:    ███░░░░░░░░░░░░░░░░░  15% (4/26)
  BACKLOG:    22
  ACTIVE:     0
  DELIVERED:  4
  Sprint:     3
  Commits 30d: 87  |  7d: 24
  Bypass:     0
  Plan:       v4.1.0
═══════════════════════════════════════════════════════════════
```
