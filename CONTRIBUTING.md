# Contributing to Uptrace Terraform Provider

Thank you for your interest in contributing! This document outlines the development workflow and setup requirements.

## Development Setup

### Prerequisites

- Go 1.23 or later
- [Task](https://taskfile.dev) - Build automation tool
- [golangci-lint](https://golangci-lint.run) - Linting tool
- [oapi-codegen](https://github.com/deepmap/oapi-codegen) - OpenAPI code generator
- [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs) - Documentation generator
- Docker and Docker Compose - For running Uptrace locally
- [gh](https://cli.github.com) - GitHub CLI (optional)

### Initial Setup

```bash
# Clone the repository
git clone https://github.com/ricCap/terraform-provider-uptrace.git
cd terraform-provider-uptrace

# Install development dependencies
task deps

# Generate API client code
task generate

# Build the provider
task build
```

## Project Structure

```
.
├── api/                    # OpenAPI specifications (source of truth)
│   └── openapi.yaml       # Main API spec
├── internal/
│   ├── client/            # Uptrace API client
│   │   ├── generated/     # Auto-generated client (never edit)
│   │   ├── client.go      # Client wrapper with high-level methods
│   │   └── client_test.go # Unit tests for client
│   ├── provider/          # Terraform provider implementation
│   │   ├── provider.go    # Provider configuration
│   │   ├── *_resource.go  # Resource implementations (CRUD)
│   │   ├── *_models.go    # Model conversion logic
│   │   ├── *_data_source.go # Data source implementations
│   │   └── *_acc_test.go  # Acceptance tests
│   ├── acceptance_tests/  # Shared test utilities
│   └── validators/        # Custom validators
├── examples/              # Terraform examples for documentation
├── docs/                  # Generated provider documentation
├── dev-env/               # Docker Compose for local Uptrace
└── tools/                 # Development tools
```

## Available Tasks

Run `task --list` to see all available tasks:

```
* build:              Build the provider binary
* deps:               Install development dependencies
* generate:           Generate API client code from OpenAPI spec
* docs:               Generate provider documentation
* lint:               Run golangci-lint
* lint:fix:           Run linters and auto-fix issues
* test:unit:          Run unit tests
* test:acc:           Run acceptance tests
* test:coverage:unit: Run unit tests with coverage
* test:coverage:acc:  Run acceptance tests with coverage
* dev:up:             Start development environment (Uptrace + dependencies)
* dev:down:           Stop development environment
* dev:logs:           Show Uptrace logs
```

## Development Workflow

### Running Tests

```bash
# Run unit tests (no external dependencies)
task test:unit

# Start local Uptrace for acceptance tests
task dev:up

# Run acceptance tests
task test:acc

# Run tests with coverage
task test:coverage:unit
task test:coverage:acc

# Stop local Uptrace
task dev:down
```

### Linting

```bash
# Run linter
task lint

# Auto-fix linting issues
task lint:fix
```

### Code Generation

```bash
# Generate API client from OpenAPI spec
task generate

# Generate provider documentation
task docs
```

**Important**: Always run `task generate` after modifying `api/openapi.yaml`. The generated client is the source of truth for API types.

## Testing Against Local Uptrace

The acceptance tests require a running Uptrace instance:

```bash
# Start Uptrace with Docker Compose
task dev:up

# Run acceptance tests
task test:acc

# View Uptrace logs
task dev:logs

# Stop Uptrace
task dev:down
```

**Local Uptrace credentials:**
- Endpoint: `http://localhost:14318/internal/v1`
- Token: `user1_secret_token`
- Project ID: `1`

## Pull Request Process

1. **Create a feature branch**:
   ```bash
   git checkout -b feat/your-feature-name
   ```

2. **Make your changes** and commit frequently:
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

3. **Run tests and linting locally**:
   ```bash
   task test:unit
   task lint
   ```

4. **Push your branch**:
   ```bash
   git push -u origin feat/your-feature-name
   ```

5. **Create a Pull Request** on GitHub:
   ```bash
   gh pr create --title "feat: your feature" --body "Description"
   ```

6. **Wait for CI to pass** - All tests and linting must pass before merge

## Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Test additions or fixes
- `refactor:` - Code refactoring
- `chore:` - Build process or auxiliary tool changes
- `ci:` - CI/CD changes

Examples:
```
feat: add dashboard resource
fix: handle 404 errors in monitor resource
docs: update installation instructions
test: add unit tests for monitor models
```

## Code Style

- Follow standard Go conventions
- Run `gofmt` and `goimports` (handled by linter)
- All comments must end with a period (godot linter)
- Document exported functions and types
- Add tests for new functionality
- Keep cyclomatic complexity under 15

### Linting Notes

The golangci-lint config is strict. Common issues:

1. **godot**: All comments must end with a period
2. **revive unused-parameter**: Rename unused params to `_`
3. **gocritic paramTypeCombine**: Combine consecutive same-type params: `func(a, b string)`
4. **dupl**: Use `//nolint:dupl` with justification for acceptable duplicates

## Adding New Features

### Adding API Endpoints

1. Update OpenAPI spec in `api/openapi.yaml`
2. Run `task generate` to regenerate client
3. Add wrapper methods in `internal/client/client.go`
4. Add unit tests in `internal/client/client_test.go`
5. Create/update provider resources in `internal/provider/`

### Adding Provider Resources

1. Create `*_resource.go` with CRUD operations
2. Create `*_models.go` for model conversion
3. Add acceptance tests in `*_acc_test.go`
4. Add examples in `examples/resources/`
5. Run `task docs` to generate documentation

## Questions or Issues?

- Open an issue for bugs or feature requests
- Tag @ricCap for questions
- Check existing issues before creating new ones

Thank you for contributing!
