package erasure_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/revenu-tech/exchangeos/internal/erasure"
)

const validPlan = `{
  "ticket": "LGPD-2026-0001",
  "subject_ref": "11111111-2222-3333-4444-555555555555",
  "approvals": ["dpo", "compliance_officer"],
  "operations": [
    {
      "table": "actors",
      "where": "id = '11111111-2222-3333-4444-555555555555'",
      "op": "redact",
      "fields": ["name", "email", "tax_id"]
    },
    {
      "table": "quote_streams",
      "where": "requester_id = '11111111-2222-3333-4444-555555555555'",
      "op": "hard_delete"
    }
  ]
}`

func TestParsePlan_Valid(t *testing.T) {
	p, err := erasure.ParsePlan([]byte(validPlan))
	require.NoError(t, err)
	assert.Equal(t, "LGPD-2026-0001", p.Ticket)
	assert.Len(t, p.Operations, 2)
	assert.True(t, p.HasRequiredApprovals())
	assert.Equal(t, "[REDACTED PER LGPD ART 18 IV LGPD-2026-0001]", p.RedactionMarker())
}

func TestParsePlan_RejectsBadTicket(t *testing.T) {
	_, err := erasure.ParsePlan([]byte(`{
		"ticket": "WRONG-0001",
		"subject_ref": "x",
		"operations": [{"table":"actors","where":"id = 'x'","op":"redact","fields":["name"]}]
	}`))
	require.Error(t, err)
	assert.ErrorIs(t, err, erasure.ErrInvalidPlan)
}

func TestParsePlan_RejectsEmptyWhere(t *testing.T) {
	_, err := erasure.ParsePlan([]byte(`{
		"ticket": "LGPD-2026-0002",
		"subject_ref": "x",
		"operations": [{"table":"actors","where":"","op":"redact","fields":["name"]}]
	}`))
	require.Error(t, err)
	assert.ErrorIs(t, err, erasure.ErrInvalidPlan)
	assert.Contains(t, err.Error(), "where required")
}

func TestParsePlan_RejectsRedactWithoutFields(t *testing.T) {
	_, err := erasure.ParsePlan([]byte(`{
		"ticket": "LGPD-2026-0003",
		"subject_ref": "x",
		"operations": [{"table":"actors","where":"id = 'x'","op":"redact","fields":[]}]
	}`))
	require.Error(t, err)
	assert.ErrorIs(t, err, erasure.ErrInvalidPlan)
}

func TestParsePlan_RejectsHardDeleteWithFields(t *testing.T) {
	_, err := erasure.ParsePlan([]byte(`{
		"ticket": "LGPD-2026-0004",
		"subject_ref": "x",
		"operations": [{"table":"quote_streams","where":"requester_id = 'x'","op":"hard_delete","fields":["name"]}]
	}`))
	require.Error(t, err)
	assert.ErrorIs(t, err, erasure.ErrInvalidPlan)
}

func TestParsePlan_RejectsUnknownOp(t *testing.T) {
	_, err := erasure.ParsePlan([]byte(`{
		"ticket": "LGPD-2026-0005",
		"subject_ref": "x",
		"operations": [{"table":"actors","where":"id = 'x'","op":"truncate"}]
	}`))
	require.Error(t, err)
	assert.True(t, errors.Is(err, erasure.ErrInvalidPlan))
}

func TestHasRequiredApprovals_MissingOne(t *testing.T) {
	p := &erasure.Plan{
		Ticket:     "LGPD-2026-0006",
		SubjectRef: "x",
		Approvals:  []string{"dpo"},
		Operations: []erasure.Operation{{Table: "x", Where: "1=1", Op: erasure.OpHardDelete}},
	}
	assert.False(t, p.HasRequiredApprovals())
}

func TestHasRequiredApprovals_CaseInsensitive(t *testing.T) {
	p := &erasure.Plan{
		Ticket:     "LGPD-2026-0007",
		SubjectRef: "x",
		Approvals:  []string{"  DPO ", "Compliance_Officer"},
		Operations: []erasure.Operation{{Table: "x", Where: "1=1", Op: erasure.OpHardDelete}},
	}
	assert.True(t, p.HasRequiredApprovals())
}
