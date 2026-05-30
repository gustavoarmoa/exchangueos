# ExchangeOS — Delivery Dashboard

> **Versao:** 1.0.0 | **Auto-gen:** `task dash-update` (delega para `.claude/scripts/generate-delivery-dashboard.sh`)
> **Live Grafana:** `https://grafana.exchangeos.revenu.tech/d/exchangeos-delivery`
> **Refresh cadence:** Auto-update hourly via cron + manual via `task dash-update`

---

## 🎯 Snapshot Executivo

```
┌──────────────────────────────────────────────────────────────────┐
│  ExchangeOS Delivery Status — Sprint 024 of 19 (auto-updated)        │
├──────────────────────────────────────────────────────────────────┤
│  Overall:      █████████████░░░░░░░  66% delivered (26/39 MS)     │
│  This sprint:  See active/ milestones                            │
│  Velocity:     Commits 7d: 70 │ 30d: 70                          │
│  Health:       🟢 ON TRACK                          │
│  Last update:  2026-05-30                                          │
└──────────────────────────────────────────────────────────────────┘
```

| KPI | Atual | Target | Status |
|-----|-------|--------|--------|
| **Milestones DELIVERED** | 0 | 26 | 🔴 0% |
| **Milestones ACTIVE** | 0 | 1-2 | 🟡 |
| **Milestones BACKLOG** | 26 | (declining) | 🟢 stable |
| **Plan version** | 4.0.0 | (stable) | 🟢 |
| **Open questions resolved** | 0 | 108 | 🔴 |
| **Sprint number** | 0 | 19 | 🟡 |
| **Days to ISO 27001 audit** | ~16 sprints × 14d = 224d | Sprint 16 | 🟡 on track |
| **Patterns documented** | 850 (em 20 catalogos) | 850 | 🟢 100% |

---

## 📅 Sprint Burndown (planejado)

> Burndown estimado baseado em milestones distribuidos por sprint conforme `master-plan.md`.

```
Sprint  | Target Milestones | Delivered | Cumulative %
--------|-------------------|-----------|-------------
   1    | MS-023a (started) | 0         | 0%
   2    | MS-023a (done)    | 1         | 4%
   3    | MS-023b           | 2         | 8%
   4    | MS-023c           | 3         | 12%
   5    | MS-023d, MS-023d2 | 5         | 19%
   6    | MS-023e           | 6         | 23%
   7    | MS-023f           | 7         | 27%
   8    | MS-023f2, MS-023g | 9         | 35%
   9    | (MS-023h ongoing) | 9         | 35%
  10    | MS-023h, MS-023i  | 11        | 42%
  11    | (MS-023i, MS-023j)| 11        | 42%
  12    | MS-023j, MS-023k, MS-023l, MS-023m | 15 | 58%
  13    | MS-023n           | 16        | 62%
  14    | MS-023o, MS-023p  | 18        | 69%
  15    | MS-023p, MS-023q  | 19        | 73%
  16    | MS-023q (DONE), MS-023r, MS-023s | 22 | 85%  ← ISO 27001 AUDIT
  17    | MS-023s, MS-023t, MS-023u | 24 | 92%
  18    | MS-023u, MS-023v, MS-023w | 25 | 96%
  19    | MS-023w, MS-023x  | 26        | 100%  ← TARGET COMPLETION
```

Gain: 100% delivery em 19 sprints (~5 meses @ 14d cada).

---

## 📊 Milestones Pipeline

### BACKLOG (26 milestones)

