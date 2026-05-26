# ExchangeOS Seeds

Idempotent SQL seeds executed (in numeric order) by `cmd/migrator seed`
after migrations are applied.

## Order matters

| # | File | Purpose | Depends on |
|---|------|---------|------------|
| 00 | `00_tenants_dev.sql` | Dev tenant (deterministic UUID) | migration 000001 |
| 01 | `01_currencies.sql` | 30 ISO 4217 currencies (18 CLS-eligible) | migration 000005 |
| 02 | `02_currency_pairs.sql` | 32 currency pairs | 01 + migration 000004 |
| 03 | `03_calendars.sql` | 6 calendars + 2026 holidays (BACEN/NYFR/BOE/TARGET2/TOKYO/TORONTO) | migration 000005 |
| 04 | `04_counterparties.sql` | 36 BIC records (CLS members + Brazilian banks) | migration 000005 |
| 05 | `05_ssi.sql` | 5 sample SSIs for dev tenant | 00 + 01 + 04 + migration 000005 |
| 06 | `06_netting_cutoffs.sql` | 24 CLS PayIn deadlines + bilateral cutoffs | migration 000004 |

## Idempotence

All seeds use `ON CONFLICT DO NOTHING` so they can be re-applied safely.

## Running

```bash
task db:seed
```

## Source-of-truth caveats

Calendars in `03_calendars.sql` are illustrative for 2026 testing only. Production
must source from authoritative feeds (BACEN, FRB, BOE, ECB) at year-end refresh.
LEIs in `04_counterparties.sql` are samples; production should pull from GLEIF.
