# Business Rules Quick Reference — RN_FX_001..050

> 50 business rules canonicos do ExchangeOS. Codificados em SHACL em .base/aasc/ontology/compliance/bacen-cambio-shapes.ttl

## Trade & Pricing (RN_FX_001..026)

| Code | Rule |
|------|------|
| RN_FX_001 | Currency pair valid (refdata + ACTIVE) |
| RN_FX_002 | Spot T+2 (default; USD/CAD T+1; USD/MXN T+1) |
| RN_FX_005 | NDF requires fixing source (PTAX/WMR/ECB) + fixing date |
| RN_FX_010 | PvP via CLS para 18 CCYs elegiveis |
| RN_FX_013 | Amendment > USD 100k requires 4-eyes |
| RN_FX_015 | NOP monitored realtime; halt se exceder limite BCB |
| RN_FX_017 | SSI obrigatoria pre-first-settlement |
| RN_FX_021 | Forward via CIP no-arbitrage (iotafinance) |
| RN_FX_026 | NUNCA float64 (decimal.Decimal obrigatorio) |

## BACEN (RN_FX_027..050)

| Code | Rule |
|------|------|
| RN_FX_027 | Operacao apenas por instituicao autorizada BCB |
| RN_FX_028 | Codigo natureza valido (95 codigos Circ 3.690) |
| RN_FX_029 | eFX limit USD 10k investimento (Res 561, 01/10/2026) |
| RN_FX_030 | VASP limit USD 100k (02/02/2026) |
| RN_FX_031 | Proibido VASP transfer para self-custody non-resident |
| RN_FX_034 | IED inbound: register SCE-IED em 30 dias |
| RN_FX_037 | IOF auto-calculado (6 aliquotas Decreto 12.499/2025) |
| RN_FX_039 | COS para SISCOAF ate 1 dia util apos decisao |
| RN_FX_046 | Penalidade ate R$ 250k por descumprimento |

Ver full list em .base/plans/02-core-domain/business-rules.md
