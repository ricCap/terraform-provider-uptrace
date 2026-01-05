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
    bot_token = var.telegram_bot_token
    chat_id   = var.telegram_chat_id
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
    metrics = [
      {
        name = "span.count"
      }
    ]
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
    metrics = [
      {
        name = "span.count"
      }
    ]
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
    metrics = [
      {
        name = "span.count"
      }
    ]
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
        alias = "$p95_latency"
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
        alias = "$p95_duration"
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
        alias = "$request_count"
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
        alias = "$request_count"
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
  yaml = <<-YAML
    schema: v2
    name: Application Overview
    grid_rows:
      - title: API Performance
        items:
          - title: Request Rate
            metrics:
              - span.count as $requests
            query:
              - per_min(sum($requests))
            where:
              - span.system
              - =
              - http
              - span.kind
              - =
              - server
          - title: P95 Latency
            metrics:
              - span.duration as $duration
            query:
              - p95($duration)
            where:
              - span.system
              - =
              - http
              - span.kind
              - =
              - server
      - title: Errors & Database
        items:
          - title: Error Rate
            metrics:
              - span.count as $errors
            query:
              - per_min(sum($errors))
            where:
              - span.status_code
              - =
              - error
          - title: Database Query Performance
            metrics:
              - span.duration as $db_duration
            query:
              - p95($db_duration)
            where:
              - db.system
              - exists
              - true
  YAML
}

# Error tracking dashboard
resource "uptrace_dashboard" "error_tracking" {
  yaml = <<-YAML
    schema: v2
    name: Error Tracking Dashboard
    grid_rows:
      - title: HTTP Errors
        items:
          - title: Total Errors
            metrics:
              - span.count as $errors
            query:
              - per_min(sum($errors))
            where:
              - span.status_code
              - =
              - error
          - title: 4xx Client Errors
            metrics:
              - span.count as $client_errors
            query:
              - per_min(sum($client_errors))
            where:
              - http.status_code
              - '>='
              - 400
              - http.status_code
              - <
              - 500
      - title: Server Errors
        items:
          - title: 5xx Server Errors
            metrics:
              - span.count as $server_errors
            query:
              - per_min(sum($server_errors))
            where:
              - http.status_code
              - '>='
              - 500
          - title: Database Errors
            metrics:
              - span.count as $db_errors
            query:
              - per_min(sum($db_errors))
            where:
              - db.system
              - exists
              - true
              - span.status_code
              - =
              - error
  YAML
}
