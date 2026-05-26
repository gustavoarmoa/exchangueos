package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ReconstituteLimit rebuilds a Limit from persisted state. Persistence-boundary only.
func ReconstituteLimit(
	id, tenantID uuid.UUID,
	limitType LimitType,
	scope, currency string,
	cap, utilised decimal.Decimal,
	version int,
) *Limit {
	return &Limit{
		id:        id,
		tenantID:  tenantID,
		limitType: limitType,
		scope:     scope,
		cap:       cap,
		utilised:  utilised,
		currency:  currency,
		version:   version,
	}
}
