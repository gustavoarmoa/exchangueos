#!/usr/bin/env bash
# Drift guard for the BACEN nature-code generated catalogue.
#
# Re-runs pkg/bacen/codegen against the canonical CSV into a tempfile and
# byte-compares against the committed pkg/bacen/codes_full.go. Exits non-zero
# on drift so CI catches "edited CSV but forgot to regenerate".
#
# Usage:
#   scripts/check-bacen-codegen.sh           # run from repo root
#   bash scripts/check-bacen-codegen.sh
#
# Exit codes:
#   0 — committed file matches a fresh codegen run
#   1 — drift detected; the script prints the diff
#   2 — environment / invocation problem
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

CSV="data/bacen/nature-codes-circ-3690-v20260101.csv"
GO_FILE="pkg/bacen/codes_full.go"

if [[ ! -f "$CSV" ]]; then
  echo "drift-check: missing CSV at $CSV" >&2
  exit 2
fi
if [[ ! -f "$GO_FILE" ]]; then
  echo "drift-check: missing committed $GO_FILE" >&2
  exit 2
fi

TMP="$(mktemp -t codes_full.go.XXXXXX)"
trap 'rm -f "$TMP"' EXIT

go run ./pkg/bacen/codegen --input "$CSV" --output "$TMP" >/dev/null

# The generated file embeds a "Generated: <UTC timestamp>" line that legitimately
# changes between runs. Strip it from both sides before comparing.
filter() {
  grep -v '^// Generated: ' "$1"
}

if ! diff -u <(filter "$GO_FILE") <(filter "$TMP") > /tmp/_bacen_codegen_drift.diff; then
  echo "✗ pkg/bacen/codes_full.go is stale vs $CSV" >&2
  echo "" >&2
  echo "  Diff (committed → fresh regeneration):" >&2
  sed 's/^/    /' /tmp/_bacen_codegen_drift.diff >&2
  echo "" >&2
  echo "  Fix: go generate ./pkg/bacen/... && git add pkg/bacen/codes_full.go" >&2
  exit 1
fi

echo "✓ pkg/bacen/codes_full.go is in sync with $CSV"
