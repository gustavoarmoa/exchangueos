package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ReconstitutePosition rebuilds a Position from persisted state. Persistence-boundary only.
func ReconstitutePosition(
	id, tenantID uuid.UUID,
	currency string,
	long, short decimal.Decimal,
	asOf time.Time,
	version int,
) *Position {
	return &Position{
		id:       id,
		tenantID: tenantID,
		currency: currency,
		long:     long,
		short:    short,
		asOf:     asOf.UTC(),
		version:  version,
	}
}
