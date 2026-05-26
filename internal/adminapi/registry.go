// Package adminapi — admin REST endpoints for inspecting + mutating every
// ExchangeOS table in the local stack. Gated by EXCHANGEOS_ENABLE_ADMIN_API
// (default false). Production NEVER enables this without the JWT scope check.
//
// Phase 2 — read-only LIST + GET endpoints under /v1/admin/{table}.
// Phase 3 — POST/PUT/DELETE per aggregate (separate file).
//
// SECURITY: this is local-dev tooling. The query/filter mechanism is column-
// allowlisted (Schema.AllowedFilters) so no SQL injection is possible via
// query params. Tenant scoping defaults to the dev tenant unless overridden by
// the X-Tenant-Id header.
package adminapi

// Schema describes one exposed table.
type Schema struct {
	// URL — path segment under /v1/admin/{URL}.
	URL string
	// Table — SQL identifier (must match a real CRDB table; we don't quote).
	Table string
	// PK — primary key column name. Used by GET /{url}/{id}.
	PK string
	// TenantColumn — if non-empty, queries auto-filter by `WHERE <col> = $tenant`
	// when the request carries X-Tenant-Id. Empty for global-scope tables.
	TenantColumn string
	// AllowedFilters — query param keys that are accepted; map maps query-param
	// → SQL column. e.g. `{"status": "status"}`. Anything else is silently ignored.
	AllowedFilters map[string]string
	// DefaultOrder — `ORDER BY` clause body (no "ORDER BY" prefix). Optional.
	DefaultOrder string
	// Mutable — when false, POST/PUT/DELETE return 405. audit_events is read-only.
	Mutable bool
}

