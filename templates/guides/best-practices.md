# Best Practices for Uptrace Terraform Provider

This guide provides recommendations for using the Uptrace Terraform provider effectively in production environments.

## Table of Contents

- [Security](#security)
- [Code Organization](#code-organization)
- [Monitor Configuration](#monitor-configuration)
- [Dashboard Design](#dashboard-design)
- [Notification Channels](#notification-channels)
- [State Management](#state-management)
- [Testing](#testing)
- [CI/CD Integration](#cicd-integration)
- [Performance](#performance)

## Security

### Never Commit Secrets

**❌ Bad:**
```hcl
provider "uptrace" {
  endpoint   = "https://uptrace.example.com/api/v1"
  token      = "user1_secret_token_hardcoded"
  project_id = 1
}
```

**✅ Good:**
```hcl
variable "uptrace_token" {
  description = "Uptrace API authentication token"
  type        = string
  sensitive   = true
}

provider "uptrace" {
  endpoint   = var.uptrace_endpoint
  token      = var.uptrace_token
  project_id = var.uptrace_project_id
}
```

### Sensitive Notification Params

Mark webhook URLs and tokens as sensitive:

```hcl
variable "slack_webhook_url" {
  description = "Slack webhook URL for alerts"
  type        = string
  sensitive   = true
}

resource "uptrace_notification_channel" "slack" {
  name = "Engineering Alerts"
  type = "slack"

  params = {
    webhookUrl = var.slack_webhook_url
  }
}
```

### Use Environment Variables in CI/CD

Set credentials via environment variables:

```bash
export UPTRACE_ENDPOINT="https://uptrace.example.com/api/v1"
export UPTRACE_TOKEN="${VAULT_UPTRACE_TOKEN}"
export UPTRACE_PROJECT_ID="1"

terraform apply -auto-approve
```

### Rotate Tokens Regularly

- Use short-lived tokens when possible
- Automate token rotation
- Audit token usage in Uptrace

## Code Organization

### File Structure

Organize your Terraform code logically:

```
monitoring/
├── main.tf                    # Provider configuration
├── variables.tf               # Input variables
├── outputs.tf                 # Outputs
├── terraform.tfvars.example   # Example values
├── channels.tf                # Notification channels
├── monitors-errors.tf         # Error monitors
├── monitors-performance.tf    # Performance monitors
├── monitors-business.tf       # Business metrics monitors
└── dashboards.tf              # Dashboards
```

### Module Structure for Reusability

Create reusable modules:

```
modules/
└── service-monitoring/
    ├── main.tf
    ├── variables.tf
    ├── outputs.tf
    └── README.md
```

**Module usage:**
```hcl
module "api_monitoring" {
  source = "./modules/service-monitoring"

  service_name    = "api-gateway"
  slack_channel   = uptrace_notification_channel.slack_critical.id
  latency_threshold = 1000000000  # 1 second
  error_threshold   = 100
}

module "worker_monitoring" {
  source = "./modules/service-monitoring"

  service_name    = "background-worker"
  slack_channel   = uptrace_notification_channel.slack_warnings.id
  latency_threshold = 10000000000  # 10 seconds
  error_threshold   = 10
}
```

### Naming Conventions

Use consistent naming:

```hcl
# Notification channels: <severity>_<service>_<type>
resource "uptrace_notification_channel" "critical_ops_slack" { }
resource "uptrace_notification_channel" "warning_eng_slack" { }

# Monitors: <severity>_<what>_<metric>
resource "uptrace_monitor" "critical_api_latency" { }
resource "uptrace_monitor" "warning_db_slow_queries" { }

# Dashboards: <service>_<purpose>
resource "uptrace_dashboard" "api_performance" { }
resource "uptrace_dashboard" "database_health" { }
```

## Monitor Configuration

### Alert Thresholds

Choose appropriate thresholds:

**Error Monitors:**
```hcl
resource "uptrace_monitor" "error_rate" {
  name = "High Error Rate"
  type = "error"

  params = {
    metrics = [{ name = "span.count" }]
    query   = "span.status_code:error"

    # Production: Tight threshold
    max_allowed_value = 10  # per interval

    # Staging: Looser threshold
    # max_allowed_value = 50
  }
}
```

**Performance Monitors:**
```hcl
resource "uptrace_monitor" "api_latency" {
  name = "API Response Time"
  type = "metric"

  params = {
    metrics = [{
      name  = "span.duration"
      alias = "$p95"
    }]

    # P95 < 2 seconds
    max_allowed_value = 2000000000

    column = "$p95"
  }
}
```

### Check Frequency

Balance alerting speed vs. noise:

```hcl
# Critical: Alert immediately
resource "uptrace_monitor" "service_down" {
  # ...
  params = {
    check_num_point = 1  # Alert on first occurrence
    # ...
  }
}

# Warning: Wait for confirmation
resource "uptrace_monitor" "elevated_latency" {
  # ...
  params = {
    check_num_point = 3  # Alert after 3 consecutive checks
    # ...
  }
}
```

### Query Optimization

Write efficient queries:

**❌ Avoid broad queries:**
```hcl
query = "span.status_code:error"  # All errors across all services
```

**✅ Be specific:**
```hcl
query = "service.name:api-gateway span.status_code:error"
```

**✅ Use multiple filters:**
```hcl
query = "service.name:api-gateway span.kind:server http.status_code:>=500"
```

### Metric Aliases

Always use descriptive aliases with $ prefix:

```hcl
params = {
  metrics = [
    {
      name  = "span.duration"
      alias = "$p95_latency"  # Clear and specific
    }
  ]

  column = "$p95_latency"
}
```

## Dashboard Design

### Dashboard Hierarchy

Organize dashboards by audience and detail level:

1. **Executive Dashboard** - High-level KPIs
   - Overall system health
   - Business metrics
   - SLO compliance

2. **Service Dashboards** - Per-service details
   - RED metrics (Rate, Errors, Duration)
   - Dependencies
   - Resource usage

3. **Troubleshooting Dashboards** - Deep dive
   - Detailed error breakdown
   - Slow queries
   - Status code distribution

### Chart Organization

Group related metrics:

```yaml
schema: v2
name: API Health Dashboard
grid_rows:
  # Row 1: Traffic
  - title: Traffic Metrics
    items:
      - title: Request Rate
        # ...
      - title: Requests by Endpoint
        # ...

  # Row 2: Performance
  - title: Latency
    items:
      - title: P50 Latency
        # ...
      - title: P95 Latency
        # ...
      - title: P99 Latency
        # ...

  # Row 3: Errors
  - title: Errors
    items:
      - title: Error Rate
        # ...
      - title: Error Types
        # ...
```

### Time Aggregations

Choose appropriate aggregations:

- **Traffic metrics**: `per_min(sum($requests))`
- **Business metrics**: `per_hour(sum($signups))`
- **Real-time monitoring**: `per_sec(sum($events))`

### Dashboard Performance

Keep dashboards fast:

**❌ Too many charts:**
```yaml
# 20+ charts = slow loading
grid_rows:
  - title: Everything
    items:
      - ...  # 20 charts
```

**✅ Focused dashboards:**
```yaml
# 6-8 charts per dashboard
# Create multiple dashboards instead
```

**❌ Broad queries:**
```yaml
where:
  - span.system  # All systems
  - exists
  - true
```

**✅ Specific filters:**
```yaml
where:
  - service.name
  - =
  - api-gateway
  - span.kind
  - =
  - server
```

## Notification Channels

### Alert Routing Strategy

Route alerts by severity and team:

```hcl
# Severity-based routing
resource "uptrace_notification_channel" "critical_pagerduty" {
  name = "Critical - PagerDuty"
  type = "webhook"
  params = { url = var.pagerduty_critical_url }
}

resource "uptrace_notification_channel" "warning_slack" {
  name = "Warnings - Slack"
  type = "slack"
  params = { webhookUrl = var.slack_warnings_webhook }
}

# Team-based routing
resource "uptrace_notification_channel" "backend_team" {
  name = "Backend Team Alerts"
  type = "slack"
  params = { webhookUrl = var.slack_backend_webhook }
}

resource "uptrace_notification_channel" "frontend_team" {
  name = "Frontend Team Alerts"
  type = "slack"
  params = { webhookUrl = var.slack_frontend_webhook }
}
```

### Fan-out Pattern

Send critical alerts to multiple channels:

```hcl
resource "uptrace_monitor" "critical_service_down" {
  name = "CRITICAL: Service Completely Down"
  type = "error"

  # Alert multiple channels simultaneously
  channel_ids = [
    uptrace_notification_channel.critical_pagerduty.id,
    uptrace_notification_channel.critical_slack.id,
    uptrace_notification_channel.oncall_telegram.id,
  ]

  params = {
    metrics = [{ name = "span.count" }]
    query   = "service.name:api-gateway span.status_code:error"
    max_allowed_value = 0  # Any error is critical
    check_num_point   = 1  # Alert immediately
  }
}
```

### Channel Testing

Test notification channels regularly:

```hcl
# Separate test channel for validation
resource "uptrace_notification_channel" "test_slack" {
  name = "Test Channel - Do Not Use for Production"
  type = "slack"
  params = { webhookUrl = var.slack_test_webhook }
}
```

## State Management

### Remote State

Always use remote state in production:

```hcl
terraform {
  backend "s3" {
    bucket = "company-terraform-state"
    key    = "monitoring/uptrace.tfstate"
    region = "us-east-1"

    dynamodb_table = "terraform-locks"
    encrypt        = true
  }
}
```

### State Locking

Prevent concurrent modifications:

- Use DynamoDB for AWS S3 backend
- Use integrated locking for Terraform Cloud
- Implement locking for custom backends

### Workspaces

Separate environments:

```bash
# Development
terraform workspace select dev
terraform apply -var-file=dev.tfvars

# Production
terraform workspace select prod
terraform apply -var-file=prod.tfvars
```

## Testing

### Validation

Test configurations before applying:

```bash
# Format check
terraform fmt -check -recursive

# Validation
terraform validate

# Plan review
terraform plan -out=plan.tfplan

# Review plan
terraform show plan.tfplan
```

### Acceptance Testing

Test in non-production first:

```bash
# Apply to staging
terraform workspace select staging
terraform apply -var-file=staging.tfvars

# Verify monitors trigger correctly
# Verify dashboards display data
# Test notification channels

# Then apply to production
terraform workspace select prod
terraform apply -var-file=prod.tfvars
```

### Monitor Testing

Verify monitor logic:

1. **Check query in Uptrace UI** - Ensure it matches expected data
2. **Test thresholds** - Trigger test scenarios
3. **Verify notifications** - Confirm channels receive alerts

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Deploy Monitoring

on:
  push:
    branches: [main]
    paths:
      - 'monitoring/**'

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - uses: hashicorp/setup-terraform@v3
        with:
          terraform_version: 1.6.0

      - name: Terraform Init
        working-directory: ./monitoring
        run: terraform init

      - name: Terraform Validate
        working-directory: ./monitoring
        run: terraform validate

      - name: Terraform Plan
        working-directory: ./monitoring
        env:
          UPTRACE_ENDPOINT: ${{ secrets.UPTRACE_ENDPOINT }}
          UPTRACE_TOKEN: ${{ secrets.UPTRACE_TOKEN }}
          UPTRACE_PROJECT_ID: ${{ secrets.UPTRACE_PROJECT_ID }}
        run: terraform plan -out=plan.tfplan

      - name: Terraform Apply
        if: github.ref == 'refs/heads/main'
        working-directory: ./monitoring
        env:
          UPTRACE_ENDPOINT: ${{ secrets.UPTRACE_ENDPOINT }}
          UPTRACE_TOKEN: ${{ secrets.UPTRACE_TOKEN }}
          UPTRACE_PROJECT_ID: ${{ secrets.UPTRACE_PROJECT_ID }}
        run: terraform apply -auto-approve plan.tfplan
```

### GitLab CI Example

```yaml
stages:
  - validate
  - plan
  - apply

variables:
  TF_ROOT: ${CI_PROJECT_DIR}/monitoring

.terraform:
  image: hashicorp/terraform:1.6
  before_script:
    - cd ${TF_ROOT}
    - terraform init

validate:
  extends: .terraform
  stage: validate
  script:
    - terraform fmt -check -recursive
    - terraform validate

plan:
  extends: .terraform
  stage: plan
  script:
    - terraform plan -out=plan.tfplan
  artifacts:
    paths:
      - ${TF_ROOT}/plan.tfplan
    expire_in: 1 day

apply:
  extends: .terraform
  stage: apply
  dependencies:
    - plan
  script:
    - terraform apply -auto-approve plan.tfplan
  only:
    - main
  when: manual
```

### Pre-commit Hooks

Catch issues before committing:

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/antonbabenko/pre-commit-terraform
    rev: v1.83.5
    hooks:
      - id: terraform_fmt
      - id: terraform_validate
      - id: terraform_docs
```

## Performance

### Provider Configuration

Optimize provider settings:

```hcl
provider "uptrace" {
  endpoint   = var.uptrace_endpoint
  token      = var.uptrace_token
  project_id = var.uptrace_project_id

  # Connection pooling handled internally
  # No additional configuration needed
}
```

### Parallel Operations

Terraform handles parallelism automatically, but you can tune it:

```bash
# Default: 10 concurrent operations
terraform apply

# Increase for faster deployments
terraform apply -parallelism=20

# Decrease if hitting rate limits
terraform apply -parallelism=5
```

### Resource Dependencies

Let Terraform manage dependencies:

```hcl
resource "uptrace_notification_channel" "slack" {
  name = "Alerts"
  type = "slack"
  params = { webhookUrl = var.slack_webhook }
}

resource "uptrace_monitor" "errors" {
  name = "Errors"
  type = "error"

  # Implicit dependency on channel
  channel_ids = [uptrace_notification_channel.slack.id]

  params = {
    metrics = [{ name = "span.count" }]
    query   = "span.status_code:error"
  }
}
```

## Common Pitfalls

### Monitor Alert Fatigue

**❌ Problem:**
```hcl
# Alerts too often
resource "uptrace_monitor" "any_error" {
  params = {
    max_allowed_value = 0     # Alerts on single error
    check_num_point   = 1     # Immediate alert
  }
}
```

**✅ Solution:**
```hcl
# Alert on sustained issues
resource "uptrace_monitor" "elevated_errors" {
  params = {
    max_allowed_value = 10    # Reasonable threshold
    check_num_point   = 3     # Confirm it's sustained
  }
}
```

### Dashboard Overload

**❌ Problem:**
```yaml
# Single dashboard with 30+ charts
grid_rows:
  - title: Everything
    items: [...30 charts...]
```

**✅ Solution:**
```yaml
# Multiple focused dashboards
# Dashboard 1: API Performance (8 charts)
# Dashboard 2: Database Health (6 charts)
# Dashboard 3: Business Metrics (7 charts)
```

### Hardcoded Values

**❌ Problem:**
```hcl
resource "uptrace_monitor" "latency" {
  params = {
    max_allowed_value = 2000000000  # What is this?
  }
}
```

**✅ Solution:**
```hcl
locals {
  latency_threshold_2s = 2 * 1000 * 1000 * 1000  # 2 seconds in nanoseconds
}

resource "uptrace_monitor" "latency" {
  params = {
    max_allowed_value = local.latency_threshold_2s
  }
}
```

## Production Checklist

Before deploying to production:

- [ ] Secrets stored securely (variables, vault, env vars)
- [ ] Remote state configured with locking
- [ ] Notification channels tested
- [ ] Monitor thresholds validated in staging
- [ ] Dashboard queries tested in Uptrace UI
- [ ] CI/CD pipeline configured
- [ ] Documentation updated
- [ ] Team trained on new monitors/dashboards
- [ ] Runbooks created for alert responses
- [ ] Backup/disaster recovery plan in place

## Resources

- [Getting Started Guide](getting-started.md)
- [Dashboard Examples](../../examples/dashboard-examples/)
- [Complete Setup Example](../../examples/complete-setup/)
- [Uptrace Documentation](https://uptrace.dev/docs/)
- [Terraform Best Practices](https://www.terraform.io/docs/language/syntax/style.html)
