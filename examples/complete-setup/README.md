# Complete Monitoring Setup Example

This example demonstrates a production-ready monitoring setup using the Uptrace Terraform provider. It includes:

- **4 Notification Channels** (Slack, PagerDuty, Telegram)
- **7 Monitors** covering errors, performance, and throughput
- **2 Dashboards** for visualization
- **Best practices** for alert routing and severity levels

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Uptrace Monitoring                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Critical   │  │   Warning    │  │  Dashboard   │     │
│  │   Monitors   │  │   Monitors   │  │   Metrics    │     │
│  └──────┬───────┘  └──────┬───────┘  └──────────────┘     │
│         │                 │                                 │
│         ▼                 ▼                                 │
│  ┌──────────────┐  ┌──────────────┐                       │
│  │  Slack #ops  │  │ Slack #warn  │                       │
│  ├──────────────┤  └──────────────┘                       │
│  │  PagerDuty   │                                          │
│  ├──────────────┤                                          │
│  │  Telegram    │                                          │
│  └──────────────┘                                          │
└─────────────────────────────────────────────────────────────┘
```

## What Gets Monitored

### Error Monitoring
1. **Critical: High Error Rate**
   - Triggers when errors > 100 in 5 minutes
   - Alerts: Slack Critical, PagerDuty, Telegram
   - Use case: Catch widespread failures immediately

2. **Warning: 4xx Errors**
   - Triggers when 4xx errors > 500 in 5 minutes
   - Alerts: Slack Warnings only
   - Use case: Track client-side issues without paging

3. **Critical: Database Errors**
   - Triggers when DB errors > 10 in 5 minutes
   - Alerts: Slack Critical, PagerDuty
   - Use case: Database health is critical

### Performance Monitoring
4. **Critical: API Latency**
   - Triggers when P95 latency > 2 seconds
   - Alerts: Slack Critical, Telegram
   - Use case: User experience degradation

5. **Warning: Slow Database Queries**
   - Triggers when DB P95 latency > 500ms
   - Alerts: Slack Warnings
   - Use case: Identify optimization opportunities

### Throughput Monitoring
6. **Warning: Traffic Spike**
   - Triggers when requests > 10,000/minute
   - Alerts: Slack Warnings
   - Use case: Capacity planning, detect attacks

7. **Critical: No Traffic**
   - Triggers when requests < 10/minute
   - Alerts: Slack Critical, PagerDuty
   - Use case: Detect complete outages

## Prerequisites

1. **Uptrace Instance**
   - Running Uptrace server
   - API token with project access
   - Project ID

2. **Notification Integrations**
   - Slack webhook URLs (2 channels recommended)
   - PagerDuty integration key
   - Telegram bot token and chat ID

3. **Terraform/OpenTofu**
   - Terraform >= 1.0 or OpenTofu >= 1.6

## Setup Instructions

### 1. Configure Credentials

```bash
# Copy the example file
cp terraform.tfvars.example terraform.tfvars

# Edit with your actual values
vim terraform.tfvars
```

### 2. Get Slack Webhooks

1. Go to https://api.slack.com/messaging/webhooks
2. Create webhook for #critical-alerts channel
3. Create webhook for #warnings channel
4. Add URLs to `terraform.tfvars`

### 3. Get PagerDuty Integration

1. In PagerDuty: Services → Your Service → Integrations
2. Add Integration → Events API V2
3. Copy Integration Key
4. Format as: `https://events.pagerduty.com/integration/{KEY}/enqueue`

### 4. Get Telegram Bot Credentials

```bash
# 1. Create bot
# Message @BotFather on Telegram: /newbot

# 2. Get chat ID
# Message @userinfobot: /start
# Add bot to group and get group chat ID
```

### 5. Apply Configuration

```bash
# Initialize Terraform
terraform init

# Review changes
terraform plan

# Apply configuration
terraform apply
```

## Customization

### Adjust Alert Thresholds

Edit `main.tf` and modify the `params` block:

```hcl
resource "uptrace_monitor" "critical_error_rate" {
  # ...
  params = {
    max_allowed_value = 50  # Lower threshold = more sensitive
    check_num_point   = 2   # More points = less noise
  }
}
```

