package pricing

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// NDFInput models a USD-settled Non-Deliverable Forward (NDF) per ISDA EMTA template.
//
// NDFs are used for currencies with capital controls or non-deliverable status —
// the trade settles in USD at fixing rather than physical delivery of the reference CCY.
// ExchangeOS supports the standard set: BRL, CNY, INR, KRW, RUB, TWD, IDR, MYR (RN_FX_005).
//
// Convention here: dealer is LONG the reference CCY (e.g. BRL) vs USD.
// `ContractRate` and `FixingRate` are both quoted as USD per 1 reference CCY's-worth
// (typical EMTA quoting: e.g. USDBRL = 5.10 means "1 USD costs 5.10 BRL", i.e. the
// reference CCY is BRL and the rate is the price-of-USD-in-BRL — which is the
// inverse of "USD per BRL"). To keep callers consistent with how brokers quote,
// we use BRL-per-USD form internally and apply the standard ISDA formula:
//
//	settlement_usd = notional_ref_ccy × (1/fixing_rate − 1/contract_rate)
//
// Sign:
//   - positive  → counterparty receives USD (favourable to dealer long-ref-CCY when
//     ref CCY APPRECIATES, i.e. fixing < contract)
//   - negative  → counterparty pays USD
//
// Example (BRL devalues): notional=1,000,000 BRL; contract=5.00; fixing=5.10
//
//	settlement = 1,000,000 × (1/5.10 − 1/5.00)
//	           = 1,000,000 × (0.196078431 − 0.20)
//	           ≈ −3,921.57 USD          (dealer pays — BRL depreciated)
//
// Example (BRL appreciates): notional=1,000,000 BRL; contract=5.00; fixing=4.90
//
//	settlement = 1,000,000 × (1/4.90 − 1/5.00)
//	           = 1,000,000 × (0.204081633 − 0.20)
//	           ≈ +4,081.63 USD          (dealer receives — BRL appreciated)
type NDFInput struct {
	NotionalReferenceCCY decimal.Decimal // notional in the non-deliverable CCY (e.g. BRL units)
	ContractRate         decimal.Decimal // ref-ccy per 1 USD at trade date (e.g. USDBRL = 5.00)
	FixingRate           decimal.Decimal // ref-ccy per 1 USD at fixing date
}

// SettlementAmount returns the USD cash payment per the ISDA EMTA NDF formula:
//
//	settlement = notional × (1/fixing − 1/contract)
//
// Returns ErrInvalidInput / ErrZeroDenominator on bad inputs.
// Result is rounded with banker's rounding at `internalScale` decimal places (8).
func SettlementAmount(in NDFInput) (decimal.Decimal, error) {
	if err := in.validate(); err != nil {
		return decimal.Zero, err
	}
	one := decimal.NewFromInt(1)
	invFixing := one.Div(in.FixingRate)
	invContract := one.Div(in.ContractRate)
	diff := invFixing.Sub(invContract)
	return in.NotionalReferenceCCY.Mul(diff).RoundBank(internalScale), nil
}

// validate enforces NDF input invariants.
func (in NDFInput) validate() error {
	if !in.NotionalReferenceCCY.IsPositive() {
		return fmt.Errorf("%w: notional_reference_ccy must be > 0", ErrInvalidInput)
	}
	if !in.ContractRate.IsPositive() {
		return fmt.Errorf("%w: contract_rate must be > 0", ErrInvalidInput)
	}
	if !in.FixingRate.IsPositive() {
		return fmt.Errorf("%w: fixing_rate must be > 0", ErrInvalidInput)
	}
	return nil
}
