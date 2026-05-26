-- ExchangeOS — 000002: fx_trades + counterparties + value_dates
-- Idempotente. NEVER float for money/rate (DECIMAL(36,18) precision).

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- counterparties (PvP + bilateral)
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS counterparties (
    counterparty_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL REFERENCES tenants(tenant_id),
    bic              STRING(11) NOT NULL,        -- ISO 9362
    lei              STRING(20),                  -- ISO 17442
    name             STRING(256) NOT NULL,
    country          STRING(2) NOT NULL,
    cls_member       BOOL NOT NULL DEFAULT false,
    cfets_member     BOOL NOT NULL DEFAULT false,
    status           STRING(32) NOT NULL DEFAULT 'ACTIVE',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    UNIQUE (tenant_id, bic),
    INDEX idx_cp_lei (lei),
    INDEX idx_cp_status (tenant_id, status)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- fx_trades — core trade entity
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS fx_trades (
    trade_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id             UUID NOT NULL REFERENCES tenants(tenant_id),
    external_ref          STRING(128),
    trade_type            STRING(16) NOT NULL,    -- SPOT | FORWARD | NDF | SWAP
    status                STRING(16) NOT NULL DEFAULT 'PENDING',
    settlement_venue      STRING(16) NOT NULL,    -- CLS | BILATERAL | CFETS

    buyer_counterparty_id UUID NOT NULL REFERENCES counterparties(counterparty_id),
    seller_counterparty_id UUID NOT NULL REFERENCES counterparties(counterparty_id),

    bought_currency       STRING(3)  NOT NULL,
    bought_amount         DECIMAL(36,18) NOT NULL CHECK (bought_amount > 0),
    sold_currency         STRING(3)  NOT NULL,
    sold_amount           DECIMAL(36,18) NOT NULL CHECK (sold_amount > 0),
    deal_rate             DECIMAL(36,18) NOT NULL CHECK (deal_rate > 0),

    trade_date            TIMESTAMPTZ NOT NULL,
    value_date            DATE NOT NULL,
    confirmed_at          TIMESTAMPTZ,
    settled_at            TIMESTAMPTZ,

    iso20022_message_id   STRING(64),
    cls_cycle_id          UUID,

    created_at            TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),

    CHECK (bought_currency <> sold_currency),
    CHECK (status IN ('PENDING','CONFIRMED','SETTLING','SETTLED','CANCELLED','REJECTED')),
    CHECK (trade_type IN ('SPOT','FORWARD','NDF','SWAP')),
    CHECK (settlement_venue IN ('CLS','BILATERAL','CFETS')),

    INDEX idx_trades_tenant_status_value (tenant_id, status, value_date),
    INDEX idx_trades_venue_cycle (settlement_venue, cls_cycle_id) STORING (status),
    INDEX idx_trades_buyer (buyer_counterparty_id),
    INDEX idx_trades_seller (seller_counterparty_id),
    INDEX idx_trades_value_date (value_date) WHERE status IN ('CONFIRMED','SETTLING'),
    INDEX idx_trades_external_ref (tenant_id, external_ref) WHERE external_ref IS NOT NULL,
    INDEX idx_trades_iso_msg (iso20022_message_id) WHERE iso20022_message_id IS NOT NULL
);

-- ─────────────────────────────────────────────────────────────────────────────
-- trade_amendments — append-only audit trail
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS trade_amendments (
    amendment_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_id       UUID NOT NULL REFERENCES fx_trades(trade_id),
    tenant_id      UUID NOT NULL REFERENCES tenants(tenant_id),
    status         STRING(16) NOT NULL DEFAULT 'PROPOSED',
    change_type    STRING(32) NOT NULL,
    before_payload JSONB NOT NULL,
    after_payload  JSONB NOT NULL,
    proposer_actor_id UUID REFERENCES actors(actor_id),
    approver_actor_id UUID REFERENCES actors(actor_id),
    proposed_at    TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    approved_at    TIMESTAMPTZ,
    applied_at     TIMESTAMPTZ,
    CHECK (status IN ('PROPOSED','APPROVED','REJECTED','APPLIED')),
    INDEX idx_amend_trade (trade_id, proposed_at DESC),
    INDEX idx_amend_status (tenant_id, status)
);

COMMIT;
