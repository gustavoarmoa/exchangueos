-- ExchangeOS — 000005: refdata — currencies + calendars + bic_records + ssis
-- Refdata tables are GLOBAL (no tenant_id) EXCEPT ssis which are per-tenant.

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- currencies — ISO 4217 catalog + CLS/CFETS eligibility flags
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS currencies (
    code            STRING(3) PRIMARY KEY,
    name            STRING(256) NOT NULL,
    minor_units     INT NOT NULL CHECK (minor_units IN (0, 2, 3)),
    cls_eligible    BOOL NOT NULL DEFAULT false,
    cfets_eligible  BOOL NOT NULL DEFAULT false,
    active          BOOL NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    INDEX idx_currencies_active (active)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- calendars — holiday set per venue (e.g. USD_NYC, EUR_TARGET2, BRL_BRASILIA)
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS calendars (
    calendar_id  STRING(64) PRIMARY KEY,
    description  STRING(256),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT current_timestamp()
);

CREATE TABLE IF NOT EXISTS calendar_holidays (
    calendar_id  STRING(64) NOT NULL REFERENCES calendars(calendar_id) ON DELETE CASCADE,
    holiday_date DATE NOT NULL,
    description  STRING(256),
    PRIMARY KEY (calendar_id, holiday_date)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- bic_records — ISO 9362 BIC catalog
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS bic_records (
    bic              STRING(11) PRIMARY KEY,
    institution_name STRING(256) NOT NULL,
    country          STRING(2) NOT NULL,
    lei              STRING(20),
    active           BOOL NOT NULL DEFAULT true,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    CHECK (length(bic) IN (8, 11)),
    INDEX idx_bic_country (country),
    INDEX idx_bic_lei (lei) WHERE lei IS NOT NULL,
    INDEX idx_bic_active (active)
);

-- ─────────────────────────────────────────────────────────────────────────────
-- ssis — Standing Settlement Instructions per (tenant, counterparty, currency)
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS ssis (
    ssi_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL REFERENCES tenants(tenant_id),
    counterparty_bic  STRING(11) NOT NULL REFERENCES bic_records(bic),
    currency          STRING(3) NOT NULL REFERENCES currencies(code),
    beneficiary_bic   STRING(11) NOT NULL,
    intermediary_bic  STRING(11),
    account_number    STRING(64),
    iban              STRING(34),
    valid_from        TIMESTAMPTZ NOT NULL,
    valid_to          TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),

    CHECK (account_number IS NOT NULL OR iban IS NOT NULL),
    CHECK (iban IS NULL OR (length(iban) BETWEEN 15 AND 34)),
    CHECK (valid_to IS NULL OR valid_to >= valid_from),

    INDEX idx_ssi_tenant_lookup (tenant_id, counterparty_bic, currency, valid_from DESC),
    -- CRDB rejects context-dependent funcs (current_timestamp/now) in INDEX
    -- PREDICATE. We index on (tenant_id, valid_to DESC) and let queries filter
    -- via `WHERE valid_to IS NULL OR valid_to > now()`. Index still covers both.
    INDEX idx_ssi_active (tenant_id, valid_to DESC)
);

COMMIT;
