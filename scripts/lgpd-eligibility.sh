#!/usr/bin/env bash
# scripts/lgpd-eligibility.sh — query CRDB for all rows referencing a data subject.
#
# Output: per-table row count + earliest occurred_at + regulatory-hold status +
# erasure eligibility (ELIGIBLE / DEFERRED_UNTIL_<date> / FROZEN_REGULATORY).
#
# This is a READ-ONLY discovery script — it does NOT mutate data. Use the
# output to build the execution plan reviewed by DPO + Compliance Officer.
#
# Usage:
#   bash scripts/lgpd-eligibility.sh <subject-uuid>
#
# Required env:
#   EXCHANGEOS_DB_DSN  — read-only role (e.g. exchangeos_dpo_ro)

set -euo pipefail

if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <subject-uuid>" >&2
    exit 2
fi

SUBJECT_ID="$1"
DSN="${EXCHANGEOS_DB_DSN:?EXCHANGEOS_DB_DSN required}"
TODAY=$(date -u +%Y-%m-%d)

# Validate UUID format (basic — 36 chars with dashes)
if [[ ! "${SUBJECT_ID}" =~ ^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$ ]]; then
    echo "ERROR: subject-id must be UUID format" >&2
    exit 2
fi

echo "═══════════════════════════════════════════════════════════════"
echo "  LGPD Eligibility Report — Subject ${SUBJECT_ID}"
echo "  Generated: ${TODAY}"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# Inline SQL — produces a synthetic eligibility table by querying each
# relevant table. Each block: name, count, earliest_at, hold_status, eligibility.
psql "${DSN}" <<SQL
WITH actors_hits AS (
    SELECT 'actors' AS table_name,
           COUNT(*) AS row_count,
           MIN(created_at) AS earliest_at,
           'NO' AS regulatory_hold,
           CASE
               WHEN COUNT(*) = 0 THEN 'N/A'
               ELSE 'ELIGIBLE_REDACT'
           END AS eligibility
    FROM actors WHERE id = '${SUBJECT_ID}' OR external_id = '${SUBJECT_ID}'
),
counterparties_hits AS (
    SELECT 'counterparties',
           COUNT(*),
           MIN(created_at),
           CASE WHEN MAX(updated_at) > now() - interval '7 years' THEN 'YES_RETENTION' ELSE 'NO' END,
           CASE
               WHEN COUNT(*) = 0 THEN 'N/A'
               WHEN MAX(updated_at) > now() - interval '7 years' THEN
                   'DEFERRED_UNTIL_' || to_char(MAX(updated_at) + interval '7 years', 'YYYY-MM-DD')
               ELSE 'ELIGIBLE_REDACT'
           END
    FROM counterparties WHERE id = '${SUBJECT_ID}'
),
fx_trades_hits AS (
    SELECT 'fx_trades',
           COUNT(*),
           MIN(trade_date),
           'YES_BACEN_10Y',
           CASE WHEN COUNT(*) = 0 THEN 'N/A' ELSE 'FROZEN_REGULATORY' END
    FROM fx_trades
    WHERE buyer_id = '${SUBJECT_ID}' OR seller_id = '${SUBJECT_ID}'
       OR initiating_actor_id = '${SUBJECT_ID}'
),
audit_events_hits AS (
    SELECT 'audit_events',
           COUNT(*),
           MIN(occurred_at),
           'YES_INTEGRITY',
           CASE WHEN COUNT(*) = 0 THEN 'N/A' ELSE 'FROZEN_REGULATORY' END
    FROM audit_events WHERE subject_id = '${SUBJECT_ID}' OR actor_id = '${SUBJECT_ID}'
),
quotes_hits AS (
    SELECT 'quotes',
           COUNT(*),
           MIN(created_at),
           'NO',
           CASE WHEN COUNT(*) = 0 THEN 'N/A' ELSE 'ELIGIBLE_HARD_DELETE' END
    FROM quotes WHERE requester_id = '${SUBJECT_ID}'
),
screening_hits AS (
    SELECT 'screening_results',
           COUNT(*),
           MIN(screened_at),
           'YES_7Y',
           CASE
               WHEN COUNT(*) = 0 THEN 'N/A'
               WHEN MAX(screened_at) > now() - interval '7 years' THEN
                   'DEFERRED_UNTIL_' || to_char(MAX(screened_at) + interval '7 years', 'YYYY-MM-DD')
               ELSE 'ELIGIBLE_REDACT'
           END
    FROM screening_results WHERE actor_id = '${SUBJECT_ID}' OR counterparty_id = '${SUBJECT_ID}'
)
SELECT * FROM actors_hits
UNION ALL SELECT * FROM counterparties_hits
UNION ALL SELECT * FROM fx_trades_hits
UNION ALL SELECT * FROM audit_events_hits
UNION ALL SELECT * FROM quotes_hits
UNION ALL SELECT * FROM screening_hits
ORDER BY 1;
SQL

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  Legend:"
echo "    ELIGIBLE_HARD_DELETE — no regulatory hold; delete rows"
echo "    ELIGIBLE_REDACT      — keep row, redact PII fields per workflow"
echo "    DEFERRED_UNTIL_<date> — retention window not yet expired"
echo "    FROZEN_REGULATORY    — regulatory record-keeping; cannot erase"
echo ""
echo "  Next: build execution plan + DPO/Compliance co-sign + execute"
echo "        per docs/security/data-lifecycle/erasure-workflow.md"
echo "═══════════════════════════════════════════════════════════════"
