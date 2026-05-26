# ExchangeOS Threat Model — STRIDE + DREAD

> Owner: Security team
> Last reviewed: 2026-05-24 (initial draft)
> Cadence: review on every major release + post-incident

## Scope

ExchangeOS API + worker + cls-cycle + eod + migrator + mq-bridge + cred-rotator,
plus their hosting (GKE Autopilot, shared CRDB hub, Vault, Kafka, OTel stack).

Out of scope: GCP physical security (inherited), sibling Revenu Platform modules
(reviewed in their own threat models).

## Method

For each architectural component we enumerate STRIDE threats:

| Letter | Threat type |
|--------|-------------|
| S | Spoofing |
| T | Tampering |
| R | Repudiation |
| I | Information disclosure |
| D | Denial of service |
| E | Elevation of privilege |

Each threat is scored DREAD (0-10 scale):

| Letter | Dimension |
|--------|-----------|
| D | Damage potential |
| R | Reproducibility |
| E | Exploitability |
| A | Affected users |
| D | Discoverability |

Risk = average of DREAD scores. > 7.0 = critical, > 5.0 = high, > 3.0 = medium, else low.

## Component: exchangeos-api (gRPC + REST)

| # | STRIDE | Threat | DREAD | Risk | Mitigation |
|---|--------|--------|-------|------|------------|
| T-1 | S | Caller spoofs counterparty BIC in trade request | D=8 R=5 E=6 A=8 D=4 → 6.2 | High | OIDC + tenant context + counterparty BIC validated against bic_records refdata + screening |
| T-2 | T | Caller alters bought/sold amounts mid-flow | D=9 R=3 E=4 A=6 D=3 → 5.0 | High | Decimal precision + amounts immutable post-construction + audit envelope |
| T-3 | R | Trader denies booking a trade | D=7 R=8 E=8 A=7 D=2 → 6.4 | High | audit_events table + outbox events + OTel trace correlation_id end-to-end |
| T-4 | I | Listing trades across tenants | D=10 R=7 E=5 A=10 D=3 → 7.0 | Critical | Every query filters tenant_id; tenant scoping enforced in Service layer |
| T-5 | D | Quote-spam DoS via /v1/quotes | D=6 R=9 E=8 A=6 D=8 → 7.4 | Critical | Rate limit at KrakenD gateway + HPA autoscale + Argo Rollouts auto-rollback |
| T-6 | E | Reading raw fx_trades via gRPC reflection without auth | D=8 R=8 E=7 A=8 D=6 → 7.4 | Critical | gRPC reflection disabled in production (Helm values) + JWT validation at gateway |

## Component: exchangeos-worker (outbox dispatcher)

| # | STRIDE | Threat | DREAD | Risk |
|---|--------|--------|-------|------|
| W-1 | T | Worker republishes already-dispatched event after restart | D=5 R=9 E=4 A=5 D=5 → 5.6 | High |
| W-2 | I | Outbox payload leaks PII via Kafka topic logs | D=6 R=4 E=5 A=4 D=3 → 4.4 | Medium |
| W-3 | E | Kafka credentials leak through pod env | D=9 R=2 E=2 A=8 D=2 → 4.6 | Medium |

**Mitigations:**
- W-1: at-least-once delivery + consumer idempotency via outbox_id dedupe table
- W-2: PII-bearing topics (`exchangeos.compliance.events`) have 7-day retention + access ACL
- W-3: WIF + External Secrets — secrets never in env files

## Component: CockroachDB hub (shared)

| # | STRIDE | Threat | DREAD | Risk |
|---|--------|--------|-------|------|
| C-1 | E | Cross-tenant SQL injection via raw concat in repo Save | D=10 R=3 E=2 A=10 D=2 → 5.4 | High |
| C-2 | I | etcd snapshot exfiltration | D=10 R=2 E=2 A=10 D=2 → 5.2 | High |

**Mitigations:**
- C-1: pgx parameterised queries everywhere (FX-CP audit); forbidigo blocks string-concat in SQL paths
- C-2: CMEK HSM encryption for etcd (Terraform `exchangeos-gke.kms_crypto_key`)

## Component: Vault + secrets

| # | STRIDE | Threat | DREAD | Risk |
|---|--------|--------|-------|------|
| V-1 | I | Stolen Vault token reads all exchangeos secrets | D=10 R=3 E=4 A=10 D=3 → 6.0 | High |
| V-2 | E | Operator misconfigures ESO ClusterSecretStore + leaks to wrong namespace | D=8 R=4 E=5 A=7 D=4 → 5.6 | High |

**Mitigations:**
- V-1: Vault token TTL ≤ 1h + automatic rotation; K8s auth method scoped to (SA, namespace)
- V-2: ESO namespace-scoped `SecretStore` (not ClusterSecretStore) + RBAC review

## Component: Kafka cluster

| # | STRIDE | Threat | DREAD | Risk |
|---|--------|--------|-------|------|
| K-1 | E | Service-A produces to Service-B's private topic | D=8 R=5 E=3 A=6 D=4 → 5.2 | High |
| K-2 | D | Producer flood saturates broker disk | D=7 R=6 E=6 A=7 D=5 → 6.2 | High |

**Mitigations:**
- K-1: SASL/SCRAM identity per service + per-topic ACL (`deploy/kafka/topics.yaml`)
- K-2: Topic-level retention + producer quota (broker config)

## Aggregate posture

- **Critical risks (>7):** 3 (T-4 tenant scoping, T-5 quote spam DoS, T-6 reflection leakage) — all mitigated
- **High risks (5..7):** 9 — all mitigated
- **Medium / low:** 4

**Top residual risk:** insider Operator with admin token in Vault. Mitigation: Vault root token sealed; daily ops use short-lived tokens with `auth/kubernetes/role/exchangeos`.

## Review cadence

- Major release → full STRIDE re-run + DREAD re-score
- Post-incident → add new threat row + mitigation
- Quarterly → review residual risks with Security + Platform leads
