# Multi-Channel Alert Routing Patterns

This example demonstrates advanced notification routing patterns for the Uptrace Terraform provider. Learn how to route alerts to the right people at the right time through the right channels.

## Patterns Demonstrated

### 1. **Severity-Based Routing** 🚨
Route alerts based on impact level:
- **SEV1 (Critical)**: PagerDuty + Slack → Page on-call immediately
- **SEV2 (High)**: Slack on-call → Notify but don't page
- **SEV3 (Medium)**: Slack engineering → Awareness only

**Use Case**: Prevent alert fatigue by matching response urgency to severity

### 2. **Team-Based Routing** 👥
Route alerts to responsible teams:
- **Backend Team**: API and service errors
- **Frontend Team**: Client-side and UI errors
- **Data Team**: Database and pipeline errors

**Use Case**: Each team sees only their relevant alerts

### 3. **Multi-Channel Fan-Out** 📢
Send critical alerts to multiple channels simultaneously:
- Slack + Telegram + PagerDuty + Custom webhook
- Ensures critical issues reach everyone
- Provides redundancy if one channel fails

**Use Case**: Revenue-impacting issues (payment failures, checkout errors)

### 4. **Environment-Specific Routing** 🌍
Different channels for different environments:
- **Production**: High sensitivity, immediate attention
- **Staging**: Lower sensitivity, more tolerance

**Use Case**: Avoid production alert fatigue from staging noise

### 5. **Time-Based Routing** ⏰
Different escalation during business hours vs after-hours:
- **Business Hours**: Slack notifications
- **After Hours**: PagerDuty for critical issues only

**Use Case**: Balance responsiveness with work-life balance

## Architecture Diagram

```
                           ┌─────────────┐
                           │   Monitor   │
                           │   Triggers  │
                           └──────┬──────┘
                                  │
                 ┌────────────────┼────────────────┐
                 │                │                │
          ┌──────▼──────┐  ┌──────▼──────┐  ┌────▼─────┐
          │  SEV1 Alert │  │  SEV2 Alert │  │SEV3 Alert│
          └──────┬──────┘  └──────┬──────┘  └────┬─────┘
                 │                │                │
         ┌───────┼───────┐        │                │
         │       │       │        │                │
    ┌────▼───┐ ┌▼────┐ ┌▼─────┐ ┌▼──────┐    ┌───▼────┐
    │PagerDuty│ │Slack│ │Telegram│ │ Slack │    │ Slack  │
    │  Page   │ │Alert│ │Mobile │ │Oncall │    │  Team  │
    └─────────┘ └─────┘ └───────┘ └───────┘    └────────┘
```

## Setup Guide

### 1. Create Slack Channels

Recommended Slack channel structure:

```
#ops-sev1-critical     → SEV1 alerts (@ channel)
#ops-sev2-oncall       → SEV2 alerts (@ here)
#ops-sev3-engineering  → SEV3 alerts (no @)
#team-backend          → Backend team
#team-frontend         → Frontend team
#team-data             → Data team
#ops-production        → Prod-only alerts
#ops-staging           → Staging-only alerts
```

### 2. Configure Webhooks

For each channel:

1. Go to Slack → Apps → Incoming Webhooks
2. Click "Add to Slack"
3. Select channel
4. Copy webhook URL
5. Add to `terraform.tfvars`

### 3. Set Up PagerDuty

```bash
# Create two PagerDuty integrations:
# 1. SEV1 (24/7 on-call)
# 2. After-hours escalation

# Format webhook URLs:
https://events.pagerduty.com/integration/{INTEGRATION_KEY}/enqueue
```

### 4. Apply Configuration

```bash
terraform init
terraform plan
terraform apply
```

## Pattern Details

### Pattern 1: Severity-Based Routing

```hcl
# SEV1: Complete outage - Page + Notify
resource "uptrace_monitor" "sev1_outage" {
  channel_ids = [
    uptrace_notification_channel.pagerduty.id,  # Page
    uptrace_notification_channel.slack_critical.id  # Also notify
  ]
  # ...
}

# SEV2: High errors - Notify only
resource "uptrace_monitor" "sev2_errors" {
  channel_ids = [
    uptrace_notification_channel.slack_oncall.id  # No page
  ]
  # ...
}
```

**Decision Matrix**:

