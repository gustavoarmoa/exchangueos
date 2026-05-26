-- ExchangeOS — 000007: risk_limits + positions
-- Money values DECIMAL(36,18); per-tenant scope.

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- risk_limits — counterparty / currency / tenor / dv01 / var caps per tenant
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS risk_limits (
    limit_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL REFERENCES tenants(tenant_id),
    limit_type   STRING(16) NOT NULL,
    scope        STRING(32) NOT NULL DEFAULT '',
    cap          DECIMAL(36,18) NOT NULL CHECK (cap > 0),
    utilised     DECIMAL(36,18) NOT NULL DEFAULT 0 CHECK (utilised >= 0),
    currency     STRING(3) NOT NULL,
    version      INT NOT NULL DEFAULT 1,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),

    CHECK (limit_type IN ('COUNTERPARTY','CURRENCY','TENOR','DV01','VAR')),
    UNIQUE (tenant_id, limit_type, scope),
    INDEX idx_limits_tenant_type (tenant_id, limit_type),
    INDEX idx_limits_breaching (tenant_id) WHERE utilised >= cap
);

-- ─────────────────────────────────────────────────────────────────────────────
-- positions — per-tenant per-currency net open position
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS positions (
    position_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    UUID NOT NULL REFERENCES tenants(tenant_id),
    currency     STRING(3) NOT NULL,
    long_amount  DECIMAL(36,18) NOT NULL DEFAULT 0 CHECK (long_amount  >= 0),
    short_amount DECIMAL(36,18) NOT NULL DEFAULT 0 CHECK (short_amount >= 0),
    net_amount   DECIMAL(36,18) NOT NULL DEFAULT 0,  -- signed; long_amount − short_amount
    as_of        TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    version      INT NOT NULL DEFAULT 1,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),

    UNIQUE (tenant_id, currency),
    INDEX idx_positions_tenant_asof (tenant_id, as_of DESC)
);

COMMIT;
