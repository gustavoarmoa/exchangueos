package domain

import (
	"time"

	"github.com/google/uuid"
)

// Reconstitute helpers — used by infrastructure repositories to rebuild aggregates
// from persisted state WITHOUT re-running the constructor's validation (which would
// reject inactive/legacy rows the database is legally allowed to hold).
//
// These are NOT part of the domain protocol — they exist solely for the persistence
// boundary and are documented as such. Application code MUST use the NewXxx
// constructors to create aggregates.

// ReconstituteCurrency rebuilds a Currency from persisted values.
func ReconstituteCurrency(code, name string, minorUnits int, cls, cfets, active bool) *Currency {
	return &Currency{
		code:          code,
		name:          name,
		minorUnits:    minorUnits,
		clsEligible:   cls,
		cfetsEligible: cfets,
		active:        active,
	}
}

// ReconstituteCalendar rebuilds a Calendar from id + holiday list.
func ReconstituteCalendar(id string, holidays []time.Time) *Calendar {
	c := &Calendar{id: id, holidays: make(map[string]struct{}, len(holidays))}
	for _, h := range holidays {
		c.holidays[keyOf(h)] = struct{}{}
	}
	return c
}

// ReconstituteBIC rebuilds a BICRecord.
func ReconstituteBIC(bic, institutionName, country, lei string, active bool) *BICRecord {
	return &BICRecord{
		bic:             bic,
		institutionName: institutionName,
		country:         country,
		lei:             lei,
		active:          active,
	}
}

// ReconstituteSSI rebuilds an SSI with a stable id (UUID from DB row).
func ReconstituteSSI(
	id, tenantID uuid.UUID,
	counterpartyBIC, currency, beneficiaryBIC, intermediaryBIC string,
	accountNumber, iban string,
	validFrom, validTo time.Time,
) *SSI {
	return &SSI{
		id:              id,
		tenantID:        tenantID,
		counterpartyBIC: counterpartyBIC,
		currency:        currency,
		beneficiaryBIC:  beneficiaryBIC,
		intermediaryBIC: intermediaryBIC,
		accountNumber:   accountNumber,
		iban:            iban,
		validFrom:       validFrom.UTC(),
		validTo:         validTo.UTC(),
	}
}
