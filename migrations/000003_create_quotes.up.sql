-- ExchangeOS — 000003: quotes + rfqs + quote_streams
-- Money/Rate use DECIMAL(36,18). NEVER float.

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- rfqs — RFQ aggregate
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS rfqs (
    rfq_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL REFERENCES tenants(tenant_id),
    requester    STRING(256) NOT NULL,
    base_ccy     STRING(3) NOT NULL,
    quote_ccy    STRING(3) NOT NULL,
    status       STRING(16) NOT NULL DEFAULT 'REQUESTED',
    version      INT NOT NULL DEFAULT 1,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),

    CHECK (base_ccy <> quote_ccy),
    CHECK (status IN ('REQUESTED','QUOTED','ACCEPTED','REJECTED','EXPIRED')),

    INDEX idx_rfqs_tenant_status (tenant_id, status, created_at DESC),
    INDEX idx_rfqs_pair (tenant_id, base_ccy, quote_ccy, created_at DESC)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- quotes — individual bid/ask snapshots tied to an RFQ or a stream
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS quotes (
    quote_id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     UUID NOT NULL REFERENCES tenants(tenant_id),
    rfq_id        UUID REFERENCES rfqs(rfq_id),   -- nullable: streams generate quotes without an RFQ
    base_ccy      STRING(3) NOT NULL,
    quote_ccy     STRING(3) NOT NULL,
    notional      DECIMAL(36,18) NOT NULL CHECK (notional > 0),
    notional_ccy  STRING(3) NOT NULL,
    bid           DECIMAL(36,18) NOT NULL CHECK (bid > 0),
    ask           DECIMAL(36,18) NOT NULL CHECK (ask > 0),
    valid_from    TIMESTAMPTZ NOT NULL,
    valid_to      TIMESTAMPTZ NOT NULL,
    venue         STRING(64),
    version       INT NOT NULL DEFAULT 1,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),

    CHECK (bid <= ask),
    CHECK (base_ccy <> quote_ccy),
    CHECK (valid_to > valid_from),
    CHECK (notional_ccy IN (base_ccy, quote_ccy)),

    INDEX idx_quotes_rfq (rfq_id),
    INDEX idx_quotes_tenant_pair_valid (tenant_id, base_ccy, quote_ccy, valid_to DESC)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- quote_streams — long-lived streaming sessions (price feeds)
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS quote_streams (
    stream_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL REFERENCES tenants(tenant_id),
    pairs        STRING[] NOT NULL,
    started_at   TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    ended_at     TIMESTAMPTZ,
    last_seen_at TIMESTAMPTZ,

    INDEX idx_streams_tenant_open (tenant_id) WHERE ended_at IS NULL
);

COMMIT;
