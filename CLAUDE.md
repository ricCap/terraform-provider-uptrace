# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a Terraform/OpenTofu provider for Uptrace (observability platform). It uses the **Terraform Plugin Framework (ProtoV6)** and auto-generates the API client from an OpenAPI specification using oapi-codegen.

## Core Architecture

### Three-Layer Design

1. **Generated Client Layer** (`internal/client/generated/`)
   - Auto-generated from `api/openapi.yaml` using oapi-codegen
   - Never edit directly - regenerate with `task generate`
   - Contains type-safe HTTP client with models

2. **Client Wrapper Layer** (`internal/client/client.go`)
   - Wraps generated client with higher-level operations
   - Handles error responses consistently via `handleErrorResponse()`
   - All methods accept `context.Context` as first parameter
   - ProjectID stored in client struct, passed to all API calls

3. **Provider Layer** (`internal/provider/`)
   - Terraform Plugin Framework resources and data sources
   - Model conversion between Terraform types and API types happens here
   - Pattern: `*_resource.go` (CRUD), `*_models.go` (conversion), `*_data_source.go` (read-only)

### Key Patterns

**Model Conversion Flow:**
```
Terraform Plan → *Model structs → API Input structs → HTTP Request
HTTP Response → API Output structs → *Model structs → Terraform State
```

**Resource Structure:**
- Each resource implements: `Resource`, `ResourceWithConfigure`, `ResourceWithImportState`
- CRUD methods: `Create()`, `Read()`, `Update()`, `Delete()`, `ImportState()`
- Models use Terraform types (`types.String`, `types.Int64`, etc.) NOT Go primitives
- API models use Go primitives and pointers for optional fields

**Dashboard Resource Exception:**
- Dashboard uses YAML-based API instead of JSON
- `CreateDashboardFromYAML()` and `UpdateDashboardFromYAML()` send `text/yaml` content type
- YAML field is required input, all other fields computed from API response

## Development Commands

### Code Generation
```bash
task generate              # Generate API client from OpenAPI spec
task docs                  # Generate Terraform provider documentation
```

**CRITICAL:** Always run `task generate` after modifying `api/openapi.yaml`. The generated client is the source of truth for API types.

### Testing
```bash
# Unit tests only (no external dependencies)
task test:unit
go test -short ./...

# Run specific test
go test -v -run TestPinDashboard ./internal/client/

# Acceptance tests (requires Uptrace instance)
task dev:up                # Start local Uptrace with Docker
task test:acc              # Run all acceptance tests
TF_ACC=1 go test -v -run TestAccDashboardResource_Basic ./internal/provider/

# Coverage
task test:coverage:unit    # Unit test coverage → coverage.html
task test:coverage:acc     # Acceptance test coverage → coverage-acc.html
task test:coverage:all     # Both coverage reports
```

**Acceptance Test Environment:**
- Endpoint: `http://localhost:14318/internal/v1`
- Token: `user1_secret_token`
- Project ID: `1`
- Set `TF_ACC=1` to enable acceptance tests

### Building & Linting
```bash
task build                 # Build provider binary
task lint                  # Run golangci-lint (strict config)
task lint:fix              # Auto-fix linting issues
```

## OpenAPI Spec Modifications

When adding/modifying API endpoints in `api/openapi.yaml`:

1. **Update OpenAPI spec** with new schemas/endpoints
2. **Run `task generate`** to regenerate client
3. **Add wrapper methods** in `internal/client/client.go`:
   ```go
   func (c *Client) GetFoo(ctx context.Context, id int64) (*generated.Foo, error) {
       resp, err := c.client.GetFooWithResponse(ctx, c.projectID, id)
       if err != nil {
           return nil, fmt.Errorf("failed to get foo: %w", err)
       }
       if !isSuccessStatus(resp.StatusCode(), http.StatusOK) {
           return nil, c.handleErrorResponse(resp.StatusCode(), resp.Body)
       }
       if resp.JSON200 == nil {
           return nil, fmt.Errorf("unexpected empty response")
       }
       return &resp.JSON200.Foo, nil
   }
   ```
4. **Add unit tests** in `internal/client/client_test.go` using `httptest.NewServer()`
5. **Create/update provider resources** in `internal/provider/`

## Testing Requirements

**Every new feature MUST include:**
1. **Unit tests** - Test client methods with mock HTTP servers
2. **Acceptance tests** - Test resources against real Uptrace instance
3. Both unit and acceptance test coverage tracked separately in CI

**Test file naming:**
- `*_test.go` - Unit tests (run with `-short` flag)
- `*_acc_test.go` - Acceptance tests (require `TF_ACC=1`)

**Common test helpers:**
- `newTestClient(server)` - Create client for unit tests
- `acceptancetests.GetTestClient()` - Get client for acceptance tests
- `acceptancetests.RandomTestName(prefix)` - Generate unique resource names

## Linting Notes

The golangci-lint config is strict. Common issues:

1. **godot**: All comments must end with a period
2. **revive unused-parameter**: Rename unused params to `_`
3. **gocritic paramTypeCombine**: Combine consecutive same-type params: `func(a, b string)`
4. **dupl**: Use `//nolint:dupl` with justification for acceptable duplicates

Generated code in `internal/client/generated/` is excluded from linting.

## Git Workflow

**Commit message format (Conventional Commits):**
```
feat: add dashboard clone functionality
fix: handle 404 errors in monitor resource
test: add unit tests for pin/unpin operations
docs: update OpenAPI spec with new endpoints
```

**All commits must include:**
```
Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

## CI/CD

GitHub Actions runs:
- **Build** - Verify compilation
- **Unit Tests** - With coverage upload to Codecov (flag: `unittests`)
- **Acceptance Tests** - Against local Uptrace (flag: `acceptancetests`)
- **Lint** - golangci-lint with strict config
- **Cloud API Tests** - If secrets configured

All checks must pass before merge.

## Common Pitfalls

1. **Don't edit generated code** - Always modify `api/openapi.yaml` and regenerate
2. **Use Terraform types in models** - `types.String` not `string`, even for required fields
3. **Pointers in API structs** - Generated API types use pointers for optional fields
4. **Context first parameter** - All methods accepting context must have it as first param
5. **Dashboard YAML quirk** - API enriches YAML with defaults, causing drift. Use `ignore_changes = ["yaml"]` or `ImportStateVerifyIgnore` in tests
6. **Test isolation** - Acceptance tests use random names via `RandomTestName()` to avoid conflicts

## Project Structure

```
api/                          # OpenAPI specifications (source of truth)
  openapi.yaml               # Main API spec
internal/
  client/
    generated/               # Auto-generated client (never edit)
    client.go                # Client wrapper with high-level methods
    client_test.go           # Unit tests for client methods
  provider/
    provider.go              # Provider configuration
    *_resource.go            # Resource implementations (CRUD)
    *_models.go              # Model conversion logic
    *_data_source.go         # Data source implementations
    *_acc_test.go            # Acceptance tests
  acceptance_tests/          # Shared test utilities
    acctest.go               # Test configuration helpers
dev-env/                     # Docker Compose for local Uptrace
examples/                    # Terraform examples for documentation
docs/                        # Generated provider documentation
```
