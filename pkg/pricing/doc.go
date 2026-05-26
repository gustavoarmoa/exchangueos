// Package pricing implements FX pricing algorithms used by ExchangeOS.
//
// HARD rules (enforced by .claude/rules/pkg-pricing.md):
//
//   - `shopspring/decimal` mandatory — NEVER float64/float32. golangci-lint
//     `forbidigo` rule blocks `float` declarations at lint time.
//   - 8 internal decimal places (rate precision). Display/quote rounding
//     done by callers via Quantize at boundary (4 decimals for majors, 2 for JPY).
//   - Banker's rounding (half-even) — see `decimal.RoundBank`.
//   - Calendar/day-count helpers under pkg/pricing/daycount/ (added in MS-023b).
//
// Algorithms (current + planned):
//
//	cip.go         ✅ Covered Interest Parity forward (iotafinance formula)
//	crossrate.go   ✅ Triangulation via shared pivot (auto-detect, auto-invert)
//	ndf.go         ✅ Non-Deliverable Forward (USD-cash-settled, ISDA EMTA)
//	tenor.go       ✅ ON / TN / SN / Spot / 1W..2Y ladder (Modified-Following)
//	ptax.go        ✅ BACEN PTAX 4-window survey + WeightedFixing/Bid/Ask + Fetcher interface
//	mtm.go         ✅ Mark-to-market revaluation (PositionMTM + PortfolioMTM)
//	points.go      ✅ Forward points helper (in cip.go: ForwardPoints)
//
// References cited inline:
//
//   - iotafinance Forward FX formula:
//     https://www.iotafinance.com/en/Tutorial-Forward-Foreign-Exchange.html
//   - BIS Triennial FX survey (rate conventions)
//   - CME FX product specs (point conventions)
//   - BACEN PTAX methodology (Resolução BCB 277/2022)
package pricing
