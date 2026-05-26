# ERD — Quote Domain

**Source migration:** `migrations/000003_create_quotes.up.sql`
**Ontology:** `.base/aasc/ontology/core/quote.ttl`

```mermaid
erDiagram
    TENANTS  ||--o{ RFQS    : "owns"
    TENANTS  ||--o{ QUOTES  : "owns"
    TENANTS  ||--o{ QUOTE_STREAMS : "owns"
    RFQS     ||--o{ QUOTES  : "rfq_id (nullable for streams)"

    RFQS {
        uuid    rfq_id PK
        uuid    tenant_id FK
        string  requester
        string  base_ccy
        string  quote_ccy
        string  status        "REQUESTED|QUOTED|ACCEPTED|REJECTED|EXPIRED"
        int     version
        timestamptz created_at
    }

    QUOTES {
        uuid    quote_id PK
        uuid    tenant_id FK
        uuid    rfq_id FK
        string  base_ccy
        string  quote_ccy
        decimal notional      "DECIMAL(36,18)"
        string  notional_ccy
        decimal bid           "DECIMAL(36,18)"
        decimal ask           "DECIMAL(36,18)"
        timestamptz valid_from
        timestamptz valid_to
        string  venue
        int     version
    }

    QUOTE_STREAMS {
        uuid       stream_id PK
        uuid       tenant_id FK
        string[]   pairs
        timestamptz started_at
        timestamptz ended_at
        timestamptz last_seen_at
    }
```

## Constraints

- `QUOTES.bid <= ask` (CHECK)
- `QUOTES.base_ccy <> quote_ccy` (RN_FX_001)
- `QUOTES.valid_to > valid_from`
- `QUOTES.notional_ccy IN (base_ccy, quote_ccy)`
- `RFQS.base_ccy <> quote_ccy`

## Indexes

- `idx_quotes_rfq (rfq_id)` — RFQ → quotes lookup
- `idx_quotes_tenant_pair_valid (tenant_id, base_ccy, quote_ccy, valid_to DESC)` — live quote feed
- `idx_streams_tenant_open (tenant_id) WHERE ended_at IS NULL` — open streams only
