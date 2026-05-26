# ExchangeOS — Master Roadmap

> **Versao:** 2.0.0 | **Sprint span:** 1-19 (MS-023) + 20-23 (MS-024)
> **Status:** MS-023 ciclo DELIVERED (26/26) | MS-024 ciclo BACKLOG (13 milestones de production hardening)

## Ciclos

- **MS-023 — Planning + scaffolding + foundations** — 26 milestones, todos em `milestones/delivered/`. Conteúdo: arquitetura, contratos, domain layers, application services, ISO 20022 toolkit, pricing engine, BACEN classifier+IOF, scaffold de Helm/Terraform/ArgoCD/Argo Rollouts/SLSA L3, governança ISO 27001 + LGPD + chaos + FinOps. **Não significa rodando em produção.**
- **MS-024 — Production hardening** — 13 milestones em `milestones/backlog/`, fechando os 13 gaps identificados no review honesto da v4.19.0 (workers que faltam, integrações live, repos persistentes, cobertura de testes completa, PR cross-repo CRDB, deploy real). **Sucesso = sistema atendendo trades reais.**

## Visao

Construir o **ExchangeOS** — modulo standalone FX da Revenu Platform — em **19 sprints** cobrindo:
1. ISO 20022 FX Business Domain completa (fxtr CLS + CFETS)
2. Pricing engine CIP nativo (iotafinance)
3. BACEN Cambio 100% (Lei 14.286 + 8 Resolucoes + IOF + VASP + eFX)
4. ISO 27001:2022 certification target (93 Annex A controls)
5. Native integration com 13 modulos da plataforma
6. Cross-platform tooling (qualquer SO)
7. Pre-Commit HARD enforcement (zero GitHub Actions desperdicado)

## Timeline (Sprint × Milestone)

```
Sprint 1-2  ████░░░░░░░░░░░░░░░░  MS-023a: Foundation & Scaffolding (F1, F2)
Sprint 3    ████░░░░░░░░░░░░░░░░  MS-023b: RefData + Pricing + Quote (F3, F4P, F4)
Sprint 4    ████░░░░░░░░░░░░░░░░  MS-023c: Trade Core (F5, F6)
Sprint 5    ████████░░░░░░░░░░░░  MS-023d: Settlement CLS + non-CLS (F7, F11)
Sprint 5    ████████░░░░░░░░░░░░  MS-023d2: CFETS Capture + Confirmation (F7F)
Sprint 6    ████░░░░░░░░░░░░░░░░  MS-023e: Risk + Position + Ledger (F8, F10)
Sprint 7    ████████░░░░░░░░░░░░  MS-023f: Compliance Core + Admin (F9, F12)
Sprint 7-8  ████████████░░░░░░░░  MS-023f2: BACEN Integration Suite (F9B)
Sprint 8    ████░░░░░░░░░░░░░░░░  MS-023g: EDA E2E (F13)
Sprint 9-10 ████████████░░░░░░░░  MS-023h: Production deploy (F14, F16)
Sprint 10-11████████████░░░░░░░░  MS-023i: Allenty Documentation (F15)
Sprint 11-12████████████░░░░░░░░  MS-023j: Ontology Suite (F15B)
Sprint 12   ████░░░░░░░░░░░░░░░░  MS-023k: Flows Suite (F15C)
Sprint 12   ████░░░░░░░░░░░░░░░░  MS-023l: ERDs Suite (F15D)
Sprint 12-13████████░░░░░░░░░░░░  MS-023m: Patterns App Layer (F15E)
Sprint 13   ████░░░░░░░░░░░░░░░░  MS-023n: Patterns Infra Layer (F15F)
Sprint 14   ████░░░░░░░░░░░░░░░░  MS-023o: Patterns DevSecOps (F15G)
Sprint 14-15████████░░░░░░░░░░░░  MS-023p: API Contracts Suite (F15H)
Sprint 15-16████████░░░░░░░░░░░░  MS-023q: IAM + ISO 27000-27005 (F15I)
Sprint 16   ████░░░░░░░░░░░░░░░░  MS-023r: Telemetry OTel (F15J)
Sprint 16-17████████░░░░░░░░░░░░  MS-023s: Local Deploy + CRUD Tests (F15K)
Sprint 17   ████░░░░░░░░░░░░░░░░  MS-023t: Local Quality Gates (F15L)
Sprint 17-18████████░░░░░░░░░░░░  MS-023u: Database Sync + Cross-Module (F15M)
Sprint 18   ████░░░░░░░░░░░░░░░░  MS-023v: Integration Audit (F15N)
Sprint 18-19████████░░░░░░░░░░░░  MS-023w: Cross-Platform Tooling (F15O)
Sprint 19   ████░░░░░░░░░░░░░░░░  MS-023x: Pre-Commit Enforcement (F15P)
─── MS-023 cycle DELIVERED ───
Sprint 20   ████████░░░░░░░░░░░░  MS-024a/b/c/h/l: Infra parity (workers + repos + CRDB hub PR)
Sprint 21   ████████░░░░░░░░░░░░  MS-024d/e/i/j: Compliance correctness + test coverage
Sprint 22   ████████░░░░░░░░░░░░  MS-024f/g: BACEN + SISCOAF submission adapters
Sprint 23   ████████░░░░░░░░░░░░  MS-024m: Production deployment close-out
Background ░░░░░░░░░░░░░░░░░░░░  MS-024k: Pattern catalogue build-out (1 pattern per PR)
```

