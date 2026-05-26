-- ExchangeOS — 000009: outbox_events for transactional outbox pattern.
--
-- The application writes domain events to outbox_events in the SAME transaction
-- as the aggregate save; a separate worker picks rows from outbox_events and
-- publishes them to Kafka, marking them dispatched. This guarantees at-least-once
-- delivery without distributed transactions.

BEGIN;

CREATE TABLE IF NOT EXISTS outbox_events (
    outbox_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL REFERENCES tenants(tenant_id),
    aggregate_type   STRING(64)  NOT NULL,   -- e.g. "Trade", "Quote", "CLSCycle"
    aggregate_id     UUID        NOT NULL,
    event_name       STRING(128) NOT NULL,   -- e.g. "trade.created.v1"
    event_payload    JSONB       NOT NULL,
    topic            STRING(128) NOT NULL,   -- target Kafka topic
    partition_key    STRING(128),            -- usually aggregate_id (string form)
    occurred_at      TIMESTAMPTZ NOT NULL,
    dispatched_at    TIMESTAMPTZ,
    attempt_count    INT         NOT NULL DEFAULT 0,
    last_error       STRING(512),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),

    INDEX idx_outbox_pending (occurred_at) WHERE dispatched_at IS NULL,
    INDEX idx_outbox_aggregate (aggregate_type, aggregate_id, occurred_at DESC),
    INDEX idx_outbox_failed (last_error, attempt_count) WHERE dispatched_at IS NULL AND attempt_count > 0
);

-- ─────────────────────────────────────────────────────────────────────────────
-- outbox_dispatched_archive — optional historical view (kept narrow for audit).
-- Worker MAY truncate dispatched rows older than retention.
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS outbox_dispatched_archive (
    outbox_id        UUID PRIMARY KEY,
    tenant_id        UUID NOT NULL,
    aggregate_type   STRING(64)  NOT NULL,
    aggregate_id     UUID        NOT NULL,
    event_name       STRING(128) NOT NULL,
    topic            STRING(128) NOT NULL,
    partition_key    STRING(128),
    occurred_at      TIMESTAMPTZ NOT NULL,
    dispatched_at    TIMESTAMPTZ NOT NULL,
    attempt_count    INT NOT NULL,

    INDEX idx_archive_agg_time (aggregate_type, aggregate_id, dispatched_at DESC)
);

COMMIT;
