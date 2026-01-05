# Notification Channel Outputs
output "notification_channels" {
  description = "Created notification channel IDs"
  value = {
    slack_critical      = uptrace_notification_channel.slack_critical.id
    slack_warnings      = uptrace_notification_channel.slack_warnings.id
    pagerduty          = uptrace_notification_channel.webhook_pagerduty.id
    telegram_oncall    = uptrace_notification_channel.telegram_oncall.id
  }
}

# Monitor Outputs
output "monitors" {
  description = "Created monitor IDs and details"
  value = {
    error_monitors = {
      critical_error_rate     = uptrace_monitor.critical_error_rate.id
      warning_4xx_errors      = uptrace_monitor.warning_4xx_errors.id
      critical_database_errors = uptrace_monitor.critical_database_errors.id
    }
    performance_monitors = {
      critical_api_latency    = uptrace_monitor.critical_api_latency.id
      warning_db_slow_queries = uptrace_monitor.warning_db_slow_queries.id
    }
    throughput_monitors = {
      warning_traffic_spike = uptrace_monitor.warning_traffic_spike.id
      critical_no_traffic   = uptrace_monitor.critical_no_traffic.id
    }
  }
}

# Dashboard Outputs
output "dashboards" {
  description = "Created dashboard IDs"
  value = {
    application_overview = uptrace_dashboard.application_overview.id
    error_tracking      = uptrace_dashboard.error_tracking.id
  }
}

# Summary
output "monitoring_summary" {
  description = "Summary of monitoring setup"
  value = {
    total_channels  = 4
    total_monitors  = 7
    total_dashboards = 2
    critical_alerts = 4
    warning_alerts  = 3
  }
}
