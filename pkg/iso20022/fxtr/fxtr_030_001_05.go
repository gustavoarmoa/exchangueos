package fxtr

import "encoding/xml"

// Namespace030 — pinned schema namespace for fxtr.030.001.05.
const Namespace030 = "urn:iso:std:iso:20022:tech:xsd:fxtr.030.001.05"

// FXSettlementNotificationV05 — fxtr.030.001.05.
// CLS variant: emitted when PvP settlement legs complete (terminal-state notification).
type FXSettlementNotificationV05 struct {
	XMLName xml.Name `xml:"urn:iso:std:iso:20022:tech:xsd:fxtr.030.001.05 FXSttlmNtfctn"`

	// Original trade reference.
	OriginalTradeIdentification TradeIdentification14 `xml:"OrgnlTradId"`

	// Settlement metadata.
	SettlementReference string      `xml:"SttlmRef"`
	SettledAt           ISODateTime `xml:"SttlmDtTm"`
	CLSCycleID          string      `xml:"ClsCyclId,omitempty"`

	// Final settled amounts (echoed for ledger reconciliation).
	BoughtSettled Amount `xml:"BghtSttldAmt"`
	SoldSettled   Amount `xml:"SldSttldAmt"`

	// Optional payment instrument identifiers (CLS-internal references).
	PayInBoughtRef string `xml:"PayInBghtRef,omitempty"`
	PayInSoldRef   string `xml:"PayInSldRef,omitempty"`
}
