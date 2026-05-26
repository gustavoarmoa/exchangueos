---
Code: RFLW.024.030.01
Domain: 024 (ExchangeOS)
Module: compliance
Version: 1.0.0
Status: DRAFT
Title: BACEN Classification + IOF Computation on Trade Booked
Traceability:
  RN: [RN_FX_028, RN_FX_037, RN_FX_039]
  Ontology: [exos:Classification, exos:IOFComputation, exos:ScreeningResult]
  Refs: [Circ 3.690 (95 nature codes), Decreto 12.499/2025 (IOF rates)]
Predecessor: RFLW.024.001.01 (Book FX Spot)
Successor: RFLW.024.020.01 (CLS daily cycle)
---

# RFLW.024.030.01 — BACEN Classification + IOF on Trade Booked

## Description

When a Trade is booked, the compliance worker (subscriber to `trade.created.v1`):

1. Resolves BACEN nature code via `pkg/bacen.Classifier`
2. Computes IOF amount via `pkg/bacen.IOFCalculator`
3. Screens counterparty against OFAC/UN/EU/COAF lists
4. If risk HIGH, emits SISCOAF COS per RN_FX_039

## Sequence

```mermaid
sequenceDiagram
    autonumber
    participant KB as Kafka outbox
    participant W as compliance worker
    participant C as ComplianceService
    participant BC as pkg/bacen.Classifier
    participant IC as pkg/bacen.IOFCalculator
    participant SC as SanctionsListProvider

    KB-->>W: trade.created.v1 {trade_id, ...}
    W->>C: ClassifyOperation(tenant, trade_id, hint)
    C->>BC: Classify(hint)
    BC-->>C: NatureCode{10002, "Importação", REMESSA}
    C-->>W: Classification persisted

    W->>C: ComputeIOF(tenant, trade_id, "DEFAULT", notional, ccy)
    C->>IC: Compute("DEFAULT", notional)
    IC-->>C: rate=0.0038, amount=38.00
    C-->>W: IOFComputation persisted

    W->>C: ScreenCounterparty(tenant, buyerBIC, lei)
    C->>SC: Query OFAC/UN/EU/COAF
    SC-->>C: hits[]
    C-->>W: ScreeningResult{risk_level}

    alt RiskHigh
        W->>W: emit SISCOAF COS (RN_FX_039)
    else clear
        Note over W: no further action
    end
```

## Error Flow

```mermaid
flowchart LR
    A[Classify] --> B{nature found?}
    B -- no --> Default[code=99999 fallback]
    B -- yes --> C[Compute IOF]
    C --> D{op_type known?}
    D -- no --> ErrIOF[/log + skip; alert ops/]
    D -- yes --> Screen[Screen]
    Screen --> Decision{risk_level}
    Decision -- HIGH --> COS[SISCOAF COS submission]
    Decision -- else --> Done
```

## Business Rules

- RN_FX_028 — 95 nature codes per Circ 3.690 (classifier ByCode + free-text fallback)
- RN_FX_037 — 6 IOF rates per Decreto 12.499/2025
- RN_FX_039 — COS for SISCOAF within 1 business day of HIGH-risk detection

## Observability

- Metric `compliance.classification.created` counter (label: nature, code)
- Metric `compliance.iof.computed` histogram (label: op_type)
- Metric `compliance.screening.hits` counter (label: risk_level)
- Alert: ScreeningResult HIGH → PagerDuty critical
