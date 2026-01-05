# Uptrace Cloud API Support

## Overview

The provider supports both self-hosted Uptrace instances and the Uptrace cloud API at `api2.uptrace.dev`.

## Cloud-Specific Fields

### Notification Channels

The `priority` field is **required** for cloud API but not used by self-hosted:

```hcl
resource "uptrace_notification_channel" "cloud_slack" {
  name = "Production Alerts"
  type = "slack"

  # Required for cloud, omit for self-hosted
  priority = ["high", "critical"]

  params = {
    webhookUrl = var.slack_webhook_url
  }
}
```

**Valid Priority Values**: To be documented after cloud testing (see Phase 9 of implementation plan).

### Monitors

The `trend_agg_func` field is **required** for cloud monitors (both metric and error types):

```hcl
resource "uptrace_monitor" "cloud_latency" {
  name = "High Latency"
  type = "metric"

  notify_everyone_by_email = false

  # Required for cloud API, optional for self-hosted v2.0.2 and earlier
  trend_agg_func = "avg"

  params = {
    metrics = [{
      name  = "span.duration"
      alias = "$latency"
    }]

    query  = "span.kind:server"
    column = "$latency"

    max_allowed_value = 1000000000
    check_num_point   = 2
  }
}
```

**Valid Aggregation Functions**: `avg`, `sum`, `min`, `max`, `p50`, `p90`, `p95`, `p99`

**Note**: For self-hosted Uptrace v2.0.2 and earlier, this field is not used and can be omitted.

## Configuration

### Cloud API
```hcl
provider "uptrace" {
  endpoint   = "https://api2.uptrace.dev/internal/v1"
  token      = var.uptrace_cloud_token
  project_id = var.uptrace_cloud_project_id
}
```

### Self-Hosted
```hcl
provider "uptrace" {
  endpoint   = "http://localhost:14318/internal/v1"
  token      = var.uptrace_local_token
  project_id = 1
}
```

**Error Monitor Example:**

```hcl
resource "uptrace_monitor" "cloud_errors" {
  name = "High Error Rate"
  type = "error"

  notify_everyone_by_email = false

  # Required for cloud API
  trend_agg_func = "sum"

  params = {
    metrics = [{
      name  = "uptrace_tracing_events"
      alias = "$events"
    }]
    query = "where span.system = 'production'"
  }
}
```

## API Differences

The cloud API (`api2.uptrace.dev`) has additional requirements compared to self-hosted instances:

1. **Notification Channels**: Must include `priority` field
2. **Monitors**: Must include `trend_agg_func` field (both metric and error types)
3. **Dashboards**: Stricter metric validation (metrics must exist in project)
4. **Query Normalization**: Cloud API automatically reformats UQL queries to canonical form

All cloud-specific fields are optional in the provider schema to maintain backward compatibility with self-hosted instances.

## Testing

To test provider against the cloud API, send test telemetry data first:

```bash
task test:cloud:send-telemetry
# Wait 30 seconds for processing
task test:cloud
```
