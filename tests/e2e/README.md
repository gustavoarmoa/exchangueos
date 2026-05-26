# ExchangeOS E2E Test Catalog

10 canonical end-to-end scenarios covering the full FX flow.

Build tag: `//go:build e2e` — run via `task test:e2e`. Each scenario uses
docker-compose stack (CRDB + Kafka + OTel + exchangeos-api) and waits for
state via `require.Eventually` (NEVER `time.Sleep`).

## Scenarios

| # | Scenario | Bounded contexts | Expected |
|---|----------|------------------|----------|
| 1 | **EUR/USD spot booking** | refdata → pricing → quote → trade | Trade SETTLED via CLS PIN2 |
| 2 | **USD/BRL NDF settlement** | refdata → quote → trade → settlement (bilateral) | NDF cash settlement in USD |
| 3 | **CFETS CNY trade capture** | cfets_capture (fxtr.031→032→033) | CaptureID + CFETSDealID assigned |
| 4 | **CFETS confirmation pairing** | cfets_confirmation (fxtr.034→035→036) | Status PAIRED |
| 5 | **Risk limit breach pre-trade** | risk → trade | Trade rejected; limit utilised unchanged |
| 6 | **Position update after settle** | trade → position | Net position matches trade net |
| 7 | **CLS daily cycle 07:00→12:00** | cls_settlement → payin → netreport | NetReport per CCY at close |
| 8 | **BACEN classification + IOF** | compliance → bacen | Classification code + IOF amount persisted |
| 9 | **Sanctions screening blocks trade** | compliance → trade | Trade rejected; ScreeningResult HIGH + RequiresCOS |
| 10 | **EOD batch full pipeline** | admin EOD → PTAX → MTM → BACEN report | EOD job COMPLETED with all 4 steps |

## Common helpers

- `tests/e2e/harness.go` — docker-compose lifecycle (Start/Stop/WaitHealthy).
- `tests/e2e/fixtures/*.json` — sample tenant + counterparty seeds.
- `tests/e2e/assertions.go` — Eventually wrappers for trade state / DB rows / Kafka topic lag.

## Running locally

```bash
task compose:up
task test:e2e
task compose:down
```

CI: same set runs nightly + on `main` merge via `.github/workflows/e2e.yml` (added in MS-023g closing).
