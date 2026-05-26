---
description: Spawn multiple agents in parallel (cross-cutting tasks)
allowed-tools: [Task]
argument-hint: <task description>
---

# /parallel

Para tarefas cross-cutting (multiple BCs/concerns), spawn agents em paralelo:

Example: `/parallel "Implement PayIn ACK feature end-to-end"`

→ Spawns simultaneamente:
- fx-domain (modelagem PayInACK)
- iso20022 (camt.063 marshaling)
- database-crdb (migration + ERD)
- kafka-flink (publish event)
- bacen-compliance (validate rules)
- observability-otel (spans + metrics)
- testing-qa (TDD tests)

Collects results → integrates → presents coherent diff.
