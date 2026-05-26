---
description: Milestone operations (list, show, start, complete)
allowed-tools: [Bash, Read, Edit]
argument-hint: <list|show|start|complete> [MS-XXX]
---

# /milestone

Manage milestones em `.base/plans/milestones/`:

- `/milestone list` — list all (BACKLOG/ACTIVE/DELIVERED)
- `/milestone show MS-023a` — show details
- `/milestone start MS-023a` — move backlog → active
- `/milestone complete MS-023a` — move active → delivered + update CHANGELOG
