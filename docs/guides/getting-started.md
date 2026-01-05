# Getting Started with the Uptrace Terraform Provider

This guide will help you get started with the Uptrace Terraform provider for managing monitors, dashboards, and notification channels.

## Prerequisites

Before you begin, ensure you have:

- [Terraform](https://www.terraform.io/downloads) or [OpenTofu](https://opentofu.org/) >= 1.0
- Access to an Uptrace instance
- Uptrace API token with appropriate permissions
- Project ID for your Uptrace project

## Installation

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    uptrace = {
      source = "registry.terraform.io/riccap/uptrace"
    }
  }
}
```

## Provider Configuration

Configure the provider with your Uptrace credentials:

```hcl
provider "uptrace" {
  endpoint   = "https://uptrace.example.com/api/v1"
  token      = var.uptrace_token
  project_id = 1
}
```

### Configuration Options

The provider supports three configuration methods (in order of precedence):

1. **Terraform variables** (as shown above)
2. **Environment variables**:
   - `UPTRACE_ENDPOINT`
   - `UPTRACE_TOKEN`
   - `UPTRACE_PROJECT_ID`
3. **Configuration file** (if supported by your Uptrace installation)

### Best Practices

- **Never commit tokens**: Use variables and mark them as sensitive
- **Use environment variables in CI/CD**: Set `UPTRACE_*` environment variables in your pipeline
- **Separate configurations**: Use different providers for different projects

Example with sensitive variables:

```hcl
variable "uptrace_token" {
  description = "Uptrace API authentication token"
  type        = string
  sensitive   = true
}

variable "uptrace_project_id" {
  description = "Uptrace project ID"
  type        = number
}

provider "uptrace" {
  endpoint   = "https://uptrace.example.com/api/v1"
  token      = var.uptrace_token
  project_id = var.uptrace_project_id
}
```

## Quick Start: Creating Your First Monitor

### Step 1: Create a Notification Channel

First, set up where alerts should be sent:

```hcl
resource "uptrace_notification_channel" "slack_alerts" {
  name = "Engineering Alerts"
  type = "slack"

  params = {
    webhookUrl = "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
  }
}
```

### Step 2: Create an Error Monitor

Monitor for application errors:

```hcl
resource "uptrace_monitor" "high_error_rate" {
  name = "High Error Rate"
  type = "error"

  channel_ids = [
    uptrace_notification_channel.slack_alerts.id
  ]

  params = {
    metrics = [
      {
        name = "span.count"
      }
    ]
    query = "span.status_code:error"
  }
}
```

### Step 3: Create a Performance Monitor

Monitor API response times:

```hcl
resource "uptrace_monitor" "slow_api_responses" {
  name = "Slow API Responses"
  type = "metric"

  channel_ids = [
    uptrace_notification_channel.slack_alerts.id
  ]

  params = {
    metrics = [
      {
        name  = "span.duration"
        alias = "$p95_duration"
      }
    ]

    query = "span.system:http span.kind:server"

    # Alert if P95 > 2 seconds
    min_allowed_value = 0
    max_allowed_value = 2000000000  # nanoseconds

    # Column to monitor
    column = "$p95_duration"

    # Alert after 2 consecutive checks
    check_num_point = 2
  }
}
```

### Step 4: Create a Dashboard

Visualize your metrics:

```hcl
resource "uptrace_dashboard" "api_overview" {
  yaml = <<-YAML
    schema: v2
    name: API Overview
    grid_rows:
      - title: Traffic & Performance
        items:
          - title: Request Rate
            metrics:
              - span.count as $requests
            query:
              - per_min(sum($requests))
            where:
              - span.system
              - =
              - http
              - span.kind
              - =
              - server
          - title: P95 Latency
            metrics:
              - span.duration as $duration
            query:
              - p95($duration)
            where:
              - span.system
              - =
              - http
      - title: Errors
        items:
          - title: Error Rate
            metrics:
              - span.count as $errors
            query:
              - per_min(sum($errors))
            where:
              - span.status_code
              - =
              - error
  YAML
}
```

### Step 5: Apply the Configuration

```bash
# Initialize Terraform
terraform init

