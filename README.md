# Terraform Provider for Uptrace

A Terraform provider for managing [Uptrace](https://uptrace.dev/) resources.

## Features

- **Monitor Management**: Create and manage metric and error monitors
- **Type-safe Client**: Auto-generated client from OpenAPI specifications
- **Full CRUD Support**: Complete lifecycle management for Uptrace resources

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25
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
  endpoint = "https://uptrace.example.com"
  token    = var.uptrace_token
}

resource "uptrace_monitor" "example" {
  name        = "High Error Rate"
  description = "Alert when error rate exceeds threshold"
  type        = "metric"

  project_id = 1

  params = {
    metrics = ["errors"]
    query   = "rate > 10"
  }

  channels = [1, 2]
}
```

## Development

### Prerequisites

- Go 1.25+
- [Task](https://taskfile.dev/)
- [golangci-lint](https://golangci-lint.run/)
- [oapi-codegen](https://github.com/deepmap/oapi-codegen)

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
