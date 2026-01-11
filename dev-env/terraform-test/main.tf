terraform {
  required_providers {
    uptrace = {
      source = "riccap/uptrace"
    }
  }
}

# ============================================================================
# Provider Configuration
# ============================================================================
# For local dev environment:
#   endpoint   = "http://localhost:14318/internal/v1"
#   token      = "user1_secret_token"
#   project_id = 1
#
# For cloud (requires metrics to exist in project):
#   endpoint   = "https://api2.uptrace.dev/internal/v1"
#   token      = var.uptrace_token
#   project_id = var.uptrace_project_id
# ============================================================================

provider "uptrace" {
  endpoint   = var.uptrace_endpoint
  token      = var.uptrace_token
  project_id = var.uptrace_project_id
}

# ============================================================================
# Error Monitor
# ============================================================================

resource "uptrace_monitor" "error_monitor" {
  name = "Terraform Test - Error Monitor"
  type = "error"

  notify_everyone_by_email = false

  params = {
    metrics = [
      {
        name  = "uptrace_tracing_events"
        alias = "$logs"
      }
    ]
    query = "sum($logs) | where span.event_name exists"
  }
}

# ============================================================================
# Metric Monitor
# Note: Requires the metric to exist in the project (send telemetry first)
# ============================================================================

# resource "uptrace_monitor" "metric_monitor" {
#   name = "Terraform Test - CPU Monitor"
#   type = "metric"
#
#   notify_everyone_by_email = false
#
#   params = {
#     metrics = [
#       {
#         name  = "system.cpu.utilization"
#         alias = "$cpu"
#       }
#     ]
#     query             = "avg($cpu) > 80"
#     max_allowed_value = 80
#     check_num_point   = 2
#   }
# }

# ============================================================================
# Data Sources
# ============================================================================

data "uptrace_monitors" "all" {
  depends_on = [uptrace_monitor.error_monitor]
}

data "uptrace_monitor" "error_monitor" {
  id = uptrace_monitor.error_monitor.id
}

# ============================================================================
# Outputs
# ============================================================================

output "error_monitor_id" {
  description = "Created error monitor ID"
  value       = uptrace_monitor.error_monitor.id
}

output "all_monitors_count" {
  description = "Total number of monitors in the project"
  value       = length(data.uptrace_monitors.all.monitors)
}

output "error_monitor_details" {
  description = "Details of the error monitor from data source"
  value = {
    id    = data.uptrace_monitor.error_monitor.id
    name  = data.uptrace_monitor.error_monitor.name
    type  = data.uptrace_monitor.error_monitor.type
    state = data.uptrace_monitor.error_monitor.state
  }
}
