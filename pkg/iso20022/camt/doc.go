// Package camt — ISO 20022 Cash Management messages (CLS PayIn cycle + NetReport).
//
// Coverage:
//
//	camt.061.001.01  PayInScheduleV01           — CLS PayIn obligations per CCY
//	camt.062.001.01  PayInScheduleAckV01        — member ack of PayIn schedule
//	camt.063.001.01  PayInScheduleCxlV01        — schedule cancellation
//	camt.088.001.02  NetReportV02               — end-of-cycle net settlement summary
package camt

import (
	"encoding/xml"

	"github.com/shopspring/decimal"
)

// Amount — ActiveCurrencyAndAmount.
type Amount struct {
	Currency string          `xml:"Ccy,attr"`
	Value    decimal.Decimal `xml:",chardata"`
}

// ISODateTime — minimal wrapper (callers normalise to RFC3339).
type ISODateTime struct {
	Value string `xml:",chardata"`
}

// PartyID — BICFI + LEI.
type PartyID struct {
	BICFI string `xml:"BICFI,omitempty"`
	LEI   string `xml:"LEI,omitempty"`
}

// ── camt.061.001.01 — PayInSchedule ────────────────────────────────────────

const Namespace061 = "urn:iso:std:iso:20022:tech:xsd:camt.061.001.01"

type PayInScheduleV01 struct {
	XMLName    xml.Name        `xml:"urn:iso:std:iso:20022:tech:xsd:camt.061.001.01 PayInSchdl"`
	ScheduleID string          `xml:"SchdlId"`
	CycleID    string          `xml:"ClsCyclId"`
	Member     PartyID         `xml:"Mmbr"`
	GeneratedAt ISODateTime    `xml:"GnrtnDtTm"`
	Lines      []PayInLine     `xml:"PayInLine"`
}

type PayInLine struct {
	Sequence       int         `xml:"Seq"`
	Currency       string      `xml:"Ccy"`
	Amount         Amount      `xml:"Amt"`
	DeadlineDtTm   ISODateTime `xml:"DdlnDtTm"`
	DeadlineBand   string      `xml:"DdlnBnd"` // "PIN1"|"PIN2"|"PIN3" (08/09/10 CET)
}

// ── camt.062.001.01 — PayInScheduleAck ─────────────────────────────────────

const Namespace062 = "urn:iso:std:iso:20022:tech:xsd:camt.062.001.01"

type PayInScheduleAckV01 struct {
	XMLName    xml.Name    `xml:"urn:iso:std:iso:20022:tech:xsd:camt.062.001.01 PayInSchdlAck"`
	ScheduleID string      `xml:"SchdlId"`
	AckBy      PartyID     `xml:"AckBy"`
	AckAt      ISODateTime `xml:"AckDtTm"`
	Status     string      `xml:"Sts"` // ACCEPTED | REJECTED
	ReasonText string      `xml:"RsnTxt,omitempty"`
}

// ── camt.063.001.01 — PayInScheduleCxl ─────────────────────────────────────

const Namespace063 = "urn:iso:std:iso:20022:tech:xsd:camt.063.001.01"

type PayInScheduleCxlV01 struct {
	XMLName       xml.Name    `xml:"urn:iso:std:iso:20022:tech:xsd:camt.063.001.01 PayInSchdlCxl"`
	ScheduleID    string      `xml:"SchdlId"`
	CancelledBy   PartyID     `xml:"CxlBy"`
	CancelledAt   ISODateTime `xml:"CxlDtTm"`
	ReasonCode    string      `xml:"RsnCd,omitempty"`
}

// ── camt.088.001.02 — NetReport ────────────────────────────────────────────

const Namespace088 = "urn:iso:std:iso:20022:tech:xsd:camt.088.001.02"

type NetReportV02 struct {
	XMLName    xml.Name      `xml:"urn:iso:std:iso:20022:tech:xsd:camt.088.001.02 NetRpt"`
	ReportID   string        `xml:"RptId"`
	CycleID    string        `xml:"ClsCyclId"`
	Member     PartyID       `xml:"Mmbr"`
	GeneratedAt ISODateTime  `xml:"GnrtnDtTm"`
	Lines      []NetLine     `xml:"NetLine"`
}

type NetLine struct {
	Currency       string  `xml:"Ccy"`
	GrossPayIn     Amount  `xml:"GrssPayIn"`
	GrossPayOut    Amount  `xml:"GrssPayOut"`
	NetSettlement  Amount  `xml:"NetSttlm"`
	TradeCount     int     `xml:"TradCnt"`
}
