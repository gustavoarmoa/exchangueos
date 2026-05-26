-- ExchangeOS seeds — 05: ssi (sample standing settlement instructions for the dev tenant)
--
-- This seed requires a dev tenant to exist. Run via:
--   psql ... -v dev_tenant_id='<UUID>' -f seeds/05_ssi.sql
--
-- For local docker-compose stack we use a deterministic UUIDv5 derived from a fixed
-- namespace + 'dev-tenant' — the migrator creates it as part of 01_tenants_dev.sql
-- (added in a follow-up; SSI rows here are a no-op if the dev tenant is absent).

BEGIN;

-- Dev tenant ID (deterministic) — must match the one inserted in seeds/00_tenants_dev.sql.
-- If not present, this block is a no-op (INSERT ... SELECT ... WHERE EXISTS).
WITH dev AS (
    SELECT tenant_id FROM tenants WHERE code = 'DEV' LIMIT 1
)
INSERT INTO ssis (tenant_id, counterparty_bic, currency, beneficiary_bic, account_number, iban, valid_from, valid_to) VALUES
    -- USD via CHASUS33 (JPMorgan)
    ((SELECT tenant_id FROM dev), 'CHASUS33', 'USD', 'CHASUS33', '1234567890', NULL,
     TIMESTAMPTZ '2026-01-01 00:00:00+00', NULL),
    -- EUR via DEUTDEFF (Deutsche Bank)
    ((SELECT tenant_id FROM dev), 'DEUTDEFF', 'EUR', 'DEUTDEFF', NULL, 'DE89370400440532013000',
     TIMESTAMPTZ '2026-01-01 00:00:00+00', NULL),
    -- GBP via HSBCGB2L (HSBC)
    ((SELECT tenant_id FROM dev), 'HSBCGB2L', 'GBP', 'HSBCGB2L', NULL, 'GB29NWBK60161331926819',
     TIMESTAMPTZ '2026-01-01 00:00:00+00', NULL),
    -- BRL via ITAUBRSP (Itaú)
    ((SELECT tenant_id FROM dev), 'ITAUBRSP', 'BRL', 'ITAUBRSP', '341-9-12345-6', NULL,
     TIMESTAMPTZ '2026-01-01 00:00:00+00', NULL),
    -- JPY via SMBCJPJT
    ((SELECT tenant_id FROM dev), 'SMBCJPJT', 'JPY', 'SMBCJPJT', '7777777', NULL,
     TIMESTAMPTZ '2026-01-01 00:00:00+00', NULL)
ON CONFLICT DO NOTHING;

COMMIT;
