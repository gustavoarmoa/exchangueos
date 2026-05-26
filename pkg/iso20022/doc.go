// Package iso20022 hosts ExchangeOS's ISO 20022 FX message toolkit.
//
// Coverage (32 FX-specific schemas):
//
//	fxtr (15) — CLS variants: 008, 013, 014, 015, 016, 017, 030
//	         — CFETS variants: 031, 032, 033, 034, 035, 036, 037, 038
//	admi  (6) — CLS system events: 002, 004, 009, 010, 011, 017
//	camt  (4) — CLS PayIn cycle: 061, 062, 063  +  NetReport: 088
//	reda  (2) — CLS settlement reference data: 060, 061
//	supporting (5) — head.001 (BAH), reda.066 (Calendar), reda.067 (SSI),
//	                head.002 (BAH ack), admi.024 (notification)
//
// IMPORTANT: ISO 20022 publishes ONLY `fxtr` for FX trade messages. There is
// NO `fxti` and NO `fxmt` namespace. Quote/Amendment are internal-only Revenu
// services (gRPC) that translate to fxtr.014/015/016 (CLS) or fxtr.031/035/036
// (CFETS) on the boundary.
//
// Submodules:
//
//	registry/    Version Registry — schema URN → version → struct factory
//	             Organisation Router — CLSBUS33 vs CFETS dispatcher
//	validator/   XSD + business-rule (SHACL-aligned) validators
//	marshaller/  Canonical XML marshal/unmarshal with namespace handling
//	fxtr/        Per-message Go structs (one file per CLS+CFETS variant)
//	admi/        Per-message Go structs (CLS system events)
//	camt/        Per-message Go structs (PayIn + NetReport)
//	reda/        Per-message Go structs (SSI + Calendar)
//
// XSDs are pinned by version in registry/sources.go and downloaded at build time
// (see scripts/download-xsd.sh) — never embedded as binary blobs in this package.
package iso20022
