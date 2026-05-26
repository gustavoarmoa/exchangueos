package registry

import (
	"strings"
	"testing"
)

func TestDefaultRegistry_HasExactly32Schemas(t *testing.T) {
	r := Default()
	got := len(r.List())
	want := 32
	if got != want {
		t.Fatalf("Default() registered %d schemas, want %d", got, want)
	}
}

func TestDefaultRegistry_CoversCLSAndCFETS(t *testing.T) {
	r := Default()
	if n := len(r.FilterByOrganisation(OrgCLS)); n < 20 {
		t.Errorf("CLS variants: got %d, want >= 20", n)
	}
	if n := len(r.FilterByOrganisation(OrgCFETS)); n != 8 {
		t.Errorf("CFETS variants: got %d, want exactly 8", n)
	}
}

func TestDescriptor_URN(t *testing.T) {
	d := Descriptor{
		Organisation: OrgCLS,
		Domain:       BusinessFXTR,
		MessageDef:   "014",
		Variant:      "001",
		Version:      "05",
		XSDSourceURL: "https://example/x.xsd",
	}
	want := "urn:iso:std:iso:20022:tech:xsd:fxtr.014.001.05"
	if got := d.URN(); got != want {
		t.Fatalf("URN: got %s, want %s", got, want)
	}
}

func TestRegistry_DuplicateRejected(t *testing.T) {
	r := New()
	d := Descriptor{
		Organisation: OrgISO, Domain: BusinessHEAD, MessageDef: "001", Variant: "001", Version: "03",
		XSDSourceURL: "https://example/h.xsd",
	}
	if err := r.Register(d); err != nil {
		t.Fatalf("first register: %v", err)
	}
	if err := r.Register(d); err == nil || !strings.Contains(err.Error(), "already registered") {
		t.Fatalf("expected duplicate err, got %v", err)
	}
}

func TestRegistry_RejectsInvalidOrganisation(t *testing.T) {
	r := New()
	d := Descriptor{
		Organisation: "BOGUS", Domain: BusinessHEAD, MessageDef: "001", Variant: "001", Version: "03",
		XSDSourceURL: "https://example/h.xsd",
	}
	if err := r.Register(d); err == nil {
		t.Fatal("expected unknown-organisation error")
	}
}

func TestRegistry_LookupByURN(t *testing.T) {
	r := Default()
	d, ok := r.LookupByURN("urn:iso:std:iso:20022:tech:xsd:fxtr.014.001.05")
	if !ok {
		t.Fatal("LookupByURN miss")
	}
	if d.Organisation != OrgCLS {
		t.Fatalf("organisation: got %s want CLSBUS33", d.Organisation)
	}
}

func TestOrganisationRouter(t *testing.T) {
	router := NewOrganisationRouter([]string{"DEUTDEFF", "CHASUS33"}, "")

	tests := []struct {
		name    string
		bic     string
		country string
		want    Organisation
		wantErr bool
	}{
		{"CLS direct", "CLSBUS33", "US", OrgCLS, false},
		{"CLS member", "DEUTDEFF", "DE", OrgCLS, false},
		{"CLS member case-insensitive", "chasus33", "US", OrgCLS, false},
		{"CFETS by prefix", "CFETSCN00", "CN", OrgCFETS, false},
		{"CFETS by country", "ICBKCNBJ", "CN", OrgCFETS, false},
		{"ISO fallback", "ITAUBRSP", "BR", OrgISO, false},
		{"empty BIC", "", "BR", "", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := router.RouteParty(tc.bic, tc.country)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err=%v wantErr=%v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Fatalf("got %s want %s", got, tc.want)
			}
		})
	}
}
