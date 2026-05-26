# ExchangeOS — production environment composition (us-east1).

terraform {
  required_version = ">= 1.5.0"
  backend "gcs" {
    bucket = "revenu-platform-tfstate"
    prefix = "exchangeos/production"
  }
}

provider "google" {
  project = var.project_id
  region  = "us-east1"
}

provider "google-beta" {
  project = var.project_id
  region  = "us-east1"
}

variable "project_id" {
  type        = string
  description = "GCP project id (e.g. revenu-platform-prod)"
}

module "network" {
  source     = "../../modules/exchangeos-network"
  project_id = var.project_id
  region     = "us-east1"
}

module "iam" {
  source     = "../../modules/exchangeos-iam"
  project_id = var.project_id
}

module "gke" {
  source           = "../../modules/exchangeos-gke"
  project_id       = var.project_id
  region           = "us-east1"
  vpc_self_link    = module.network.vpc_self_link
  subnet_self_link = module.network.subnet_self_link
}

output "cluster_endpoint" { value = module.gke.cluster_endpoint }
output "service_account"  { value = module.iam.service_account_email }
