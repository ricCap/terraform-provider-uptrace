terraform {
  required_providers {
    uptrace = {
      source = "registry.terraform.io/riccap/uptrace"
    }
  }
}

provider "uptrace" {
  endpoint   = var.uptrace_endpoint
  token      = var.uptrace_token
  project_id = var.uptrace_project_id
}

# ============================================================================
# Notification Channels
# ============================================================================

# Critical alerts go to Slack and PagerDuty
resource "uptrace_notification_channel" "slack_critical" {
  name = "Critical Alerts - Slack"
  type = "slack"

  params = {
    webhookUrl = var.slack_critical_webhook
  }
}

resource "uptrace_notification_channel" "webhook_pagerduty" {
  name = "PagerDuty Integration"
  type = "webhook"

  params = {
    url = var.pagerduty_webhook_url
  }
}

# Warning alerts go to team Slack channel
resource "uptrace_notification_channel" "slack_warnings" {
  name = "Warning Alerts - Slack"
  type = "slack"

  params = {
    webhookUrl = var.slack_warnings_webhook
  }
}

# Telegram for mobile notifications
resource "uptrace_notification_channel" "telegram_oncall" {
  name = "On-Call Mobile Alerts"
  type = "telegram"

  params = {
    botToken = var.telegram_bot_token
    chatId   = var.telegram_chat_id
  }
}

# ============================================================================
# Error Monitoring
# ============================================================================

# Critical: High error rate across all services
resource "uptrace_monitor" "critical_error_rate" {
  name = "Critical: High Error Rate"
  type = "error"

  channel_ids = [
    uptrace_notification_channel.slack_critical.id,
    uptrace_notification_channel.webhook_pagerduty.id,
    uptrace_notification_channel.telegram_oncall.id,
  ]

  params = {
    metrics = []
    query   = "span.status_code:error"

    # Alert if error count > 100 in 5 minutes
    min_allowed_value = 0
    max_allowed_value = 100

    # Check every 5 minutes
    check_num_point = 1

    # Notification settings
    notify_everyone_by_email = false
  }
}

# Warning: 4xx errors increasing
resource "uptrace_monitor" "warning_4xx_errors" {
  name = "Warning: High 4xx Error Rate"
  type = "error"

  channel_ids = [
    uptrace_notification_channel.slack_warnings.id,
  ]

  params = {
    metrics = []
    query   = "span.status_code:>=400 span.status_code:<500"

    min_allowed_value = 0
    max_allowed_value = 500

    check_num_point = 1

    notify_everyone_by_email = false
  }
}

# Critical: Database errors
resource "uptrace_monitor" "critical_database_errors" {
  name = "Critical: Database Errors"
  type = "error"

  channel_ids = [
    uptrace_notification_channel.slack_critical.id,
    uptrace_notification_channel.webhook_pagerduty.id,
  ]

  params = {
    metrics = []
    query   = "db.system:* span.status_code:error"

    min_allowed_value = 0
    max_allowed_value = 10

    check_num_point = 1

    notify_everyone_by_email = true
  }
}

# ============================================================================
# Performance Monitoring
# ============================================================================

# Critical: API response time degradation
resource "uptrace_monitor" "critical_api_latency" {
  name = "Critical: API Response Time >2s"
  type = "metric"

  channel_ids = [
    uptrace_notification_channel.slack_critical.id,
    uptrace_notification_channel.telegram_oncall.id,
  ]

  params = {
    metrics = [
      {
        name  = "span.duration"
        alias = "p95_latency"
      }
    ]

    query = "span.system:http span.kind:server"

    # Alert if P95 latency > 2 seconds
    min_allowed_value = 0
    max_allowed_value = 2000000000 # 2 seconds in nanoseconds

    check_num_point = 2 # Alert if true for 2 consecutive checks

    # Column to monitor
    column = "$p95_latency"

    notify_everyone_by_email = false
  }
}

