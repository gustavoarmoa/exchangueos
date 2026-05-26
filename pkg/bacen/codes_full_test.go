package bacen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllNatureCodes_Populated(t *testing.T) {
	require.NotEmpty(t, AllNatureCodes, "generated catalogue must not be empty")
	assert.GreaterOrEqual(t, len(AllNatureCodes), 40, "expected at least 40 codes from CSV")
}

func TestAllNatureCodes_CategoriesCovered(t *testing.T) {
	cats := CountByCategory()
	for _, expected := range []string{
		"COMERCIAL", "SERVICOS", "CAPITAL", "TRANSFERENCIAS",
		"TURISMO", "CARTAO", "RENDA", "DERIVATIVOS", "OUTROS", "VASP",
	} {
		assert.Positive(t, cats[expected], "category %s missing from catalogue", expected)
	}
}

func TestFullByCode_KnownCodes(t *testing.T) {
	cases := []struct {
		code      string
		category  string
		direction string
	}{
		{"10000", "COMERCIAL", "INGRESSO"},
		{"10200", "COMERCIAL", "REMESSA"},
		{"20000", "CAPITAL", "INGRESSO"},
		{"70000", "DERIVATIVOS", "BIDIRECTIONAL"},
		{"90001", "VASP", "REMESSA"},
	}
	for _, c := range cases {
		t.Run(c.code, func(t *testing.T) {
			full, ok := FullByCode(c.code)
			require.True(t, ok, "code %s must be present", c.code)
			assert.Equal(t, c.category, full.Category)
			assert.Equal(t, c.direction, full.Direction)
			assert.True(t, full.Active)
		})
	}
}

func TestClassifierByCode_ResolvesViaGeneratedCatalogue(t *testing.T) {
	c := NewClassifier()

	// 20000 is in generated catalogue but NOT in legacy builtin.
	n, ok := c.ByCode("20000")
	require.True(t, ok, "must resolve via AllNatureCodes")
	assert.Equal(t, "20000", n.Code)
	assert.Equal(t, NatureIngresso, n.Nature) // INGRESSO direction → ingresso nature
	assert.Contains(t, n.Description, "Investimento")
}

func TestClassifierByCode_LegacyFallbackStillWorks(t *testing.T) {
	c := NewClassifier()

	// 60001 is legacy-only (conversao between currencies), not in generated catalogue.
	n, ok := c.ByCode("60001")
	require.True(t, ok)
	assert.Equal(t, NatureConversao, n.Nature)
}

func TestDeriveNature_BidirectionalCollapsesToConversao(t *testing.T) {
	assert.Equal(t, NatureConversao, deriveNature("BIDIRECTIONAL"))
	assert.Equal(t, NatureIngresso, deriveNature("INGRESSO"))
	assert.Equal(t, NatureRemessa, deriveNature("REMESSA"))
	assert.Equal(t, NatureRemessa, deriveNature("UNKNOWN"))
}
