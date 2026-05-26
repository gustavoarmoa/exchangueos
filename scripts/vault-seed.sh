#!/usr/bin/env bash
# scripts/vault-seed.sh — bootstrap ExchangeOS secrets in Vault.
#
# Run ONCE per environment by the platform operator after Vault is provisioned.
# Idempotent — re-running with the same VAULT_TOKEN updates existing secrets.
#
# Required env:
#   VAULT_ADDR     e.g. https://vault.revenu.internal:8200
#   VAULT_TOKEN    operator token with write on secret/data/exchangeos/*
#   EXCHANGEOS_DB_DSN     full CRDB DSN (verify-full + cert paths)
#   OIDC_CLIENT_SECRET    14-char rotating M2M secret (KeycloakOS)
#
# Optional:
#   KAFKA_BROKERS         comma-separated brokers (only for production)
#
# Usage:
#   VAULT_ADDR=... VAULT_TOKEN=... EXCHANGEOS_DB_DSN=... \
#     OIDC_CLIENT_SECRET=... ./scripts/vault-seed.sh

set -euo pipefail

: "${VAULT_ADDR:?VAULT_ADDR is required}"
: "${VAULT_TOKEN:?VAULT_TOKEN is required}"
: "${EXCHANGEOS_DB_DSN:?EXCHANGEOS_DB_DSN is required}"
: "${OIDC_CLIENT_SECRET:?OIDC_CLIENT_SECRET is required}"

if ! command -v vault >/dev/null 2>&1; then
    echo "ERROR: vault CLI not installed (https://developer.hashicorp.com/vault/install)" >&2
    exit 2
fi

echo "→ Writing secret/data/exchangeos/db"
vault kv put secret/data/exchangeos/db \
    dsn="${EXCHANGEOS_DB_DSN}"

echo "→ Writing secret/data/exchangeos/oidc"
vault kv put secret/data/exchangeos/oidc \
    client_id=exchangeos-api \
    client_secret="${OIDC_CLIENT_SECRET}"

if [[ -n "${KAFKA_BROKERS:-}" ]]; then
    echo "→ Writing secret/data/exchangeos/kafka"
    vault kv put secret/data/exchangeos/kafka \
        brokers="${KAFKA_BROKERS}" \
        client_id=exchangeos-worker
fi

echo "→ Applying minimum read-only policy for the External Secrets Operator"
cat > /tmp/exchangeos-readonly.hcl <<'HCL'
path "secret/data/exchangeos/*" {
  capabilities = ["read", "list"]
}
HCL
vault policy write exchangeos-readonly /tmp/exchangeos-readonly.hcl
rm -f /tmp/exchangeos-readonly.hcl

echo "→ Binding K8s ServiceAccount → policy via Kubernetes auth method"
vault write auth/kubernetes/role/exchangeos \
    bound_service_account_names=exchangeos \
    bound_service_account_namespaces=exchangeos \
    policies=exchangeos-readonly \
    ttl=1h

echo ""
echo "✅ Vault seeded for ExchangeOS"
echo ""
echo "Next steps:"
echo "  1. Verify: vault kv get secret/data/exchangeos/db"
echo "  2. Deploy External Secrets Operator + ClusterSecretStore pointing at this Vault"
echo "  3. Apply Helm chart — ExternalSecret CRs will materialise Secrets in the namespace"
