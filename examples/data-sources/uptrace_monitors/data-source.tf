# List all monitors in the project
data "uptrace_monitors" "all" {}

output "all_monitors_count" {
  description = "Total number of monitors"
  value       = length(data.uptrace_monitors.all.monitors)
}

output "all_monitor_names" {
  description = "Names of all monitors"
  value       = [for m in data.uptrace_monitors.all.monitors : m.name]
}

# Filter monitors by type
data "uptrace_monitors" "metric_monitors" {
  type = "metric"
}

output "metric_monitors_count" {
  description = "Number of metric monitors"
  value       = length(data.uptrace_monitors.metric_monitors.monitors)
}

# Filter monitors by state
data "uptrace_monitors" "firing_monitors" {
  state = "firing"
}

output "firing_monitors" {
  description = "List of currently firing monitors"
  value = [
    for m in data.uptrace_monitors.firing_monitors.monitors : {
      name  = m.name
      type  = m.type
      state = m.state
    }
  ]
}

# Filter monitors by name (substring match)
data "uptrace_monitors" "cpu_monitors" {
  name = "CPU"
}

output "cpu_monitors" {
  description = "Monitors with 'CPU' in the name"
  value = [
    for m in data.uptrace_monitors.cpu_monitors.monitors : {
      id   = m.id
      name = m.name
      type = m.type
    }
  ]
}

# Combine multiple filters
data "uptrace_monitors" "critical_metric_monitors" {
  type  = "metric"
  state = "firing"
  name  = "critical"
}

output "critical_alerts" {
  description = "Critical metric monitors currently firing"
  value       = data.uptrace_monitors.critical_metric_monitors.monitors
}

# Use in locals for processing
locals {
  all_monitors = data.uptrace_monitors.all.monitors

  # Group monitors by type
  monitors_by_type = {
    for m in local.all_monitors :
    m.type => m...
  }

  # Get IDs of all firing monitors
  firing_monitor_ids = [
    for m in local.all_monitors :
    m.id if m.state == "firing"
  ]

  # Count monitors by state
  monitors_by_state = {
    for m in local.all_monitors :
    m.state => length([for mon in local.all_monitors : mon if mon.state == m.state])...
  }
}

output "monitors_summary" {
  description = "Summary of monitors by type and state"
  value = {
    total              = length(local.all_monitors)
    by_type            = { for k, v in local.monitors_by_type : k => length(v) }
    firing_count       = length(local.firing_monitor_ids)
    firing_monitor_ids = local.firing_monitor_ids
  }
}

# Reference monitors created in the same configuration
resource "uptrace_monitor" "example" {
  name = "Example Monitor"
  type = "metric"

  params = {
    metrics = [{
      name  = "system.cpu.utilization"
      alias = "$cpu"
    }]
    query             = "avg($cpu) > 90"
    max_allowed_value = 90
  }
}

# List all monitors including the one just created
data "uptrace_monitors" "with_new_monitor" {
  depends_on = [uptrace_monitor.example]
}

# Verify the new monitor appears in the list
output "new_monitor_exists" {
  description = "Check if the new monitor is in the list"
  value = anytrue([
    for m in data.uptrace_monitors.with_new_monitor.monitors :
    m.name == "Example Monitor"
  ])
}
