# ADR 0006 — Workload Identity Federation over JSON service-account keys

- Status: Accepted
- Date: 2026-05-24

## Context

Pods need GCP API credentials (logging, monitoring, KMS decrypt, Secret Manager read). The legacy approach mounts a JSON service-account key as a Kubernetes Secret. JSON keys are static long-lived credentials that:

- Frequently leak into source repos / Slack / logs
- Don't rotate without operator action
- Can't be revoked without breaking running workloads

## Decision

**Workload Identity Federation (WIF).** Bind the K8s ServiceAccount `exchangeos` in namespace `exchangeos` to a GCP ServiceAccount via `roles/iam.workloadIdentityUser`. Pods authenticate using their ambient K8s identity — no secret material to manage.

Terraform binding in `deploy/terraform/modules/exchangeos-iam/main.tf`:

```hcl
resource "google_service_account_iam_member" "wif" {
  service_account_id = google_service_account.exchangeos.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project_id}.svc.id.goog[${var.namespace}/${var.k8s_service_account}]"
}
```

ServiceAccount annotation in Helm values:

```yaml
serviceAccount:
  annotations:
    iam.gke.io/gcp-service-account: exchangeos@revenu-platform.iam.gserviceaccount.com
```

## Consequences

### Positive

- **Zero secrets to manage** — no JSON keys ever generated
- **Automatic rotation** — short-lived OIDC tokens, refreshed by metadata server
- **Per-namespace + per-SA scoping** — granular least-privilege
- **Revocation is instant** — delete the IAM binding

### Negative

- **GKE-specific** — won't work on bare-metal K8s; not a concern for us
- **First-time setup more complex** — Workload Identity Pool + provider config required
- **Less universal than JSON keys** — third-party tools sometimes only document JSON

### Mitigations

- Terraform module captures the WIF setup; new modules copy the same pattern
- Onboarding doc references this ADR

## Alternatives considered

- **JSON service-account keys** — rejected: 90% of secret-leak incidents in the wild
- **Vault → GCP dynamic credentials** — overlapping with WIF; extra moving piece
- **AWS-style instance metadata** — GCP equivalent IS WIF on GKE