// schemas is the source-of-truth for which tables are exposed. Maps URL → Schema.
var schemas = map[string]Schema{
	// ── refdata (mutable but rare) ──────────────────────────────────────────
	"tenants": {
		URL: "tenants", Table: "tenants", PK: "tenant_id",
		AllowedFilters: map[string]string{"status": "status", "country": "country", "code": "code"},
		DefaultOrder:   "code ASC",
		Mutable:        true,
	},
	"actors": {
		URL: "actors", Table: "actors", PK: "actor_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"status": "status", "type": "type"},
		DefaultOrder:   "display_name ASC",
		Mutable:        true,
	},
	"currencies": {
		URL: "currencies", Table: "currencies", PK: "code",
		AllowedFilters: map[string]string{"active": "active", "cls_eligible": "cls_eligible", "cfets_eligible": "cfets_eligible"},
		DefaultOrder:   "code ASC",
		Mutable:        true,
	},
	"currency-pairs": {
		// PK is composite (base_ccy, quote_ccy). For LIST + filter the composite
		// is fine; GET/PUT/DELETE by id is best-effort using base_ccy alone
		// (returns the first matching pair). For per-pair mutation, prefer the
		// dedicated currency_pairs domain endpoints when they land in MS-024.
		URL: "currency-pairs", Table: "currency_pairs", PK: "base_ccy",
		AllowedFilters: map[string]string{"active": "active", "cls_eligible": "cls_eligible", "cfets_eligible": "cfets_eligible", "base_ccy": "base_ccy", "quote_ccy": "quote_ccy"},
		DefaultOrder:   "base_ccy ASC, quote_ccy ASC",
		Mutable:        true,
	},
	"calendars": {
		URL: "calendars", Table: "calendars", PK: "calendar_id",
		AllowedFilters: map[string]string{},
		DefaultOrder:   "calendar_id ASC",
		Mutable:        true,
	},
	"calendar-holidays": {
		// Composite PK (calendar_id, holiday_date). Same caveat as currency_pairs.
		URL: "calendar-holidays", Table: "calendar_holidays", PK: "calendar_id",
		AllowedFilters: map[string]string{"calendar_id": "calendar_id"},
		DefaultOrder:   "calendar_id ASC, holiday_date ASC",
		Mutable:        true,
	},
	"bic-records": {
		URL: "bic-records", Table: "bic_records", PK: "bic",
		AllowedFilters: map[string]string{"active": "active", "country": "country"},
		DefaultOrder:   "bic ASC",
		Mutable:        true,
	},
	"ssis": {
		URL: "ssis", Table: "ssis", PK: "ssi_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"currency": "currency", "counterparty_bic": "counterparty_bic"},
		DefaultOrder:   "valid_from DESC",
		Mutable:        true,
	},
	"netting-cutoffs": {
		// Composite PK (venue, currency, band). Same caveat as currency_pairs.
		URL: "netting-cutoffs", Table: "netting_cutoffs", PK: "venue",
		AllowedFilters: map[string]string{"venue": "venue", "currency": "currency", "band": "band"},
		DefaultOrder:   "venue ASC, currency ASC, band ASC",
		Mutable:        true,
	},
	"counterparties": {
		URL: "counterparties", Table: "counterparties", PK: "counterparty_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"status": "status", "country": "country", "bic": "bic", "cls_member": "cls_member", "cfets_member": "cfets_member"},
		DefaultOrder:   "bic ASC",
		Mutable:        true,
	},

	// ── domain aggregates (mutable via app services in Phase 3) ─────────────
	"fx-trades": {
		URL: "fx-trades", Table: "fx_trades", PK: "trade_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"status": "status", "trade_type": "trade_type", "settlement_venue": "settlement_venue", "external_ref": "external_ref"},
		DefaultOrder:   "trade_date DESC",
		Mutable:        true,
	},
	"trade-amendments": {
		URL: "trade-amendments", Table: "trade_amendments", PK: "amendment_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"status": "status", "trade_id": "trade_id", "change_type": "change_type"},
		DefaultOrder:   "proposed_at DESC",
		Mutable:        true,
	},
	"rfqs": {
		URL: "rfqs", Table: "rfqs", PK: "rfq_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"status": "status", "base_ccy": "base_ccy", "quote_ccy": "quote_ccy"},
		DefaultOrder:   "created_at DESC",
		Mutable:        true,
	},
	"quotes": {
		URL: "quotes", Table: "quotes", PK: "quote_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"rfq_id": "rfq_id", "base_ccy": "base_ccy", "quote_ccy": "quote_ccy", "venue": "venue"},
		DefaultOrder:   "valid_to DESC",
		Mutable:        true,
	},
	"quote-streams": {
		URL: "quote-streams", Table: "quote_streams", PK: "stream_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{},
		DefaultOrder:   "started_at DESC",
		Mutable:        true,
	},
	"cls-cycles": {
		URL: "cls-cycles", Table: "cls_cycles", PK: "cycle_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"status": "status", "cycle_date": "cycle_date"},
		DefaultOrder:   "cycle_date DESC",
		Mutable:        true,
	},
	"cls-cycle-trades": {
		// Composite PK (cycle_id, trade_id). GET/PUT/DELETE by id uses cycle_id
		// alone; use the trade_id filter on LIST to disambiguate.
		URL: "cls-cycle-trades", Table: "cls_cycle_trades", PK: "cycle_id",
		AllowedFilters: map[string]string{"cycle_id": "cycle_id", "trade_id": "trade_id"},
		DefaultOrder:   "attached_at DESC",
		Mutable:        true,
	},
	"payin-instructions": {
		URL: "payin-instructions", Table: "payin_instructions", PK: "instruction_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"status": "status", "currency": "currency", "band": "band", "cycle_id": "cycle_id"},
		DefaultOrder:   "deadline ASC",
		Mutable:        true,
	},
	"net-reports": {
		URL: "net-reports", Table: "net_reports", PK: "report_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"currency": "currency", "cycle_id": "cycle_id"},
		DefaultOrder:   "generated_at DESC",
		Mutable:        true,
	},
	"risk-limits": {
		URL: "risk-limits", Table: "risk_limits", PK: "limit_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"limit_type": "limit_type", "scope": "scope", "currency": "currency"},
		DefaultOrder:   "limit_type ASC, scope ASC",
		Mutable:        true,
	},
	"positions": {
		URL: "positions", Table: "positions", PK: "position_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"currency": "currency"},
		DefaultOrder:   "currency ASC",
		Mutable:        true,
	},
	"classifications": {
		URL: "classifications", Table: "classifications", PK: "classification_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"code": "code", "nature": "nature", "trade_id": "trade_id"},
		DefaultOrder:   "created_at DESC",
		Mutable:        true,
	},
	"iof-computations": {
		URL: "iof-computations", Table: "iof_computations", PK: "iof_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"operation_type": "operation_type", "trade_id": "trade_id"},
		DefaultOrder:   "computed_at DESC",
		Mutable:        true,
	},
	"bacen-reports": {
		URL: "bacen-reports", Table: "bacen_reports", PK: "report_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"status": "status", "report_type": "report_type", "reference_date": "reference_date"},
		DefaultOrder:   "reference_date DESC",
		Mutable:        true,
	},
	"screening-results": {
		URL: "screening-results", Table: "screening_results", PK: "screening_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"risk_level": "risk_level", "counterparty_bic": "counterparty_bic"},
		DefaultOrder:   "screened_at DESC",
		Mutable:        true,
	},
	"system-events": {
		URL: "system-events", Table: "system_events", PK: "event_id",
		AllowedFilters: map[string]string{"code": "code", "component": "component"},
		DefaultOrder:   "at DESC",
		Mutable:        true,
	},
	"eod-jobs": {
		URL: "eod-jobs", Table: "eod_jobs", PK: "job_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"status": "status", "business_date": "business_date"},
		DefaultOrder:   "business_date DESC",
		Mutable:        true,
	},
	"outbox-events": {
		URL: "outbox-events", Table: "outbox_events", PK: "outbox_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"aggregate_type": "aggregate_type", "event_name": "event_name", "topic": "topic"},
		DefaultOrder:   "occurred_at DESC",
		Mutable:        true,
	},
	"outbox-archive": {
		URL: "outbox-archive", Table: "outbox_dispatched_archive", PK: "outbox_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"aggregate_type": "aggregate_type", "topic": "topic"},
		DefaultOrder:   "dispatched_at DESC",
		Mutable:        false, // archive is append-only; never accept mutations via admin
	},

	// ── audit + immutable ──────────────────────────────────────────────────
	"audit-events": {
		URL: "audit-events", Table: "audit_events", PK: "event_id", TenantColumn: "tenant_id",
		AllowedFilters: map[string]string{"event_type": "event_type", "source": "source", "correlation_id": "correlation_id"},
		DefaultOrder:   "occurred_at DESC",
		Mutable:        false, // ISO 27001 + LGPD: NEVER deletable
	},
}

// SchemaByURL looks up a Schema by URL path segment, returning false if not registered.
func SchemaByURL(url string) (Schema, bool) {
	s, ok := schemas[url]
	return s, ok
}

// AllSchemas returns all registered URLs (for diagnostics + smoke).
func AllSchemas() []Schema {
	out := make([]Schema, 0, len(schemas))
	for _, s := range schemas {
		out = append(out, s)
	}
	return out
}
