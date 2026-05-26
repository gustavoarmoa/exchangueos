// Package outbox implements the transactional outbox pattern for at-least-once
// event delivery without distributed transactions.
//
// Pattern:
//
//  1. Domain layer records DomainEvents into the aggregate (RecordEvent).
//  2. Repository.Save persists aggregate state + outbox rows in the SAME DB tx.
//  3. A background worker (cmd/worker) reads pending outbox rows in commit order,
//     publishes to Kafka, then marks them dispatched (and optionally archives).
//
// Why a separate package?
//   - Modules stay free of Kafka client coupling — they depend on `outbox.Store`.
//   - The Kafka client choice (Sarama / franz-go / kgo) is a single replacement
//     point in `pkg/outbox/kafka` (not in every bounded context).
//
// Schema: migrations/000009_create_outbox.up.sql (outbox_events + dispatched archive).
package outbox
