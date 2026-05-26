// Package domain — Quote/RFQ bounded context.
//
// Two collaborating aggregates:
//
//   - RFQ — Request-For-Quote initiated by a trader; transitions
//     REQUESTED → QUOTED → ACCEPTED|REJECTED|EXPIRED.
//   - Quote — a streamable price (bid/ask) with a validity window.
//
// Accepted quotes hand off to trade via a domain event consumed by the
// trade application layer (no cross-aggregate transaction).
package domain
