# Glossary — ExchangeOS Ubiquitous Language

> Quick reference cacheable. Loaded via @context/glossary.md em CLAUDE.md.

## FX Domain

- **Spot** — FX trade settled em T+2 (default; T+1 USD/CAD)
- **Forward** — FX trade settled em T+N futuro
- **Swap** — 2 trades simultaneos (near leg + far leg)
- **NDF (Non-Deliverable Forward)** — cash-settled em USD (BRL/CNY/INR/KRW)
- **Pip** — smallest unit (4 decimals for majors, 2 for JPY)
- **NOP (Net Open Position)** — exposicao realtime per CCY pair
- **MTM (Mark-to-Market)** — revaluation EOD via fixing source
- **CIP (Covered Interest Parity)** — no-arbitrage forward formula
- **PTAX** — BACEN daily fixing rate (4 windows survey 10/11/12/13h SP)

## Standards

- **ISO 20022 fxtr** — Foreign Exchange Trade messages (15 messages CLS + CFETS)
- **CLS** — Continuous Linked Settlement Bank (BIC: CLSBUS33, 18 CCYs eligible)
- **CFETS** — China Foreign Exchange Trade System (PTPP)
- **PvP** — Payment vs Payment (atomic settlement)
- **SWIFT MT** — Legacy financial messaging (MT300/MT304/MT202)

## BACEN

- **DEC** — Declaracao Eletronica de Cambio (Lei 14.286/2021, obrigatoria > USD 10K)
- **SCE-IED** — Sistema Capitais Estrangeiros — Investimento Estrangeiro Direto
- **SCE-Credito** — Sucessor RDE-ROF (credito externo)
- **SCE-CBE** — Capitais Brasileiros no Exterior
- **SISCOAF** — COS (Comunicacao Operacao Suspeita) — PLD/FT
- **IOF** — Imposto Operacoes Financeiras (cambio: 6 aliquotas Decreto 12.499/2025)
- **PTAX** — BACEN fixing rate
- **VASP/PSAV** — Virtual Asset Service Provider (regulacao 02/02/2026)
- **eFX** — Electronic Foreign Exchange (Res 561/2026, USD 10K limit IPs)
- **Lei 14.286/2021** — Novo Marco Cambial
- **Res BCB 277-561** — Resolucoes operacionais cambio

## Infrastructure

- **CRDB** — CockroachDB (shared hub TLS pattern: cockroachdb/modules/exchangeos/)
- **GKE Autopilot** — Google Kubernetes Engine managed
- **WIF** — Workload Identity Federation (zero JSON keys)
- **CMEK** — Customer-Managed Encryption Keys (Cloud KMS HSM)
- **mTLS** — Mutual TLS (cert-based authn)
- **OIDC** — OpenID Connect (Keycloak)
- **Vault SPI** — Keycloak Service Provider Interface para Vault secrets

## Standards-Absorbed

- **FIBO** — Financial Industry Business Ontology (canonical base)
- **ISO 27001:2022** — ISMS Requirements (93 Annex A controls)
- **SLSA L3** — Supply chain Level 3 (provenance + attestation + non-falsifiable)
- **OWL 2 DL** — Web Ontology Language Description Logic profile

## Internal

- **Allenty** — Revenu enterprise planning framework
- **RN_FX_NNN** — Business rule code (50 rules total)
- **RFLW.024.NNN.NN** — Flow code (Revenu FLoW, domain 024 = ExchangeOS)
- **FX-* patterns** — 850 patterns em 20 catalogs
- **BC (Bounded Context)** — 14 BCs em ExchangeOS
