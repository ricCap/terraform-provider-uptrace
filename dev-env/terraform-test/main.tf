terraform {
  required_providers {
    uptrace = {
      source  = "riccap/uptrace"
      version = "~> 0.1"
    }
  }
}

provider "uptrace" {
  endpoint   = "https://api2.uptrace.dev/internal/v1"
  token      = var.uptrace_token
  project_id = var.uptrace_project_id
}

# ============================================================================
# Monitors (without notification channels to avoid API version issues)
# ============================================================================

# Error monitor - Test cloud API (also needs monitor_trend_aggregation_func)
resource "uptrace_monitor" "api_errors" {
  name = "Terraform Test - API Error Rate"
  type = "error"

  notify_everyone_by_email = false

  params = {
    metrics = [
      {
        name = "span.count"
      }
    ]
    query                          = "span.status_code:error"
    monitor_trend_aggregation_func = "sum"
    max_allowed_value              = 100
    check_num_point                = 1
  }
}

# Metric monitor - Test cloud API with monitor_trend_aggregation_func
resource "uptrace_monitor" "high_latency" {
  name = "Terraform Test - High Latency"
  type = "metric"

  notify_everyone_by_email = false

  params = {
    metrics = [
      {
        name  = "span.duration"
        alias = "$latency"
      }
    ]
    query                          = "span.kind:server"
    column                         = "$latency"
    monitor_trend_aggregation_func = "avg"
    max_allowed_value              = 5000000000 # 5 seconds in nanoseconds
    check_num_point                = 2
  }
}

# Notification channel - Test cloud API with priority
resource "uptrace_notification_channel" "test_webhook" {
  name = "Terraform Test - Webhook"
  type = "webhook"

  priority = ["high", "critical"]

  params = {
    url = "https://example.com/webhook"
  }
}

# ============================================================================
# Data Sources
# ============================================================================

# Fetch all monitors
data "uptrace_monitors" "all" {
  depends_on = [
    uptrace_monitor.api_errors,
    uptrace_monitor.high_latency,
    uptrace_notification_channel.test_webhook
  ]
}

# Fetch specific monitor
data "uptrace_monitor" "error_monitor" {
  id = uptrace_monitor.api_errors.id
}

# ============================================================================
# Outputs
# ============================================================================

output "monitors" {
  description = "Created monitor IDs"
  value = {
    error_monitor   = uptrace_monitor.api_errors.id
    latency_monitor = uptrace_monitor.high_latency.id
  }
}

output "notification_channel" {
  description = "Created notification channel"
  value = {
    id       = uptrace_notification_channel.test_webhook.id
    name     = uptrace_notification_channel.test_webhook.name
    priority = uptrace_notification_channel.test_webhook.priority
  }
}

output "all_monitors_count" {
  description = "Total number of monitors in the project"
  value       = length(data.uptrace_monitors.all.monitors)
}

output "error_monitor_details" {
  description = "Details of the error monitor from data source"
  value = {
    id   = data.uptrace_monitor.error_monitor.id
    name = data.uptrace_monitor.error_monitor.name
    type = data.uptrace_monitor.error_monitor.type
  }
}
