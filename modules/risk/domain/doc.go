// Package domain — Risk bounded context.
//
// Models per-tenant trading limits (counterparty / currency / tenor / DV01 / VaR).
// Pre-trade limit checks return (allowed, breached_list, explanation).
// Cite RN_FX_015 — NOP monitored realtime; halt if exceeds BCB cap.
package domain
