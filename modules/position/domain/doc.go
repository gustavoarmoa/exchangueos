// Package domain — Position bounded context.
//
// Tracks per (tenant, currency) net open position aggregated from trades.
// Updated incrementally via UpdateFromTrade (reacts to trade.created events).
package domain
