# Performance Best Practices — ExchangeOS Claude Code

> Otimizacao para fluxos rapidos + context window eficiente + parallelismo.

## 1. Parallel Tool Batching (CRITICAL)

**Sempre que multiplas operacoes independentes:** chame em **paralelo** numa unica message (multiplos `<tool_use>` blocks).

✅ **BOM:**
```
Single message with:
- Bash(go test ./modules/trade/...)
- Bash(go test ./modules/cls_settlement/...)
- Bash(go test ./pkg/pricing/...)
- Read(.base/plans/milestones/active/MS-023b.md)
```

❌ **RUIM:**
```
Message 1: Bash(go test ./modules/trade/...)
[wait]
Message 2: Bash(go test ./modules/cls_settlement/...)
[wait]
...
```

**Gain:** 3-5x faster em multi-step tasks.

## 2. Subagent Parallel Execution

Para tarefas cross-cutting (ex: nova feature toca domain + DB + Kafka + tests), **spawn multiple subagents em paralelo** em uma message:

```
Tools: [
  Task(subagent_type=fx-domain, prompt=...),
  Task(subagent_type=database-crdb, prompt=...),
  Task(subagent_type=kafka-flink, prompt=...),
  Task(subagent_type=testing-qa, prompt=...)
]
```

**Gain:** Isolated contexts per agent (no token waste) + parallel execution.

## 3. Context Window Optimization

### Cache .base/plans/ index in context
- Master `index.md` lido 1x (5min cache)
- Per-workstream `index.md` lazy-load via glob match

### Selective reads
- NUNCA `Read` arquivo grande monolitico (10K+ linhas) sem `offset`/`limit`
- Use `Grep` com `output_mode: files_with_matches` para descoberta
- Use `Glob` antes de `Read` (find specific files first)

### Avoid re-reads
- Se acabou de editar arquivo, NAO Read de novo (edit cache valida)

## 4. Hook Performance

| Hook | SLO | Optimizations |
|------|-----|---------------|
| `SessionStart` | < 500ms | Defer heavy checks; show loading hint |
| `UserPromptSubmit` | < 100ms | Lightweight pattern match only |
| `PreToolUse` | < 200ms | Block-only patterns; no I/O slow |
| `PostToolUse` | < 300ms | Log async + non-blocking |
| `Stop`/`SessionEnd` | < 1s | Cleanup pode demorar mas non-critical |

## 5. Caching Strategy

| Cache | Path | TTL |
|-------|------|-----|
| Bash audit | `.claude/cache/bash-audit.log` | 30d (cleanupPeriodDays) |
| Test results | `~/.cache/go-test-results` | 24h |
| Build artifacts | `~/.cache/go-build` | 24h |
| Trivy DB | `~/.cache/trivy` | 1d daily refresh |
| MCP memory | `.claude/memory/mcp-memory.json` | infinite (knowledge graph) |
| Session log | `.claude/memory/sessions.log` | infinite |
| Statusline | (computed every render) | n/a — keep fast |

## 6. Tools Allowlist Optimization

Settings.json `allow` list reduz prompts por permissao (faster UX).
Tools usados frequentemente DEVEM estar em `allow` ou `ask` (vs `deny` ou unlisted).

## 7. Smart Defaults

- `permissions.defaultMode: "acceptEdits"` — auto-accept edits em files
- `cleanupPeriodDays: 30` — cache pruning automatico
- `enableAllProjectMcpServers: true` — todos MCP servers de .mcp.json auto-loaded

## 8. Subagent Tool Restrictions

Cada agent em `.claude/agents/<name>.md` define `tools: [Read, Edit, Write, Bash, Grep, Glob]` — restringe ao minimo necessario.

**Anti-pattern:** subagent com TODAS as tools (overhead + risk).

## 9. Memory Hierarchy (CLAUDE.md modular)

Use `@path/to/file` imports em CLAUDE.md para split de memory:

```markdown
# CLAUDE.md

@docs/quick-rules.md
@.base/plans/index.md
@.claude/context/glossary.md
```

**Gain:** Modular memory loadable per session need.

## 10. Specs/Contract-First Workflow

`specs/<feature>.md` = contract que mantem agents honest.
Especifica EXACTLY what to build → reduces back-and-forth → faster delivery.
