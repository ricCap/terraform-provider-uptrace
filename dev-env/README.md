# Uptrace Local Development Environment

This directory contains a Docker Compose setup for running Uptrace locally to test the Terraform provider.

## Quick Start

1. **Start Uptrace** (from project root):
   ```bash
   task dev:up
   ```

   This will start all services and wait for Uptrace to be ready.

2. **Access the Uptrace UI (optional):**
   - URL: http://localhost:14318
   - Email: `admin@uptrace.local`
   - Password: `SomeRandomPassword`

3. **Configure OpenTofu test:**
   ```bash
   cp dev-env/terraform-test/terraform.tfvars.example dev-env/terraform-test/terraform.tfvars
   # The API token is already set to: user1_secret_token
   ```

4. **Test the provider:**
   ```bash
   task dev:install-local  # Build and install provider locally
   task dev:test           # Run tofu plan
   ```

5. **Apply the test configuration:**
   ```bash
   cd dev-env/terraform-test
   tofu apply
   ```

## Available Tasks

From the project root, you can run:

- `task dev:up` - Start Uptrace and all dependencies
- `task dev:down` - Stop all services
- `task dev:reset` - Stop and remove all data (fresh start)
- `task dev:logs` - View Uptrace logs
- `task dev:install-local` - Build and install provider locally
- `task dev:test` - Run terraform plan against local Uptrace

## Services

- **Uptrace UI**: http://localhost:14318
- **Uptrace gRPC**: localhost:14317
- **OpenTelemetry Collector gRPC**: localhost:4317
- **OpenTelemetry Collector HTTP**: localhost:4318
- **PostgreSQL**: localhost:5432 (user: uptrace, password: uptrace)
- **ClickHouse**: localhost:9000 (user: uptrace, password: uptrace)
- **ClickHouse HTTP**: localhost:8123

## Project Details

- **Project ID**: 1
- **Project Name**: Test Project
- **Project Token**: project1_secret_token (for OTLP ingestion)
- **API Endpoint**: http://localhost:14318/internal/v1
- **User API Token**: user1_secret_token (for Terraform provider)

## Stopping the Environment

```bash
task dev:down
```

To remove all data and start fresh:
```bash
task dev:reset
```

## Troubleshooting

**Services not starting:**
```bash
task dev:logs
# or
cd dev-env && docker compose logs
```

**Reset everything:**
```bash
task dev:reset
task dev:up
```

**Check service health:**
```bash
cd dev-env && docker compose ps
```

**View all available tasks:**
```bash
task --list
```
