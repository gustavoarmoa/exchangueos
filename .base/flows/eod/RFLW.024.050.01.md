---
Code: RFLW.024.050.01
Domain: 024 (ExchangeOS)
Module: eod
Version: 1.0.0
Status: DRAFT
Title: End-of-Day Batch (PTAX → MTM → Position Snapshot → BACEN Report)
Traceability:
  RN: [RN_FX_028, RN_FX_037]
  Ontology: [admin:EODJob, refdata:SpotRate, position:Position, compliance:BACENReport]
Predecessor: RFLW.024.020.01 (CLS daily cycle close)
Successor: —
---

# RFLW.024.050.01 — End-of-Day Batch

## Description

Daily batch orchestrated by `cmd/eod` (CronJob 23:00 UTC weekdays). Runs 4
canonical steps with idempotent step-mark tracking on the EODJob aggregate.
If any step fails the job moves to FAILED with the step as failure_reason.

## Sequence

```mermaid
sequenceDiagram
    autonumber
    participant CRON as eod CronJob
    participant AS as AdminService
    participant PE as PTAX Fetcher (OLINDA)
    participant PR as PricingEngine (MTM)
    participant PS as PositionService
    participant CS as ComplianceService
    participant KB as Kafka outbox

    CRON->>AS: TriggerEOD(tenant, business_date)
    AS-->>CRON: EODJob{status:PENDING}
    CRON->>AS: StartEOD(job_id)
    Note over AS,KB: emits admin.eod_started

    CRON->>PE: FetchPTAX(business_date)
    PE-->>CRON: PTAX{4 windows, WeightedFixing}
    CRON->>AS: MarkStep("PTAX")

    CRON->>PR: PositionMTM per position
    PR-->>CRON: P&L per (tenant, currency)
    CRON->>AS: MarkStep("MTM")

    CRON->>PS: List positions + snapshot to position_snapshots
    CRON->>AS: MarkStep("POSITION_SNAPSHOT")

    CRON->>CS: SubmitBACENReport(SISBACEN, business_date, payload)
    CS-->>CRON: BACENReport{status:PENDING}
    CRON->>AS: MarkStep("BACEN_REPORT")

    CRON->>AS: CompleteEOD(job_id)
    Note over AS,KB: emits admin.eod_completed
```

## Error Flow

```mermaid
flowchart TB
    A[Step fails] --> B[FailEOD reason=step]
    B --> C[admi.004 EOD_FAILED event]
    C --> D[PagerDuty critical]
    D --> Retry{Next run<br/>(idempotent)}
    Retry -- yes --> Resume[MarkStep skips done steps]
    Retry -- no --> Manual[Manual intervention]
```

## Business Rules

- Idempotent: re-running an EODJob with already-completed steps via `MarkStep` is a no-op
- One EOD per (tenant, business_date) — UNIQUE constraint on eod_jobs
- RN_FX_028 — BACEN classifications must be in place before SISBACEN submission
- RN_FX_037 — IOF computations frozen at EOD timestamp

## Observability

- Metric `admin.eod_job.duration` histogram per step
- Metric `admin.eod_job.status` counter (label: status)
- Grafana panel: "Last 30 EOD runs by tenant"
- Alert: any FAILED EOD → critical
