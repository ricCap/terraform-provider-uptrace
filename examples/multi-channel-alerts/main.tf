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
# Pattern 1: Severity-Based Routing
# Different channels for different severity levels
# ============================================================================

resource "uptrace_notification_channel" "sev1_pagerduty" {
  name = "SEV1 - PagerDuty"
  type = "webhook"

  params = {
    url = var.pagerduty_url
  }
}

resource "uptrace_notification_channel" "sev2_slack" {
  name = "SEV2 - Slack Oncall"
  type = "slack"

  params = {
    webhookUrl = var.slack_oncall_webhook
  }
}

resource "uptrace_notification_channel" "sev3_slack" {
  name = "SEV3 - Slack Engineering"
  type = "slack"

  params = {
    webhookUrl = var.slack_engineering_webhook
  }
}

# SEV1: Critical - Page immediately
resource "uptrace_monitor" "sev1_complete_outage" {
  name = "SEV1: Complete Service Outage"
  type = "metric"

  channel_ids = [
    uptrace_notification_channel.sev1_pagerduty.id,
    uptrace_notification_channel.sev2_slack.id, # Also notify Slack
  ]

  params = {
    metrics = [
      {
        name  = "span.count"
        alias = "requests"
      }
    ]

    query = "span.system:http span.kind:server"

    # No traffic for 5 minutes
    min_allowed_value = 1
    max_allowed_value = 999999

    check_num_point = 3

    column = "$requests"

    notify_everyone_by_email = true
  }
}

# SEV2: High - Slack on-call channel
resource "uptrace_monitor" "sev2_high_error_rate" {
  name = "SEV2: High Error Rate"
  type = "error"

  channel_ids = [
    uptrace_notification_channel.sev2_slack.id,
  ]

  params = {
    metrics = []
    query   = "span.status_code:error"

    max_allowed_value = 100

    check_num_point = 2

    notify_everyone_by_email = false
  }
}

# SEV3: Medium - Engineering team awareness
resource "uptrace_monitor" "sev3_elevated_latency" {
  name = "SEV3: Elevated API Latency"
  type = "metric"

  channel_ids = [
    uptrace_notification_channel.sev3_slack.id,
  ]

  params = {
    metrics = [
      {
        name  = "span.duration"
        alias = "latency"
      }
    ]

    query = "span.system:http span.kind:server"

    max_allowed_value = 1000000000 # 1 second

    check_num_point = 3

    column = "$latency"

    notify_everyone_by_email = false
  }
}

# ============================================================================
# Pattern 2: Team-Based Routing
# Different channels for different teams
# ============================================================================

resource "uptrace_notification_channel" "team_backend" {
  name = "Team: Backend Engineering"
  type = "slack"

  params = {
    webhookUrl = var.slack_backend_webhook
  }
}

resource "uptrace_notification_channel" "team_frontend" {
  name = "Team: Frontend Engineering"
  type = "slack"

  params = {
    webhookUrl = var.slack_frontend_webhook
  }
}

resource "uptrace_notification_channel" "team_data" {
  name = "Team: Data Engineering"
  type = "slack"

  params = {
    webhookUrl = var.slack_data_webhook
  }
}

# Backend team: API errors
resource "uptrace_monitor" "backend_api_errors" {
  name = "Backend: API Server Errors"
  type = "error"

  channel_ids = [
    uptrace_notification_channel.team_backend.id,
  ]

  params = {
    metrics = []
    query   = "service.name:api-server span.status_code:error"

    max_allowed_value = 50

    check_num_point = 2

    notify_everyone_by_email = false
  }
}

# Frontend team: Client-side errors
resource "uptrace_monitor" "frontend_client_errors" {
  name = "Frontend: Client-Side Errors"
  type = "error"

  channel_ids = [
    uptrace_notification_channel.team_frontend.id,
  ]

  params = {
    metrics = []
    query   = "service.name:web-app span.status_code:error"

    max_allowed_value = 100

    check_num_point = 2

    notify_everyone_by_email = false
  }
}

# Data team: Database and pipeline errors
resource "uptrace_monitor" "data_pipeline_errors" {
  name = "Data: Pipeline Errors"
  type = "error"

  channel_ids = [
    uptrace_notification_channel.team_data.id,
  ]

  params = {
    metrics = []
    query   = "service.name:data-pipeline span.status_code:error"

    max_allowed_value = 10

    check_num_point = 1

    notify_everyone_by_email = false
  }
}

