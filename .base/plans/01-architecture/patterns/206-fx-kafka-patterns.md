# FX-KP-* — Kafka Patterns (60 patterns)

ExchangeOS Kafka production-grade patterns. Built on franz-go/kgo.

## Catalog (representative)

| # | Title | Status |
|---|-------|--------|
| FX-KP-001 | Idempotent producer + acks=all + zstd | ✅ |
| FX-KP-002 | Topic naming `exchangeos.<bc>.events` + per-bc isolation | ✅ |
| FX-KP-003 | partition_key = aggregate_id ensures per-aggregate ordering | ✅ |
| FX-KP-004 | ACL policy per service identity | ✅ |
| FX-KP-005 | Consumer group config + manual offset commit | ⏳ |

---

## FX-KP-001 — Idempotent producer + acks=all + zstd

**Context:** Production-grade publisher reliability + throughput.

**Problem:** Default Kafka client config has weak durability + uncompressed payloads.

**Solution:** Configure:
- `RequiredAcks=AllISRAcks` (leader + all in-sync replicas)
- `ProducerBatchCompression=ZstdCompression` (best ratio + cpu balance)
- Idempotent producer enabled (auto-dedupe on broker retry)
- `ProducerLinger=5ms` (small batch window for throughput)
- `ProducerBatchMaxBytes=1<<20` (1 MiB)
- `MaxBufferedRecords=10_000`

**Example:** `pkg/outbox/kafka/publisher.go:New(cfg)` — the full kgo.NewClient option list.

**Anti-pattern:** Default `kgo.NewClient(kgo.SeedBrokers(...))` — acks=leader, no compression.

**Related:** FX-EDA-002, FX-KP-002.

---

## FX-KP-002 — Topic naming `exchangeos.<bc>.events`

**Context:** Topic catalogue for the platform.

**Problem:** Ad-hoc topic names (e.g. `trade-events`, `trades_v2`) lead to discoverability hell + ACL drift.

**Solution:** `exchangeos.<bounded_context>.<topic_class>` lowercase + dot-separated. Versioning via `.v2` suffix when format-breaking. Compacted feeds (e.g. refdata) use `<topic>.compacted` suffix.

**Example:** `deploy/kafka/topics.yaml` defines 14 topics following this scheme:
- `exchangeos.trade.events`
- `exchangeos.quote.events`
- `exchangeos.cls_settlement.events`
- `exchangeos.refdata.spot_rates` (cleanup_policy: compact)

**Anti-pattern:** `trade.events`, `trades`, `fx.trade`.

**Related:** FX-KP-003, FX-KP-004.

---

## FX-KP-003 — partition_key = aggregate_id

**Context:** Per-aggregate ordering across consumers.

**Problem:** Without a partition key, events for the same trade can land on different partitions, breaking causal ordering downstream.

**Solution:** Always set `Record.Key = aggregateID.String()` (or partition_key column). Kafka hashes to a stable partition.

**Example:** `pkg/outbox/Dispatch` defaults `key = AggregateID.String()` when PartitionKey field is empty.

**Anti-pattern:** Letting Kafka round-robin without a key — settled-event arrives before created-event on a slow consumer.

**Related:** FX-EDA-004, FX-KP-001.

---

## FX-KP-004 — ACL policy per service identity

**Context:** Multi-tenant cluster shared with sibling Revenu Platform modules.

**Problem:** Without ACLs, any service can produce/consume any topic — supply-chain risk + audit failure.

**Solution:** Per-service Kafka ACLs scoped to the topics each service actually produces/consumes. SASL/SCRAM identity + ACL definition in IaC.

**Example:** `deploy/kafka/topics.yaml` `acls:` section maps each of 5 service identities (api/worker/cls-cycle/eod/mq-bridge) to explicit `produce:` and `consume:` topic lists.

**Anti-pattern:** Single shared SASL user with `ALLOW *` on `*`.

**Related:** FX-KP-002, FX-DS-008 (WIF for GCP — same principle for cluster identities).
