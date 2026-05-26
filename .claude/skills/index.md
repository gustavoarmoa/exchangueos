# `.claude/skills/` — Slash Commands Especializados

> **Skills** = slash commands invocaveis no Claude Code (ex: `/fx-trade-book`)
> Cada skill encapsula um fluxo recorrente do projeto

## Catalogo (6 skills iniciais)

| Skill | Slash command | Purpose |
|-------|---------------|---------|
| `fx-trade-book` | `/fx-trade-book` | Book a new FX trade end-to-end com validacoes + audit |
| `fx-pricing-test` | `/fx-pricing-test` | Run pricing golden test cases (BIS/CME/PTAX historico) |
| `bacen-compliance-check` | `/bacen-compliance-check` | Validate operacao FX against ALL BACEN rules (DEC + SCE-IED + IOF + sanctions + COS) |
| `ontology-validate` | `/ontology-validate` | Run SHACL validation completa em `.base/aasc/ontology/` |
| `integration-audit` | `/integration-audit` | Run 4-vector × 13-module integration audit (§20 do plano) |
| `cost-savings-report` | `/cost-savings-report` | Generate weekly GitHub Actions cost savings report (lefthook telemetry) |

## Skill Definition Template

```markdown
---
name: <skill-name>
description: <one-line>
allowed-tools: [Bash, Read, Edit, Grep, ...]
---

# Skill: <Name>

## Trigger

Slash command: `/<skill-name>` [args]

## Description
<o que faz>

## Workflow
1. <step 1>
2. <step 2>
...

## Examples

`/<skill-name>` → ...
`/<skill-name> <arg>` → ...
```

## Como Adicionar Nova Skill

1. Criar `.claude/skills/<name>.md` com template
2. Update este `index.md`
3. Invocar via slash command no Claude Code
