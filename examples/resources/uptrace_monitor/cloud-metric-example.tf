# Metric monitor configured for Uptrace cloud
resource "uptrace_monitor" "cloud_latency" {
  name = "Cloud API Latency Monitor"
  type = "metric"

  notify_everyone_by_email = false

  # Trend aggregation function is required for cloud API
  # Omit this field for self-hosted Uptrace v2.0.2 and earlier
  # Valid values: avg, sum, min, max, p50, p90, p95, p99
  trend_agg_func = "avg"

  params = {
    metrics = [{
      name  = "span.duration"
      alias = "$latency"
    }]

    query  = "span.kind:server"
    column = "$latency"

    max_allowed_value = 2000000000 # 2 seconds in nanoseconds
    check_num_point   = 2
  }
}
