// Package domain — Admin bounded context.
//
// Two aggregates:
//
//	SystemEvent — operational notifications mapped to admi.002/004/009/010/011/017
//	EODJob      — end-of-day batch (PTAX fixing, MTM, position snapshot, BACEN reports)
package domain
