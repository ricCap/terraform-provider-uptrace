# Uptrace Configuration
variable "uptrace_endpoint" {
  description = "Uptrace API endpoint"
  type        = string
}

variable "uptrace_token" {
  description = "Uptrace API token"
  type        = string
  sensitive   = true
}

variable "uptrace_project_id" {
  description = "Uptrace project ID"
  type        = number
}

# Severity-Based Channels
variable "pagerduty_url" {
  description = "PagerDuty webhook URL for SEV1 alerts"
  type        = string
  sensitive   = true
}

variable "slack_oncall_webhook" {
  description = "Slack webhook for on-call team (SEV2)"
  type        = string
  sensitive   = true
}

variable "slack_engineering_webhook" {
  description = "Slack webhook for engineering team (SEV3)"
  type        = string
  sensitive   = true
}

# Team-Based Channels
variable "slack_backend_webhook" {
  description = "Slack webhook for backend team"
  type        = string
  sensitive   = true
}

variable "slack_frontend_webhook" {
  description = "Slack webhook for frontend team"
  type        = string
  sensitive   = true
}

variable "slack_data_webhook" {
  description = "Slack webhook for data team"
  type        = string
  sensitive   = true
}

# Fan-Out Channels
variable "slack_critical_webhook" {
  description = "Slack webhook for critical alerts"
  type        = string
  sensitive   = true
}

variable "telegram_bot_token" {
  description = "Telegram bot token"
  type        = string
  sensitive   = true
}

variable "telegram_chat_id" {
  description = "Telegram chat ID"
  type        = string
  sensitive   = true
}

variable "custom_webhook_url" {
  description = "Custom webhook URL for integrations"
  type        = string
  sensitive   = true
}

# Environment-Specific Channels
variable "slack_production_webhook" {
  description = "Slack webhook for production alerts"
  type        = string
  sensitive   = true
}

variable "slack_staging_webhook" {
  description = "Slack webhook for staging alerts"
  type        = string
  sensitive   = true
}

# Time-Based Channels
variable "slack_business_hours_webhook" {
  description = "Slack webhook for business hours"
  type        = string
  sensitive   = true
}

variable "pagerduty_after_hours_url" {
  description = "PagerDuty URL for after-hours on-call"
  type        = string
  sensitive   = true
}
