# 09 — Compliance

> **Workstream:** Compliance
> **Versao:** 1.0.0
> **Status:** DRAFT

## Conteudo

| Documento | Status | Descricao |
|-----------|--------|-----------|
| `bacen-cambio-coverage.md` | TODO | BACEN Regulatory Coverage 100% (§2.6 monolitico) — Lei 14.286 + 8 Resolucoes + Circulares + IOF + VASP + eFX |
| `lei-14286-novo-marco-cambial.md` | TODO | Lei nº 14.286/2021 (Novo Marco Legal) detalhado |
| `resolucoes-bcb.md` | TODO | Resolucoes BCB 277/278/279/280/281/337/348/539/561 detalhadas |
| `circular-3978-pldft.md` | TODO | Circular BCB 3.978/2020 (PLD/FT) — COAF + SISCOAF |
| `circular-3690-classificacao.md` | TODO | Circular BCB 3.690/2013 — 95 codigos de natureza |
| `iof-decreto-12499.md` | TODO | Decreto 12.499/2025 — Aliquotas IOF 2025-2026 |
| `vasp-ativos-virtuais.md` | TODO | Regulacao Ativos Virtuais (vigencia 02/02/2026) — VASP/PSAV |
| `efx-resolucao-561.md` | TODO | Resolucao 561/2026 — eFX (vigencia 01/10/2026) |
| `pix-internacional-nexus.md` | TODO | PIX Internacional Project Nexus (BIS Innovation Hub) — roadmap |
| `rmcci-articles.md` | TODO | RMCCI (Regulamento Mercado Cambio + Capitais Internacionais) artigos |
| `pims-lgpd.md` | TODO | LGPD + ISO 27018 PII protection (`pims/`) |
| `bacen/` | TODO | BACEN-specific docs (sce-ied, sce-credito, sce-cbe manuals) |
| `certification/` | TODO | Certification path docs |
| `legal/` | TODO | Legal references |
| `pims/` | TODO | PIMS (Privacy Information Management System) |

## BACEN Marco Cambial Atual

| Norma | Vigencia | Escopo |
|-------|----------|--------|
| **Lei nº 14.286/2021** | 31/12/2022 | Novo Marco Cambial |
| **Resolucao BCB nº 277/2022** | 31/12/2022 | Mercado de cambio + instituicoes autorizadas + eFX |
| **Resolucao BCB nº 278/2022** | 31/12/2022 | Credito externo + IED |
| **Resolucao BCB nº 279/2022** | 31/12/2022 | Capital brasileiro no exterior (CBE) |
| **Resolucao BCB nº 280/2022** | 31/12/2022 | Definicao de residente vs nao residente |
| **Resolucao BCB nº 281/2022** | 31/12/2022 | Disposicoes transitorias |
| **Resolucao BCB nº 337/2023** | 22/08/2023 | Alteracoes em Res 277 |
| **Resolucao BCB nº 348/2023** | 01/11/2023 | Cambio simbolico revogacoes |
| **Resolucao BCB nº 539/2025** | 18/12/2025 | Atualizacoes complementares |
| **Resolucao BCB nº 561/2026** | **01/10/2026** | eFX restricao a instituicoes autorizadas + USD 10k limit |
| **Regulacao Ativos Virtuais (VASP)** | **02/02/2026** | PSAV no mercado cambio + USD 100k limit + self-custody bans |
| **Decreto nº 12.499/2025** | 2025 | Aliquotas IOF cambio |
| **Circular BCB nº 3.690/2013** | atualizada 11/2023 | 95 codigos de natureza |
| **Circular BCB nº 3.978/2020** | 23/01/2020 | PLD/FT — COAF + SISCOAF |

## Sistemas BACEN Integrados

- **Sistema Cambio** — registro contrato cambio
- **SISBACEN** — Information System
- **SCE-IED** (sucessor RDE-IED) — Investimento Estrangeiro Direto
- **SCE-Credito** (sucessor RDE-ROF) — Credito Externo
- **SCE-CBE** — Capitais Brasileiros no Exterior
- **SISCOAF** — COS (Comunicacao Operacao Suspeita)
- **SISCOMEX** — Comercio Exterior
- **OLINDA API** — PTAX + dados publicos

## Sources

- §2.6 (BACEN Regulatory Coverage 100%) do [snapshot monolitico v3.11.7](../_archive/allenty-v3.11.7-monolithic-plan.md)
- Pattern de referencia: [LedgerOS 09-compliance](../../../../ledgeros/.base/plans/09-compliance/)
