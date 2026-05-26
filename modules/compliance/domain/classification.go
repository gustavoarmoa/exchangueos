package domain

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Nature — direction of the FX operation per BACEN classification.
type Nature string

const (
	NatureRemessa   Nature = "REMESSA"   // outflow (BRL → foreign CCY)
	NatureIngresso  Nature = "INGRESSO"  // inflow (foreign CCY → BRL)
	NatureConversao Nature = "CONVERSAO" // FX between two foreign CCYs
)

// Classification — BACEN nature code attached to a trade.
// Cite RN_FX_028 (95 codes per Circ 3.690).
type Classification struct {
	id          uuid.UUID
	tradeID     uuid.UUID
	tenantID    uuid.UUID
	code        string // 4-digit BACEN nature code (e.g. "32101", "63010")
	description string
	nature      Nature
}

// NewClassificationInput parameterises construction.
type NewClassificationInput struct {
	TenantID    uuid.UUID
	TradeID     uuid.UUID
	Code        string
	Description string
	Nature      Nature
}

// NewClassification validates and constructs.
func NewClassification(in NewClassificationInput) (*Classification, error) {
	if in.TenantID == uuid.Nil || in.TradeID == uuid.Nil {
		return nil, fmt.Errorf("%w: tenant_id + trade_id required", ErrInvalidInput)
	}
	code := strings.TrimSpace(in.Code)
	if len(code) < 4 || len(code) > 6 {
		return nil, fmt.Errorf("%w: code must be 4-6 digits, got %q", ErrInvalidInput, in.Code)
	}
	for _, r := range code {
		if r < '0' || r > '9' {
			return nil, fmt.Errorf("%w: code must be numeric, got %q", ErrInvalidInput, in.Code)
		}
	}
	switch in.Nature {
	case NatureRemessa, NatureIngresso, NatureConversao:
	default:
		return nil, fmt.Errorf("%w: nature must be REMESSA|INGRESSO|CONVERSAO, got %q", ErrInvalidInput, in.Nature)
	}
	if in.Description == "" {
		return nil, fmt.Errorf("%w: description required", ErrInvalidInput)
	}
	return &Classification{
		id:          uuid.New(),
		tenantID:    in.TenantID,
		tradeID:     in.TradeID,
		code:        code,
		description: in.Description,
		nature:      in.Nature,
	}, nil
}

func (c *Classification) ID() uuid.UUID       { return c.id }
func (c *Classification) TenantID() uuid.UUID { return c.tenantID }
func (c *Classification) TradeID() uuid.UUID  { return c.tradeID }
func (c *Classification) Code() string        { return c.code }
func (c *Classification) Description() string { return c.description }
func (c *Classification) Nature() Nature      { return c.nature }
