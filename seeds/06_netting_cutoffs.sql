-- ExchangeOS seeds — 06: netting_cutoffs (CLS PayIn deadline bands per currency)
--
-- CLS daily cycle (Europe/Zurich = CET):
--   07:00  Cycle opens
--   PIN1   08:00  First PayIn deadline (Asia-Pacific currencies)
--   PIN2   09:00  Second PayIn deadline (European currencies)
--   PIN3   10:00  Third PayIn deadline (Americas currencies)
--   12:00  Cycle closes
--
-- CCY-to-band mapping reflects market convention (Asia → PIN1, Europe → PIN2, Americas → PIN3).
-- Bilateral cutoffs are illustrative; production sources from counterparty SLAs.

BEGIN;

-- ── CLS PayIn deadlines per currency ────────────────────────────────────────
INSERT INTO netting_cutoffs (venue, currency, band, cutoff_time_cet, description, active) VALUES
    -- PIN1 — Asia-Pacific
    ('CLS', 'JPY', 'PIN1', '08:00:00', 'CLS PayIn 1 — Asia-Pacific (JPY)',         true),
    ('CLS', 'AUD', 'PIN1', '08:00:00', 'CLS PayIn 1 — Asia-Pacific (AUD)',         true),
    ('CLS', 'NZD', 'PIN1', '08:00:00', 'CLS PayIn 1 — Asia-Pacific (NZD)',         true),
    ('CLS', 'HKD', 'PIN1', '08:00:00', 'CLS PayIn 1 — Asia-Pacific (HKD)',         true),
    ('CLS', 'SGD', 'PIN1', '08:00:00', 'CLS PayIn 1 — Asia-Pacific (SGD)',         true),
    ('CLS', 'KRW', 'PIN1', '08:00:00', 'CLS PayIn 1 — Asia-Pacific (KRW)',         true),

    -- PIN2 — Europe/Africa
    ('CLS', 'EUR', 'PIN2', '09:00:00', 'CLS PayIn 2 — Europe (EUR)',               true),
    ('CLS', 'GBP', 'PIN2', '09:00:00', 'CLS PayIn 2 — Europe (GBP)',               true),
    ('CLS', 'CHF', 'PIN2', '09:00:00', 'CLS PayIn 2 — Europe (CHF)',               true),
    ('CLS', 'NOK', 'PIN2', '09:00:00', 'CLS PayIn 2 — Europe (NOK)',               true),
    ('CLS', 'SEK', 'PIN2', '09:00:00', 'CLS PayIn 2 — Europe (SEK)',               true),
    ('CLS', 'DKK', 'PIN2', '09:00:00', 'CLS PayIn 2 — Europe (DKK)',               true),
    ('CLS', 'HUF', 'PIN2', '09:00:00', 'CLS PayIn 2 — Europe (HUF)',               true),
    ('CLS', 'ILS', 'PIN2', '09:00:00', 'CLS PayIn 2 — Europe (ILS)',               true),
    ('CLS', 'ZAR', 'PIN2', '09:00:00', 'CLS PayIn 2 — Africa (ZAR)',               true),

    -- PIN3 — Americas
    ('CLS', 'USD', 'PIN3', '10:00:00', 'CLS PayIn 3 — Americas (USD)',             true),
    ('CLS', 'CAD', 'PIN3', '10:00:00', 'CLS PayIn 3 — Americas (CAD)',             true),
    ('CLS', 'MXN', 'PIN3', '10:00:00', 'CLS PayIn 3 — Americas (MXN)',             true),

    -- ── Bilateral default cutoffs (per CCY EOD on local exchange) ───────────
    ('BILATERAL', 'USD', 'EOD', '22:00:00', 'NY EOD cutoff',                       true),
    ('BILATERAL', 'BRL', 'EOD', '21:00:00', 'B3 EOD cutoff (BRT 17:00 = 21:00 CET)', true),
    ('BILATERAL', 'EUR', 'EOD', '17:00:00', 'TARGET2 EOD cutoff (CET 17:00)',      true),
    ('BILATERAL', 'GBP', 'EOD', '17:00:00', 'CHAPS EOD cutoff (UK 16:00 = 17:00 CET)', true),
    ('BILATERAL', 'JPY', 'EOD', '08:00:00', 'BOJ-NET EOD cutoff (JST 16:00 = next day 08:00 CET)', true),

    -- ── CFETS (China onshore CNY) ──────────────────────────────────────────
    ('CFETS', 'CNY', 'EOD', '11:00:00', 'CFETS EOD cutoff (Shanghai 17:00 = 11:00 CET in summer)', true)
ON CONFLICT (venue, currency, band) DO NOTHING;

COMMIT;
