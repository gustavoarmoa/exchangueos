-- ExchangeOS seeds — 02: currency_pairs
-- 30+ major pairs:
--   - 17 G10 vs USD (CLS-eligible both sides)
--   - 5 high-volume CLS crosses (EURGBP, EURCHF, EURJPY, GBPJPY, AUDJPY)
--   - 1 CNY pair (USDCNY — CFETS eligible)
--   - 6 EM/NDF pairs (USDBRL, USDINR, USDIDR, USDKRW, USDMXN, USDZAR)
--
-- spot_days:
--   - USDCAD = 1 (industry convention)
--   - USDMXN = 1 (industry convention)
--   - USDTRY = 1 (historical)
--   - all others = 2 (FX spot default)

BEGIN;

INSERT INTO currency_pairs (base_ccy, quote_ccy, spot_days, cls_eligible, cfets_eligible, active, description) VALUES
    -- ── G10 vs USD ──────────────────────────────────────────────────────────
    ('EUR', 'USD', 2, true,  false, true, 'Euro / US Dollar (most-traded G10 pair)'),
    ('GBP', 'USD', 2, true,  false, true, 'Sterling / US Dollar (Cable)'),
    ('USD', 'JPY', 2, true,  false, true, 'US Dollar / Japanese Yen'),
    ('USD', 'CHF', 2, true,  false, true, 'US Dollar / Swiss Franc'),
    ('AUD', 'USD', 2, true,  false, true, 'Australian Dollar / US Dollar (Aussie)'),
    ('NZD', 'USD', 2, true,  false, true, 'New Zealand Dollar / US Dollar (Kiwi)'),
    ('USD', 'CAD', 1, true,  false, true, 'US Dollar / Canadian Dollar (T+1)'),
    ('USD', 'NOK', 2, true,  false, true, 'US Dollar / Norwegian Krone'),
    ('USD', 'SEK', 2, true,  false, true, 'US Dollar / Swedish Krona'),
    ('USD', 'DKK', 2, true,  false, true, 'US Dollar / Danish Krone'),
    ('USD', 'SGD', 2, true,  false, true, 'US Dollar / Singapore Dollar'),
    ('USD', 'HKD', 2, true,  false, true, 'US Dollar / Hong Kong Dollar'),
    ('USD', 'ZAR', 2, true,  false, true, 'US Dollar / South African Rand'),
    ('USD', 'MXN', 1, true,  false, true, 'US Dollar / Mexican Peso (T+1)'),
    ('USD', 'ILS', 2, true,  false, true, 'US Dollar / Israeli New Shekel'),
    ('USD', 'HUF', 2, true,  false, true, 'US Dollar / Hungarian Forint'),
    ('USD', 'KRW', 2, true,  false, true, 'US Dollar / South Korean Won'),

    -- ── G10 crosses (high CLS volume) ───────────────────────────────────────
    ('EUR', 'GBP', 2, true,  false, true, 'Euro / Sterling'),
    ('EUR', 'CHF', 2, true,  false, true, 'Euro / Swiss Franc'),
    ('EUR', 'JPY', 2, true,  false, true, 'Euro / Japanese Yen'),
    ('GBP', 'JPY', 2, true,  false, true, 'Sterling / Japanese Yen'),
    ('AUD', 'JPY', 2, true,  false, true, 'Australian Dollar / Japanese Yen'),
    ('EUR', 'NOK', 2, true,  false, true, 'Euro / Norwegian Krone'),
    ('EUR', 'SEK', 2, true,  false, true, 'Euro / Swedish Krona'),

    -- ── China (CFETS) ───────────────────────────────────────────────────────
    ('USD', 'CNY', 2, false, true,  true, 'US Dollar / Chinese Yuan onshore (CFETS)'),
    ('USD', 'CNH', 2, false, false, true, 'US Dollar / Chinese Yuan offshore (HK)'),

    -- ── Emerging market / NDF universe ──────────────────────────────────────
    ('USD', 'BRL', 2, false, false, true, 'US Dollar / Brazilian Real (NDF)'),
    ('USD', 'INR', 2, false, false, true, 'US Dollar / Indian Rupee (NDF)'),
    ('USD', 'IDR', 2, false, false, true, 'US Dollar / Indonesian Rupiah (NDF)'),
    ('USD', 'PHP', 2, false, false, true, 'US Dollar / Philippine Peso (NDF)'),
    ('USD', 'TWD', 2, false, false, true, 'US Dollar / New Taiwan Dollar (NDF)'),
    ('USD', 'MYR', 2, false, false, true, 'US Dollar / Malaysian Ringgit (NDF)'),

    -- ── BRL crosses (LatAm institutional flow) ──────────────────────────────
    ('EUR', 'BRL', 2, false, false, true, 'Euro / Brazilian Real (cross via USD)')
ON CONFLICT (base_ccy, quote_ccy) DO NOTHING;

COMMIT;
