# Read an existing monitor by ID
data "uptrace_monitor" "existing" {
  id = "123"
}

# Use the data source outputs
output "monitor_name" {
  description = "The name of the monitor"
  value       = data.uptrace_monitor.existing.name
}

output "monitor_state" {
  description = "Current state of the monitor (open, firing, paused)"
  value       = data.uptrace_monitor.existing.state
}

output "monitor_type" {
  description = "Monitor type (metric or error)"
  value       = data.uptrace_monitor.existing.type
}

# Common pattern: Reference a monitor created elsewhere
resource "uptrace_monitor" "example" {
  name = "Example CPU Monitor"
  type = "metric"

  params = {
    metrics = [{
      name  = "system.cpu.utilization"
      alias = "$cpu"
    }]
    query             = "avg($cpu) > 90"
    max_allowed_value = 90
    check_num_point   = 2
  }
}

# Read it back via data source
data "uptrace_monitor" "example" {
  id = uptrace_monitor.example.id
}

# Access monitor details
output "example_monitor_params" {
  description = "Monitor parameters"
  value       = data.uptrace_monitor.example.params
}

output "example_created_at" {
  description = "When the monitor was created"
  value       = data.uptrace_monitor.example.created_at
}