| Status | Code | Sprint | Title | Owner | Blocked by |
|--------|------|--------|-------|-------|------------|
| 🔵 BACKLOG | MS-023a | 1-2 | Foundation & Scaffolding | Platform | None (entry) |
| 🔵 BACKLOG | MS-023b | 3 | RefData + Pricing + Quote | Platform | MS-023a |
| 🔵 BACKLOG | MS-023c | 4 | Trade Core | Platform | MS-023b |
| 🔵 BACKLOG | MS-023d | 5 | Settlement CLS + non-CLS | Platform | MS-023c |
| 🔵 BACKLOG | MS-023d2 | 5 | CFETS Capture + Confirmation | Platform | MS-023c |
| 🔵 BACKLOG | MS-023e | 6 | Risk + Position + Ledger | Platform | MS-023d, d2 |
| 🔵 BACKLOG | MS-023f | 7 | Compliance Core + Admin | Platform | MS-023e |
| 🔵 BACKLOG | MS-023f2 | 7-8 | BACEN Integration Suite | Compliance | MS-023f |
| 🔵 BACKLOG | MS-023g | 8 | EDA E2E | Platform | MS-023d/d2/e/f/f2 |
| 🔵 BACKLOG | MS-023h | 9-10 | Production | Platform + SRE | MS-023g |
| 🔵 BACKLOG | MS-023i | 10-11 | Allenty Documentation | Platform | MS-023h |
| 🔵 BACKLOG | MS-023j | 11-12 | Ontology Suite v1.2.0 | Ontology | MS-023i |
| 🔵 BACKLOG | MS-023k | 12 | Flows Suite | Platform | MS-023j |
| 🔵 BACKLOG | MS-023l | 12 | ERDs Suite | Platform | MS-023j |
| 🔵 BACKLOG | MS-023m | 12-13 | Patterns App Layer | Platform | MS-023i |
| 🔵 BACKLOG | MS-023n | 13 | Patterns Infra Layer | Platform + SRE | MS-023m |
| 🔵 BACKLOG | MS-023o | 14 | Patterns DevSecOps + IaC | SRE | MS-023n |
| 🔵 BACKLOG | MS-023p | 14-15 | API Contracts Suite | Platform | MS-023o |
| 🔵 BACKLOG | MS-023q | 15-16 | IAM + ISO 27000-27005 | Security | MS-023p |
| 🔵 BACKLOG | MS-023r | 16 | Telemetry OTel | SRE | MS-023q |
| 🔵 BACKLOG | MS-023s | 16-17 | Local Deploy + CRUD Tests | Platform | MS-023r |
| 🔵 BACKLOG | MS-023t | 17 | Local Quality Gates | DevSecOps | MS-023s |
| 🔵 BACKLOG | MS-023u | 17-18 | Database Sync + Cross-Module | Platform | MS-023t |
| 🔵 BACKLOG | MS-023v | 18 | Integration Verification | Platform | MS-023u |
| 🔵 BACKLOG | MS-023w | 18-19 | Cross-Platform Tooling | DevSecOps | MS-023v |
| 🔵 BACKLOG | MS-023x | 19 | Pre-Commit HARD Enforcement | DevSecOps | MS-023w |

### ACTIVE (target: max 2 ativos por sprint)

| Code | Started | ETA | Owner | Progress | Blockers |
|------|---------|-----|-------|----------|----------|
| _(vazio — plano ainda em DRAFT)_ | | | | | |

### DELIVERED

| Code | Sprint | Delivered date | Closed by | Retrospective |
|------|--------|----------------|-----------|---------------|
| _(vazio — 0/26 entregues)_ | | | | |

---

## 🔥 Health Indicators

### DORA Metrics (target)

| Metric | Target | Atual | Status |
|--------|--------|-------|--------|
| **Deployment Frequency** | Daily (dev) / Weekly (prod) | N/A (pre-MVP) | 🟡 |
| **Lead Time for Changes** | < 24h dev → prod | N/A | 🟡 |
| **MTTR (Mean Time to Restore)** | < 4h CRITICAL / < 24h HIGH | N/A | 🟡 |
| **Change Failure Rate** | < 5% | N/A | 🟡 |

### Quality Gates Status (Pos-MVP)

| Gate | Threshold | Atual | Status |
|------|-----------|-------|--------|
| Domain coverage | >= 80% | N/A | 🟡 |
| Application coverage | >= 70% | N/A | 🟡 |
| SAST findings (HIGH) | 0 | N/A | 🟡 |
| Vulnerability scan (HIGH/CRITICAL) | 0 | N/A | 🟡 |
| Cosign signed images | 100% | N/A | 🟡 |
| SBOM published | 100% | N/A | 🟡 |
| SLSA Level | L3 | N/A | 🟡 |

### SLI/SLO (Production target)

| SLI | SLO | Error Budget (30d) |
|-----|-----|---------------------|
| RFQ latency p95 | < 50ms | 99.5% — budget 0.5% |
| Trade booking p95 | < 200ms | 99.5% — budget 0.5% |
| CLS submission p95 | < 500ms | 99.0% — budget 1.0% |
| PayIn deadline adherence | 100% | zero tolerance |
| API availability | 99.95% | budget 0.05% (~22min/month) |
| CLS daily cycle success | 100% | zero tolerance |
| SISCOAF filing SLA | 100% | zero tolerance |
| Audit log integrity | 100% | zero tolerance |

---

## 🚨 Risks & Blockers Dashboard

### CRITICAL (acao imediata)

| Risk | Open Question | Owner | Decision needed by |
|------|---------------|-------|---------------------|
| BACEN licenca categoria indefinida | 3c | Business | Antes Sprint 1 |
| Cloud provider escolha | 4a | Infra | Antes Sprint 1 |
| Shared CRDB hub TLS adoption | 10a | Platform | Antes Sprint 1 |

### HIGH (acao no proximo sprint)

| Risk | Mitigation | Owner |
|------|-----------|-------|
| Multi-region CRDB cost | Single-region MVP + multi-region v2 | Infra |
| AccountOS legacy migration | 6 meses backward compat documentado | Platform |
| ISO 27001 audit gap | Monthly gap analysis + Sprint 16 prep | Security |
| Pre-commit hooks > 30s | SLO monitoring + cache + impact analysis | DevSecOps |

