# ExchangeOS IAM module — Workload Identity Federation (zero JSON keys) +
# least-privilege roles per binary.

terraform {
  required_version = ">= 1.5.0"
  required_providers {
    google = { source = "hashicorp/google", version = "~> 5.30" }
  }
}

variable "project_id"     { type = string }
variable "namespace"      {
  type    = string
  default = "exchangeos"
}
variable "k8s_service_account" {
  type    = string
  default = "exchangeos"
}

# Google service account bound 1:1 to the K8s service account via WIF.
resource "google_service_account" "exchangeos" {
  project      = var.project_id
  account_id   = "exchangeos"
  display_name = "ExchangeOS application"
  description  = "Workload identity for exchangeos.svc.exchangeos in K8s"
}

# WIF binding: K8s SA → GCP SA.
resource "google_service_account_iam_member" "wif" {
  service_account_id = google_service_account.exchangeos.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project_id}.svc.id.goog[${var.namespace}/${var.k8s_service_account}]"
}

# Least-privilege role assignments.
locals {
  roles = [
    "roles/cloudtrace.agent",
    "roles/monitoring.metricWriter",
    "roles/logging.logWriter",
    "roles/cloudprofiler.agent",
    "roles/cloudkms.cryptoKeyEncrypterDecrypter",
    "roles/secretmanager.secretAccessor",
  ]
}

resource "google_project_iam_member" "exchangeos" {
  for_each = toset(local.roles)
  project  = var.project_id
  role     = each.value
  member   = "serviceAccount:${google_service_account.exchangeos.email}"
}

output "service_account_email" { value = google_service_account.exchangeos.email }
