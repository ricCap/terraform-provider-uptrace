# Terraform Provider for Uptrace

A Terraform provider for managing [Uptrace](https://uptrace.dev/) monitoring resources - monitors, dashboards, and notification channels.

[![Tests](https://github.com/riccap/terraform-provider-uptrace/actions/workflows/test.yml/badge.svg)](https://github.com/riccap/terraform-provider-uptrace/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/riccap/terraform-provider-uptrace)](https://goreportcard.com/report/github.com/riccap/terraform-provider-uptrace)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Monitor Management**: Create metric and error monitors with configurable thresholds and alerts
- **Dashboard Creation**: YAML-based dashboards with flexible grid layouts and visualizations
- **Notification Channels**: Support for Slack, Telegram, Mattermost, and generic webhooks
- **Data Sources**: Query individual monitors or filter multiple monitors by criteria
- **Type-Safe Client**: Auto-generated client from OpenAPI specifications
- **Full CRUD Support**: Complete lifecycle management for all resources
- **Import Support**: Import existing Uptrace resources into Terraform state

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) or [OpenTofu](https://opentofu.org/) >= 1.0
- [Uptrace](https://uptrace.dev/) instance with API access
- Uptrace API token and project ID

## Building The Provider

```bash
go build -o terraform-provider-uptrace
```

Or using Task:

```bash
task build
```

## Quick Start

### Installation

```hcl
terraform {
  required_providers {
    uptrace = {
      source  = "registry.terraform.io/riccap/uptrace"
      version = "~> 0.1"
    }
  }
}

provider "uptrace" {
  endpoint   = "https://uptrace.example.com/api/v1"
  token      = var.uptrace_token
  project_id = var.uptrace_project_id
}
```

### Configuration

The provider can be configured via:

1. **HCL variables** (shown above)
2. **Environment variables**:
   ```bash
   export UPTRACE_ENDPOINT="https://uptrace.example.com/api/v1"
   export UPTRACE_TOKEN="your-api-token"
   export UPTRACE_PROJECT_ID="1"
   ```

### Example Usage

**Create a Notification Channel:**

```hcl
resource "uptrace_notification_channel" "slack_alerts" {
  name = "Engineering Alerts"
  type = "slack"

  params = {
    webhookUrl = var.slack_webhook_url
  }
}
```

**Create an Error Monitor:**

```hcl
resource "uptrace_monitor" "api_errors" {
  name = "High Error Rate"
  type = "error"

  channel_ids = [uptrace_notification_channel.slack_alerts.id]

  params = {
    metrics = [
      {
        name = "span.count"
      }
    ]
    query = "service.name:api-gateway span.status_code:error"

    max_allowed_value = 100
    check_num_point   = 1
  }
}
```

**Create a Performance Monitor:**

```hcl
resource "uptrace_monitor" "api_latency" {
  name = "API Response Time >2s"
  type = "metric"

  channel_ids = [uptrace_notification_channel.slack_alerts.id]

  params = {
    metrics = [
      {
        name  = "span.duration"
        alias = "$p95_latency"
      }
    ]

    query = "span.system:http span.kind:server"

    max_allowed_value = 2000000000  # 2 seconds in nanoseconds
    column            = "$p95_latency"
    check_num_point   = 2
  }
}
```

**Create a Dashboard:**

```hcl
resource "uptrace_dashboard" "api_overview" {
  yaml = <<-YAML
    schema: v2
    name: API Overview
    grid_rows:
      - title: Performance
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

## Documentation

- **[Getting Started Guide](docs/guides/getting-started.md)** - Complete introduction with examples
- **[Best Practices Guide](docs/guides/best-practices.md)** - Production recommendations
- **[Resource Documentation](docs/)** - Detailed resource and data source reference
- **[Example Configurations](examples/)** - Real-world usage examples

## Examples

### Complete Examples

- **[Dashboard Examples](examples/dashboard-examples/)** - 8 dashboard patterns (RED metrics, database performance, business metrics, etc.)
- **[Multi-Channel Alerts](examples/multi-channel-alerts/)** - Advanced notification routing patterns
- **[Complete Setup](examples/complete-setup/)** - Full monitoring stack with channels, monitors, and dashboards

### Resource Examples

- **[Monitor Resource](examples/resources/uptrace_monitor/)** - Metric and error monitor configurations
- **[Dashboard Resource](examples/resources/uptrace_dashboard/)** - Dashboard YAML examples
- **[Notification Channel Resource](examples/resources/uptrace_notification_channel/)** - Channel configurations for all supported types

## Development

### Prerequisites

- Go 1.21+
- [Task](https://taskfile.dev/) - Task runner for development commands
- [golangci-lint](https://golangci-lint.run/) - Go linter
- [oapi-codegen](https://github.com/deepmap/oapi-codegen) - OpenAPI code generator
- [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs) - Documentation generator
- [Docker](https://www.docker.com/) - For running local Uptrace instance and acceptance tests

### Setup

```bash
# Install development dependencies
task deps

# Generate API client from OpenAPI spec
task generate

# Build the provider
task build

# Run unit tests
task test:unit

# Run linters
task lint

# Start local Uptrace instance for testing
task dev:up

# Run acceptance tests (requires running Uptrace instance)
task test:acc

# Stop local Uptrace instance
task dev:down

# Generate documentation
task docs
```

### Available Tasks

Run `task --list` to see all available tasks:

```
* build:          Build the provider binary
* deps:           Install development dependencies
* dev:down:       Stop development environment
* dev:up:         Start development environment (Uptrace + dependencies)
* docs:           Generate provider documentation
* generate:       Generate API client code from OpenAPI spec
* lint:           Run golangci-lint
* test:acc:       Run acceptance tests
* test:unit:      Run unit tests
* test:           Run all tests
```

### Project Structure

```
.
├── api/                    # OpenAPI specifications
├── internal/
│   ├── client/            # Uptrace API client
│   ├── provider/          # Terraform provider implementation
│   └── validators/        # Custom validators
├── examples/              # Example configurations
├── docs/                  # Generated documentation
└── tools/                 # Development tools
```

## Testing

### Unit Tests

Run unit tests without requiring an Uptrace instance:

```bash
task test:unit
```

### Acceptance Tests

Acceptance tests create real resources in an Uptrace instance.

**Option 1: Use local Uptrace instance (recommended for development):**

```bash
# Start local Uptrace with Docker Compose
task dev:up

# Run acceptance tests
task test:acc

# Stop local instance when done
task dev:down
```

**Option 2: Use existing Uptrace instance:**

```bash
export UPTRACE_ENDPOINT="https://uptrace.example.com/api/v1"
export UPTRACE_TOKEN="your-api-token"
export UPTRACE_PROJECT_ID="1"

task test:acc
```

### Linting

```bash
task lint
```

## Contributing

Contributions are welcome! Here's how to get started:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/my-feature`
3. **Make your changes**
4. **Run tests**: `task test`
5. **Run linters**: `task lint`
6. **Commit your changes**: `git commit -m 'feat: add my feature'`
7. **Push to your fork**: `git push origin feature/my-feature`
8. **Open a Pull Request**

### Development Guidelines

- Follow existing code style and patterns
- Add tests for new features
- Update documentation as needed
- Keep commits focused and atomic
- Write clear commit messages

### Generating Documentation

Documentation is auto-generated from code and examples:

```bash
task docs
```

This generates resource documentation from:
- Schema definitions in `internal/provider/`
- Examples in `examples/resources/` and `examples/data-sources/`

## Support

- **[GitHub Issues](https://github.com/riccap/terraform-provider-uptrace/issues)** - Bug reports and feature requests
- **[GitHub Discussions](https://github.com/riccap/terraform-provider-uptrace/discussions)** - Questions and community support
- **[Uptrace Documentation](https://uptrace.dev/docs/)** - Uptrace platform documentation
- **[Getting Started Guide](docs/guides/getting-started.md)** - Provider usage guide

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release notes and version history.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- Built with [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)
- API client generated with [oapi-codegen](https://github.com/deepmap/oapi-codegen)
- Designed for [Uptrace](https://uptrace.dev/) observability platform
