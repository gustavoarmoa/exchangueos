# 04 — DSL & Compiler

> **Workstream:** DSL & Compiler
> **Versao:** 1.0.0
> **Status:** OUT-OF-SCOPE MVP — Roadmap futuro (v2)

## Escopo

DSL para descrever operacoes FX de forma declarativa (ex: book trade + register DEC + post ledger em uma sintaxe unificada). **Out-of-scope no MVP** — implementacao depende de DSL framework consolidado do LedgerOS (MS-044 ja delivered).

## Roadmap Futuro (v2+)

| Documento | Status | Descricao |
|-----------|--------|-----------|
| `fx-dsl-syntax.md` | FUTURE | Sintaxe DSL FX (sample: `trade USDBRL spot 100K with itau via cls`) |
| `parser.md` | FUTURE | Parser ANTLR ou peg |
| `ast.md` | FUTURE | AST nodes |
| `codegen.md` | FUTURE | Go code generation (handler + saga + tests) |
| `dsl-to-ontology-bridge.md` | FUTURE | Bridge para ontology IRIs |

## Dependencia

- LedgerOS MS-044 (Ontology 100% Implementation) — **delivered Sprint Q3 2026**
- DSL framework do LedgerOS reutilizado

## Sources

- Pattern de referencia: [LedgerOS 04-dsl-compiler](../../../../ledgeros/.base/plans/04-dsl-compiler/)
- ADR-049 (LedgerOS): Domain Bridge Pattern para projection DSL ↔ ontology
