-- ExchangeOS — 000008: compliance + admin tables.

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- classifications — 95-code BACEN nature classification per trade
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS classifications (
    classification_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id          UUID NOT NULL REFERENCES tenants(tenant_id),
    trade_id           UUID NOT NULL REFERENCES fx_trades(trade_id) ON DELETE CASCADE,
    code               STRING(6) NOT NULL,
    description        STRING(256) NOT NULL,
    nature             STRING(16) NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    CHECK (nature IN ('REMESSA','INGRESSO','CONVERSAO')),
    CHECK (length(code) BETWEEN 4 AND 6),
    INDEX idx_class_trade (trade_id),
    INDEX idx_class_tenant_code (tenant_id, code)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- iof_computations — tax per trade per Decreto 12.499/2025
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS iof_computations (
    iof_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(tenant_id),
    trade_id        UUID NOT NULL REFERENCES fx_trades(trade_id) ON DELETE CASCADE,
    operation_type  STRING(32) NOT NULL,
    notional        DECIMAL(36,18) NOT NULL CHECK (notional > 0),
    notional_ccy    STRING(3) NOT NULL,
    rate            DECIMAL(36,18) NOT NULL CHECK (rate >= 0 AND rate <= 1),
    iof_amount      DECIMAL(36,18) NOT NULL CHECK (iof_amount >= 0),
    computed_at     TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    INDEX idx_iof_trade (trade_id),
    INDEX idx_iof_op_type (tenant_id, operation_type, computed_at DESC)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- bacen_reports — submission tracker
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS bacen_reports (
    report_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL REFERENCES tenants(tenant_id),
    report_type      STRING(16) NOT NULL,
    reference_date   DATE NOT NULL,
    payload_hash     STRING(128) NOT NULL,
    status           STRING(16) NOT NULL DEFAULT 'PENDING',
    submitted_at     TIMESTAMPTZ,
    responded_at     TIMESTAMPTZ,
    rejection_reason STRING(512),
    version          INT NOT NULL DEFAULT 1,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    CHECK (report_type IN ('SISBACEN','BCB-CCS','BCB-CAMBIO')),
    CHECK (status IN ('PENDING','SUBMITTED','ACCEPTED','REJECTED')),
    INDEX idx_bacen_tenant_status (tenant_id, status, reference_date DESC),
    INDEX idx_bacen_hash (payload_hash)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- screening_results — OFAC/UN/EU/COAF outcomes per counterparty per screen run
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS screening_results (
    screening_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL REFERENCES tenants(tenant_id),
    counterparty_bic STRING(11) NOT NULL,
    lei              STRING(20),
    risk_level       STRING(8) NOT NULL,
    hits             STRING[] NOT NULL DEFAULT '{}',
    screened_at      TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    CHECK (risk_level IN ('LOW','MEDIUM','HIGH')),
    INDEX idx_screen_cp (tenant_id, counterparty_bic, screened_at DESC),
    INDEX idx_screen_high (tenant_id, screened_at DESC) WHERE risk_level = 'HIGH'
);

-- ─────────────────────────────────────────────────────────────────────────────
-- system_events — admi.x correlated operational events
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS system_events (
    event_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code             STRING(32) NOT NULL,
    component        STRING(128) NOT NULL,
    description      STRING(512),
    at               TIMESTAMPTZ NOT NULL,
    iso20022_ref     STRING(128),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    INDEX idx_sysev_time (at DESC),
    INDEX idx_sysev_code (code, at DESC)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- eod_jobs — end-of-day batch orchestration tracker
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS eod_jobs (
    job_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(tenant_id),
    business_date   DATE NOT NULL,
    status          STRING(16) NOT NULL DEFAULT 'PENDING',
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    failure_reason  STRING(512),
    steps_done      STRING[] NOT NULL DEFAULT '{}',
    version         INT NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    UNIQUE (tenant_id, business_date),
    CHECK (status IN ('PENDING','RUNNING','COMPLETED','FAILED')),
    INDEX idx_eod_status (tenant_id, status, business_date DESC)
);

COMMIT;
