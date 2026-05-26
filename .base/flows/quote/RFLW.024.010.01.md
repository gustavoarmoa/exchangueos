---
Code: RFLW.024.010.01
Domain: 024 (ExchangeOS)
Module: quote
Version: 1.0.0
Status: DRAFT
Title: RFQ Streaming Lifecycle (Requested → Quoted → Accepted)
Traceability:
  RN: [RN_FX_001]
  Ontology: [exos:RFQ, exos:Quote, exos:Accepted]
Predecessor: —
Successor: RFLW.024.001.01 (Book FX Spot via CLS)
---

# RFLW.024.010.01 — RFQ Streaming Lifecycle

## Description

Trader requests a quote via RFQ; pricing engine streams multiple quotes; trader
picks one + accepts. Acceptance emits `rfq.accepted.v1` which downstream Trade
flow (RFLW.024.001.01) reacts to.

## Sequence

```mermaid
sequenceDiagram
    autonumber
    participant T as Trader
    participant API as exchangeos-api
    participant Q as QuoteService
    participant P as PricingEngine
    participant KB as Kafka outbox

    T->>API: POST /v1/rfq (EUR/USD, 1M notional)
    API->>Q: CreateRFQ(...)
    Q-->>API: RFQ{status:REQUESTED, version:1}
    Note over Q,KB: emits rfq.requested.v1

    loop streamed quotes
        API->>P: GetMidRate(EUR,USD)
        P-->>API: mid+halfSpread
        API->>Q: GetQuote(...)
        Q-->>API: Quote{id, bid, ask, ttl:10s}
        API->>Q: AttachQuoteToRFQ(rfqID, quoteID)
        Q-->>API: RFQ{status:QUOTED, version:N+1}
        Note over Q,KB: rfq.quoted.v1
        API-->>T: WS stream quote
    end

    T->>API: POST /v1/rfq/:id/accept {quote_id}
    API->>Q: AcceptRFQ(rfqID, quoteID, actor)
    Q-->>API: RFQ{status:ACCEPTED}
    Note over Q,KB: rfq.accepted.v1 → triggers RFLW.024.001.01
```

## Error Flow

```mermaid
flowchart LR
    A[POST accept] --> B{quote_id in RFQ?}
    B -- no --> E1[/400 quote_id not in RFQ/]
    B -- yes --> C{RFQ status QUOTED?}
    C -- no --> E2[/409 invalid transition/]
    C -- yes --> OK[emit accepted + trigger Trade]
```

## Business Rules

- RN_FX_001 — currency pair validated by domain.NewRFQ

## Observability

- OTel span `rfq.lifecycle` covering CreateRFQ → all AttachQuote → Accept
- Metric `rfq.accepted.v1` counter

## Related Patterns

- FX-EDA-* (outbox dispatcher)
- FX-API-* (cursor pagination for RFQ list)
