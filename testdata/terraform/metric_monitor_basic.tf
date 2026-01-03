variable "monitor_name" {
  type = string
}

resource "uptrace_monitor" "test" {
  name = var.monitor_name
  type = "metric"

  notify_everyone_by_email = false

  params = {
    metrics = [
      {
        name  = "system.cpu.utilization"
        alias = "$cpu"
      }
    ]
    query             = "avg($cpu) > 80"
    max_allowed_value = 80
    check_num_point   = 2
  }
}

output "monitor_id" {
  value = uptrace_monitor.test.id
}

output "monitor_state" {
  value = uptrace_monitor.test.state
}