### Add More Channels

```hcl
resource "uptrace_notification_channel" "mattermost" {
  name = "Mattermost Alerts"
  type = "mattermost"

  params = {
    webhookUrl = var.mattermost_webhook
  }
}
```

### Modify Monitor Queries

Use Uptrace query language:

```hcl
params = {
  # Monitor specific service
  query = "service.name:api-gateway span.status_code:error"

  # Monitor specific environment
  query = "deployment.environment:production span.status_code:error"

  # Complex conditions
  query = "span.system:http span.kind:server http.status_code:>=500"
}
```

## Dashboard YAML Structure

Dashboards use YAML configuration. Structure:

```yaml
table:
  - schema_version: v2
    row:
      - title: Chart Title
        columns:
          - span.system: [http]    # Filter criteria
            span.kind: [server]
        chart:
          - type: line             # line, bar, area
            metric: span.duration|p95
            legend: Label
            color: blue            # Optional
```

## Troubleshooting

### Monitors Not Triggering

1. **Check Query Results**
   - Go to Uptrace UI → Metrics
   - Run the same query as your monitor
   - Verify data exists

2. **Verify Thresholds**
   - Check `min_allowed_value` and `max_allowed_value`
   - Ensure they match your scale (nanoseconds for duration)

3. **Check Time Windows**
   - `check_num_point` determines how many consecutive points must breach
   - Increase for less noise, decrease for faster alerting

### Channels Not Receiving Alerts

1. **Test Webhooks**
   ```bash
   # Test Slack webhook
   curl -X POST -H 'Content-type: application/json' \
     --data '{"text":"Test from Terraform"}' \
     YOUR_SLACK_WEBHOOK_URL
   ```

2. **Verify Channel Assignment**
   - Check monitor's `channel_ids` array
   - Ensure channel IDs are correct

3. **Check Channel Status**
   ```bash
   terraform state show uptrace_notification_channel.slack_critical
   ```

### Dashboard Not Showing Data

1. **Verify YAML Syntax**
   - YAML is indentation-sensitive
   - Use 2-space indentation
   - Ensure proper nesting

2. **Check Filter Criteria**
   - Dashboard `columns` filters must match actual span attributes
   - Verify attribute names in Uptrace UI

## Cost Considerations

This setup creates:
- 4 notification channels (free)
- 7 monitors (~7 checks/minute depending on configuration)
- 2 dashboards (free)

**Uptrace Costs**: Based on your Uptrace pricing tier
**External Services**: Slack/PagerDuty/Telegram have their own limits

## Security Best Practices

1. **Never commit `terraform.tfvars`**
   ```bash
   echo "terraform.tfvars" >> .gitignore
   ```

2. **Use Terraform Backends**
   ```hcl
   terraform {
     backend "s3" {
       # Your backend config
     }
   }
   ```

3. **Rotate Tokens Regularly**
   - Update `uptrace_token` every 90 days
   - Rotate webhook URLs if exposed

4. **Use Workspace Separation**
   ```bash
   terraform workspace new production
   terraform workspace new staging
   ```

## Next Steps

1. **Fine-tune Thresholds**
   - Monitor alert frequency for 1 week
   - Adjust thresholds to reduce noise

2. **Add More Monitors**
   - Memory usage
   - CPU utilization
   - Custom business metrics

3. **Create On-Call Schedules**
   - Configure PagerDuty schedules
   - Set up escalation policies

4. **Document Runbooks**
   - Create runbooks for each monitor
   - Link in monitor descriptions

## Related Examples

- [Multi-Channel Alerts](../multi-channel-alerts/) - Complex notification patterns
- [Dashboard Examples](../dashboard-examples/) - More dashboard configurations
- [Best Practices Guide](../../docs/guides/best-practices.md) - Production recommendations

## Support

- [Uptrace Documentation](https://uptrace.dev/docs/)
- [Provider Documentation](../../docs/)
- [GitHub Issues](https://github.com/riccap/tofu-uptrace-provider/issues)
