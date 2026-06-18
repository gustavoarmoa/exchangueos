-- ExchangeOS — 000012: CFETS Trade Capture (fxtr 031-033) + Confirmation (fxtr 034-038)
-- tables. Persistence layer for modules/cfets_capture + modules/cfets_confirmation.

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- cfets_captures — one row per CFETS Trade Capture submission
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS cfets_captures (
    capture_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES tenants(tenant_id),
    trade_id          UUID NOT NULL REFERENCES fx_trades(trade_id) ON DELETE CASCADE,
    submitter_ref     STRING(64) NOT NULL,
    cfets_deal_id     STRING(64),
    status            STRING(16) NOT NULL,
    submitted_at      TIMESTAMPTZ,
    ack_at            TIMESTAMPTZ,
    notified_at       TIMESTAMPTZ,
    rejection_reason  STRING(512),
    version           INT NOT NULL DEFAULT 1,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    CHECK (status IN ('DRAFT','SUBMITTED','ACK','REJECTED','NOTIFIED')),
    UNIQUE (tenant_id, submitter_ref),
    INDEX idx_cfets_cap_trade (trade_id),
    INDEX idx_cfets_cap_status (tenant_id, status, created_at DESC),
    INDEX idx_cfets_cap_deal (cfets_deal_id) WHERE cfets_deal_id IS NOT NULL
);

-- ─────────────────────────────────────────────────────────────────────────────
-- cfets_confirmations — one row per CFETS Trade Confirmation flow
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS cfets_confirmations (
    confirmation_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES tenants(tenant_id),
    trade_id          UUID NOT NULL REFERENCES fx_trades(trade_id) ON DELETE CASCADE,
    cfets_deal_id     STRING(64) NOT NULL,
    status            STRING(16) NOT NULL,
    requested_at      TIMESTAMPTZ NOT NULL,
    confirmed_at      TIMESTAMPTZ,
    rejection_reason  STRING(512),
    version           INT NOT NULL DEFAULT 1,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    CHECK (status IN ('CONFIRMING','CONFIRMED','UNPAIRED','REJECTED')),
    UNIQUE (tenant_id, cfets_deal_id),
    INDEX idx_cfets_conf_trade (trade_id),
    INDEX idx_cfets_conf_status (tenant_id, status, requested_at DESC)
);

COMMIT;
