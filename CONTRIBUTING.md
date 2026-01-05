# Contributing to Uptrace Terraform Provider

Thank you for your interest in contributing! This document outlines the development workflow and setup requirements.

## Development Setup

### Prerequisites

- Go 1.23 or later
- [Task](https://taskfile.dev) - Build automation tool
- [golangci-lint](https://golangci-lint.run) - Linting tool
- Docker and Docker Compose - For running Uptrace locally
- [gh](https://cli.github.com) - GitHub CLI (optional, for project management)

### Initial Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/ricCap/terraform-provider-uptrace.git
   cd terraform-provider-uptrace
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Generate API client code:
   ```bash
   task generate
   ```

## Development Workflow

### Running Tests

```bash
# Run all tests
task test

# Run only unit tests
task test:unit

# Run only acceptance tests (requires running Uptrace instance)
task test:acc

# Start local Uptrace for testing
task dev:up

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

### Building

```bash
# Build the provider
task build

# Install locally for testing
task install
```

### Code Generation

```bash
# Generate API client from OpenAPI spec
task generate:client

# Generate provider documentation
task docs
```

## Pull Request Process

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** and commit frequently:
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

3. **Run tests and linting locally**:
   ```bash
   task test
   task lint
   ```

4. **Push your branch**:
   ```bash
   git push origin feature/your-feature-name
   ```

5. **Create a Pull Request** on GitHub

6. **Wait for CI to pass** - All tests and linting must pass before merge

7. **Request review** - Tag the maintainer for review

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

## GitHub Workflows Setup

### Documentation Workflow

The documentation workflow (`docs.yml`) auto-generates provider documentation.

**Required Setup**:
1. Go to **Settings → Actions → General**
2. Scroll to **Workflow permissions**
3. ✓ Check **"Allow GitHub Actions to create and approve pull requests"**

Without this, the workflow will fail with:
```
GitHub Actions is not permitted to create or approve pull requests
```

### Auto-assign to Project Workflow

The auto-assign workflow adds new issues/PRs to the GitHub Project board.

**Note**: This workflow currently doesn't work for user-level projects because `GITHUB_TOKEN` lacks project permissions.

**Options**:
1. **Manual assignment**: Use `gh project item-add` to manually add items
2. **Use PAT** (if you want automation):
   - Create a Personal Access Token with `project` and `repo` scopes
   - Add it as repository secret `PROJECT_TOKEN`
   - Update workflow to use `${{ secrets.PROJECT_TOKEN }}`

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

Default Uptrace endpoint: `http://localhost:14318`

## Code Style

- Follow standard Go conventions
- Run `gofmt` and `goimports` (handled by linter)
- Document exported functions and types
- Add tests for new functionality
- Keep cyclomatic complexity under 15

## Questions or Issues?

- Open an issue for bugs or feature requests
- Tag @ricCap for questions
- Check existing issues before creating new ones

Thank you for contributing!
