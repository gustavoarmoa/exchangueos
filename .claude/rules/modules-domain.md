---
glob: "modules/**/domain/**/*.go"
---

# Rule: modules/*/domain/*.go (loaded on glob match)

## TDD Obrigatorio
Toda mudanca em domain layer requer test que falha primeiro (Red), implementacao minima (Green), refactor (Refactor).

## Invariants
- Aggregate root e o unico ponto de entrada
- Reference outras aggregates por ID (NUNCA por ponteiro)
- 1 aggregate = 1 transacao DB
- Zero imports de infrastructure

## Required Patterns
- Construtor `NewXxx(...) (*Xxx, error)` com validation
- Method receivers em pointer para aggregate root
- `version` field para optimistic concurrency
- Domain events via `RecordEvent(e)` (flushed via outbox)

## Money & Rates
- `decimal.Decimal` obrigatorio (NUNCA float64)
- Specifications composable (And/Or/Not)

## Cite Source
- Adicionar/modificar RN_FX_NNN: cite no comment + add SHACL shape em `.base/aasc/ontology/compliance/bacen-cambio-shapes.ttl`
