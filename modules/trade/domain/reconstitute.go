package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ReconstituteFXTrade rebuilds an FXTrade from persisted values (no validation).
// Use ONLY in infrastructure repositories. Application code must call NewFXTrade.
func ReconstituteFXTrade(
	id, tenantID uuid.UUID,
	externalRef string,
	tradeType TradeType,
	status TradeStatus,
	venue SettlementVenue,
	buyerBIC, sellerBIC string,
	boughtCurrency string, boughtAmount decimal.Decimal,
	soldCurrency string, soldAmount decimal.Decimal,
	dealRate decimal.Decimal,
	tradeDate, valueDate time.Time,
	version int,
) *FXTrade {
	return &FXTrade{
		id:             id,
		tenantID:       tenantID,
		externalRef:    externalRef,
		tradeType:      tradeType,
		status:         status,
		venue:          venue,
		buyerBIC:       buyerBIC,
		sellerBIC:      sellerBIC,
		boughtCurrency: boughtCurrency,
		boughtAmount:   boughtAmount,
		soldCurrency:   soldCurrency,
		soldAmount:     soldAmount,
		dealRate:       dealRate,
		tradeDate:      tradeDate.UTC(),
		valueDate:      valueDate.UTC(),
		version:        version,
	}
}
