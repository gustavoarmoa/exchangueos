# ExchangeOS Segregation of Duties (SoD) Matrix

> ISO 27001 A.5.3 + A.8.2 compliance evidence
> Owner: Security + Compliance teams
> Last reviewed: 2026-05-24

## Roles defined

| Role | Description | Identity source |
|------|-------------|-----------------|
| `platform-engineer` | Day-to-day ops: deploy, troubleshoot, scale | KeycloakOS group `revenu-platform:platform-engineer` |
| `security-officer` | Security policy, audit, incident lead | KeycloakOS group `revenu-platform:security-officer` |
| `compliance-officer` | BACEN/SISCOAF interactions, regulatory reports | KeycloakOS group `revenu-platform:compliance-officer` |
| `trader` | Books trades via API (M2M client) | KeycloakOS M2M client `exchangeos-trader-<tenant>` |
| `operator` | Read-only production observability | KeycloakOS group `revenu-platform:operator` |
| `auditor` | Read-only access to all audit + compliance data | KeycloakOS group `revenu-platform:auditor` |
| `dba` | CRDB cluster operations (shared hub) | Cross-platform role; not exchangeos-scoped |

## Critical action matrix

Row = action. Column = role. Cell = ✓ allowed / ✗ denied / Δ requires 4-eyes.

| Action | platform-eng | security | compliance | trader | operator | auditor | dba |
|--------|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| **Code + Deploy** |  |  |  |  |  |  |  |
| Merge PR to main | ✓ Δ | ✓ Δ | ✗ | ✗ | ✗ | ✗ | ✗ |
| Sign release with Cosign keyless | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ |
| Trigger ArgoCD sync | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ |
| Abort canary rollout | ✓ | ✓ | ✗ | ✗ | ✓ | ✗ | ✗ |
| EMERGENCY_BYPASS git hook | ✓ Δ + log | ✓ Δ + log | ✗ | ✗ | ✗ | ✗ | ✗ |
| **Production data** |  |  |  |  |  |  |  |
| Read fx_trades (any tenant) | ✗ | ✓ | ✓ | ✗ | ✗ | ✓ | ✗ |
| Read fx_trades (own tenant only) | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ |
| Book FX trade | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ |
| Cancel trade | ✗ | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ |
| Mark trade settled (manual override) | ✗ | ✓ Δ | ✗ | ✗ | ✗ | ✗ | ✗ |
| Mutate refdata (currencies/calendars) | ✓ Δ | ✓ Δ | ✗ | ✗ | ✗ | ✗ | ✗ |
| **Risk + Compliance** |  |  |  |  |  |  |  |
| Update risk limits | ✓ Δ | ✓ Δ | ✗ | ✗ | ✗ | ✗ | ✗ |
| Override breach (allow trade above limit) | ✗ | ✓ Δ | ✓ Δ | ✗ | ✗ | ✗ | ✗ |
| Mark BACEN report submitted | ✗ | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ |
| Add hits to ScreeningResult | ✗ | ✓ | ✓ | ✗ | ✗ | ✗ | ✗ |
| Submit SISCOAF COS | ✗ | ✗ | ✓ Δ | ✗ | ✗ | ✗ | ✗ |
| **Secrets + Identity** |  |  |  |  |  |  |  |
| Rotate M2M client_secret | ✓ Δ | ✓ Δ | ✗ | ✗ | ✗ | ✗ | ✗ |
| Read Vault `secret/data/exchangeos/*` | ✗ | ✓ Δ | ✗ | ✗ | ✗ | ✗ | ✗ |
| Modify Vault policies | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ |
| Issue Workload Identity bindings | ✓ Δ | ✓ Δ | ✗ | ✗ | ✗ | ✗ | ✗ |
| **Database** |  |  |  |  |  |  |  |
| Connect to CRDB hub as exchangeos_app | ✗ (via app only) | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ |
| Run ad-hoc UPDATE/DELETE in production | ✗ | ✗ | ✗ | ✗ | ✗ | ✗ | ✓ Δ |
| Restore from backup | ✗ | ✓ | ✗ | ✗ | ✗ | ✗ | ✓ Δ |
| **Observability** |  |  |  |  |  |  |  |
| Read Grafana dashboards | ✓ | ✓ | ✓ | ✗ | ✓ | ✓ | ✓ |
| Modify alert rules | ✓ | ✓ | ✗ | ✗ | ✗ | ✗ | ✗ |
| Read audit_events (cross-tenant) | ✗ | ✓ | ✓ | ✗ | ✗ | ✓ | ✗ |

## Δ (4-eyes) enforcement

Critical actions marked Δ require a second authorised approver. Mechanism:

- **PR merges:** GitHub branch protection requires 2 approvers + signed commits
- **EMERGENCY_BYPASS:** `scripts/git-hooks-wrapper.sh` records to `.git/audit-bypass.log` + Slack alert; daily review by security
- **Risk limit override / BACEN report submission:** approval ticket in compliance system; system records approver_actor_id alongside trade_amendments.approver_actor_id
- **Vault secret read / policy change:** Vault audit log + Slack `#vault-audit` channel

## Conflict-of-duties analysis

The following pairs MUST be held by different individuals (enforced via group exclusivity in KeycloakOS):

- `trader` ∩ `compliance-officer` → forbidden (trader can't classify own trades)
- `trader` ∩ `security-officer` → forbidden (can't override own breaches)
- `dba` ∩ `platform-engineer` → discouraged but allowed for emergency response (logged)
- `auditor` ∩ any other role → forbidden (auditor must be independent)

## Review cadence

- Quarterly review of role memberships + matrix accuracy
- Post-incident: audit log review for SoD violations
- Annual: external auditor review
