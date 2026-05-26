// modules/exchangeos-budget/main.tf — GCP Billing budgets + Pub/Sub alert routing.
//
// Creates 3 budgets per environment with progressive alert thresholds (50/80/100%).
// Alerts route via Pub/Sub → Cloud Function → Slack #exchangeos-finops channel.
// Cost-allocation labels (module=exchangeos + env + bc) are enforced via Helm
// values + Terraform module conventions — see docs/operations/cost-allocation.md.

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

variable "billing_account" {
  description = "GCP billing account ID (XXXXXX-XXXXXX-XXXXXX)"
  type        = string
}

variable "project_id" {
  description = "GCP project ID hosting exchangeos workloads"
  type        = string
}

variable "env" {
  description = "Environment label (dev|staging|production)"
  type        = string
  validation {
    condition     = contains(["dev", "staging", "production"], var.env)
    error_message = "env must be dev, staging, or production."
  }
}

variable "monthly_budget_usd" {
  description = "Monthly budget cap in USD for this env"
  type        = number
}

variable "pubsub_topic" {
  description = "Pub/Sub topic for budget alerts (created if not provided)"
  type        = string
  default     = ""
}

locals {
  topic_name = var.pubsub_topic != "" ? var.pubsub_topic : google_pubsub_topic.budget_alerts[0].name
}

// ── Pub/Sub topic for budget alerts (created if caller didn't provide one) ────
resource "google_pubsub_topic" "budget_alerts" {
  count   = var.pubsub_topic == "" ? 1 : 0
  project = var.project_id
  name    = "exchangeos-budget-alerts-${var.env}"

  labels = {
    module = "exchangeos"
    env    = var.env
    type   = "finops"
  }
}

// ── Main budget: overall exchangeos spend ────────────────────────────────────
resource "google_billing_budget" "exchangeos_total" {
  billing_account = var.billing_account
  display_name    = "exchangeos-${var.env}-total"

  budget_filter {
    projects = ["projects/${var.project_id}"]
    labels = {
      module = "exchangeos"
      env    = var.env
    }
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units         = tostring(var.monthly_budget_usd)
    }
  }

  threshold_rules {
    threshold_percent = 0.5
    spend_basis       = "CURRENT_SPEND"
  }

  threshold_rules {
    threshold_percent = 0.8
    spend_basis       = "CURRENT_SPEND"
  }

  threshold_rules {
    threshold_percent = 1.0
    spend_basis       = "CURRENT_SPEND"
  }

  threshold_rules {
    threshold_percent = 1.2 // overshoot — page finops + engineering lead
    spend_basis       = "CURRENT_SPEND"
  }

  // Forecasted overspend alert — catches trajectory before it hits 100%
  threshold_rules {
    threshold_percent = 1.0
    spend_basis       = "FORECASTED_SPEND"
  }

  all_updates_rule {
    pubsub_topic                     = "projects/${var.project_id}/topics/${local.topic_name}"
    disable_default_iam_recipients   = false
    monitoring_notification_channels = []
  }
}

// ── Per-service sub-budgets — early warning when one service dominates ──────
resource "google_billing_budget" "exchangeos_compute" {
  billing_account = var.billing_account
  display_name    = "exchangeos-${var.env}-compute"

  budget_filter {
    projects = ["projects/${var.project_id}"]
    services = ["services/6F81-5844-456A"] // Compute Engine (GKE Autopilot)
    labels = {
      module = "exchangeos"
      env    = var.env
    }
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units         = tostring(floor(var.monthly_budget_usd * 0.6)) // 60% of total
    }
  }

  threshold_rules {
    threshold_percent = 0.8
    spend_basis       = "CURRENT_SPEND"
  }

  all_updates_rule {
    pubsub_topic = "projects/${var.project_id}/topics/${local.topic_name}"
  }
}

resource "google_billing_budget" "exchangeos_storage" {
  billing_account = var.billing_account
  display_name    = "exchangeos-${var.env}-storage"

  budget_filter {
    projects = ["projects/${var.project_id}"]
    services = [
      "services/95FF-2EF5-5EA1", // Cloud Storage
      "services/A1E8-BE35-7EBC", // Cloud KMS
    ]
    labels = {
      module = "exchangeos"
      env    = var.env
    }
  }

  amount {
    specified_amount {
      currency_code = "USD"
      units         = tostring(floor(var.monthly_budget_usd * 0.15))
    }
  }

  threshold_rules {
    threshold_percent = 0.8
    spend_basis       = "CURRENT_SPEND"
  }

  all_updates_rule {
    pubsub_topic = "projects/${var.project_id}/topics/${local.topic_name}"
  }
}

// ── Outputs ──────────────────────────────────────────────────────────────────
output "budget_id_total" {
  value = google_billing_budget.exchangeos_total.id
}

output "pubsub_topic" {
  value = local.topic_name
}
