-- ExchangeOS seeds — 00: tenants (idempotent)
--
-- DEV is the deterministic UUIDv5-shaped primary tenant. Other 4 are sample
-- multi-tenant fixtures so admin-query endpoints have ≥ 5 rows + we can prove
-- tenant isolation by querying with a non-DEV tenant_id.

BEGIN;

INSERT INTO tenants (tenant_id, code, name, country, status, metadata) VALUES
    ('00000000-0000-5000-8000-000000000001', 'DEV',     'ExchangeOS Dev Tenant',          'BR', 'ACTIVE',    '{"primary":true}'::JSONB),
    ('00000000-0000-5000-8000-000000000002', 'BANK-BR-01', 'Sample Brazilian Bank #1',    'BR', 'ACTIVE',    '{}'::JSONB),
    ('00000000-0000-5000-8000-000000000003', 'BANK-US-01', 'Sample US Investment Bank',   'US', 'ACTIVE',    '{}'::JSONB),
    ('00000000-0000-5000-8000-000000000004', 'BANK-EU-01', 'Sample EU Universal Bank',    'DE', 'ACTIVE',    '{}'::JSONB),
    ('00000000-0000-5000-8000-000000000005', 'SANDBOX',    'Sandbox / Integration Tests', 'BR', 'SUSPENDED', '{"purpose":"sandbox"}'::JSONB)
ON CONFLICT (code) DO NOTHING;

COMMIT;
