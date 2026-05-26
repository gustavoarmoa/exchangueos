# ERD — CLS Settlement Domain (CLS Cycles + PayIn + NetReport)

**Source migration:** `migrations/000006_create_settlement.up.sql`
**Ontology:** `.base/aasc/ontology/core/cls_settlement.ttl`

```mermaid
erDiagram
    TENANTS    ||--o{ CLS_CYCLES         : "owns"
    CLS_CYCLES ||--o{ CLS_CYCLE_TRADES   : "cycle_id"
    FX_TRADES  ||--o{ CLS_CYCLE_TRADES   : "trade_id"
    CLS_CYCLES ||--o{ PAYIN_INSTRUCTIONS : "cycle_id"
    CLS_CYCLES ||--o{ NET_REPORTS        : "cycle_id"

    CLS_CYCLES {
        uuid    cycle_id PK
        uuid    tenant_id FK
        date    cycle_date            "UNIQUE per tenant"
        string  status                "OPEN|PAY_IN_WINDOW|SETTLING|CLOSED|FAILED"
        timestamptz opened_at         "07:00 CET"
        timestamptz pin1_deadline     "08:00 CET (Asia)"
        timestamptz pin2_deadline     "09:00 CET (Europe)"
        timestamptz pin3_deadline     "10:00 CET (Americas)"
        timestamptz scheduled_close   "12:00 CET"
        timestamptz closed_at
        string  failure_reason
        int     version
    }

    CLS_CYCLE_TRADES {
        uuid cycle_id FK
        uuid trade_id FK
        timestamptz attached_at
    }

    PAYIN_INSTRUCTIONS {
        uuid    instruction_id PK
        uuid    tenant_id FK
        uuid    cycle_id FK
        string  currency
        decimal amount             "DECIMAL(36,18)"
        string  band               "PIN1|PIN2|PIN3"
        timestamptz deadline
        string  status             "PENDING|SUBMITTED|CONFIRMED|FAILED"
        timestamptz submitted_at
        timestamptz confirmed_at
        string  failure_reason
        int     version
    }

    NET_REPORTS {
        uuid    report_id PK
        uuid    tenant_id FK
        uuid    cycle_id FK
        string  currency
        decimal gross_pay_in       "DECIMAL(36,18)"
        decimal gross_pay_out      "DECIMAL(36,18)"
        decimal net_settlement     "signed; > 0 receivable"
        int     trade_count
        timestamptz generated_at
    }
```

## Constraints

- `CLS_CYCLES.status` enum CHECK + `pin1 < pin2 < pin3 < scheduled_close`
- `UNIQUE (tenant_id, cycle_date)` — one cycle per business date
- `PAYIN_INSTRUCTIONS.amount > 0`
- `PAYIN_INSTRUCTIONS.band` ∈ {PIN1,PIN2,PIN3}
- `NET_REPORTS.gross_pay_in / gross_pay_out >= 0`
- `UNIQUE (cycle_id, currency)` on NET_REPORTS

## Indexes

- `idx_cycles_open (tenant_id, cycle_date) WHERE status IN (OPEN,PAY_IN_WINDOW,SETTLING)` — partial for scheduler
- `idx_payin_cycle_ccy (cycle_id, currency)` — NetReport aggregation
- `idx_payin_status_deadline (tenant_id, status, deadline)` — deadline monitor
