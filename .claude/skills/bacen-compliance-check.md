---
name: bacen-compliance-check
description: Validate operacao FX contra TODAS as regras BACEN (DEC + SCE-* + IOF + sanctions + COS)
allowed-tools: [Bash, Read, Grep]
---

# Skill: /bacen-compliance-check

## Trigger
`/bacen-compliance-check <trade_id_or_request_yaml>`

## Workflow (paralelo via `bacen-compliance` agent)
1. Validate classificacao codigo natureza (95 codigos Circ 3.690)
2. Check DEC required (> USD 10K per Lei 14.286/2021)
3. Calculate IOF (Decreto 12.499/2025 — 6 aliquotas)
4. Sanctions screen (OFAC SDN + ONU + BCB lista + FATF jurisdicoes risco)
5. Check SCE-IED registration trigger (if IED inbound)
6. Check SCE-Credito registration (if credito externo)
7. Check SCE-CBE thresholds (if residente com ativos exterior)
8. VASP rules check (USD 100k + self-custody bans)
9. eFX limit check (USD 10k para IPs)
10. Generate compliance report (PASS/FAIL per rule + required actions)
