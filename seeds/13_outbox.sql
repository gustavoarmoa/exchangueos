-- ExchangeOS seeds — 13: outbox_events + outbox_dispatched_archive
-- 5 outbox_events split across pending (will be drained by worker) + dispatched
-- (sample row + archive copy for ARCHIVE table seed).
--
-- WARNING: the 2 dispatched rows ship with dispatched_at IN THE PAST so the
-- worker doesn't double-publish. The 3 pending rows will fire to Kafka once
-- the worker boots (or it just runs through them and they end up archived).

BEGIN;

-- ─── 3 pending outbox_events ───────────────────────────────────────────────
INSERT INTO outbox_events (outbox_id, tenant_id, aggregate_type, aggregate_id, event_name, event_payload, topic, partition_key, occurred_at, dispatched_at, attempt_count) VALUES
    ('16161616-0000-5000-8000-000000000001', '00000000-0000-5000-8000-000000000001',
     'Trade', '33333333-0000-5000-8000-000000000001',
     'trade.created.v1',
     '{"trade_id":"33333333-0000-5000-8000-000000000001","trade_type":"SPOT","pair":"EURUSD","amount":"1000000"}'::JSONB,
     'exchangeos.trade.events', '33333333-0000-5000-8000-000000000001',
     '2026-05-23T14:35:00Z', NULL, 0),

    ('16161616-0000-5000-8000-000000000002', '00000000-0000-5000-8000-000000000001',
     'Quote', '88888888-0000-5000-8000-000000000004',
     'quote.created.v1',
     '{"quote_id":"88888888-0000-5000-8000-000000000004","pair":"EURUSD","bid":"1.0847","ask":"1.0853"}'::JSONB,
     'exchangeos.quote.events', '88888888-0000-5000-8000-000000000004',
     '2026-05-25T11:00:00Z', NULL, 0),

    ('16161616-0000-5000-8000-000000000003', '00000000-0000-5000-8000-000000000001',
     'CLSCycle', '66666666-0000-5000-8000-000000000003',
     'cls.payin_window_opened.v1',
     '{"cycle_id":"66666666-0000-5000-8000-000000000003","band":"PIN1"}'::JSONB,
     'exchangeos.cls_settlement.events', '66666666-0000-5000-8000-000000000003',
     '2026-05-26T06:00:00Z', NULL, 0)
ON CONFLICT (outbox_id) DO NOTHING;

-- ─── 2 dispatched outbox_events (kept for read-API testing) ────────────────
INSERT INTO outbox_events (outbox_id, tenant_id, aggregate_type, aggregate_id, event_name, event_payload, topic, partition_key, occurred_at, dispatched_at, attempt_count) VALUES
    ('16161616-0000-5000-8000-000000000004', '00000000-0000-5000-8000-000000000001',
     'Position', 'dddddddd-0000-5000-8000-000000000001',
     'position.updated.v1',
     '{"currency":"EUR","net_amount":"1000000"}'::JSONB,
     'exchangeos.position.events', 'EUR',
     '2026-05-25T11:00:00Z', '2026-05-25T11:00:05Z', 1),

    ('16161616-0000-5000-8000-000000000005', '00000000-0000-5000-8000-000000000001',
     'RiskLimit', 'cccccccc-0000-5000-8000-000000000001',
     'risk.reserve.v1',
     '{"limit_type":"COUNTERPARTY","scope":"DEUTDEFF","utilised":"3000000"}'::JSONB,
     'exchangeos.risk.events', 'DEUTDEFF',
     '2026-05-25T11:00:00Z', '2026-05-25T11:00:10Z', 1)
ON CONFLICT (outbox_id) DO NOTHING;

-- ─── 5 outbox_dispatched_archive rows ──────────────────────────────────────
INSERT INTO outbox_dispatched_archive (outbox_id, tenant_id, aggregate_type, aggregate_id, event_name, topic, partition_key, occurred_at, dispatched_at, attempt_count) VALUES
    ('17171717-0000-5000-8000-000000000001', '00000000-0000-5000-8000-000000000001', 'Trade',    '33333333-0000-5000-8000-000000000001', 'trade.settled.v1',       'exchangeos.trade.events',     '33333333-0000-5000-8000-000000000001', '2026-05-22T10:00:00Z', '2026-05-22T10:00:05Z', 1),
    ('17171717-0000-5000-8000-000000000002', '00000000-0000-5000-8000-000000000001', 'Trade',    '33333333-0000-5000-8000-000000000004', 'trade.settled.v1',       'exchangeos.trade.events',     '33333333-0000-5000-8000-000000000004', '2026-05-22T10:00:01Z', '2026-05-22T10:00:06Z', 1),
    ('17171717-0000-5000-8000-000000000003', '00000000-0000-5000-8000-000000000001', 'CLSCycle', '66666666-0000-5000-8000-000000000001', 'cls.closed.v1',          'exchangeos.cls_settlement.events', '66666666-0000-5000-8000-000000000001', '2026-05-22T10:00:30Z', '2026-05-22T10:00:35Z', 1),
    ('17171717-0000-5000-8000-000000000004', '00000000-0000-5000-8000-000000000001', 'EOD',      '15151515-0000-5000-8000-000000000001', 'eod.completed.v1',       'exchangeos.admin.eod_jobs',   '15151515-0000-5000-8000-000000000001', '2026-05-22T22:35:00Z', '2026-05-22T22:35:05Z', 1),
    ('17171717-0000-5000-8000-000000000005', '00000000-0000-5000-8000-000000000001', 'PayIn',    'aaaaaaaa-0000-5000-8000-000000000001', 'payin.confirmed.v1',     'exchangeos.payin.events',     'aaaaaaaa-0000-5000-8000-000000000001', '2026-05-23T05:55:30Z', '2026-05-23T05:55:35Z', 1)
ON CONFLICT (outbox_id) DO NOTHING;

COMMIT;
