---
name: fx-trade-book
description: Book a new FX trade end-to-end (pre-trade validations + book + post-trade)
allowed-tools: [Bash, Read, Edit, Grep, Glob]
---

# Skill: /fx-trade-book

## Trigger
`/fx-trade-book <ccy_pair> <amount> <counterparty> <type>`

## Workflow
1. **Pre-trade validation** (paralelo via agents)
   - `bacen-compliance` → DEC required check + sanctions + classificacao codigo natureza
   - `iam-security` → tenant resolution + scope check (`exchangeos:trade:write`)
   - `pricing-quant` → get current CIP forward rate via PTAX
   - `database-crdb` → AccountOS multi-CCY balance check
2. **Book trade** (`fx-domain` agent)
   - Build FXTrade aggregate via `domain.NewFXTrade(...)`
   - Validate invariants (rate > 0, value_date >= trade_date, NDF requires fixing)
   - Save via repository + outbox event
3. **Post-trade** (paralelo via agents)
   - `cls-settlement` ou `cfets-confirmation` (rota via Organisation Router)
   - `bacen-compliance` → register DEC + IOF calc + IED if applicable
   - `observability-otel` → emit trade.booked event + span
   - `kafka-flink` → publish `exchangeos.trade.booked` topic
4. **Validate audit** (`testing-qa`)
   - Audit log Merkle chain valid

## Example
`/fx-trade-book USDBRL 100000 itau spot`