# ============================================================================
# Pattern 3: Multi-Channel Fan-Out
# Send to multiple channels for critical issues
# ============================================================================

resource "uptrace_notification_channel" "fanout_slack" {
  name = "Fanout: Slack"
  type = "slack"

  params = {
    webhookUrl = var.slack_critical_webhook
  }
}

resource "uptrace_notification_channel" "fanout_telegram" {
  name = "Fanout: Telegram"
  type = "telegram"

  params = {
    botToken = var.telegram_bot_token
    chatId   = var.telegram_chat_id
  }
}

resource "uptrace_notification_channel" "fanout_webhook" {
  name = "Fanout: Custom Webhook"
  type = "webhook"

  params = {
    url = var.custom_webhook_url
  }
}

# Critical alert goes to all channels
resource "uptrace_monitor" "critical_payment_failures" {
  name = "CRITICAL: Payment Processing Failures"
  type = "error"

  # Fan out to ALL channels
  channel_ids = [
    uptrace_notification_channel.fanout_slack.id,
    uptrace_notification_channel.fanout_telegram.id,
    uptrace_notification_channel.fanout_webhook.id,
    uptrace_notification_channel.sev1_pagerduty.id, # Also page
  ]

  params = {
    metrics = []
    query   = "service.name:payment-service span.status_code:error"

    max_allowed_value = 5

    check_num_point = 1

    notify_everyone_by_email = true
  }
}

# ============================================================================
# Pattern 4: Environment-Specific Routing
# Different channels for production vs staging
# ============================================================================

resource "uptrace_notification_channel" "prod_alerts" {
  name = "Production Alerts"
  type = "slack"

  params = {
    webhookUrl = var.slack_production_webhook
  }
}

resource "uptrace_notification_channel" "staging_alerts" {
  name = "Staging Alerts"
  type = "slack"

  params = {
    webhookUrl = var.slack_staging_webhook
  }
}

# Production monitor
resource "uptrace_monitor" "prod_error_rate" {
  name = "Production: Error Rate"
  type = "error"

  channel_ids = [
    uptrace_notification_channel.prod_alerts.id,
  ]

  params = {
    metrics = []
    query   = "deployment.environment:production span.status_code:error"

    max_allowed_value = 50

    check_num_point = 2

    notify_everyone_by_email = false
  }
}

# Staging monitor (less sensitive)
resource "uptrace_monitor" "staging_error_rate" {
  name = "Staging: Error Rate"
  type = "error"

  channel_ids = [
    uptrace_notification_channel.staging_alerts.id,
  ]

  params = {
    metrics = []
    query   = "deployment.environment:staging span.status_code:error"

    max_allowed_value = 200

    check_num_point = 5

    notify_everyone_by_email = false
  }
}

# ============================================================================
# Pattern 5: Time-Based Routing (Business Hours vs On-Call)
# Use different channels based on urgency
# ============================================================================

resource "uptrace_notification_channel" "business_hours" {
  name = "Business Hours - Slack"
  type = "slack"

  params = {
    webhookUrl = var.slack_business_hours_webhook
  }
}

resource "uptrace_notification_channel" "after_hours" {
  name = "After Hours - PagerDuty"
  type = "webhook"

  params = {
    url = var.pagerduty_after_hours_url
  }
}

# During business hours: Slack is fine
resource "uptrace_monitor" "degraded_performance" {
  name = "Degraded Performance (Business Hours Alert)"
  type = "metric"

  channel_ids = [
    uptrace_notification_channel.business_hours.id,
  ]

  params = {
    metrics = [
      {
        name  = "span.duration"
        alias = "latency"
      }
    ]

    query = "span.system:http span.kind:server"

    max_allowed_value = 3000000000 # 3 seconds

    check_num_point = 5

    column = "$latency"

    notify_everyone_by_email = false
  }
}

# After hours: Critical issues page on-call
resource "uptrace_monitor" "after_hours_outage" {
  name = "Critical Outage (Page On-Call)"
  type = "error"

  channel_ids = [
    uptrace_notification_channel.after_hours.id,
    uptrace_notification_channel.business_hours.id, # Also post to Slack
  ]

  params = {
    metrics = []
    query   = "span.status_code:error"

    max_allowed_value = 100

    check_num_point = 2

    notify_everyone_by_email = true
  }
}
