package fxtr

import "encoding/xml"

// Namespace017 — pinned schema namespace for fxtr.017.001.05.
const Namespace017 = "urn:iso:std:iso:20022:tech:xsd:fxtr.017.001.05"

// FXTradeStatusReportV05 — fxtr.017.001.05.
// CLS variant: per-trade status update (e.g. PAIRED, SETTLED, RESCINDED, REJECTED).
type FXTradeStatusReportV05 struct {
	XMLName xml.Name `xml:"urn:iso:std:iso:20022:tech:xsd:fxtr.017.001.05 FXTradStsRpt"`

	// Original trade reference.
	OriginalTradeIdentification TradeIdentification14 `xml:"OrgnlTradId"`

	// Status itself.
	Status         FXTradeStatusCode `xml:"Sts"`
	StatusReason   string            `xml:"StsRsn,omitempty"`
	StatusDateTime ISODateTime       `xml:"StsDtTm"`

	// Optional cycle context.
	CLSCycleID string `xml:"ClsCyclId,omitempty"`
}

// FXTradeStatusCode — enumerated CLS lifecycle states (subset; full XSD has more).
type FXTradeStatusCode string

const (
	StatusReceived  FXTradeStatusCode = "RCVD"
	StatusPaired    FXTradeStatusCode = "PAIR"
	StatusUnpaired  FXTradeStatusCode = "UPRD"
	StatusSettling  FXTradeStatusCode = "STGN" // settling
	StatusSettled   FXTradeStatusCode = "SETT"
	StatusRescinded FXTradeStatusCode = "RESC"
	StatusRejected  FXTradeStatusCode = "REJT"
)
