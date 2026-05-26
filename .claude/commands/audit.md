---
description: Run integration audit (4 vetores × 13 modulos)
allowed-tools: [Bash, Task]
---

# /audit

Run full integration audit:
1. `/integration-audit` skill (4-vector × 13-module matrix)
2. Verify pkg/integration/<module>/ structure
3. Kafka ACL consistency
4. CDC consumer registry validation
5. Generate HTML report em .base/plans/00-governance/audits/
