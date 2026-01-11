# Uptrace Cloud API Support

## Overview

The provider supports both self-hosted Uptrace instances and the Uptrace cloud API at `api2.uptrace.dev`.

## Known Issues

The Uptrace Cloud API has stricter validation requirements than self-hosted. See the GitHub issues for current status:

| Issue | Description | Status |
|-------|-------------|--------|
| [#53](https://github.com/ricCap/terraform-provider-uptrace/issues/53) | Query single quotes normalized to double quotes causes state drift | Open |
| [#54](https://github.com/ricCap/terraform-provider-uptrace/issues/54) | Cloud API uses `priorities` (plural) not `priority` (singular) | Open |
| [#55](https://github.com/ricCap/terraform-provider-uptrace/issues/55) | Cloud API stricter validation requirements | Open |

## Cloud-Specific Requirements

### Monitors

The cloud API has additional validation requirements:

1. **`trend_agg_func` is required** - Must specify aggregation function
2. **`query` cannot be empty** - Must provide a valid UQL query
3. **Metrics must exist** - Referenced metrics must exist in the project
4. **`alias` is required** - Each metric must have an alias

```hcl
resource "uptrace_monitor" "cloud_monitor" {
  name = "High Latency"
  type = "metric"

  notify_everyone_by_email = false

  # Required for cloud API
  trend_agg_func = "avg"

  params = {
    metrics = [{
      name  = "span.duration"
      alias = "$latency"  # Required
    }]

    query             = "span.kind:server"  # Cannot be empty
    column            = "$latency"
    max_allowed_value = 1000000000
    check_num_point   = 2
  }
}
```

**Valid Aggregation Functions**: `avg`, `sum`, `min`, `max`, `p50`, `p90`, `p95`, `p99`

### Notification Channels

The `priority` field is required for cloud API but currently not functional due to a field name mismatch ([#54](https://github.com/ricCap/terraform-provider-uptrace/issues/54)).

**Workaround**: Use self-hosted Uptrace for notification channel testing until the issue is resolved.

### Query Normalization

The cloud API normalizes queries to canonical form, which can cause state drift ([#53](https://github.com/ricCap/terraform-provider-uptrace/issues/53)):

- Single quotes `'` are converted to double quotes `"`
- Query formatting may change

**Workaround**: Use double quotes in queries:

```hcl
# Instead of:
query = "where span.status_code = 'error'"

# Use:
query = "where span.status_code = \"error\""
```

Or use `lifecycle` to ignore changes:

```hcl
resource "uptrace_monitor" "example" {
  # ... configuration ...

  lifecycle {
    ignore_changes = [params]
  }
}
```

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

## Testing Against Cloud API

Before creating monitors on cloud, ensure the referenced metrics exist in your project by sending telemetry data first.

**Recommended approach**: Use the local dev environment for development and testing:

```bash
# Start local Uptrace
task dev:up

# Test your configuration
cd dev-env/terraform-test
tofu apply

# Stop when done
task dev:down
```