# Warning: Database query slow
resource "uptrace_monitor" "warning_db_slow_queries" {
  name = "Warning: Slow Database Queries >500ms"
  type = "metric"

  channel_ids = [
    uptrace_notification_channel.slack_warnings.id,
  ]

  params = {
    metrics = [
      {
        name  = "span.duration"
        alias = "p95_duration"
      }
    ]

    query = "db.system:* span.kind:client"

    min_allowed_value = 0
    max_allowed_value = 500000000 # 500ms in nanoseconds

    check_num_point = 3

    column = "$p95_duration"

    notify_everyone_by_email = false
  }
}

# ============================================================================
# Throughput Monitoring
# ============================================================================

# Warning: Unusual traffic spike
resource "uptrace_monitor" "warning_traffic_spike" {
  name = "Warning: Traffic Spike Detected"
  type = "metric"

  channel_ids = [
    uptrace_notification_channel.slack_warnings.id,
  ]

  params = {
    metrics = [
      {
        name  = "span.count"
        alias = "request_count"
      }
    ]

    query = "span.system:http span.kind:server"

    # Alert if request count > 10000/minute
    min_allowed_value = 0
    max_allowed_value = 10000

    check_num_point = 2

    column = "$request_count"

    notify_everyone_by_email = false
  }
}

# Critical: Traffic dropped to zero
resource "uptrace_monitor" "critical_no_traffic" {
  name = "Critical: No Traffic Detected"
  type = "metric"

  channel_ids = [
    uptrace_notification_channel.slack_critical.id,
    uptrace_notification_channel.webhook_pagerduty.id,
  ]

  params = {
    metrics = [
      {
        name  = "span.count"
        alias = "request_count"
      }
    ]

    query = "span.system:http span.kind:server"

    # Alert if request count < 10/minute
    min_allowed_value = 10
    max_allowed_value = 999999999

    check_num_point = 3

    column = "$request_count"

    notify_everyone_by_email = true
  }
}

# ============================================================================
# Dashboards
# ============================================================================

# Main application overview dashboard
resource "uptrace_dashboard" "application_overview" {
  name = "Application Overview"

  yaml_config = <<-YAML
    table:
      - schema_version: v2
        row:
          - title: API Performance
            columns:
              - span.system: [http]
                span.kind: [server]
            chart:
              - type: line
                metric: span.duration|p95
                legend: P95 Latency
              - type: line
                metric: span.count|per_min
                legend: Requests/min
          - title: Error Rate
            columns:
              - span.status_code: [error]
            chart:
              - type: line
                metric: span.count|per_min
                legend: Errors/min
                color: red
      - row:
          - title: Database Performance
            columns:
              - db.system: [*]
            chart:
              - type: line
                metric: span.duration|p95
                legend: DB Query P95
              - type: line
                metric: span.count|per_min
                legend: Queries/min
          - title: External API Calls
            columns:
              - span.kind: [client]
                span.system: [http]
            chart:
              - type: line
                metric: span.duration|p95
                legend: External API P95
  YAML
}

# Error tracking dashboard
resource "uptrace_dashboard" "error_tracking" {
  name = "Error Tracking Dashboard"

  yaml_config = <<-YAML
    table:
      - schema_version: v2
        row:
          - title: Total Errors
            columns:
              - span.status_code: [error]
            chart:
              - type: line
                metric: span.count|per_min
                legend: Errors/min
                color: red
          - title: 4xx Errors
            columns:
              - span.status_code: [>=400, <500]
            chart:
              - type: line
                metric: span.count|per_min
                legend: 4xx/min
                color: orange
      - row:
          - title: 5xx Errors
            columns:
              - span.status_code: [>=500]
            chart:
              - type: line
                metric: span.count|per_min
                legend: 5xx/min
                color: red
          - title: Database Errors
            columns:
              - db.system: [*]
                span.status_code: [error]
            chart:
              - type: line
                metric: span.count|per_min
                legend: DB Errors/min
                color: purple
  YAML
}
