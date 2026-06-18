package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Reconstitute helpers — persistence boundary only. They bypass the validation
// performed by the Newxxx constructors because the row is presumed to have been
// validated when first persisted. Never call from application logic.

// ReconstituteClassification hydrates a Classification from a row.
func ReconstituteClassification(
	id, tenantID, tradeID uuid.UUID,
	code, description string,
	nature Nature,
) *Classification {
	return &Classification{
		id:          id,
		tenantID:    tenantID,
		tradeID:     tradeID,
		code:        code,
		description: description,
		nature:      nature,
	}
}

// ReconstituteIOFComputation hydrates an IOFComputation from a row.
func ReconstituteIOFComputation(
	id, tenantID, tradeID uuid.UUID,
	operationType string,
	notional decimal.Decimal,
	notionalCCY string,
	rate, iofAmount decimal.Decimal,
	computedAt time.Time,
) *IOFComputation {
	return &IOFComputation{
		id:            id,
		tenantID:      tenantID,
		tradeID:       tradeID,
		operationType: operationType,
		notional:      notional,
		notionalCCY:   notionalCCY,
		rate:          rate,
		iofAmount:     iofAmount,
		computedAt:    computedAt.UTC(),
	}
}

// ReconstituteBACENReport hydrates a BACENReport from a row.
func ReconstituteBACENReport(
	id, tenantID uuid.UUID,
	reportType ReportType,
	referenceDate time.Time,
	payloadHash string,
	status ReportStatus,
	submittedAt, respondedAt time.Time,
	rejectionReason string,
	version int,
) *BACENReport {
	return &BACENReport{
		id:              id,
		tenantID:        tenantID,
		reportType:      reportType,
		referenceDate:   referenceDate.UTC(),
		payloadHash:     payloadHash,
		status:          status,
		submittedAt:     submittedAt,
		respondedAt:     respondedAt,
		rejectionReason: rejectionReason,
		version:         version,
	}
}

// ReconstituteScreeningResult hydrates a ScreeningResult from a row.
func ReconstituteScreeningResult(
	id, tenantID uuid.UUID,
	counterpartyBIC, lei string,
	hits []string,
	riskLevel RiskLevel,
	screenedAt time.Time,
) *ScreeningResult {
	return &ScreeningResult{
		id:              id,
		tenantID:        tenantID,
		counterpartyBIC: counterpartyBIC,
		lei:             lei,
		hits:            append([]string(nil), hits...),
		riskLevel:       riskLevel,
		screenedAt:      screenedAt.UTC(),
	}
}
