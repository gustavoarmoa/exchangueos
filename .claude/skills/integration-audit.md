---
name: integration-audit
description: Run 4-vector × 13-module integration audit (Kafka × DB × gRPC × Sync × 13 modulos)
allowed-tools: [Bash, Read, Grep]
---

# Skill: /integration-audit

## Trigger
`/integration-audit`

## Workflow
1. Verify `pkg/integration/<module>/` para cada um dos 13 modulos
2. Validate Kafka ACLs Terraform vs Topics catalog
3. Validate CDC consumer registry consistency
4. Validate gRPC service discovery + LB config
5. Validate schema evolution policy compliance (proto v1/v2 + Avro BACKWARD)
6. Validate saga compensation matrix
7. Validate integration test coverage
8. Generate HTML report em `.base/plans/00-governance/audits/integration-audit-<date>.html`
9. Post summary no Slack #exchangeos-platform
