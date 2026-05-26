package fxtr

import "encoding/xml"

// Namespace013 — pinned schema namespace for fxtr.013.001.04.
const Namespace013 = "urn:iso:std:iso:20022:tech:xsd:fxtr.013.001.04"

// FXTradePositionStatementV04 — fxtr.013.001.04.
// CLS variant: periodic statement of open trade positions per (member, currency) pair.
type FXTradePositionStatementV04 struct {
	XMLName xml.Name `xml:"urn:iso:std:iso:20022:tech:xsd:fxtr.013.001.04 FXTradPosSttmnt"`

	StatementID    string      `xml:"StmtId"`
	MemberBIC      string      `xml:"MmbBICFI"`
	StatementDate  ISODate     `xml:"StmtDt"`
	GeneratedAt    ISODateTime `xml:"GnrtnDtTm"`

	// Per-currency aggregate positions (long/short legs + net).
	Positions []PositionLine `xml:"Pos"`
}

// PositionLine — one currency-pair position summary.
type PositionLine struct {
	Currency        string `xml:"Ccy"`
	LongAmount      Amount `xml:"LngAmt,omitempty"`
	ShortAmount     Amount `xml:"ShrtAmt,omitempty"`
	NetAmount       Amount `xml:"NetAmt"`
	OpenTradeCount  int    `xml:"OpenTradCnt"`
}
