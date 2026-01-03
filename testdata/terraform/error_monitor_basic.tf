variable "monitor_name" {
  type = string
}

resource "uptrace_monitor" "test" {
  name = var.monitor_name
  type = "error"

  notify_everyone_by_email = false

  params = {
    metrics = [
      {
        name  = "uptrace_tracing_events"
        alias = "$logs"
      }
    ]
    query = "sum($logs) | where span.event_name exists"
  }
}

output "monitor_id" {
  value = uptrace_monitor.test.id
}

output "monitor_state" {
  value = uptrace_monitor.test.state
}
