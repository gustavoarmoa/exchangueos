package fxtr

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestAmount_RoundTripPreservesPrecision(t *testing.T) {
	v, err := decimal.NewFromString("1234567.890123456789")
	if err != nil {
		t.Fatalf("decimal seed: %v", err)
	}
	in := struct {
		XMLName xml.Name `xml:"Amt"`
		Amount
	}{Amount: Amount{Currency: "USD", Value: v}}

	raw, err := xml.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(raw), `Ccy="USD"`) {
		t.Fatalf("missing Ccy attribute: %s", raw)
	}
	if !strings.Contains(string(raw), "1234567.890123456789") {
		t.Fatalf("decimal precision lost: %s", raw)
	}

	var out struct {
		XMLName xml.Name `xml:"Amt"`
		Amount
	}
	if err := xml.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Currency != "USD" {
		t.Fatalf("Ccy: got %s want USD", out.Currency)
	}
	if !out.Value.Equal(v) {
		t.Fatalf("Value: got %s want %s", out.Value, v)
	}
}

func TestRate_RoundTrip(t *testing.T) {
	in := struct {
		XMLName xml.Name `xml:"Rate"`
		Rate
	}{Rate: Rate{Value: decimal.RequireFromString("1.087654321")}}

	raw, err := xml.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out struct {
		XMLName xml.Name `xml:"Rate"`
		Rate
	}
	if err := xml.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !out.Value.Equal(in.Value) {
		t.Fatalf("Rate: got %s want %s", out.Value, in.Value)
	}
}

func TestISODate_RoundTrip(t *testing.T) {
	in := struct {
		XMLName xml.Name `xml:"Dt"`
		ISODate
	}{ISODate: ISODate{Time: time.Date(2026, 5, 24, 0, 0, 0, 0, time.UTC)}}

	raw, err := xml.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(raw), "2026-05-24") {
		t.Fatalf("expected ISO date in output: %s", raw)
	}

	var out struct {
		XMLName xml.Name `xml:"Dt"`
		ISODate
	}
	if err := xml.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !out.Time.Equal(in.Time) {
		t.Fatalf("Date: got %v want %v", out.Time, in.Time)
	}
}

func TestISODateTime_RoundTrip(t *testing.T) {
	in := struct {
		XMLName xml.Name `xml:"DtTm"`
		ISODateTime
	}{ISODateTime: ISODateTime{Time: time.Date(2026, 5, 24, 14, 30, 45, 0, time.UTC)}}

	raw, err := xml.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out struct {
		XMLName xml.Name `xml:"DtTm"`
		ISODateTime
	}
	if err := xml.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !out.Time.Equal(in.Time) {
		t.Fatalf("DateTime: got %v want %v", out.Time, in.Time)
	}
}

func TestAmount_EmptyValueDecodesAsZero(t *testing.T) {
	raw := []byte(`<Amt Ccy="EUR"></Amt>`)
	var a Amount
	if err := xml.Unmarshal(raw, &a); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !a.Value.IsZero() {
		t.Fatalf("expected zero, got %s", a.Value)
	}
	if a.Currency != "EUR" {
		t.Fatalf("expected EUR, got %s", a.Currency)
	}
}
