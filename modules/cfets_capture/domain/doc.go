// Package domain — CFETS Trade Capture bounded context.
//
// Lifecycle:
//
//	DRAFT → SUBMITTED → ACK         (CFETS accepted capture; deal_id assigned)
//	                 → REJECTED     (CFETS rejected; reason captured)
//	ACK → NOTIFIED (counterparty side notified; informational)
//
// Maps to fxtr.031 (submit) + fxtr.032 (ack) + fxtr.033 (notification).
package domain
