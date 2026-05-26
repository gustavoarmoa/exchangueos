package fxtr

import "encoding/xml"

// Namespace016 — pinned schema namespace for fxtr.016.001.05.
const Namespace016 = "urn:iso:std:iso:20022:tech:xsd:fxtr.016.001.05"

// FXTradeCancellationConfirmationV05 — fxtr.016.001.05.
// CLS variant: confirms cancellation of a previously captured (and not-yet-settled) trade.
//
// Lifecycle constraint: cancellation is forbidden once trade.status == SETTLING (see
// modules/trade/domain.FXTrade.Cancel) — emitting this message for a settling/settled
// trade will be rejected at the boundary.
type FXTradeCancellationConfirmationV05 struct {
	XMLName xml.Name `xml:"urn:iso:std:iso:20022:tech:xsd:fxtr.016.001.05 FXTradCxlConf"`

	// Original trade reference (links back to the fxtr.014 that established the trade).
	OriginalTradeIdentification TradeIdentification14 `xml:"OrgnlTradId"`

	// Side identification carried forward for audit.
	TradingSideIdentification      SideIdentification `xml:"TradgSdIdr"`
	CounterpartySideIdentification SideIdentification `xml:"CtrPtySdIdr"`

	// Cancellation metadata.
	CancellationReason string      `xml:"CxlRsn,omitempty"`
	CancelledAt        ISODateTime `xml:"CxlDtTm"`

	// Optional audit / 4-eyes (cancellations above threshold typically require 4-eyes).
	AuditTrail *AuditTrail `xml:"AdtTrl,omitempty"`
}
