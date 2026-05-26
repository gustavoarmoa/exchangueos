// Package domain — CFETS Trade Confirmation bounded context.
//
// Lifecycle:
//
//	CONFIRMING → CONFIRMED       (CFETS returns PAIRED status)
//	          → UNPAIRED         (CFETS returns UNPAIRED status; awaiting counterparty)
//	          → REJECTED         (CFETS rejected the confirmation)
//
// Maps to fxtr.034 (request) + fxtr.035 (confirmation) + fxtr.036 (status).
package domain
