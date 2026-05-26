---
Code: RFLW.024.040.01
Domain: 024 (ExchangeOS)
Module: cfets_capture
Version: 1.0.0
Status: DRAFT
Title: CFETS Trade Capture (fxtr.031 → fxtr.032 → fxtr.033)
Traceability:
  RN: [RN_FX_001]
  ISO20022: [fxtr.031.001.02, fxtr.032.001.02, fxtr.033.001.02]
Predecessor: RFLW.024.001.01 (Book FX Trade)
Successor: RFLW.024.041.01 (CFETS Confirmation; not implemented yet)
---

# RFLW.024.040.01 — CFETS Trade Capture

## Description

For trades involving CNY (USDCNY), ExchangeOS captures the trade with CFETS
(China Foreign Exchange Trade System) via PTPP:

1. Submit `fxtr.031.001.02` (Trade Capture Request)
2. Wait for `fxtr.032.001.02` (Ack — SUCC|REJT + CFETSDealID on success)
3. CFETS forwards `fxtr.033.001.02` to the counterparty (informational notification)

## Sequence

```mermaid
sequenceDiagram
    autonumber
    participant TS as TradeService
    participant CC as CFETSCaptureService
    participant CFETS as CFETS PTPP
    participant KB as Kafka outbox

    TS->>CC: Create(tenant, trade_id, submitter_ref)
    CC-->>TS: Capture{status:DRAFT, version:1}

    CC->>CC: Submit(at:now)
    Note over CC,KB: cfets_capture.submitted.v1
    CC->>CFETS: fxtr.031.001.02 (Trade Capture Request)

    CFETS-->>CC: fxtr.032.001.02 {status:SUCC, cfetsDealID}
    CC->>CC: Ack(at, cfetsDealID)
    Note over CC,KB: cfets_capture.acked.v1

    CFETS-->>CC: fxtr.033.001.02 (counterparty notified)
    CC->>CC: NotifyCounterparty(at)
    Note over CC,KB: cfets_capture.notified.v1
```

## Error Flow

```mermaid
flowchart TB
    Submit[Submit] --> Wait[Wait for fxtr.032]
    Wait --> Status{Ack status}
    Status -- SUCC --> Ack[transition to ACK]
    Status -- REJT --> Reject[transition to REJECTED + record reason]
    Status -- timeout --> Timeout[/log + alert; manual retry/]
```

## Business Rules

- Currency pair must include CNY (USDCNY etc.) — RN_FX_001 validation upstream

## Observability

- Metric `cfets_capture.submitted.v1` counter
- Metric `cfets_capture.acked.v1` / `cfets_capture.rejected.v1` (label: reason)
- Histogram for ack latency (submit → ack)
