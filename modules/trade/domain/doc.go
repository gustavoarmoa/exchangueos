// Package domain — FXTrade bounded context (DDD aggregate root pattern).
//
// Conventions (enforced by .claude/rules/modules-domain.md):
//
//   - Aggregate root is the SOLE entry point — references to other aggregates by ID only.
//   - 1 aggregate = 1 transaction. No cross-aggregate writes in a single transaction.
//   - Zero imports of infrastructure (no SQL, no HTTP, no Kafka here).
//   - Constructor returns (*FXTrade, error) with validation.
//   - Pointer-receiver methods on root for state mutation.
//   - `version` field for optimistic concurrency.
//   - Domain events emitted via `RecordEvent(...)`, flushed by infrastructure outbox.
//   - Money/rate use shopspring/decimal — NEVER float.
//
// Business rules cited inline: RN_FX_001..050 (SHACL shapes in
// .base/aasc/ontology/compliance/bacen-cambio-shapes.ttl).
package domain
