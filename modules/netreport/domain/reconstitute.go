package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ReconstituteNetReport hydrates a NetReport from persisted state.
//
// USAGE: persistence boundary only. Bypasses NewNetReport validation +
// derived-field calculation (netSettlement) because the row stores both
// the gross sides AND the precomputed net. Caller is responsible for
// passing a consistent triplet.
func ReconstituteNetReport(
	id, tenantID, cycleID uuid.UUID,
	currency string,
	grossPayIn, grossPayOut, netSettlement decimal.Decimal,
	tradeCount int,
	generatedAt time.Time,
) *NetReport {
	return &NetReport{
		id:            id,
		tenantID:      tenantID,
		cycleID:       cycleID,
		currency:      currency,
		grossPayIn:    grossPayIn,
		grossPayOut:   grossPayOut,
		netSettlement: netSettlement,
		tradeCount:    tradeCount,
		generatedAt:   generatedAt.UTC(),
	}
}
