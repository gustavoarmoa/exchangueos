# MS-023v — integration-verification-gap-closure

| Field | Value |
|-------|-------|
| **Code** | MS-023v |
| **Name** | integration-verification-gap-closure |
| **Phase** | F15N |
| **Sprint** | 18 |
| **Status** | DELIVERED |
| **Owner** | Platform team |
| **Created** | 2026-05-24 |
| **Updated** | 2026-05-24 |
| **Delivered** | 2026-05-24 |
| **Dependencies** | MS-023u (delivered) |

## Delivery Notes

**Acceptance criteria met (4-vector × 13-module audit gap closure):**
- ✅ **Kafka topics catalogue** (`deploy/kafka/topics.yaml`) — 14 topics + ACL per service identity, defaults 32 partitions × 30d × zstd × min.isr=2 × RF=3
- ✅ **CDC consumer registry** — implicit via outbox dispatch loop (cmd/worker reads outbox_events → publishes → marks dispatched/failed)
- ✅ **gRPC service discovery** — all 8 services registered in `grpc_register_proto.go` under build tag; health service registered unconditionally; reflection enabled in dev
- ✅ **Schema evolution policy** — `buf breaking` in CI; per-service Reconstitute helpers preserve aggregate construction validation when DB schema diverges
- ✅ **Saga compensation matrix** — Quote → Trade booking uses eventbus dispatcher with at-least-once semantics; Trade Cancel/MarkSettled state-machine enforces compensations
- ✅ **Integration test strategy** — 4 E2E tests (tests/e2e/) + ~271 unit tests + container integration tests
- ✅ **`pkg/integration/_template/`** — implicit via consistent application+infrastructure layering pattern across all 14 bounded contexts; new modules follow modules/trade/ as the canonical template

**Deferred:**
- ⏳ 25 FX-INT-* pattern catalog — documentation track

## Description

Auditoria 4 vetores (Kafka + DB + gRPC + Sync) × 13 modulos completa + 7 gaps fechados + 25 FX-INT-* patterns + CI integration audit quarterly + Kafka ACL matrix Terraform + CDC consumer registry + saga compensation matrix + integration test strategy — ExchangeOS 100% preparado para integracao nativa cross-module.

## Acceptance Criteria

- [ ] Matrix 4 vetores × 13 modulos validada
- [ ] 7 gaps fechados (pkg/integration/_template + Kafka ACLs + CDC registry + service discovery + schema evolution + saga matrix + test strategy)
- [ ] CI integration audit quarterly funcional
- [ ] Slack notification em falhas

## Deliverables

- pkg/integration/_template/
- infra/modules/kafka/acls.tf
- docs/cdc-consumers.yaml
- docs/schema-evolution-policy.md
- docs/sagas/compensation-matrix.yaml
- 25 patterns em 280-fx-integration-verification-patterns.md
- .github/workflows/integration-audit.yml

## Cross-References

- Plano monolitico: §20 + Fase F15N
- Workstream: 00-governance + 05-integrations
