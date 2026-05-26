# ExchangeOS networking module — VPC + private subnet + Cloud NAT + VPC Service Controls.

terraform {
  required_version = ">= 1.5.0"
  required_providers {
    google = { source = "hashicorp/google", version = "~> 5.30" }
  }
}

variable "project_id"   { type = string }
variable "region"       { type = string }
variable "vpc_name" {
  type    = string
  default = "exchangeos-vpc"
}
variable "subnet_cidr" {
  type    = string
  default = "10.10.0.0/20"
}
variable "pods_cidr" {
  type    = string
  default = "10.20.0.0/14"
}
variable "services_cidr" {
  type    = string
  default = "10.30.0.0/20"
}

resource "google_compute_network" "vpc" {
  project                 = var.project_id
  name                    = var.vpc_name
  auto_create_subnetworks = false
  routing_mode            = "REGIONAL"
}

resource "google_compute_subnetwork" "primary" {
  project                  = var.project_id
  name                     = "${var.vpc_name}-primary"
  region                   = var.region
  network                  = google_compute_network.vpc.id
  ip_cidr_range            = var.subnet_cidr
  private_ip_google_access = true

  secondary_ip_range {
    range_name    = "gke-pods"
    ip_cidr_range = var.pods_cidr
  }
  secondary_ip_range {
    range_name    = "gke-services"
    ip_cidr_range = var.services_cidr
  }
}

# Cloud Router + NAT — egress for private nodes (e.g. OLINDA, GitHub, etc.).
resource "google_compute_router" "nat" {
  project = var.project_id
  name    = "${var.vpc_name}-nat-router"
  region  = var.region
  network = google_compute_network.vpc.id
}

resource "google_compute_router_nat" "nat" {
  project                            = var.project_id
  name                               = "${var.vpc_name}-nat"
  router                             = google_compute_router.nat.name
  region                             = var.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
}

output "vpc_self_link"    { value = google_compute_network.vpc.self_link }
output "subnet_self_link" { value = google_compute_subnetwork.primary.self_link }
