---
description: Build all binaries (api, worker, cls-cycle, eod, mq-bridge, cred-rotator, migrator)
allowed-tools: [Bash]
---

# /build

Build cross-platform via Task (Taskfile.yml):

```bash
task build
```

Builds:
- bin/api
- bin/worker
- bin/cls-cycle
- bin/eod
- bin/mq-bridge
- bin/cred-rotator
- bin/migrator
