package bacen

//go:generate go run ./codegen --input ../../data/bacen/nature-codes-circ-3690-v20260101.csv --output ./codes_full.go

import (
	"fmt"
	"strings"
)

// Nature mirrors compliance.domain.Nature without an import cycle.
type Nature string

const (
	NatureRemessa   Nature = "REMESSA"
	NatureIngresso  Nature = "INGRESSO"
	NatureConversao Nature = "CONVERSAO"
)

// NatureCode is one row of the BACEN catalog.
type NatureCode struct {
	Code        string
	Description string
	Nature      Nature
}

// builtin is the seed set (~20 most-common codes). Production deployments load
// the full 95-code list from refdata at boot; this seed is sufficient for tests
// + smoke runs.
var builtin = []NatureCode{
	// Mercadorias (goods)
	{"10001", "Exportação de mercadorias", NatureIngresso},
	{"10002", "Importação de mercadorias", NatureRemessa},
	{"10005", "Devolução de exportação", NatureRemessa},
	{"10006", "Devolução de importação", NatureIngresso},

	// Serviços (services)
	{"20001", "Receita de serviços técnicos", NatureIngresso},
	{"20002", "Pagamento de serviços técnicos", NatureRemessa},
	{"20010", "Royalties — recebimento", NatureIngresso},
	{"20011", "Royalties — pagamento", NatureRemessa},

	// Capital (capital flows)
	{"30001", "Investimento estrangeiro direto — ingresso", NatureIngresso},
	{"30002", "Retorno de IED — remessa", NatureRemessa},
	{"30010", "Empréstimo externo — ingresso", NatureIngresso},
	{"30011", "Empréstimo externo — pagamento de principal", NatureRemessa},
	{"30012", "Empréstimo externo — pagamento de juros", NatureRemessa},

	// Transferências unilaterais (transfers)
	{"40001", "Manutenção de residentes — remessa", NatureRemessa},
	{"40002", "Manutenção de residentes — ingresso", NatureIngresso},

	// Turismo + cartão (travel + cards)
	{"50001", "Viagens internacionais — turismo", NatureRemessa},
	{"50002", "Cartão de crédito internacional", NatureRemessa},

	// Outros
	{"60001", "Conversão entre moedas estrangeiras", NatureConversao},
	{"63010", "Operação financeira — derivativo", NatureRemessa},
	{"99999", "Outros — código residual", NatureRemessa},
}

// Classifier resolves nature codes by code or by free-text hint.
type Classifier struct {
	byCode map[string]NatureCode
	// keyword index — first-match wins (order matters).
	keywords []keywordRule
}

type keywordRule struct {
	keyword string
	code    string
}

// NewClassifier constructs a Classifier seeded with the builtin catalog plus
// optional `extra` rows (from refdata).
func NewClassifier(extra ...NatureCode) *Classifier {
	c := &Classifier{
		byCode: make(map[string]NatureCode, len(builtin)+len(extra)),
		keywords: []keywordRule{
			{"export", "10001"},
			{"import", "10002"},
			{"service", "20002"},
			{"royalt", "20011"},
			{"investment", "30001"},
			{"loan", "30010"},
			{"interest", "30012"},
			{"travel", "50001"},
			{"card", "50002"},
			{"cross", "60001"},
			{"derivative", "63010"},
		},
	}
	for _, n := range builtin {
		c.byCode[n.Code] = n
	}
	for _, n := range extra {
		c.byCode[n.Code] = n
	}
	return c
}

// ByCode returns the NatureCode for an exact code.
//
// Resolution order:
//  1. Generated catalogue (AllNatureCodes — full 46-code BACEN Circ 3.690 set),
//     compressed to the slim NatureCode shape via deriveNature().
//  2. Constructor-supplied extra rows + builtin seed (legacy path).
// Returns false on miss.
func (c *Classifier) ByCode(code string) (NatureCode, bool) {
	code = strings.TrimSpace(code)
	if full, ok := AllNatureCodes[code]; ok && full.Active {
		return NatureCode{
			Code:        full.Code,
			Description: full.DescriptionPT,
			Nature:      deriveNature(full.Direction),
		}, true
	}
	n, ok := c.byCode[code]
	return n, ok
}

// deriveNature maps the generated catalogue's Direction (INGRESSO/REMESSA/BIDIRECTIONAL)
// to the slim Nature type. BIDIRECTIONAL collapses to CONVERSAO (commercial intent
// matches the swap/conversion class of operations in legacy callers).
func deriveNature(direction string) Nature {
	switch direction {
	case "INGRESSO":
		return NatureIngresso
	case "REMESSA":
		return NatureRemessa
	case "BIDIRECTIONAL":
		return NatureConversao
	default:
		return NatureRemessa
	}
}

// Classify resolves a nature from a free-text hint (case-insensitive substring match).
// Returns ErrUnknown if no rule matches.
func (c *Classifier) Classify(hint string) (NatureCode, error) {
	h := strings.ToLower(strings.TrimSpace(hint))
	if h == "" {
		return NatureCode{}, fmt.Errorf("%w: hint is empty", ErrUnknown)
	}
	for _, rule := range c.keywords {
		if strings.Contains(h, rule.keyword) {
			return c.byCode[rule.code], nil
		}
	}
	return NatureCode{}, fmt.Errorf("%w: no rule for %q", ErrUnknown, hint)
}

// All returns the seeded set (in catalog order).
func (c *Classifier) All() []NatureCode {
	out := make([]NatureCode, 0, len(c.byCode))
	for _, n := range builtin {
		out = append(out, n)
	}
	return out
}
