---
description: Show ExchangeOS project status (git + active milestone + cache health)
allowed-tools: [Bash, Read]
---

# /status

Quick health check do ExchangeOS project:

1. Git status + current branch
2. Active milestone (read .base/plans/milestones/active/)
3. Last 5 commits
4. Lefthook status (installed?)
5. Local stack status (docker compose ps)
6. Cache health (.claude/cache/ size)
