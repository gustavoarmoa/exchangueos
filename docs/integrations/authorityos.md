# Integration — ExchangeOS ↔ AuthorityOS

> Owner: Compliance team
> Compatible since: ExchangeOS v4.11.0 (compliance domain delivered)
> Status: 🟡 Spec — AuthorityOS adapter pending

## Purpose

AuthorityOS centralises Revenu Platform's interactions with regulatory
authorities (BACEN, SISCOAF, COAF, CVM). ExchangeOS generates raw BACEN
reports + screening events; AuthorityOS owns the secure submission channel,
authority retries, response normalisation, and the audit trail required by
external regulators.

## Direction

```
ExchangeOS ── compliance.report_ready.v1 ────▶ AuthorityOS  (Kafka topic)
            ── compliance.cos_required.v1 ─▶
                                                  │
                                                  ▼
                                      submits to BACEN/SISCOAF
                                                  │
AuthorityOS ── regulator.response_received.v1 ─▶ ExchangeOS
            ── regulator.policy_updated.v1 ───▶
```

## Events ExchangeOS produces

### `compliance.report_ready.v1`

Emitted when a BACENReport transitions to PENDING (the payload is finalised
and ready for submission). AuthorityOS picks up + submits.

```json
{
  "event_id": "uuid",
  "occurred_at": "RFC3339",
  "tenant_id": "uuid",
  "report_id": "uuid",
  "report_type": "SISBACEN|BCB-CCS|BCB-CAMBIO",
  "reference_date": "2026-05-26",
  "payload_hash": "sha256...",
  "payload_location": "s3://exchangeos-reports/<tenant>/<report_id>.xml"
}
```

### `compliance.cos_required.v1`

Emitted when a ScreeningResult is HIGH and RequiresCOS() is true (RN_FX_039).
AuthorityOS submits the SISCOAF COS within 1 business day.

## Events ExchangeOS consumes

### `regulator.response_received.v1`

AuthorityOS reports the regulator's response. ExchangeOS updates BACENReport:

- `MarkAccepted(at)` on success
- `MarkRejected(at, reason)` on rejection

### `regulator.policy_updated.v1`

AuthorityOS notifies of regulatory changes (e.g. new BACEN nature code,
updated IOF rate). ExchangeOS:

- Logs the event in admin.system_events
- Surfaces alert to Compliance team for manual update of `pkg/bacen` catalog
  (cannot auto-update without code review)

## Sync RPCs

### ExchangeOS → AuthorityOS

Async-only. No sync RPCs (regulators are slow + heavy; everything queues).

### AuthorityOS → ExchangeOS

```protobuf
service ComplianceQuery {
  // Read-only view for AuthorityOS to fetch full BACEN report payload.
  rpc GetReport(GetReportRequest) returns (GetReportResponse);

  // Validate that a tenant has authority to operate before AuthorityOS submits.
  rpc CheckAuthorityScope(CheckAuthorityScopeRequest) returns (CheckAuthorityScopeResponse);
}
```

## Failure semantics

- **AuthorityOS down:** reports queue indefinitely in outbox. Regulators have
  deadlines (DEC: USD 10k threshold immediate; COS: 1 business day) — alert
  ops if outbox lag > 4h on these topics.
- **Regulator rejects:** ExchangeOS marks REJECTED + Compliance team manually
  remediates via ComplianceService.SubmitBACENReport with corrected payload.
- **Repeated rejection:** escalate to BACEN technical liaison.

## Open questions

- [ ] BACEN payload format: native ISO 20022 XML or SISBACEN-specific?
- [ ] LGPD: AuthorityOS sees PII via screening events — separate retention policy?
- [ ] SISCOAF COS template: which fields can ExchangeOS pre-populate vs which need Compliance review?
- [ ] Cross-tenant regulator account: is the cluster a single registered entity at BACEN, or per-tenant?