| Severity | Criteria | Channels | Example |
|----------|----------|----------|---------|
| SEV1 | Service down, data loss | PagerDuty + Slack + Telegram | No traffic, payment failures |
| SEV2 | Degraded service | Slack on-call | High error rate, slow responses |
| SEV3 | Minor issues | Slack engineering | Elevated latency, warnings |

### Pattern 2: Team-Based Routing

```hcl
# Backend owns API errors
resource "uptrace_monitor" "api_errors" {
  channel_ids = [uptrace_notification_channel.team_backend.id]
  params = {
    query = "service.name:api-server span.status_code:error"
  }
}

# Frontend owns client errors
resource "uptrace_monitor" "client_errors" {
  channel_ids = [uptrace_notification_channel.team_frontend.id]
  params = {
    query = "service.name:web-app span.status_code:error"
  }
}
```

**Query Patterns by Team**:

```hcl
# Backend Team
query = "service.name:api-* span.status_code:error"
query = "span.system:http span.kind:server http.status_code:>=500"

# Frontend Team
query = "service.name:web-* span.status_code:error"
query = "span.system:browser error.type:*"

# Data Team
query = "service.name:*-pipeline span.status_code:error"
query = "db.system:* span.status_code:error"
```

### Pattern 3: Fan-Out

Use for business-critical alerts that must reach everyone:

```hcl
resource "uptrace_monitor" "payment_failures" {
  name = "CRITICAL: Payment Failures"

  # Send to ALL channels
  channel_ids = [
    uptrace_notification_channel.slack.id,
    uptrace_notification_channel.pagerduty.id,
    uptrace_notification_channel.telegram.id,
    uptrace_notification_channel.webhook.id,
  ]
}
```

**When to Fan-Out**:
- ✅ Payment processing failures
- ✅ Authentication system down
- ✅ Database connection lost
- ❌ Minor performance degradation
- ❌ Non-critical warnings

### Pattern 4: Environment-Specific

```hcl
# Production: Low tolerance
resource "uptrace_monitor" "prod_errors" {
  params = {
    query = "deployment.environment:production"
    max_allowed_value = 50
    check_num_point = 2  # Alert after 2 checks
  }
}

# Staging: Higher tolerance
resource "uptrace_monitor" "staging_errors" {
  params = {
    query = "deployment.environment:staging"
    max_allowed_value = 200
    check_num_point = 5  # Alert after 5 checks
  }
}
```

### Pattern 5: Time-Based

Implement using alert escalation:

```hcl
# Business hours (9-5): Slack
resource "uptrace_monitor" "degraded_perf" {
  channel_ids = [
    uptrace_notification_channel.slack_business.id
  ]
  params = {
    max_allowed_value = 3000000000  # 3 seconds (tolerant)
    check_num_point = 5
  }
}

# After hours: Only page for critical
resource "uptrace_monitor" "critical_outage" {
  channel_ids = [
    uptrace_notification_channel.pagerduty.id,
    uptrace_notification_channel.slack.id  # Also log
  ]
  params = {
    max_allowed_value = 100
    check_num_point = 1  # Immediate
  }
}
```

**Note**: Uptrace doesn't have built-in time-based routing. Implement in PagerDuty schedules or use separate monitors with different sensitivities.

## Best Practices

### 1. Channel Naming Convention

Use consistent naming:

```
[Severity/Team/Purpose] - [Platform]

Examples:
- "SEV1 - PagerDuty"
- "Backend Team - Slack"
- "Production Alerts - Slack"
```

### 2. Alert Grouping

Group related monitors with same channels:

```hcl
locals {
  critical_channels = [
    uptrace_notification_channel.pagerduty.id,
    uptrace_notification_channel.slack_critical.id,
  ]

  warning_channels = [
    uptrace_notification_channel.slack_warnings.id,
  ]
}

resource "uptrace_monitor" "critical_1" {
  channel_ids = local.critical_channels
}

resource "uptrace_monitor" "critical_2" {
  channel_ids = local.critical_channels
}
```

### 3. Avoid Alert Fatigue

```hcl
# ❌ BAD: Too sensitive
params = {
  max_allowed_value = 1  # Single error triggers
  check_num_point = 1   # Immediate alert
}

# ✅ GOOD: Balanced
params = {
  max_allowed_value = 50     # Tolerate some errors
  check_num_point = 2         # Wait for trend
}
```

