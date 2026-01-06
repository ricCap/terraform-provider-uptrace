# Notification channel configured for Uptrace cloud
resource "uptrace_notification_channel" "cloud_slack" {
  name = "Cloud Production Alerts"
  type = "slack"

  # Priority is required for cloud API
  # Omit this field for self-hosted Uptrace
  priority = ["high", "critical"]

  params = {
    webhookUrl = var.slack_webhook_url
  }
}
