# Metric monitor example
resource "uptrace_monitor" "high_cpu" {
  name = "High CPU Usage"
  type = "metric"

  notify_everyone_by_email = false
  channel_ids              = [1, 2]

  params = {
    metrics = [
      {
        name  = "system.cpu.utilization"
        alias = "cpu_usage"
      }
    ]
    query              = "avg(cpu_usage) > 90"
    column             = "cpu_usage"
    max_allowed_value  = 90
    grouping_interval  = 60000
    check_num_point    = 3
    nulls_mode         = "allow"
  }
}

# Error monitor example
resource "uptrace_monitor" "api_errors" {
  name = "API Error Rate"
  type = "error"

  notify_everyone_by_email = true
  team_ids                 = [1]

  params = {
    metrics = [
      {
        name = "uptrace_tracing_logs"
      }
    ]
    query = "severity:ERROR AND service:api"
  }
}
