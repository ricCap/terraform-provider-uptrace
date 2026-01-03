terraform {
  required_providers {
    uptrace = {
      source  = "registry.terraform.io/riccap/uptrace"
      version = "~> 0.1"
    }
  }
}

provider "uptrace" {
  endpoint   = "http://localhost:14318/internal/v1"
  token      = var.uptrace_token
  project_id = 1
}

# Test metric monitor
resource "uptrace_monitor" "high_cpu" {
  name = "High CPU Usage - Test"
  type = "metric"

  notify_everyone_by_email = false

  params = {
    metrics = [
      {
        name  = "system.cpu.utilization"
        alias = "$cpu_usage"
      }
    ]
    query             = "avg($cpu_usage) > 90"
    max_allowed_value = 90
    grouping_interval = 60000
    check_num_point   = 3
    nulls_mode        = "allow"
  }
}

# Test error monitor
resource "uptrace_monitor" "api_errors" {
  name = "API Error Monitor - Updated"
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

output "metric_monitor_id" {
  value       = uptrace_monitor.high_cpu.id
  description = "ID of the metric monitor"
}

output "metric_monitor_state" {
  value       = uptrace_monitor.high_cpu.state
  description = "State of the metric monitor"
}

output "error_monitor_id" {
  value       = uptrace_monitor.api_errors.id
  description = "ID of the error monitor"
}

output "error_monitor_state" {
  value       = uptrace_monitor.api_errors.state
  description = "State of the error monitor"
}
