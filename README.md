# Terraform Provider for Uptrace

A Terraform provider for managing [Uptrace](https://uptrace.dev/) resources.

## Features

- **Monitor Resource**: Create and manage metric and error monitors with full CRUD support
- **Dashboard Resource**: Manage YAML-based dashboards with grid layouts and visualizations
- **Monitor Data Source**: Read individual monitor configurations
- **Monitors Data Source**: Query and filter multiple monitors by type, state, or name
- **Type-safe Client**: Auto-generated client from OpenAPI specifications
- **Full Lifecycle Management**: Complete CRUD operations for all resources

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24
- [Uptrace](https://uptrace.dev/) instance with API access

## Building The Provider

```bash
go build -o terraform-provider-uptrace
```

Or using Task:

```bash
task build
```

## Using the Provider

```hcl
terraform {
  required_providers {
    uptrace = {
      source  = "riccap/uptrace"
      version = "~> 0.1"
    }
  }
}

provider "uptrace" {
  endpoint   = "https://uptrace.example.com"
  token      = var.uptrace_token
  project_id = 1
}

# Or using environment variables:
# export UPTRACE_ENDPOINT="https://uptrace.example.com"
# export UPTRACE_TOKEN="your-token"
# export UPTRACE_PROJECT_ID="1"

# Metric monitor
resource "uptrace_monitor" "high_cpu" {
  name = "High CPU Usage"
  type = "metric"

  notify_everyone_by_email = false
  channel_ids              = [1, 2]

  params = {
    metrics = [
      {
        name  = "system.cpu.utilization"
        alias = "cpu_usage"
      }
    ]
    query             = "avg(cpu_usage) > 90"
    max_allowed_value = 90
  }
}

# Dashboard
resource "uptrace_dashboard" "overview" {
  name = "System Overview"

  yaml = <<-YAML
    schema: v2
    grid_rows:
      - title: Metrics
        items:
          - title: CPU Usage
            metrics:
              - system.cpu.utilization as $cpu
            query:
              - avg($cpu)
  YAML
}
```

## Development

### Prerequisites

- Go 1.24+
- [Task](https://taskfile.dev/)
- [golangci-lint](https://golangci-lint.run/)
- [oapi-codegen](https://github.com/deepmap/oapi-codegen)
- [Docker](https://www.docker.com/) (for running acceptance tests)

### Setup

```bash
# Install dependencies
task deps

# Generate client code
task generate

# Run tests
task test

# Run acceptance tests (requires Uptrace instance)
task testacc
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

```bash
task test
```

### Acceptance Tests

Acceptance tests create real resources. Set up environment variables:

```bash
export UPTRACE_ENDPOINT="https://uptrace.example.com"
export UPTRACE_TOKEN="your-token"
export UPTRACE_PROJECT_ID="1"

task testacc
```

## Documentation

Documentation is auto-generated using [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs):

```bash
task docs
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.
