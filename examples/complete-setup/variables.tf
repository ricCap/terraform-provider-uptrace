# Uptrace Connection
variable "uptrace_endpoint" {
  description = "Uptrace API endpoint"
  type        = string
  default     = "https://uptrace.example.com/api/v1"
}

variable "uptrace_token" {
  description = "Uptrace API authentication token"
  type        = string
  sensitive   = true
}

variable "uptrace_project_id" {
  description = "Uptrace project ID"
  type        = number
}

# Notification Channels
variable "slack_critical_webhook" {
  description = "Slack webhook URL for critical alerts"
  type        = string
  sensitive   = true
}

variable "slack_warnings_webhook" {
  description = "Slack webhook URL for warning alerts"
  type        = string
  sensitive   = true
}

variable "pagerduty_webhook_url" {
  description = "PagerDuty webhook URL for incidents"
  type        = string
  sensitive   = true
}

variable "telegram_bot_token" {
  description = "Telegram bot token for notifications"
  type        = string
  sensitive   = true
}

variable "telegram_chat_id" {
  description = "Telegram chat ID for notifications"
  type        = string
  sensitive   = true
}
