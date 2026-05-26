// Package marshaller wraps canonical ISO 20022 XML marshal/unmarshal.
//
// Responsibilities:
//   - Wrap any business message in a head.001 Business Application Header (BAH)
//   - Apply canonical namespace prefixes for the registered URN
//   - Detect/reject schema-version drift on unmarshal
//
// Heavy XSD validation lives in pkg/iso20022/validator/.
package marshaller

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/revenu-tech/exchangeos/pkg/iso20022/registry"
)

// Envelope is the canonical outer document combining BAH + business message.
type Envelope struct {
	XMLName xml.Name    `xml:"urn:iso:std:iso:20022:tech:xsd:head.001.001.03 BusinessMessage"`
	Header  BAH         `xml:"AppHdr"`
	Body    interface{} `xml:",any"`
}

// BAH — minimal head.001 fields used by ExchangeOS messages.
// Full schema covered in head/head_001_001_03.go (generated from XSD).
type BAH struct {
	From          string `xml:"Fr>FIId>FinInstnId>BICFI"`
	To            string `xml:"To>FIId>FinInstnId>BICFI"`
	BizMsgIdr     string `xml:"BizMsgIdr"`
	MsgDefIdr     string `xml:"MsgDefIdr"`
	CreDt         string `xml:"CreDt"`
	Signature     string `xml:"Sgntr,omitempty"`
}

// MarshalOptions controls XML output.
type MarshalOptions struct {
	Indent     string // e.g. "  " for pretty-print; empty for production
	WriteBOM   bool   // SWIFT compatibility may require BOM
}

// Marshal serialises body with a BAH envelope using the descriptor's URN as MsgDefIdr.
func Marshal(desc registry.Descriptor, header BAH, body interface{}, opts MarshalOptions) ([]byte, error) {
	if desc.MessageDef == "" || desc.Domain == "" {
		return nil, fmt.Errorf("marshaller: descriptor missing required fields")
	}
	header.MsgDefIdr = fmt.Sprintf("%s.%s.%s.%s",
		desc.Domain, desc.MessageDef, desc.Variant, desc.Version)

	env := Envelope{Header: header, Body: body}

	var buf bytes.Buffer
	if opts.WriteBOM {
		buf.Write([]byte{0xEF, 0xBB, 0xBF})
	}
	buf.WriteString(xml.Header)

	enc := xml.NewEncoder(&buf)
	if opts.Indent != "" {
		enc.Indent("", opts.Indent)
	}
	if err := enc.Encode(env); err != nil {
		return nil, fmt.Errorf("marshaller: encode: %w", err)
	}
	if err := enc.Flush(); err != nil {
		return nil, fmt.Errorf("marshaller: flush: %w", err)
	}
	return buf.Bytes(), nil
}

// Unmarshal parses an Envelope, populating header + body. `body` must be a non-nil pointer.
// Returns the resolved Descriptor (looked up via Registry.LookupByURN on MsgDefIdr) and an error.
func Unmarshal(reg *registry.Registry, raw []byte, body interface{}) (registry.Descriptor, BAH, error) {
	if reg == nil {
		return registry.Descriptor{}, BAH{}, fmt.Errorf("marshaller: nil registry")
	}

	var env struct {
		XMLName xml.Name
		Header  BAH    `xml:"AppHdr"`
		BodyRaw []byte `xml:",innerxml"`
	}
	if err := xml.NewDecoder(bytes.NewReader(raw)).Decode(&env); err != nil && err != io.EOF {
		return registry.Descriptor{}, BAH{}, fmt.Errorf("marshaller: decode envelope: %w", err)
	}

	urn := "urn:iso:std:iso:20022:tech:xsd:" + env.Header.MsgDefIdr
	desc, ok := reg.LookupByURN(urn)
	if !ok {
		return registry.Descriptor{}, env.Header, fmt.Errorf("marshaller: unknown URN %s", urn)
	}

	if body != nil {
		if err := xml.Unmarshal(env.BodyRaw, body); err != nil {
			return desc, env.Header, fmt.Errorf("marshaller: decode body: %w", err)
		}
	}
	return desc, env.Header, nil
}
