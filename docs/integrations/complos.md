# Integration — ExchangeOS ↔ ComplOS

> Owner: Compliance team
> Compatible since: future (ComplOS not in initial delivery)
> Status: ⏳ Out of scope for initial release

## Purpose

ComplOS is a platform-wide compliance orchestrator that handles cross-module
compliance workflows (KYC across products, sanctions list refresh, regulatory
correspondence). ExchangeOS owns FX-specific compliance (BACEN classification,
IOF, SISCOAF COS, OFAC screening) — see `modules/compliance/` and `pkg/bacen/`.

ComplOS would be the upstream feed for OFAC/UN/EU/COAF sanctions lists +
distribute regulatory policy updates platform-wide.

## Direction (future)

```
ComplOS    ── sanctions.list_updated.v1 ──▶ ExchangeOS  (refresh local cache)
            ── policy.bacen_code_added.v1 ▶
            ── kyc.actor_status_changed.v1 ▶
                                                  │
                                                  ▼
                                  ExchangeOS updates local rules

ExchangeOS ── compliance.report_ready.v1 ─▶ ComplOS  (cross-module aggregation)
            ── compliance.cos_required.v1 ▶
```

## Events ExchangeOS would consume

### `sanctions.list_updated.v1`

ComplOS refreshes one of: OFAC SDN, UN 1267, EU restrictive measures, COAF list.
ExchangeOS invalidates its local screening cache + re-screens active counterparties.

### `policy.bacen_code_added.v1`

ComplOS notifies of new BACEN nature codes (Circ 3.690 updates). ExchangeOS
surfaces alert to engineering to update `pkg/bacen/classifier.go` — does NOT
auto-update (code review required for regulatory data).

### `kyc.actor_status_changed.v1`

ComplOS reports KYC re-verification result. ExchangeOS may need to mark a
counterparty as `cls_member = false` if KYC expired.

## Events ExchangeOS would produce

Already emitted to `exchangeos.compliance.events` (7d retention per PII policy):

- `compliance.classification_created.v1`
- `compliance.iof_computed.v1`
- `compliance.report_ready.v1`
- `compliance.cos_required.v1` (HIGH-risk screening)

## Sync RPCs

### ComplOS → ExchangeOS

```protobuf
service ComplianceQuery {
  rpc ListReportsForPeriod(...) returns (...);
  rpc GetScreeningHistory(...) returns (...);
}
```

Cross-module compliance reports may need to join exchangeos data with sibling
data (paymentos, accountos). Read-only via auditor role (`exchangeos_auditor`).

### ExchangeOS → ComplOS

```protobuf
service SanctionsScreening {
  rpc ScreenParty(...) returns (...);
}
```

ExchangeOS would call ComplOS instead of an in-process stub when ScreenCounterparty
runs. 5s timeout + cache hit/miss telemetry.

## Failure semantics

- **ComplOS down:** ExchangeOS falls back to last cached sanctions list (max age 24h);
  trades continue but with stale screening. Alert ops after 1h.
- **Stale sanctions:** regulatory exposure if list updated and we missed it;
  mitigation = strict 24h freshness SLO + cache-staleness alert.

## Open questions

- [ ] When will ComplOS exist?
- [ ] Does ExchangeOS's local screening stub stay as fallback OR get fully replaced?
- [ ] Cross-module COS: does ComplOS aggregate or does each module submit independently to SISCOAF?
- [ ] BACEN code update cadence: how do we coordinate code review + production deploy?
