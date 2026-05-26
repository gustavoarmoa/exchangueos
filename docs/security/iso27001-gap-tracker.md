# ISO 27001:2022 Gap Closure Tracker

> Companion to `iso27001-controls-mapping.md`. Lists every 🟡 partial + ⏳ deferred control
> with owner + action + target date. Reviewed monthly by Security + Compliance.

Last reviewed: 2026-05-24
Cert audit target: Sprint 16 (TBD calendar date)
Target gap-close: ≥ 90% implemented before audit

## 🟡 Partial (5)

| # | Control | Gap | Owner | Action | ETA |
|---|---------|-----|-------|--------|-----|
| A.5.1 | Information security policies | Formal ISMS policy doc absent (cited inline in CLAUDE.md) | Security Officer | Author standalone ISMS Policy + sign by CISO | 2026-Q3 |
| A.5.4 | Management responsibilities | Sign-off table empty in go-live-checklist.md | Platform Lead | Populate on first go-live execution | First prod deploy |
| A.5.7 | Threat intelligence | Reactive scans only; no proactive feeds | Security Officer | Subscribe to BACEN security bulletins + Mandiant FinServ feed | 2026-Q3 |
| A.5.10 | Acceptable use of information | AUP cited in CLAUDE.md not standalone | HR + Security | Author standalone AUP signed at onboarding | 2026-Q3 |
| A.8.11 | Data masking | Dev/staging carry real PII shape (synthetic seeds) | Platform | Add data-masking pipeline for backups before dev refresh | 2026-Q4 |

## ⏳ Deferred (18)

### People controls (HR-owned)

| # | Control | Owner | Notes |
|---|---------|-------|-------|
| A.6.1 | Screening | HR | Pre-employment checks — outside module |
| A.6.2 | Terms and conditions of employment | HR | Standard employment contract template |
| A.6.3 | Information security awareness | HR + Security | Annual mandatory training program |
| A.6.4 | Disciplinary process | HR | Existing HR process |
| A.6.5 | Responsibilities after termination | HR + Security | Offboarding playbook including cert rotation + account deactivation |
| A.6.6 | Confidentiality / NDA | Legal | DPA / NDA templates per role |
| A.6.7 | Remote working | HR | Workforce remote-work policy |
| A.5.19..A.5.23 | Supplier relationships | Procurement | Vendor due diligence + DPA reviews — 5 controls |

### Module-scoped deferred

| # | Control | Owner | Action | ETA |
|---|---------|-------|--------|-----|
| A.5.35 | Independent review | Security + External | Annual ISO 27001 audit engagement | Sprint 16 |
| A.8.1 | User endpoint devices | IT | Endpoint MDM (outside this module's scope) | — |
| A.8.23 | Web filtering | Platform networking | Egress allowlist via Cloud NAT | 2026-Q3 |
| A.8.34 | Protection of information systems during audit testing | Security | Document audit-window procedures (read-only Vault token, etc.) | Pre-audit |

## Roll-up by quarter

| Quarter | Targeted | Closure pace required |
|---------|----------|----------------------|
| 2026-Q3 | 4 (A.5.1, A.5.7, A.5.10, A.8.23) | 1.3/month |
| 2026-Q4 | 2 (A.8.11 + cert prep wrap) | 0.7/month |
| Pre-audit | Independent review engagement (A.5.35) | Critical path |

## Review cadence

- **Monthly:** Security + Compliance leads update this tracker, move items between 🟡/⏳/✅
- **Quarterly:** Platform Lead reviews aggregate; presents to CISO
- **Pre-audit (T-90):** All 🟡 controls promoted to ✅ OR explicitly accepted as residual risk by CISO sign-off
