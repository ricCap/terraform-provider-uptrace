# Slack notification channel
resource "uptrace_notification_channel" "slack" {
  name = "Engineering Alerts"
  type = "slack"

  params = {
    webhookUrl = var.slack_webhook_url
  }
}

# Webhook notification channel
resource "uptrace_notification_channel" "webhook" {
  name = "Custom Webhook"
  type = "webhook"

  params = {
    url = "https://example.com/webhook"
  }
}

# Telegram notification channel
resource "uptrace_notification_channel" "telegram" {
  name = "Telegram Alerts"
  type = "telegram"

  params = {
    botToken = var.telegram_bot_token
    chatId   = "-1001234567890"
  }
}

# Mattermost notification channel
resource "uptrace_notification_channel" "mattermost" {
  name = "Mattermost Alerts"
  type = "mattermost"

  params = {
    webhookUrl = var.mattermost_webhook_url
  }
}

# Channel with condition expression
# NOTE: The valid condition syntax for notification channels is not yet documented.
# The condition field is optional and may filter which alerts trigger this channel.
# Uncomment and test with your specific use case.
#
# resource "uptrace_notification_channel" "conditional_slack" {
#   name      = "Critical Alerts Only"
#   type      = "slack"
#   condition = "your_condition_here"
#
#   params = {
#     webhookUrl = var.slack_webhook_url
#   }
# }

# Use notification channel in a monitor
resource "uptrace_monitor" "high_error_rate" {
  name = "High Error Rate"
  type = "error"

  channel_ids = [
    uptrace_notification_channel.slack.id,
    uptrace_notification_channel.telegram.id,
  ]

  params = {
    metrics = []
    query   = "error"
  }
}
