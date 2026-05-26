# ExchangeOS ERDs Suite

Per bounded-context Entity-Relationship Diagrams in Mermaid `erDiagram` format,
synchronised with `migrations/*.sql` and `.base/aasc/ontology/core/*.ttl`.

## Layout

```
.base/erds/
├── README.md
├── domain/
│   ├── erd-trade-domain.md            ✅ representative
│   ├── erd-quote-domain.md            ⏳
│   ├── erd-settlement-domain.md       ⏳
│   ├── erd-risk-position-domain.md    ⏳
│   ├── erd-compliance-admin-domain.md ⏳
│   └── erd-cfets-domain.md            ⏳
└── sql/
    ├── 01-tenants-ddl.sql             (mirror of migrations 000001)
    ├── 02-trade-ddl.sql               (mirror of 000002)
    └── ...
```

## Synchronisation rule

When `migrations/00000N_*.up.sql` changes, the matching ERD `erd-*-domain.md`
**MUST** be updated in the same PR. Lefthook `pre-commit` enforces this via
a glob check (TODO MS-023n).

## Status

- ✅ erd-trade-domain.md — representative (tenants + counterparties + fx_trades + trade_amendments + actors + audit_events)
- ⏳ 5 remaining domain ERDs + 9 sql/ DDL mirrors
