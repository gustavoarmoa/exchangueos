# ERD — Compliance + Admin Domain

**Source migration:** `migrations/000008_create_compliance_admin.up.sql`
**Ontology:** `.base/aasc/ontology/core/compliance.ttl`

```mermaid
erDiagram
    TENANTS    ||--o{ CLASSIFICATIONS   : "owns"
    TENANTS    ||--o{ IOF_COMPUTATIONS  : "owns"
    TENANTS    ||--o{ BACEN_REPORTS     : "owns"
    TENANTS    ||--o{ SCREENING_RESULTS : "owns"
    TENANTS    ||--o{ EOD_JOBS          : "owns"
    FX_TRADES  ||--o{ CLASSIFICATIONS   : "trade_id"
    FX_TRADES  ||--o{ IOF_COMPUTATIONS  : "trade_id"

    CLASSIFICATIONS {
        uuid    classification_id PK
        uuid    tenant_id FK
        uuid    trade_id FK
        string  code              "BACEN 4-6 digit"
        string  description
        string  nature            "REMESSA|INGRESSO|CONVERSAO"
    }

    IOF_COMPUTATIONS {
        uuid    iof_id PK
        uuid    tenant_id FK
        uuid    trade_id FK
        string  operation_type
        decimal notional          "DECIMAL(36,18) > 0"
        string  notional_ccy
        decimal rate              "fraction (<= 1)"
        decimal iof_amount        "notional × rate, 2-decimal banker"
        timestamptz computed_at
    }

    BACEN_REPORTS {
        uuid    report_id PK
        uuid    tenant_id FK
        string  report_type       "SISBACEN|BCB-CCS|BCB-CAMBIO"
        date    reference_date
        string  payload_hash      "sha256 hex"
        string  status            "PENDING|SUBMITTED|ACCEPTED|REJECTED"
        timestamptz submitted_at
        timestamptz responded_at
        string  rejection_reason
        int     version
    }

    SCREENING_RESULTS {
        uuid    screening_id PK
        uuid    tenant_id FK
        string  counterparty_bic
        string  lei
        string  risk_level        "LOW|MEDIUM|HIGH"
        string[] hits
        timestamptz screened_at
    }

    SYSTEM_EVENTS {
        uuid    event_id PK
        string  code              "STARTUP|SHUTDOWN|CYCLE_OPEN|..."
        string  component
        string  description
        timestamptz at
        string  iso20022_ref
    }

    EOD_JOBS {
        uuid    job_id PK
        uuid    tenant_id FK
        date    business_date     "UNIQUE per tenant"
        string  status            "PENDING|RUNNING|COMPLETED|FAILED"
        timestamptz started_at
        timestamptz completed_at
        string  failure_reason
        string[] steps_done
        int     version
    }
```

## Constraints

- `CLASSIFICATIONS.nature` enum + `code` length 4-6
- `IOF_COMPUTATIONS.rate` 0..1 + `iof_amount >= 0`
- `BACEN_REPORTS.{report_type,status}` enum CHECKs
- `SCREENING_RESULTS.risk_level` enum
- `EOD_JOBS UNIQUE (tenant_id, business_date)` + `status` enum

## Indexes

- `idx_iof_op_type (tenant_id, operation_type, computed_at DESC)` — IOF audit by op
- `idx_bacen_tenant_status (tenant_id, status, reference_date DESC)` — submission queue
- `idx_screen_high (tenant_id, screened_at DESC) WHERE risk_level = 'HIGH'` — partial for SISCOAF triage
- `idx_eod_status (tenant_id, status, business_date DESC)` — EOD orchestrator monitor
