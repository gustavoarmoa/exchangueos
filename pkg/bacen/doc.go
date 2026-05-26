// Package bacen — BACEN regulatory utilities.
//
// Contents:
//
//	classifier.go  — Nature-code classifier (Circ 3.690, 95 codes). Implements the
//	                 most common 20 codes; the long tail (rare/historical codes) is
//	                 loaded from refdata on demand.
//	iof.go         — IOFCalculator with 6 rates per Decreto 12.499/2025.
//	dec.go         — DEC submission helper (deferred; envelope structure documented).
//
// Cite: RN_FX_028 (95 codes), RN_FX_037 (IOF rates), RN_FX_039 (SISCOAF COS).
package bacen
