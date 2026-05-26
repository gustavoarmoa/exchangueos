-- ExchangeOS — 000004: currency_pairs + netting_cutoffs
-- Depends on currencies (created in 000005 conceptually; the migrator runs in
-- file order, so this file is ordered before 000005). We use deferred FK
-- semantics by NOT declaring the FK here — currencies table will hold the
-- canonical codes and applications validate.
--
-- Actually we DO declare the FK because in CockroachDB the schema_migrations
-- pin is monotonic; the operator runs `migrator up` once and the validator
-- requires referenced tables exist. To keep this file independent, we use a
-- LOOSE check constraint instead of FK and leave a comment for future
-- tightening.

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- currency_pairs — quoted FX pairs (e.g. EUR/USD)
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS currency_pairs (
    base_ccy        STRING(3) NOT NULL,
    quote_ccy       STRING(3) NOT NULL,
    spot_days       INT NOT NULL DEFAULT 2 CHECK (spot_days IN (0, 1, 2)),  -- T+0/T+1/T+2
    cls_eligible    BOOL NOT NULL DEFAULT false,
    cfets_eligible  BOOL NOT NULL DEFAULT false,
    active          BOOL NOT NULL DEFAULT true,
    description     STRING(256),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    PRIMARY KEY (base_ccy, quote_ccy),
    CHECK (base_ccy <> quote_ccy),
    INDEX idx_pairs_active (active),
    INDEX idx_pairs_cls (cls_eligible) WHERE cls_eligible = true
);

-- ─────────────────────────────────────────────────────────────────────────────
-- netting_cutoffs — CLS PayIn deadlines per currency
-- Times stored as TIME (HH:MM:SS) in CET — caller converts to local zone.
-- One row per (venue, currency, band).
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS netting_cutoffs (
    venue           STRING(16) NOT NULL,         -- e.g. 'CLS','BILATERAL','CFETS'
    currency        STRING(3) NOT NULL,
    band            STRING(8) NOT NULL,          -- e.g. 'PIN1','PIN2','PIN3'
    cutoff_time_cet TIME NOT NULL,
    description     STRING(256),
    active          BOOL NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT current_timestamp(),
    PRIMARY KEY (venue, currency, band),
    CHECK (band IN ('PIN1','PIN2','PIN3','EOD')),
    INDEX idx_cutoffs_venue_ccy (venue, currency)
);

COMMIT;
