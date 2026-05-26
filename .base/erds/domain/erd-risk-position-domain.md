# ERD — Risk + Position Domain

**Source migration:** `migrations/000007_create_risk_position.up.sql`
**Ontology:** (planned in `.base/aasc/ontology/core/risk.ttl` + `position.ttl`)

```mermaid
erDiagram
    TENANTS ||--o{ RISK_LIMITS : "owns"
    TENANTS ||--o{ POSITIONS   : "owns"

    RISK_LIMITS {
        uuid    limit_id PK
        uuid    tenant_id FK
        string  limit_type        "COUNTERPARTY|CURRENCY|TENOR|DV01|VAR"
        string  scope             "BIC/CCY/tenor/'' for portfolio-wide"
        decimal cap               "DECIMAL(36,18) > 0"
        decimal utilised          ">= 0"
        string  currency          "ISO 4217"
        int     version
    }

    POSITIONS {
        uuid    position_id PK
        uuid    tenant_id FK
        string  currency
        decimal long_amount       "DECIMAL(36,18) >= 0"
        decimal short_amount      ">= 0"
        decimal net_amount        "signed: long - short"
        timestamptz as_of
        int     version
    }
```

## Constraints

- `RISK_LIMITS.limit_type` enum CHECK
- `RISK_LIMITS UNIQUE (tenant_id, limit_type, scope)` — one limit per scope
- `RISK_LIMITS.cap > 0` + `utilised >= 0`
- `POSITIONS.long/short_amount >= 0`
- `POSITIONS UNIQUE (tenant_id, currency)` — one position per CCY

## Indexes

- `idx_limits_breaching (tenant_id) WHERE utilised >= cap` — partial for breach observability
- `idx_limits_tenant_type (tenant_id, limit_type)` — Find() lookup
- `idx_positions_tenant_asof (tenant_id, as_of DESC)` — most-recent positions
