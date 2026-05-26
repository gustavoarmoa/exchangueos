# ExchangeOS — CLAUDE.local.md (Personal Overrides — Gitignored)

> **NAO commitar este arquivo.** Esta no `.gitignore`.
> Use para overrides pessoais que NAO devem chegar ao repo.

## Personal Setup

Exemplos (adicione conforme necessario):

- `MY_EMAIL=tech@revenu.com.br`
- `LOCAL_CRDB_VERSION=v24.3.32` (override do default)
- Editor preferences
- Personal slash command shortcuts
- Local-only test fixtures paths

## Personal Debugging Notes

(persistente entre sessoes — coloque insights pessoais que voce nao quer perder)

## Personal Agent Preferences

Se voce prefere certos agents para certas tarefas (override do default):

- Para BACEN: prefer agent `bacen-compliance` em vez de delegate ao `fx-domain`
- Para pricing: sempre invoke `pricing-quant` + `fx-domain` em paralelo

## Personal Permissions Override

Override de permissoes especificas em `settings.local.json` (tambem gitignored).