## Major Milestones (resumo)

| Milestone | Sprint | Foco |
|-----------|--------|------|
| MS-023a Foundation | 1-2 | Repo + scaffolding + proto + DB schemas + CI |
| MS-023b Pricing + Quote | 3 | CIP formula + PTAX + RFQ funcional |
| MS-023c Trade Core | 4 | Spot/Forward/Swap/NDF + amendments |
| MS-023d Settlement CLS | 5 | PvP CLS completo + 18 CCYs + non-CLS gross |
| MS-023d2 CFETS | 5 | fxtr 031-038 end-to-end |
| MS-023e Risk + Position + Ledger | 6 | NOP + limits + MTM + dual-ledger multi-CCY |
| MS-023f Compliance + Admin | 7 | DEC + sanctions + admi |
| MS-023f2 BACEN Suite | 7-8 | 95 codigos + Sistema Cambio + SISBACEN + SCE-* + SISCOAF + IOF |
| MS-023g EDA E2E | 8 | Saga completa + MQ bridge |
| MS-023h Production | 9-10 | Deploy K8s + 900+ tests + performance |
| MS-023i Documentation | 10-11 | Allenty docs completa |
| MS-023j-x Patterns + Quality + Integration | 11-19 | 850 patterns + ERDs + flows + IAM + OTel + CRUD + TDD/E2E + Sync + Audit + Cross-platform + Pre-commit enforcement |
| **MS-024a..m Production Hardening** | **20-23** | **Workers (erasure/archiver/cred-rotator) + live sanctions + full BACEN codes + DEC/SCE/SISCOAF submission + postgres repos + integration test suite + 7 E2E scenarios + pattern build-out + CRDB hub cross-repo PR + real GKE deploy with Litmus/Chaos Mesh** |

## Princip de Sequencia

1. **F1 (Foundation)** desbloqueia F2 + F14 + F15A em paralelo
2. **F2 (ISO 20022 toolkit)** desbloqueia F3-F12 (mas em ordem F3 → F4P → F4 → F5)
3. **F11 (Counterparty Adapters)** paralelo a F6-F10
4. **F13 (EDA Saga)** depende de F5+F7+F10+F11
5. **F15A-P (Documentation, Patterns, Quality)** rodam em paralelo com F2-F13
6. **F16 (Testing)** e a fase final
7. Workstreams 15B-P sao sub-fases incrementais por documentacao/patterns/quality

## Risk Highlights

- **CRITICAL:** Open question 3c (categoria licenca BACEN) — DEFINE escopo MVP
- **HIGH:** ISO 27001 certification target Sprint 16 — depende de gap analysis completo
- **HIGH:** AccountOS + PaymentOS ainda em legacy inline insecure — backward compat 6 meses
- **MEDIUM:** Multi-region CRDB custo justifica apenas com CLS production access

Ver [`risk-register.md`](../00-governance/risk-register.md) para 90+ riscos identificados.

## Status

Ver [`status-dashboard.md`](./status-dashboard.md) para tracking em tempo real.

## Sources

- §4 (Plano de Fases) + §5 (Milestones & Timeline) do [snapshot monolitico v3.11.7](../_archive/allenty-v3.11.7-monolithic-plan.md)
