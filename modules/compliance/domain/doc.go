// Package domain — Compliance bounded context.
//
// Four aggregates:
//
//	Classification    — operation nature code (Circ 3.690, 95 codes)
//	IOFComputation    — tax computation (Decreto 12.499/2025, 6 rates)
//	BACENReport       — SISBACEN/CCS/CAMBIO submission tracker
//	ScreeningResult   — OFAC/UN/EU/COAF sanctions screening outcome
//
// Cite: RN_FX_028 (95 nature codes), RN_FX_037 (IOF), RN_FX_039 (SISCOAF COS).
package domain