# Preview changes
terraform plan

# Apply changes
terraform apply
```

## Understanding Monitor Types

### Error Monitors

Error monitors count error occurrences based on query filters:

```hcl
resource "uptrace_monitor" "database_errors" {
  name = "Database Errors"
  type = "error"

  channel_ids = [uptrace_notification_channel.slack_alerts.id]

  params = {
    metrics = [
      {
        name = "span.count"
      }
    ]
    query = "db.system:postgresql span.status_code:error"

    # Alert if more than 10 errors in check interval
    min_allowed_value = 0
    max_allowed_value = 10

    # Check every interval
    check_num_point = 1
  }
}
```

### Metric Monitors

Metric monitors track numerical values with aggregations:

```hcl
resource "uptrace_monitor" "high_cpu" {
  name = "High CPU Usage"
  type = "metric"

  channel_ids = [uptrace_notification_channel.slack_alerts.id]

  params = {
    metrics = [
      {
        name  = "system.cpu.utilization"
        alias = "$cpu_usage"
      }
    ]

    query = "service.name:my-service"

    # Alert if CPU > 90%
    max_allowed_value = 90

    # Column to monitor (must match alias)
    column = "$cpu_usage"

    # Alert after 3 consecutive high readings
    check_num_point = 3
  }
}
```

## Notification Channel Types

### Slack

```hcl
resource "uptrace_notification_channel" "slack" {
  name = "Slack Alerts"
  type = "slack"

  params = {
    webhookUrl = "https://hooks.slack.com/services/..."
  }
}
```

### Generic Webhook

```hcl
resource "uptrace_notification_channel" "webhook" {
  name = "Custom Webhook"
  type = "webhook"

  params = {
    url = "https://example.com/webhook"
  }
}
```

### Telegram

```hcl
resource "uptrace_notification_channel" "telegram" {
  name = "Telegram Alerts"
  type = "telegram"

  params = {
    bot_token = "123456:ABC-DEF..."
    chat_id   = "-1001234567890"
  }
}
```

### Mattermost

```hcl
resource "uptrace_notification_channel" "mattermost" {
  name = "Mattermost Alerts"
  type = "mattermost"

  params = {
    webhookUrl = "https://mattermost.example.com/hooks/..."
  }
}
```

## Dashboard YAML Format

Dashboards use a YAML-based configuration:

### Basic Structure

```yaml
schema: v2
name: Dashboard Name
grid_rows:
  - title: Row Title
    items:
      - title: Chart Title
        metrics:
          - metric.name as $alias
        query:
          - aggregation($alias)
        where:
          - attribute
          - operator
          - value
```

### Common Metrics

**Request Rate:**
```yaml
metrics:
  - span.count as $requests
query:
  - per_min(sum($requests))
```

**Latency Percentiles:**
```yaml
metrics:
  - span.duration as $duration
query:
  - p95($duration)  # or p50, p99
```

**Error Rate:**
```yaml
metrics:
  - span.count as $errors
query:
  - per_min(sum($errors))
where:
  - span.status_code
  - =
  - error
```

### Common Filters (where clauses)

**By Service:**
```yaml
where:
  - service.name
  - =
  - my-service
```

**By HTTP Method:**
```yaml
where:
  - http.method
  - =
  - GET
```

**By Status Code Range:**
```yaml
where:
  - http.status_code
  - '>='
  - 500
```

**Multiple Conditions (AND):**
```yaml
where:
  - span.system
  - =
  - http
  - span.kind
  - =
  - server
```

## Data Sources

### Query Single Monitor

```hcl
data "uptrace_monitor" "existing" {
  id = 123
}

