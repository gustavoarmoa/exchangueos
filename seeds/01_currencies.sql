-- ExchangeOS seeds — 01: currencies (ISO 4217 with CLS/CFETS eligibility)
-- 18 CLS-eligible currencies + BRL + 8 emerging-market currencies (NDF universe).
-- Idempotent (ON CONFLICT DO NOTHING).

BEGIN;

INSERT INTO currencies (code, name, minor_units, cls_eligible, cfets_eligible, active) VALUES
    -- ── 18 CLS-eligible (Continuous Linked Settlement) ──────────────────────
    ('AUD', 'Australian Dollar',              2, true,  false, true),
    ('CAD', 'Canadian Dollar',                2, true,  false, true),
    ('CHF', 'Swiss Franc',                    2, true,  false, true),
    ('DKK', 'Danish Krone',                   2, true,  false, true),
    ('EUR', 'Euro',                           2, true,  false, true),
    ('GBP', 'Pound Sterling',                 2, true,  false, true),
    ('HKD', 'Hong Kong Dollar',               2, true,  false, true),
    ('HUF', 'Hungarian Forint',               2, true,  false, true),
    ('ILS', 'Israeli New Shekel',             2, true,  false, true),
    ('JPY', 'Japanese Yen',                   0, true,  false, true),
    ('KRW', 'South Korean Won',               0, true,  false, true),
    ('MXN', 'Mexican Peso',                   2, true,  false, true),
    ('NOK', 'Norwegian Krone',                2, true,  false, true),
    ('NZD', 'New Zealand Dollar',             2, true,  false, true),
    ('SEK', 'Swedish Krona',                  2, true,  false, true),
    ('SGD', 'Singapore Dollar',               2, true,  false, true),
    ('USD', 'United States Dollar',           2, true,  false, true),
    ('ZAR', 'South African Rand',             2, true,  false, true),
    -- ── Brazilian Real (NDF universe — not CLS-eligible) ────────────────────
    ('BRL', 'Brazilian Real',                 2, false, false, true),
    -- ── CFETS / China onshore + offshore ────────────────────────────────────
    ('CNY', 'Chinese Yuan Renminbi (onshore)', 2, false, true,  true),
    ('CNH', 'Chinese Yuan Renminbi (offshore)', 2, false, false, true),
    -- ── Other emerging market (typically NDF) ───────────────────────────────
    ('INR', 'Indian Rupee',                   2, false, false, true),
    ('IDR', 'Indonesian Rupiah',              2, false, false, true),
    ('MYR', 'Malaysian Ringgit',              2, false, false, true),
    ('PHP', 'Philippine Peso',                2, false, false, true),
    ('RUB', 'Russian Ruble',                  2, false, false, true),
    ('THB', 'Thai Baht',                      2, false, false, true),
    ('TWD', 'New Taiwan Dollar',              2, false, false, true),
    -- ── Middle East / 3-decimal currencies ──────────────────────────────────
    ('BHD', 'Bahraini Dinar',                 3, false, false, true),
    ('KWD', 'Kuwaiti Dinar',                  3, false, false, true),
    ('OMR', 'Omani Rial',                     3, false, false, true)
ON CONFLICT (code) DO NOTHING;

COMMIT;
