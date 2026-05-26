---
description: Invoke a specialized subagent (15 disponiveis)
allowed-tools: [Task]
argument-hint: <agent-name> <prompt>
---

# /agent

Invoca subagent especializado:

`/agent fx-domain "Refactor FXTrade aggregate to support multi-leg swap"`
`/agent bacen-compliance "Validate trade 12345 against all BACEN rules"`
`/agent pricing-quant "Add new NDF formula for INR"`

Agents disponiveis: fx-domain, iso20022, bacen-compliance, pricing-quant, cls-settlement, cfets-confirmation, ontology-shacl, database-crdb, kafka-flink, iam-security, observability-otel, testing-qa, devsecops-cicd, infra-k8s-terraform, cross-platform.
