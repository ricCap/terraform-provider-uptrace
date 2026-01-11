# Terraform Provider for Uptrace

A Terraform provider for managing [Uptrace](https://uptrace.dev/) monitoring resources - monitors, dashboards, and notification channels.

[![Terraform Registry](https://img.shields.io/badge/Terraform-Registry-purple?logo=terraform)](https://registry.terraform.io/providers/riccap/uptrace)
[![OpenTofu Registry](https://img.shields.io/badge/OpenTofu-Registry-blue?logo=opentofu)](https://search.opentofu.org/provider/riccap/uptrace)
[![Tests](https://github.com/riccap/terraform-provider-uptrace/actions/workflows/test.yml/badge.svg)](https://github.com/riccap/terraform-provider-uptrace/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/riccap/terraform-provider-uptrace/branch/main/graph/badge.svg)](https://codecov.io/gh/riccap/terraform-provider-uptrace)
[![Go Report Card](https://goreportcard.com/badge/github.com/riccap/terraform-provider-uptrace)](https://goreportcard.com/report/github.com/riccap/terraform-provider-uptrace)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> This provider is in early development. Any feedback is welcome (just create an issue)!

## Features

- **Monitor Management**: Create metric and error monitors with configurable thresholds and alerts
- **Dashboard Creation**: YAML-based dashboards with flexible grid layouts and visualizations
- **Notification Channels**: Support for Slack, Telegram, Mattermost, and generic webhooks
- **Data Sources**: Query individual monitors or filter multiple monitors by criteria
- **Full CRUD Support**: Complete lifecycle management for all resources
- **Import Support**: Import existing Uptrace resources into Terraform state

## Cloud API Support

The provider supports both **Uptrace Cloud** (`api2.uptrace.dev`) and **Self-Hosted Uptrace**. See the [Cloud API Guide](docs/guides/cloud-api.md) for details on cloud-specific fields.

### Feature Compatibility

| Feature | Self-Hosted | Uptrace Cloud | Notes |
|---------|-------------|---------------|-------|
| **Monitor Resource** | ✅ Full Support | ⚠️ Partial | Cloud requires `trend_agg_func`, non-empty `query`, and existing metrics. Query normalization may cause state drift ([#53](https://github.com/ricCap/terraform-provider-uptrace/issues/53)). |
| **Dashboard Resource** | ✅ Full Support | ✅ Full Support | YAML-based dashboards work identically across both platforms. |
| **Notification Channels** | ✅ Full Support | ❌ Not Functional | Cloud API uses `priorities` (plural) field name ([#54](https://github.com/ricCap/terraform-provider-uptrace/issues/54)). |
| **Data Sources** | ✅ Full Support | ✅ Full Support | Query individual monitors or filter by criteria. |
| **Monitor `trend_agg_func`** | ⚠️ Optional | ✅ Required | Required for cloud API ([#55](https://github.com/ricCap/terraform-provider-uptrace/issues/55)). Valid values: `avg`, `sum`, `min`, `max`, `p50`, `p90`, `p95`, `p99`. |
| **Import Support** | ✅ Full Support | ✅ Full Support | Import existing resources into Terraform state. |

**Legend:**
- ✅ **Full Support** - Feature works as expected
- ⚠️ **Partial/Optional** - Feature available with limitations
- ❌ **Not Functional** - Feature not yet working (API limitation)

## Installation

```hcl
terraform {
  required_providers {
    uptrace = {
      source  = "riccap/uptrace"
      version = "~> 0.3"
    }
  }
}

provider "uptrace" {
  endpoint   = "https://uptrace.example.com/api/v1"
  token      = var.uptrace_token
  project_id = var.uptrace_project_id
}
```

The provider can also be configured via environment variables:

```bash
export UPTRACE_ENDPOINT="https://uptrace.example.com/api/v1"
export UPTRACE_TOKEN="your-api-token"
export UPTRACE_PROJECT_ID="1"
```

## Example Usage

```hcl
# Create a notification channel
resource "uptrace_notification_channel" "slack" {
  name = "Engineering Alerts"
  type = "slack"
  params = { webhookUrl = var.slack_webhook_url }
}

# Create an error monitor
resource "uptrace_monitor" "api_errors" {
  name        = "High Error Rate"
  type        = "error"
  channel_ids = [uptrace_notification_channel.slack.id]

  params = {
    metrics           = [{ name = "span.count" }]
    query             = "service.name:api-gateway span.status_code:error"
    max_allowed_value = 100
    check_num_point   = 1
  }
}
```

## Documentation

- **[Terraform Registry Docs](https://registry.terraform.io/providers/riccap/uptrace/latest/docs)** - Full provider documentation
- **[OpenTofu Registry Docs](https://search.opentofu.org/provider/riccap/uptrace)** - OpenTofu documentation
- **[Getting Started Guide](docs/guides/getting-started.md)** - Complete introduction
- **[Cloud API Guide](docs/guides/cloud-api.md)** - Uptrace Cloud configuration and known issues
- **[Example Configurations](examples/)** - Real-world usage examples

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## Support

- **[GitHub Issues](https://github.com/riccap/terraform-provider-uptrace/issues)** - Bug reports and feature requests
- **[Uptrace Documentation](https://uptrace.dev/docs/)** - Uptrace platform documentation

## License

MIT License - see [LICENSE](LICENSE) for details.