output "monitor_name" {
  value = data.uptrace_monitor.existing.name
}
```

### Query Multiple Monitors

```hcl
data "uptrace_monitors" "error_monitors" {
  type = "error"
}

output "error_monitor_count" {
  value = length(data.uptrace_monitors.error_monitors.monitors)
}
```

## Importing Existing Resources

### Import a Monitor

```bash
terraform import uptrace_monitor.existing <monitor_id>
```

### Import a Dashboard

```bash
terraform import uptrace_dashboard.existing <dashboard_id>
```

### Import a Notification Channel

```bash
terraform import uptrace_notification_channel.existing <channel_id>
```

## Common Patterns

### Multiple Notification Channels

Route different severity levels to different channels:

```hcl
resource "uptrace_monitor" "critical_error" {
  name = "Critical: Service Down"
  type = "error"

  channel_ids = [
    uptrace_notification_channel.pagerduty.id,
    uptrace_notification_channel.slack_critical.id,
    uptrace_notification_channel.telegram_oncall.id,
  ]

  params = {
    metrics = [{ name = "span.count" }]
    query   = "span.status_code:error"
  }
}
```

### Conditional Monitoring

Monitor different services with different thresholds:

```hcl
resource "uptrace_monitor" "api_latency" {
  name = "API Latency"
  type = "metric"

  channel_ids = [uptrace_notification_channel.slack_alerts.id]

  params = {
    metrics = [{
      name  = "span.duration"
      alias = "$duration"
    }]

    query = "service.name:api-gateway span.kind:server"

    max_allowed_value = 1000000000  # 1 second
    column            = "$duration"
    check_num_point   = 2
  }
}

resource "uptrace_monitor" "background_job_latency" {
  name = "Background Job Latency"
  type = "metric"

  channel_ids = [uptrace_notification_channel.slack_alerts.id]

  params = {
    metrics = [{
      name  = "span.duration"
      alias = "$duration"
    }]

    query = "service.name:worker span.kind:internal"

    max_allowed_value = 10000000000  # 10 seconds
    column            = "$duration"
    check_num_point   = 3
  }
}
```

## Troubleshooting

### Monitor Not Triggering

1. **Check the query**: Verify it matches your data in Uptrace UI
2. **Check thresholds**: Ensure `min_allowed_value` and `max_allowed_value` are correct
3. **Check check_num_point**: High values require multiple consecutive violations

### Dashboard Shows No Data

1. **Verify attribute names**: Check they match your telemetry data
2. **Check where clauses**: Ensure filters match existing data
3. **Test in Uptrace UI**: Build the query there first

### Notification Channel Not Working

1. **Test webhook URL**: Verify it's accessible
2. **Check channel status**: Look at the `status` computed attribute
3. **Verify params**: Ensure all required params are set correctly

### Validation Errors

**"metric alias must start with the dollar sign":**
```hcl
# ❌ Wrong
alias = "duration"

# ✅ Correct
alias = "$duration"
```

**"at least one metric is required":**
```hcl
# ❌ Wrong (error monitor)
params = {
  metrics = []
  query   = "span.status_code:error"
}

# ✅ Correct
params = {
  metrics = [{ name = "span.count" }]
  query   = "span.status_code:error"
}
```

## Next Steps

- Review [Best Practices Guide](best-practices.md) for production recommendations
- Explore [Dashboard Examples](../../examples/dashboard-examples/) for visualization patterns
- Check [Resource Examples](../../examples/resources/) for detailed configurations
- Read the [Uptrace Documentation](https://uptrace.dev/docs/) for query syntax and metrics

## Getting Help

- [GitHub Issues](https://github.com/riccap/tofu-uptrace-provider/issues) - Report bugs or request features
- [Uptrace Community](https://github.com/uptrace/uptrace/discussions) - General Uptrace questions
- [Provider Documentation](../../docs/) - Detailed resource reference
