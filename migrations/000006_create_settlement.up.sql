-- ExchangeOS — 000006: cls_cycles + payin_instructions + net_reports
-- Models the CLS daily PvP cycle (07:00 open / 08-09-10 PIN deadlines / 12:00 close).

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- cls_cycles — one row per (tenant, business_date)
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS cls_cycles (
    cycle_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES tenants(tenant_id),
    cycle_date        DATE NOT NULL,
    status            STRING(20) NOT NULL DEFAULT 'OPEN',
    opened_at         TIMESTAMPTZ NOT NULL,
    pin1_deadline     TIMESTAMPTZ NOT NULL,
    pin2_deadline     TIMESTAMPTZ NOT NULL,
    pin3_deadline     TIMESTAMPTZ NOT NULL,
    scheduled_close   TIMESTAMPTZ NOT NULL,
    closed_at         TIMESTAMPTZ,
    failure_reason    STRING(256),
    version           INT NOT NULL DEFAULT 1,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),

    CHECK (status IN ('OPEN','PAY_IN_WINDOW','SETTLING','CLOSED','FAILED')),
    CHECK (pin1_deadline < pin2_deadline AND pin2_deadline < pin3_deadline AND pin3_deadline < scheduled_close),

    UNIQUE (tenant_id, cycle_date),
    INDEX idx_cycles_status (tenant_id, status, cycle_date DESC),
    INDEX idx_cycles_open (tenant_id, cycle_date) WHERE status IN ('OPEN','PAY_IN_WINDOW','SETTLING')
);

-- ─────────────────────────────────────────────────────────────────────────────
-- cls_cycle_trades — many-to-many: which trades belong to which cycle
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS cls_cycle_trades (
    cycle_id    UUID NOT NULL REFERENCES cls_cycles(cycle_id) ON DELETE CASCADE,
    trade_id    UUID NOT NULL REFERENCES fx_trades(trade_id)  ON DELETE RESTRICT,
    attached_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    PRIMARY KEY (cycle_id, trade_id),
    INDEX idx_cct_trade (trade_id)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- payin_instructions
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS payin_instructions (
    instruction_id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID NOT NULL REFERENCES tenants(tenant_id),
    cycle_id        UUID NOT NULL REFERENCES cls_cycles(cycle_id) ON DELETE RESTRICT,
    currency        STRING(3) NOT NULL,
    amount          DECIMAL(36,18) NOT NULL CHECK (amount > 0),
    band            STRING(4) NOT NULL,
    deadline        TIMESTAMPTZ NOT NULL,
    status          STRING(16) NOT NULL DEFAULT 'PENDING',
    submitted_at    TIMESTAMPTZ,
    confirmed_at    TIMESTAMPTZ,
    failure_reason  STRING(256),
    version         INT NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),

    CHECK (band IN ('PIN1','PIN2','PIN3')),
    CHECK (status IN ('PENDING','SUBMITTED','CONFIRMED','FAILED')),

    INDEX idx_payin_cycle_ccy (cycle_id, currency),
    INDEX idx_payin_status_deadline (tenant_id, status, deadline)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- net_reports — one per (cycle, currency)
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS net_reports (
    report_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        UUID NOT NULL REFERENCES tenants(tenant_id),
    cycle_id         UUID NOT NULL REFERENCES cls_cycles(cycle_id) ON DELETE CASCADE,
    currency         STRING(3) NOT NULL,
    gross_pay_in     DECIMAL(36,18) NOT NULL CHECK (gross_pay_in  >= 0),
    gross_pay_out    DECIMAL(36,18) NOT NULL CHECK (gross_pay_out >= 0),
    net_settlement   DECIMAL(36,18) NOT NULL,  -- signed; positive = receivable
    trade_count      INT NOT NULL CHECK (trade_count >= 0),
    generated_at     TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),

    UNIQUE (cycle_id, currency),
    INDEX idx_netrep_tenant_cycle (tenant_id, cycle_id)
);

COMMIT;
