-- ExchangeOS — 000012 down: rollback CFETS tables.

BEGIN;

DROP TABLE IF EXISTS cfets_confirmations CASCADE;
DROP TABLE IF EXISTS cfets_captures CASCADE;

COMMIT;
