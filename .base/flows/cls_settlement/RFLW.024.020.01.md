---
Code: RFLW.024.020.01
Domain: 024 (ExchangeOS)
Module: cls_settlement
Version: 1.0.0
Status: DRAFT
Title: CLS Daily Cycle Lifecycle (07:00 Open → 12:00 Close)
Traceability:
  RN: [RN_FX_010]
  Ontology: [exos:CLSCycle, exos:Open, exos:Closed]
Predecessor: RFLW.024.001.01 (Book FX Spot)
Successor: RFLW.024.030.01 (BACEN compliance)
---

# RFLW.024.020.01 — CLS Daily Cycle Lifecycle

## Description

Scheduled cycle orchestration anchored to Europe/Zurich CET timezone:

| CET    | Event |
|--------|-------|
| 07:00  | Cycle OPEN — trades attach |
| 08:00  | PIN1 deadline (Asia-Pacific PayIns) |
| 09:00  | PIN2 deadline (Europe PayIns) |
| 10:00  | PIN3 deadline (Americas PayIns) — enter SETTLING |
| 12:00  | Cycle CLOSED — NetReport emitted |

## Sequence

```mermaid
sequenceDiagram
    autonumber
    participant CRON as cls-cycle (cron)
    participant SS as SettlementService
    participant TS as TradeService
    participant PS as PayInService
    participant NS as NetReportService
    participant KB as Kafka outbox

    CRON->>SS: OpenCycle(tenant, 2026-05-22)
    SS-->>CRON: Cycle{status:OPEN}
    Note over SS,KB: cls_cycle.opened.v1

    loop attach trades during OPEN window
        TS->>SS: AttachTrade(cycleID, tradeID)
    end

    CRON->>SS: EnterPayInWindow(at:08:00 CET)
    Note over SS,KB: cls_cycle.payin_opened.v1

    par PIN1/2/3 in parallel by CCY
        PS->>SS: SubmitPayIn (PIN1, JPY/AUD/...)
        PS->>SS: SubmitPayIn (PIN2, EUR/GBP/CHF/...)
        PS->>SS: SubmitPayIn (PIN3, USD/CAD/MXN)
    end

    CRON->>SS: EnterSettling(at:10:00 CET)
    Note over SS,KB: cls_cycle.settling.v1

    CRON->>NS: Generate NetReport per (cycle, currency)
    NS-->>CRON: 18 NetReports

    CRON->>SS: CloseCycle(at:12:00 CET)
    Note over SS,KB: cls_cycle.closed.v1
```

## Error Flow

```mermaid
flowchart TB
    A[PIN deadline reached] --> B{All PayIns confirmed?}
    B -- yes --> Settle[enter SETTLING]
    B -- no --> Fail[FailCycle reason="PIN deadline missed"]
    Fail --> Notify[admi.004 system event]
```

## Business Rules

- RN_FX_010 — PvP for 18 CLS-eligible currencies
- Deadline ordering enforced by `cls_cycles` CHECK pin1 < pin2 < pin3 < scheduled_close

## Observability

- Metric `cls_cycle.opened.v1` / `cls_cycle.closed.v1` counters
- Per-PIN-band PayIn latency histogram
- Grafana dashboard panel: deadline burn-down

## Compliance Notes

- CLS daily cycle is monitored by BACEN for CCY exposure (DEC reporting if BRL on either leg).
