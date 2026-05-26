-- ExchangeOS seeds — 03: calendars (BACEN BRL, NYFR USD, BOE GBP, TARGET2 EUR) — 2026
-- Holidays are illustrative for testing only. PRODUCTION must source from authoritative
-- feeds (BACEN ENDPOINT, Federal Reserve calendar, BOE calendar, TARGET2 calendar) at
-- each year-end refresh.

BEGIN;

-- ── Calendar definitions ────────────────────────────────────────────────────
INSERT INTO calendars (calendar_id, description) VALUES
    ('BACEN_BRL',  'Banco Central do Brasil — BRL settlement calendar'),
    ('NYFR_USD',   'New York Federal Reserve — USD settlement calendar'),
    ('BOE_GBP',    'Bank of England — GBP settlement calendar'),
    ('TARGET2_EUR','TARGET2 (Eurosystem) — EUR settlement calendar'),
    ('TOKYO_JPY',  'Bank of Japan — JPY settlement calendar'),
    ('TORONTO_CAD','Bank of Canada — CAD settlement calendar')
ON CONFLICT (calendar_id) DO NOTHING;

-- ── BACEN_BRL — 2026 fixed holidays ─────────────────────────────────────────
INSERT INTO calendar_holidays (calendar_id, holiday_date, description) VALUES
    ('BACEN_BRL','2026-01-01','Confraternização Universal'),
    ('BACEN_BRL','2026-02-16','Carnaval segunda'),
    ('BACEN_BRL','2026-02-17','Carnaval terça'),
    ('BACEN_BRL','2026-04-03','Sexta-feira Santa'),
    ('BACEN_BRL','2026-04-21','Tiradentes'),
    ('BACEN_BRL','2026-05-01','Dia do Trabalho'),
    ('BACEN_BRL','2026-06-04','Corpus Christi'),
    ('BACEN_BRL','2026-09-07','Independência'),
    ('BACEN_BRL','2026-10-12','N. Sra. Aparecida'),
    ('BACEN_BRL','2026-11-02','Finados'),
    ('BACEN_BRL','2026-11-15','Proclamação da República'),
    ('BACEN_BRL','2026-12-25','Natal')
ON CONFLICT (calendar_id, holiday_date) DO NOTHING;

-- ── NYFR_USD — 2026 Federal holidays (observed) ─────────────────────────────
INSERT INTO calendar_holidays (calendar_id, holiday_date, description) VALUES
    ('NYFR_USD','2026-01-01','New Year''s Day'),
    ('NYFR_USD','2026-01-19','Martin Luther King Jr. Day'),
    ('NYFR_USD','2026-02-16','Presidents'' Day'),
    ('NYFR_USD','2026-05-25','Memorial Day'),
    ('NYFR_USD','2026-06-19','Juneteenth'),
    ('NYFR_USD','2026-07-03','Independence Day (observed — July 4 is Saturday)'),
    ('NYFR_USD','2026-09-07','Labor Day'),
    ('NYFR_USD','2026-10-12','Columbus Day'),
    ('NYFR_USD','2026-11-11','Veterans Day'),
    ('NYFR_USD','2026-11-26','Thanksgiving'),
    ('NYFR_USD','2026-12-25','Christmas')
ON CONFLICT (calendar_id, holiday_date) DO NOTHING;

-- ── BOE_GBP — 2026 UK bank holidays ─────────────────────────────────────────
INSERT INTO calendar_holidays (calendar_id, holiday_date, description) VALUES
    ('BOE_GBP','2026-01-01','New Year''s Day'),
    ('BOE_GBP','2026-04-03','Good Friday'),
    ('BOE_GBP','2026-04-06','Easter Monday'),
    ('BOE_GBP','2026-05-04','Early May bank holiday'),
    ('BOE_GBP','2026-05-25','Spring bank holiday'),
    ('BOE_GBP','2026-08-31','Summer bank holiday'),
    ('BOE_GBP','2026-12-25','Christmas'),
    ('BOE_GBP','2026-12-28','Boxing Day (substitute — Dec 26 is Sat)')
ON CONFLICT (calendar_id, holiday_date) DO NOTHING;

-- ── TARGET2_EUR — 2026 closing days ─────────────────────────────────────────
INSERT INTO calendar_holidays (calendar_id, holiday_date, description) VALUES
    ('TARGET2_EUR','2026-01-01','New Year''s Day'),
    ('TARGET2_EUR','2026-04-03','Good Friday'),
    ('TARGET2_EUR','2026-04-06','Easter Monday'),
    ('TARGET2_EUR','2026-05-01','Labour Day'),
    ('TARGET2_EUR','2026-12-25','Christmas'),
    ('TARGET2_EUR','2026-12-26','Boxing Day')
ON CONFLICT (calendar_id, holiday_date) DO NOTHING;

-- ── TOKYO_JPY — 2026 banking holidays (subset) ──────────────────────────────
INSERT INTO calendar_holidays (calendar_id, holiday_date, description) VALUES
    ('TOKYO_JPY','2026-01-01','New Year''s Day'),
    ('TOKYO_JPY','2026-01-02','Bank holiday'),
    ('TOKYO_JPY','2026-01-12','Coming of Age Day'),
    ('TOKYO_JPY','2026-02-11','Foundation Day'),
    ('TOKYO_JPY','2026-04-29','Showa Day'),
    ('TOKYO_JPY','2026-05-04','Greenery Day'),
    ('TOKYO_JPY','2026-05-05','Children''s Day'),
    ('TOKYO_JPY','2026-12-31','Bank holiday')
ON CONFLICT (calendar_id, holiday_date) DO NOTHING;

-- ── TORONTO_CAD — 2026 statutory holidays (subset) ──────────────────────────
INSERT INTO calendar_holidays (calendar_id, holiday_date, description) VALUES
    ('TORONTO_CAD','2026-01-01','New Year''s Day'),
    ('TORONTO_CAD','2026-02-16','Family Day'),
    ('TORONTO_CAD','2026-04-03','Good Friday'),
    ('TORONTO_CAD','2026-05-18','Victoria Day'),
    ('TORONTO_CAD','2026-07-01','Canada Day'),
    ('TORONTO_CAD','2026-09-07','Labour Day'),
    ('TORONTO_CAD','2026-10-12','Thanksgiving'),
    ('TORONTO_CAD','2026-12-25','Christmas')
ON CONFLICT (calendar_id, holiday_date) DO NOTHING;

COMMIT;
