# ExchangeOS GKE Autopilot module — opinionated defaults for production FX workloads.

terraform {
  required_version = ">= 1.5.0"
  required_providers {
    google      = { source = "hashicorp/google", version = "~> 5.30" }
    google-beta = { source = "hashicorp/google-beta", version = "~> 5.30" }
  }
}

variable "project_id" { type = string }
variable "region"     { type = string }
variable "cluster_name" {
  type    = string
  default = "exchangeos"
}
variable "release_channel" {
  type    = string
  default = "REGULAR"
}
variable "vpc_self_link" { type = string }
variable "subnet_self_link" { type = string }
variable "master_ipv4_cidr_block" {
  type    = string
  default = "10.100.0.0/28"
}
variable "binary_authorization_evaluation_mode" {
  type    = string
  default = "PROJECT_SINGLETON_POLICY_ENFORCE"
}

# GKE Autopilot — managed by Google, hardened by default.
resource "google_container_cluster" "this" {
  provider           = google-beta
  project            = var.project_id
  name               = var.cluster_name
  location           = var.region
  enable_autopilot   = true
  deletion_protection = true

  release_channel { channel = var.release_channel }

  network    = var.vpc_self_link
  subnetwork = var.subnet_self_link

  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = false
    master_ipv4_cidr_block  = var.master_ipv4_cidr_block
  }

  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }

  binary_authorization {
    evaluation_mode = var.binary_authorization_evaluation_mode
  }

  database_encryption {
    state    = "ENCRYPTED"
    key_name = google_kms_crypto_key.gke_etcd.id
  }

  master_authorized_networks_config {
    cidr_blocks {
      cidr_block   = "10.0.0.0/8"
      display_name = "internal-vpn"
    }
  }
}

# CMEK for etcd at-rest encryption.
resource "google_kms_key_ring" "exchangeos" {
  project  = var.project_id
  name     = "exchangeos-${var.region}"
  location = var.region
}

resource "google_kms_crypto_key" "gke_etcd" {
  name            = "gke-etcd"
  key_ring        = google_kms_key_ring.exchangeos.id
  rotation_period = "7776000s"   # 90 days
  purpose         = "ENCRYPT_DECRYPT"

  version_template {
    algorithm        = "GOOGLE_SYMMETRIC_ENCRYPTION"
    protection_level = "HSM"
  }
}

output "cluster_name"     { value = google_container_cluster.this.name }
output "cluster_endpoint" { value = google_container_cluster.this.endpoint }
output "kms_key_id"       { value = google_kms_crypto_key.gke_etcd.id }
