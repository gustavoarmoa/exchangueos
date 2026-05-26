// Package domain — CLS settlement bounded context.
//
// Models the CLS Bank daily PvP cycle:
//
//	07:00 CET  Cycle OPEN
//	08:00 CET  PIN1 deadline (Asia-Pacific currencies)
//	09:00 CET  PIN2 deadline (Europe currencies)
//	10:00 CET  PIN3 deadline (Americas currencies)
//	12:00 CET  Cycle CLOSED  (or FAILED)
//
// Reference: CLS Operating Procedures Manual + RN_FX_010.
package domain