### 4. Document Runbooks

Add monitor descriptions:

```hcl
resource "uptrace_monitor" "payment_failures" {
  name = "CRITICAL: Payment Failures"

  # Document response process
  # Note: Description field not yet supported in provider
  # Add to Uptrace UI manually or via API
}
```

### 5. Test Channels

Before deploying:

```bash
# Test Slack
curl -X POST -H 'Content-type: application/json' \
  --data '{"text":"Test alert"}' \
  YOUR_WEBHOOK_URL

# Test PagerDuty
curl -X POST https://events.pagerduty.com/v2/enqueue \
  -H 'Content-Type: application/json' \
  -d '{
    "routing_key": "YOUR_KEY",
    "event_action": "trigger",
    "payload": {
      "summary": "Test alert",
      "severity": "critical",
      "source": "terraform-test"
    }
  }'
```

## Common Patterns

### Weekly Digest (Not Yet Supported)

Future pattern for low-priority alerts:

```hcl
# Weekly summary of warnings
resource "uptrace_notification_channel" "weekly_digest" {
  name = "Weekly Digest"
  type = "webhook"

  params = {
    url = var.digest_webhook
  }

  # Future: schedule field
  # schedule = "0 9 * * MON"  # Every Monday at 9am
}
```

### Conditional Routing

Use query filters for smart routing:

```hcl
# Route database errors to DBA team
resource "uptrace_monitor" "database_team" {
  channel_ids = [uptrace_notification_channel.dba_team.id]
  params = {
    query = "db.system:* span.status_code:error"
  }
}

# Route cache errors to infrastructure team
resource "uptrace_monitor" "infra_team" {
  channel_ids = [uptrace_notification_channel.infra_team.id]
  params = {
    query = "cache.system:* span.status_code:error"
  }
}
```

## Troubleshooting

### Alerts Not Routing Correctly

1. **Verify channel IDs**:
   ```bash
   terraform state show uptrace_notification_channel.slack_critical
   ```

2. **Check monitor configuration**:
   ```bash
   terraform state show uptrace_monitor.sev1_outage
   ```

3. **Test webhook**:
   ```bash
   curl -X POST -H 'Content-type: application/json' \
     --data '{"text":"Test"}' YOUR_WEBHOOK
   ```

### Too Many Alerts

1. **Increase thresholds**:
   ```hcl
   max_allowed_value = 100  # Was 50
   ```

2. **Add more check points**:
   ```hcl
   check_num_point = 3  # Was 1
   ```

3. **Adjust queries**:
   ```hcl
   # More specific
   query = "service.name:critical-service span.status_code:error"
   # Instead of
   query = "span.status_code:error"
   ```

### Missing Alerts

1. **Check monitor is triggering**:
   - View in Uptrace UI
   - Verify query returns data

2. **Verify channel status**:
   ```bash
   terraform state show uptrace_notification_channel.slack
   # Check status field
   ```

3. **Test with lower threshold**:
   ```hcl
   max_allowed_value = 1  # Temporarily very sensitive
   ```

## Cost Optimization

- **PagerDuty**: Limited free tier, pay per incident
  - Use only for SEV1
  - Consider PagerDuty schedules to reduce unnecessary pages

- **Slack**: Free for basic webhooks
  - Unlimited channels and messages
  - Use liberally

- **Telegram**: Free
  - Good for mobile notifications
  - No rate limits for bots

## Related Examples

- [Complete Setup](../complete-setup/) - Full monitoring stack
- [Dashboard Examples](../dashboard-examples/) - Visualization configs
- [Getting Started Guide](../../docs/guides/getting-started.md) - Basic setup

## Next Steps

1. **Customize for Your Org**
   - Map teams to Slack channels
   - Define severity levels
   - Set escalation policies

2. **Add Runbooks**
   - Document response procedures
   - Link from monitor names/descriptions
   - Train team on escalation

3. **Refine Over Time**
   - Monitor alert frequency
   - Adjust thresholds based on feedback
   - Reduce noise, increase signal

## Support

- [Uptrace Documentation](https://uptrace.dev/docs/)
- [Provider Issues](https://github.com/riccap/tofu-uptrace-provider/issues)
- [Slack Best Practices](https://api.slack.com/messaging/webhooks)
- [PagerDuty Integration](https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTgw-events-api-v2-overview)
