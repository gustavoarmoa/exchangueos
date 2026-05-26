---
Code: RFLW.024.002.01
Domain: 024 (ExchangeOS)
Module: trade
Version: 1.0.0
Status: DRAFT
Title: Settle FX Trade via CLS PvP
Traceability:
  RN: [RN_FX_010, RN_FX_026]
  ISO20022: [fxtr.030.001.05, camt.088.001.02]
  Ontology: [exost:FXTrade, exost:Settled, exosc:CLSCycle]
Predecessor: RFLW.024.001.01 (Book FX Spot)
Successor: RFLW.024.030.01 (BACEN compliance)
---

# RFLW.024.002.01 — Settle FX Trade via CLS PvP

## Description

After the CLS cycle closes for the trade's value_date, CLS reports settlement
via fxtr.030 + camt.088. ExchangeOS reflects the SETTLING → SETTLED transition
on the FXTrade aggregate and emits trade.settled.v1.

## Sequence

```mermaid
sequenceDiagram
    autonumber
    participant CLS as CLS Bank
    participant API as exchangeos-api
    participant TS as TradeService
    participant SS as SettlementService
    participant KB as Kafka outbox

    CLS->>API: fxtr.030 SettlementNotification
    API->>SS: GetNetReport(cycle_id)
    SS-->>API: camt.088 lines
    API->>TS: MarkSettling(trade_id)
    TS-->>API: Trade{status:SETTLING, version+1}
    Note over TS,KB: emits trade.settling.v1

    API->>TS: MarkSettled(trade_id, settlement_ref=fxtr.030.OurTradRef)
    TS-->>API: Trade{status:SETTLED, version+1}
    Note over TS,KB: emits trade.settled.v1 → downstream BACEN report
```

## Error Flow

```mermaid
flowchart TB
    A[fxtr.030 received] --> B{Trade in SETTLING?}
    B -- no --> Err[/409 invalid transition + alert ops/]
    B -- yes --> C{settlement_ref present?}
    C -- no --> Err2[/400 missing ref/]
    C -- yes --> OK[transition + emit + downstream BACEN report]
```

## Business Rules

- RN_FX_010 — CLS PvP for 18 eligible CCYs
- RN_FX_026 — decimal precision preserved through settlement amounts
- Cancellation forbidden after SETTLING (domain enforces)

## Observability

- Metric `trade.settled.v1` counter (label: settlement_venue)
- Span links: settle span linked to original book span via trade_id
- Settlement latency histogram (trade_date → settled_at)
