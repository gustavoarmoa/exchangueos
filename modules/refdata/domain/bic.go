package domain

import (
	"fmt"
	"strings"
)

// BICRecord captures the ISO 9362 BIC plus identifying metadata for a financial institution.
type BICRecord struct {
	bic             string
	institutionName string
	country         string
	lei             string
	active          bool
}

// NewBICRecord constructs and validates a BICRecord.
//
// BIC structure (ISO 9362):
//
//	chars 1-4:  bank prefix (alpha)
//	chars 5-6:  ISO 3166 country code (alpha)
//	chars 7-8:  location code (alphanum)
//	chars 9-11: optional branch code (alphanum) — when present length = 11
func NewBICRecord(bic, institutionName, country, lei string) (*BICRecord, error) {
	bic = strings.ToUpper(strings.TrimSpace(bic))
	country = strings.ToUpper(strings.TrimSpace(country))
	switch len(bic) {
	case 8, 11:
	default:
		return nil, fmt.Errorf("%w: BIC must be 8 or 11 chars, got %d", ErrInvalidInput, len(bic))
	}
	if !isAlpha(bic[:4]) {
		return nil, fmt.Errorf("%w: BIC prefix must be alpha", ErrInvalidInput)
	}
	if !isAlpha(bic[4:6]) {
		return nil, fmt.Errorf("%w: BIC country segment must be alpha", ErrInvalidInput)
	}
	if !isAlnum(bic[6:8]) {
		return nil, fmt.Errorf("%w: BIC location segment must be alphanumeric", ErrInvalidInput)
	}
	if len(bic) == 11 && !isAlnum(bic[8:11]) {
		return nil, fmt.Errorf("%w: BIC branch segment must be alphanumeric", ErrInvalidInput)
	}
	if institutionName == "" {
		return nil, fmt.Errorf("%w: institution_name required", ErrInvalidInput)
	}
	if len(country) != 2 {
		return nil, fmt.Errorf("%w: country must be ISO 3166 alpha-2", ErrInvalidInput)
	}
	if lei != "" && len(lei) != 20 {
		return nil, fmt.Errorf("%w: LEI must be 20 chars (ISO 17442) when present", ErrInvalidInput)
	}
	return &BICRecord{
		bic:             bic,
		institutionName: institutionName,
		country:         country,
		lei:             lei,
		active:          true,
	}, nil
}

func (b *BICRecord) BIC() string             { return b.bic }
func (b *BICRecord) InstitutionName() string { return b.institutionName }
func (b *BICRecord) Country() string         { return b.country }
func (b *BICRecord) LEI() string             { return b.lei }
func (b *BICRecord) IsActive() bool          { return b.active }
func (b *BICRecord) Deactivate()             { b.active = false }
func (b *BICRecord) Activate()               { b.active = true }

func isAlpha(s string) bool {
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return len(s) > 0
}

func isAlnum(s string) bool {
	for _, r := range s {
		if (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return false
		}
	}
	return len(s) > 0
}
