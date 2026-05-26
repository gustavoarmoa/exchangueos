#!/usr/bin/env bash
# scripts/audit-evidence-bundle.sh — package ISO 27001 evidence for auditor handoff.
#
# Produces a timestamped tarball at .audit-bundles/exchangeos-evidence-YYYYMMDD-HHMM.tar.gz
# containing every evidence artefact referenced from the controls mapping +
# interview prep documents. SAFE — read-only collection.
#
# Usage:
#   bash scripts/audit-evidence-bundle.sh                  # default destination
#   OUT_DIR=/tmp bash scripts/audit-evidence-bundle.sh     # custom destination

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT_DIR="${OUT_DIR:-${ROOT}/.audit-bundles}"
STAMP="$(date -u +%Y%m%d-%H%M)"
BUNDLE_NAME="exchangeos-evidence-${STAMP}"
STAGE="${OUT_DIR}/${BUNDLE_NAME}"

mkdir -p "${STAGE}"

echo "📦 Staging audit evidence at ${STAGE}"
echo ""

copy() {
    local src="$1"
    local dst="$2"
    if [[ -e "${ROOT}/${src}" ]]; then
        mkdir -p "${STAGE}/$(dirname "${dst}")"
        cp -R "${ROOT}/${src}" "${STAGE}/${dst}"
        echo "  ✓ ${src}"
    else
        echo "  ⚠ MISSING: ${src} (mentioned by controls mapping)" >&2
    fi
}

# ── Top-level governance ─────────────────────────────────────────────────────
copy "CLAUDE.md"                                "01-governance/CLAUDE.md"
copy ".base/plans/index.md"                     "01-governance/plan-index.md"
copy ".base/plans/version.md"                   "01-governance/version-history.md"
copy ".base/plans/CHANGELOG.md"                 "01-governance/CHANGELOG.md"

# ── Security policy + risk ───────────────────────────────────────────────────
copy "docs/security/iso27001-controls-mapping.md" "02-security/controls-mapping.md"
copy "docs/security/iso27001-gap-tracker.md"      "02-security/gap-tracker.md"
copy "docs/security/threat-model-stride.md"       "02-security/threat-model.md"
copy "docs/security/sod-matrix.md"                "02-security/sod-matrix.md"
copy "docs/security/incident-response.md"         "02-security/incident-response.md"
copy "docs/security/dr-runbook.md"                "02-security/dr-runbook.md"
copy "docs/security/audit-interview-prep.md"      "02-security/audit-interview-prep.md"
copy "docs/security/drills"                       "02-security/drills"

# ── Operations + runbooks ────────────────────────────────────────────────────
copy "docs/operations/runbook-index.md"           "03-operations/runbook-index.md"
copy "docs/operations/canary-runbook.md"          "03-operations/canary-runbook.md"
copy "docs/operations/go-live-checklist.md"       "03-operations/go-live-checklist.md"
copy "docs/operations/crdb-hub-tls-pr.md"         "03-operations/crdb-hub-tls-pr.md"

# ── Technical controls evidence ──────────────────────────────────────────────
copy "lefthook.yml"                              "04-technical/lefthook-3tier.yml"
copy ".golangci.yml"                             "04-technical/golangci-strict.yml"
copy "scripts/git-hooks-wrapper.sh"              "04-technical/git-hooks-wrapper.sh"
copy "scripts/vault-seed.sh"                     "04-technical/vault-seed.sh"
copy ".github/workflows/ci.yml"                  "04-technical/ci.yml"
copy ".github/workflows/security.yml"            "04-technical/security-scans.yml"
copy ".github/workflows/slsa-attestation.yml"    "04-technical/slsa-attestation.yml"
copy "deploy/terraform/modules"                  "04-technical/terraform-modules"
copy "deploy/helm/exchangeos/values.yaml"        "04-technical/helm-values.yaml"
copy "deploy/k8s/argo-rollouts/api-rollout.yaml" "04-technical/argo-rollouts-canary.yaml"
copy "deploy/k8s/cert-manager/cluster-issuer.yaml" "04-technical/cert-manager.yaml"
copy "deploy/kafka/topics.yaml"                  "04-technical/kafka-topics-acls.yaml"

# ── Compliance code path ─────────────────────────────────────────────────────
copy "pkg/bacen"                                 "05-compliance/pkg-bacen"
copy "pkg/iso20022/registry/sources.go"          "05-compliance/iso20022-32-schemas.go"
copy "modules/compliance/domain"                 "05-compliance/compliance-domain"
copy "migrations/000008_create_compliance_admin.up.sql" "05-compliance/migration-compliance.sql"

# ── Delivered milestones (proof of work) ─────────────────────────────────────
copy ".base/plans/milestones/delivered"          "06-delivery/all-26-milestones"
copy ".base/plans/roadmap/delivery-dashboard.md" "06-delivery/dashboard.md"

# ── Tests as evidence ────────────────────────────────────────────────────────
echo "  → counting tests..."
(
    cd "${ROOT}"
    {
        echo "# Test inventory (collected $(date -u))"
        echo ""
        echo "## Go tests"
        echo ""
        find . -name '*_test.go' -not -path '*/proto/gen/*' -not -path '*/.cache/*' | sort
        echo ""
        echo "## Test count by package"
        grep -rh '^func Test' --include='*_test.go' --exclude-dir=proto --exclude-dir=.cache \
            | wc -l | xargs printf '  Go test funcs: %s\n'
        grep -rh '^func Benchmark' --include='*_test.go' --exclude-dir=proto --exclude-dir=.cache \
            | wc -l | xargs printf '  Benchmark funcs: %s\n'
    } > "${STAGE}/06-delivery/test-inventory.md"
)
echo "  ✓ test inventory generated"

# ── Pack ─────────────────────────────────────────────────────────────────────
echo ""
echo "📦 Creating tarball..."
TAR_PATH="${OUT_DIR}/${BUNDLE_NAME}.tar.gz"
tar -czf "${TAR_PATH}" -C "${OUT_DIR}" "${BUNDLE_NAME}"
SIZE=$(du -h "${TAR_PATH}" | awk '{print $1}')
SHA=$(shasum -a 256 "${TAR_PATH}" | awk '{print $1}')

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "  Audit evidence bundle ready"
echo "═══════════════════════════════════════════════════════════════"
echo "  Path:   ${TAR_PATH}"
echo "  Size:   ${SIZE}"
echo "  SHA256: ${SHA}"
echo ""
echo "Next steps:"
echo "  1. Upload to secure file-share for auditor"
echo "  2. Email SHA256 to auditor SEPARATELY (integrity check)"
echo "  3. Schedule walkthrough using docs/security/audit-interview-prep.md"
echo "═══════════════════════════════════════════════════════════════"
