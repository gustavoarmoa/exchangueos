// Package domain — PayIn bounded context (CLS daily PayIn cycle).
//
// PayInInstruction lifecycle:
//
//	PENDING → SUBMITTED → CONFIRMED (settled by counterparty before deadline)
//	                    → FAILED    (missed deadline / rejected)
//
// Each instruction belongs to a CLSCycle (by cycle_id) and a deadline band
// (PIN1/PIN2/PIN3) per netting_cutoffs.
package domain
