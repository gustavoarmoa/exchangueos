---
name: fx-pricing-test
description: Run pricing golden test cases (BIS/CME/PTAX historical data)
allowed-tools: [Bash, Read, Grep]
---

# Skill: /fx-pricing-test

## Trigger
`/fx-pricing-test [--scope <forward|ndf|crossrate|mtm|all>]`

## Workflow
1. Invoke `pricing-quant` agent
2. Run `go test -v ./pkg/pricing/... -run TestGolden`
3. Compare results against `tests/integration/golden/pricing/*.json`
4. Validate property-based tests (gopter): inverse, round-trip, cross-rate consistency
5. Report: passing/failing per scope + edge cases

## Examples
`/fx-pricing-test` → run all 80+ golden tests
`/fx-pricing-test --scope ndf` → only NDF (BRL onshore PTAX D-1 + offshore D-2)
`/fx-pricing-test --scope crossrate` → only USD pivot triangulation
