#!/usr/bin/env bash
# scripts/smoke-prod.sh — production smoke validation suite.
#
# Hits every public health + version + smoke endpoint and gates the canary
# promotion. Returns 0 only when all checks pass; non-zero exits identify
# the first failing check so an Argo Rollouts AnalysisTemplate or CI job
# can abort early.
#
# Env (with defaults for local compose):
#   EXCHANGEOS_BASE_URL    default https://api.exchangeos.revenu.tech
#   EXCHANGEOS_GRPC_TARGET default grpc.exchangeos.revenu.tech:443
#   STRICT                 if "true", a 404 on /v1/trades/<unknown> is treated
#                          as success (correct error code propagation)
#                          rather than failure. Default: true.

set -euo pipefail

BASE_URL="${EXCHANGEOS_BASE_URL:-https://api.exchangeos.revenu.tech}"
GRPC_TARGET="${EXCHANGEOS_GRPC_TARGET:-grpc.exchangeos.revenu.tech:443}"
STRICT="${STRICT:-true}"

OK="✅"
FAIL="❌"
fails=0

check() {
    local label="$1"; shift
    if "$@" >/tmp/smoke.out 2>&1; then
        echo "  ${OK}  ${label}"
    else
        echo "  ${FAIL} ${label}" >&2
        sed 's/^/      /' /tmp/smoke.out >&2 || true
        fails=$((fails + 1))
    fi
}

http_status_eq() {
    local path="$1" want="$2"
    local got
    got=$(curl -sS -o /dev/null -w '%{http_code}' --max-time 10 "${BASE_URL}${path}")
    [[ "$got" == "$want" ]] || { echo "got $got want $want"; return 1; }
}

http_json_field_present() {
    local path="$1" field="$2"
    curl -sS --max-time 10 "${BASE_URL}${path}" | grep -q "\"${field}\""
}

echo "── ExchangeOS production smoke ──"
echo "  base URL : ${BASE_URL}"
echo "  gRPC    : ${GRPC_TARGET}"
echo "  STRICT  : ${STRICT}"
echo ""

# ── HTTP health probes ──────────────────────────────────────────────────────
check "GET /healthz returns 200"            http_status_eq /healthz 200
check "GET /readyz returns 200"             http_status_eq /readyz 200
check "GET /version exposes service+version" http_json_field_present /version service

# ── HTTP smoke endpoints ────────────────────────────────────────────────────
check "GET /v1/refdata/currencies?active_only=true returns 200" \
    http_status_eq "/v1/refdata/currencies?active_only=true" 200
check "GET /v1/refdata/currencies exposes 'currencies' field" \
    http_json_field_present "/v1/refdata/currencies?active_only=true" currencies

# ── Error-code propagation (correct 404/400 across gRPC↔HTTP boundary) ──────
if [[ "$STRICT" == "true" ]]; then
    check "GET /v1/trades/<unknown UUID> returns 404 (correct ErrNotFound mapping)" \
        http_status_eq "/v1/trades/00000000-0000-0000-0000-000000000001" 404
    check "GET /v1/trades/not-a-uuid returns 400 (correct InvalidArgument mapping)" \
        http_status_eq "/v1/trades/not-a-uuid" 400
fi

# ── gRPC health (requires grpcurl) ──────────────────────────────────────────
# Auto-detect TLS vs plaintext. Production targets end in :443 (system trust roots).
# Local/non-prod (localhost, 127.0.0.1, *.local) use -plaintext.
# Override via EXCHANGEOS_GRPC_FLAGS for custom mTLS.
GRPC_FLAGS=""
case "${GRPC_TARGET}" in
    *:443)                            GRPC_FLAGS="" ;;
    localhost:*|127.0.0.1:*|*.local:*) GRPC_FLAGS="-plaintext" ;;
    *)                                GRPC_FLAGS="${EXCHANGEOS_GRPC_FLAGS:-}" ;;
esac
if command -v grpcurl >/dev/null 2>&1; then
    check "gRPC health.Check returns SERVING" bash -c "
        grpcurl ${GRPC_FLAGS} -max-time 5 ${GRPC_TARGET} grpc.health.v1.Health/Check 2>/dev/null \
            | grep -q '\"SERVING\"'
    "
else
    echo "  ⚠  grpcurl not installed — skipping gRPC health check (install via: brew install grpcurl)"
fi

echo ""
if (( fails > 0 )); then
    echo "${FAIL} ${fails} check(s) failed — abort canary promotion."
    exit 1
fi
echo "${OK} All smoke checks passed. Safe to promote canary."