### MEDIUM

| Risk | Mitigation |
|------|-----------|
| OTel cardinality explosion | CI validates < 10K attr sets per metric |
| Saga compensation failure | Manual review queue + 4-eyes |
| Cosign verify deps lento | Cache 30d TTL |

---

## 💰 Cost Savings Dashboard (post-MVP)

| Metric | Source | Cadence |
|--------|--------|---------|
| GitHub Actions minutos saved | Lefthook telemetry | Weekly |
| Estimated $$ saved (ubuntu @ $0.008/min) | Computed | Weekly |
| Bypass count per developer | `.git/audit-bypass.log` | Weekly |
| Local hooks executed | Lefthook telemetry | Weekly |
| GCP infra cost (per env) | GCP Billing API | Daily |
| CockroachDB Dedicated cost | CRDB Cloud API | Daily |
| Kafka MSK cost | AWS/GCP Billing | Daily |

Detail in [`08-security/iso27004-fx-security-metrics.md`](../08-security/iso27004-fx-security-metrics.md).

---

## 📈 Velocity Tracking

> Velocity = milestones delivered per sprint (target ~1.5 average)

```
Sprint:  1  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16 17 18 19
Planned: 0  1  1  1  2  1  1  2  0  2  1  4  1  2  1  3  3  3  2 → 26
Actual:  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -  -
```

Update apos cada sprint via `task dash-update`.

---

## 🎯 ISO 27001 Certification Roadmap

| Phase | Sprint | Milestone | Status |
|-------|--------|-----------|--------|
| **Foundation (ISMS scope + gap analysis)** | 1-5 | MS-023a..e | BACKLOG |
| **Domain Core (controls implementation)** | 6-9 | MS-023f..h | BACKLOG |
| **Integration (compliance suite + evidence)** | 10-15 | MS-023i..q | BACKLOG |
| **Certification (internal audit + external audit)** | 16+ | MS-023q | BACKLOG |

**Certification target:** Sprint 16 audit-ready + Sprint 17 external auditor.

---

## 🤖 Auto-Update

Este dashboard e atualizado via:

1. **Script:** `task dash-update` → executa `.claude/scripts/generate-delivery-dashboard.sh`
2. **Cron:** Hourly via GitHub Actions `.github/workflows/dashboard-update.yml`
3. **Manual:** Editar e commitar diretamente quando milestone muda status

Inputs:
- `.base/plans/milestones/{backlog,active,delivered}/*.md` — milestone status
- `.git/audit-bypass.log` — bypass count
- `.claude/memory/sessions.log` — session frequency
- GCP Billing API (when MVP+)
- Prometheus SLI/SLO metrics (when MVP+)

Outputs:
- Este arquivo (`delivery-dashboard.md`)
- `.base/plans/00-governance/audits/dashboard-snapshot-<date>.html`
- Grafana panels via API
- Slack weekly summary

---

## 📊 Live Grafana Dashboards (post-MVP)

| Dashboard | URL | Audience |
|-----------|-----|----------|
| ExchangeOS Delivery | `grafana.exchangeos.revenu.tech/d/exchangeos-delivery` | Leadership + PM |
| Sprint Velocity | `.../d/exchangeos-velocity` | Engineering manager |
| DORA Metrics | `.../d/exchangeos-dora` | SRE + Platform |
| Quality Gates | `.../d/exchangeos-quality` | DevSecOps |
| SLI/SLO Burn | `.../d/exchangeos-slo` | SRE + on-call |
| Cost Savings | `.../d/exchangeos-cost-savings` | FinOps + Leadership |

Provisioned via `docker/grafana/dashboards/exchangeos-delivery.json`.

---

## 📝 Retrospective Template

Cada milestone DELIVERED gera entry em `closures/` com retrospective:

```markdown
# Retrospective: MS-XXX

## What went well
- ...

## What didn't
- ...

## Actions
- [ ] Pattern catalog update
- [ ] ADR update if architectural decision made
- [ ] Open question resolved (X)
```

---

## 🔗 Cross-References

- [`master-plan.md`](./master-plan.md) — Full 19-sprint roadmap
- [`status-dashboard.md`](./status-dashboard.md) — Per-workstream progress
- [`../milestones/`](../milestones/) — Milestone files
- [`../milestones/index.md`](../milestones/index.md) — Milestones catalog
- [`../00-governance/open-questions.md`](../00-governance/open-questions.md) — 108 open questions
- [`../00-governance/risk-register.md`](../00-governance/risk-register.md) — 90+ risks
- [`../10-quality/sli-slo-catalog.md`](../10-quality/sli-slo-catalog.md) — SLI/SLO definitions
- [`../08-security/iso27004-fx-security-metrics.md`](../08-security/iso27004-fx-security-metrics.md) — Security metrics
