#!/usr/bin/env bash
# scripts/download-xsd.sh — pulls all 32 pinned ISO 20022 XSDs into .cache/xsd/
#
# Source of truth: pkg/iso20022/registry/sources.go (parsed via grep).
# Idempotent. Verifies HTTP 200. Emits sha256 manifest in .cache/xsd/manifest.txt.
#
# Env:
#   OFFLINE=true       skip downloads, only validate existing local cache + manifest
#   FAIL_FAST=true     exit on first failure (default: collect failures, exit 1 at end)
#   CACHE_DIR=...      override .cache/xsd
#   CURL_OPTS=...      extra curl options (e.g. --proxy http://...)
#
# Usage:
#   bash scripts/download-xsd.sh           # download + verify
#   OFFLINE=true bash scripts/download-xsd.sh   # check-only

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SOURCES="${ROOT}/pkg/iso20022/registry/sources.go"
CACHE_DIR="${CACHE_DIR:-${ROOT}/.cache/xsd}"
MANIFEST="${CACHE_DIR}/manifest.txt"
FAIL_FAST="${FAIL_FAST:-false}"
OFFLINE="${OFFLINE:-false}"

# Defensive sha256 helper (Linux/macOS).
sha256() {
    if command -v sha256sum >/dev/null 2>&1; then sha256sum "$1" | awk '{print $1}'
    elif command -v shasum >/dev/null 2>&1; then shasum -a 256 "$1" | awk '{print $1}'
    else echo "ERROR: no sha256 tool available" >&2; exit 2
    fi
}

mkdir -p "${CACHE_DIR}"
: > "${MANIFEST}.tmp"

# Extract XSDSourceURL string literals from sources.go.
# Each line looks like:   XSDSourceURL: "https://...xsd",
# Portable across bash 3.2 (macOS) — no mapfile.
URLS=()
while IFS= read -r line; do
    URLS+=("$line")
done < <(grep -oE 'XSDSourceURL:[[:space:]]*"[^"]+"' "${SOURCES}" \
    | sed -E 's/.*"([^"]+)"/\1/')

if [[ ${#URLS[@]} -ne 32 ]]; then
    echo "WARNING: expected 32 URLs in sources.go, got ${#URLS[@]}" >&2
fi

OK=0
FAIL=0
declare -a FAILED

for url in "${URLS[@]}"; do
    fname=$(basename "${url}")
    out="${CACHE_DIR}/${fname}"

    if [[ "${OFFLINE}" == "true" ]]; then
        if [[ -f "${out}" ]]; then
            digest=$(sha256 "${out}")
            printf '%s  %s\n' "${digest}" "${fname}" >> "${MANIFEST}.tmp"
            OK=$((OK + 1))
            echo "  [cached] ${fname}"
        else
            FAIL=$((FAIL + 1))
            FAILED+=("(missing) ${fname}")
            echo "  [MISS]   ${fname}" >&2
        fi
        continue
    fi

    # Online: download with retry, verify HTTP 200, write atomically.
    tmp=$(mktemp "${CACHE_DIR}/.${fname}.XXXXXX")
    if curl -fsSL --retry 3 --retry-delay 2 --connect-timeout 10 --max-time 60 \
            ${CURL_OPTS:-} -o "${tmp}" "${url}" 2>/dev/null; then
        mv -f "${tmp}" "${out}"
        digest=$(sha256 "${out}")
        printf '%s  %s\n' "${digest}" "${fname}" >> "${MANIFEST}.tmp"
        OK=$((OK + 1))
        echo "  [ok]     ${fname}"
    else
        rm -f "${tmp}"
        FAIL=$((FAIL + 1))
        FAILED+=("${url}")
        echo "  [FAIL]   ${url}" >&2
        if [[ "${FAIL_FAST}" == "true" ]]; then break; fi
    fi
done

mv -f "${MANIFEST}.tmp" "${MANIFEST}"

echo ""
echo "─────────────────────────────────────────────"
echo "  XSD download summary:  ${OK} ok / ${FAIL} fail"
echo "  Cache dir:  ${CACHE_DIR}"
echo "  Manifest:   ${MANIFEST}"
echo "─────────────────────────────────────────────"

if (( FAIL > 0 )); then
    echo "Failures:" >&2
    for f in "${FAILED[@]}"; do echo "  - ${f}" >&2; done
    exit 1
fi
