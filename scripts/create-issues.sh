#!/bin/bash
set -e

# GitHub Issue Creation Script for tofu-uptrace-provider
# This script creates labels, milestones, and all planned issues for the provider roadmap
# Requires: GitHub CLI (gh) - Install with: brew install gh
# Usage: ./scripts/create-issues.sh

REPO="riccap/tofu-uptrace-provider"

echo "🚀 Setting up GitHub issues for $REPO"
echo ""

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo "❌ GitHub CLI (gh) is not installed"
    echo "Install it with: brew install gh"
    exit 1
fi

# Check if authenticated
if ! gh auth status &> /dev/null; then
    echo "❌ Not authenticated with GitHub CLI"
    echo "Run: gh auth login"
    exit 1
fi

echo "✅ GitHub CLI is installed and authenticated"
echo ""

# Create labels
echo "📌 Creating labels..."

gh label create "priority:critical" --color "d73a4a" --description "Critical priority - blocking issues" --repo $REPO --force 2>/dev/null || true
gh label create "priority:high" --color "ff9900" --description "High priority - important features" --repo $REPO --force 2>/dev/null || true
gh label create "priority:medium" --color "ffcc00" --description "Medium priority - nice to have" --repo $REPO --force 2>/dev/null || true
gh label create "priority:low" --color "00cc66" --description "Low priority - future enhancements" --repo $REPO --force 2>/dev/null || true

gh label create "testing" --color "8b5cf6" --description "Test infrastructure" --repo $REPO --force 2>/dev/null || true
gh label create "resource" --color "0969da" --description "New resources" --repo $REPO --force 2>/dev/null || true
gh label create "data-source" --color "00c2e0" --description "Data sources" --repo $REPO --force 2>/dev/null || true
gh label create "documentation" --color "6e7781" --description "Docs and examples" --repo $REPO --force 2>/dev/null || true
gh label create "enhancement" --color "10b981" --description "Improvements to existing features" --repo $REPO --force 2>/dev/null || true
gh label create "ci/cd" --color "1f77d0" --description "CI/CD improvements" --repo $REPO --force 2>/dev/null || true

gh label create "blocked" --color "d73a4a" --description "Blocked by dependencies" --repo $REPO --force 2>/dev/null || true
gh label create "needs-research" --color "fbca04" --description "Requires API research" --repo $REPO --force 2>/dev/null || true
gh label create "good-first-issue" --color "7057ff" --description "Good for new contributors" --repo $REPO --force 2>/dev/null || true
gh label create "enterprise" --color "8b5cf6" --description "Enterprise edition only" --repo $REPO --force 2>/dev/null || true

echo "✅ Labels created"
echo ""

# Create milestones
echo "🎯 Creating milestones..."

# Note: Due dates need to be adjusted based on when you run this script
SPRINT1_DATE="2026-02-01T00:00:00Z"
SPRINT2_DATE="2026-02-15T00:00:00Z"
SPRINT3_DATE="2026-03-15T00:00:00Z"
SPRINT4_DATE="2026-04-01T00:00:00Z"
SPRINT5_DATE="2026-05-01T00:00:00Z"

gh api repos/$REPO/milestones -f title="v0.2.0 - Testing Foundation" -f description="Implement comprehensive testing infrastructure" -f due_on="$SPRINT1_DATE" 2>/dev/null || echo "Milestone 'v0.2.0 - Testing Foundation' may already exist"
gh api repos/$REPO/milestones -f title="v0.3.0 - Data Sources" -f description="Add data sources for read-only access" -f due_on="$SPRINT2_DATE" 2>/dev/null || echo "Milestone 'v0.3.0 - Data Sources' may already exist"
gh api repos/$REPO/milestones -f title="v0.4.0 - Core Resources" -f description="Implement notification channels and dashboards" -f due_on="$SPRINT3_DATE" 2>/dev/null || echo "Milestone 'v0.4.0 - Core Resources' may already exist"
gh api repos/$REPO/milestones -f title="v0.5.0 - Documentation" -f description="Complete provider documentation" -f due_on="$SPRINT4_DATE" 2>/dev/null || echo "Milestone 'v0.5.0 - Documentation' may already exist"
gh api repos/$REPO/milestones -f title="v0.6.0 - Advanced Features" -f description="Advanced monitoring features" -f due_on="$SPRINT5_DATE" 2>/dev/null || echo "Milestone 'v0.6.0 - Advanced Features' may already exist"

echo "✅ Milestones created"
echo ""

# Create issues
echo "📝 Creating issues..."

# Issue #1: Unit Test Suite
gh issue create --repo $REPO \
  --title "Implement Unit Test Suite" \
  --label "testing,enhancement,priority:critical" \
  --milestone "v0.2.0 - Testing Foundation" \
  --body "## Description

Implement comprehensive unit tests for all provider components using Go's testing package and testify assertions.

## Scope

### Unit tests for \`internal/provider/monitor_models.go\` (conversion functions)
- \`TestPlanToMonitorInput_MetricMonitor\`
- \`TestPlanToMonitorInput_ErrorMonitor\`
- \`TestMonitorToState_MetricMonitor\`
- \`TestMonitorToState_ErrorMonitor\`
- \`TestConvertMetricParams_AllFields\`
- \`TestConvertErrorParams_AllFields\`
- Null/optional field handling tests

### Unit tests for \`internal/client/client.go\` (API wrapper methods)
- \`TestNew_ValidConfig\`
- \`TestNew_ValidationErrors\`
- \`TestListMonitors_Success\`
- \`TestGetMonitor_Success\`
- \`TestCreateMonitor_Success\`
- \`TestUpdateMonitor_Success\`
- \`TestDeleteMonitor_Success\`
- HTTP error handling tests (401, 404, 500)

### Unit tests for \`internal/provider/provider.go\`
- Schema validation
- Configuration from env vars
- Configuration from Terraform config

## Files to Create
- \`/internal/provider/monitor_models_test.go\`
- \`/internal/provider/monitor_resource_test.go\` (unit tests section)
- \`/internal/client/client_test.go\`
- \`/internal/provider/provider_test.go\`

## Acceptance Criteria
- [ ] Minimum 80% code coverage for conversion logic
- [ ] All HTTP client methods tested with mock server
- [ ] Tests run in < 5 seconds
- [ ] CI passes with \`go test -race -cover ./...\`

## Related
Part of Sprint 1 - Testing Foundation"

echo "✅ Created Issue #1: Unit Test Suite"

# Issue #2: Acceptance Test Framework
gh issue create --repo $REPO \
  --title "Implement Acceptance Test Framework" \
  --label "testing,enhancement,priority:critical" \
  --milestone "v0.2.0 - Testing Foundation" \
  --body "## Description

Implement acceptance tests using \`terraform-plugin-testing\` framework running against Docker Uptrace instance.

## Scope

### Create test helper package \`/internal/acctest/\`
- \`acctest.go\` - PreCheck, provider factory, test utilities
- \`docker.go\` - Docker Compose lifecycle management
- \`fixtures.go\` - Reusable test data
- \`provider.go\` - Test provider factories

### Acceptance tests for monitor resource
- \`TestAccMonitorResource_MetricBasic\` - Create/Read/Delete
- \`TestAccMonitorResource_MetricUpdate\` - Full CRUD lifecycle
- \`TestAccMonitorResource_ErrorBasic\` - Error monitor CRUD
- \`TestAccMonitorResource_Import\` - Import functionality
- \`TestAccMonitorResource_Disappears\` - Resource deletion handling
- \`TestAccMonitorResource_RequiredFields\` - Validation
- \`TestAccMonitorResource_TypeChangeRequiresReplace\` - Force new

### Test fixtures in \`/testdata/\`
- Terraform configurations
- JSON monitor definitions
- Test data for various scenarios

## Files to Create
- \`/internal/acctest/acctest.go\`
- \`/internal/acctest/docker.go\`
- \`/internal/acctest/fixtures.go\`
- \`/internal/acctest/provider.go\`
- \`/internal/provider/monitor_resource_test.go\` (acceptance tests)
- \`/testdata/terraform/*.tf\`
- \`/testdata/monitors/*.json\`
- \`/docker-compose.test.yml\`

## Acceptance Criteria
- [ ] All CRUD operations tested end-to-end
- [ ] Tests run against real Docker Uptrace instance
- [ ] Both metric and error monitors tested
- [ ] Import functionality validated
- [ ] Tests complete in < 10 minutes
- [ ] CI integration with GitHub Actions

## Related
Part of Sprint 1 - Testing Foundation"

echo "✅ Created Issue #2: Acceptance Test Framework"

# Issue #3: CI/CD Test Automation
gh issue create --repo $REPO \
  --title "Add CI/CD Test Automation" \
  --label "ci/cd,testing,priority:high" \
  --milestone "v0.2.0 - Testing Foundation" \
  --body "## Description

Update GitHub Actions workflows to run unit and acceptance tests automatically.

## Scope

### Modify \`.github/workflows/test.yml\`
- Separate jobs for unit and acceptance tests
- Docker Compose setup for acceptance tests
- Coverage reporting to Codecov
- Test result artifacts

### Add new tasks to \`Taskfile.yml\`
- \`task test:unit\` - Run unit tests only
- \`task test:acc\` - Run acceptance tests with Docker
- \`task test:acc:metric\` - Test metric monitors only
- \`task test:acc:error\` - Test error monitors only
- \`task test:all\` - Run complete test suite
- \`task test:coverage\` - Generate coverage reports

## Files to Modify
- \`/.github/workflows/test.yml\`
- \`/Taskfile.yml\`

## Acceptance Criteria
- [ ] Tests run on every PR and push to main
- [ ] Separate unit test job (fast, < 2 min)
- [ ] Separate acceptance test job (slower, < 10 min)
- [ ] Coverage reports generated and uploaded
- [ ] Failed tests block PR merge

## Dependencies
- Issue #1 (Unit tests)
- Issue #2 (Acceptance tests)

## Related
Part of Sprint 1 - Testing Foundation"

echo "✅ Created Issue #3: CI/CD Test Automation"

# Issue #4: Monitor Data Source
gh issue create --repo $REPO \
  --title "Implement Monitor Data Source (Read-Only)" \
  --label "data-source,feature,priority:high" \
  --milestone "v0.3.0 - Data Sources" \
  --body "## Description

Implement \`data \"uptrace_monitor\"\` data source for read-only access to existing monitors.

## Scope
- Create \`/internal/provider/monitor_data_source.go\`
- Schema with filter by ID or name
- Read operation using existing client
- Unit and acceptance tests
- Documentation and examples

## Usage Example
\`\`\`hcl
data \"uptrace_monitor\" \"existing\" {
  id = \"123\"
}

output \"monitor_state\" {
  value = data.uptrace_monitor.existing.state
}
\`\`\`

## Files to Create
- \`/internal/provider/monitor_data_source.go\`
- \`/internal/provider/monitor_data_source_test.go\`
- \`/examples/data-sources/uptrace_monitor/data-source.tf\`

## Acceptance Criteria
- [ ] Read monitor by ID
- [ ] Optional filter by name
- [ ] All monitor attributes available
- [ ] Full test coverage
- [ ] Generated documentation

## Related
Part of Sprint 2 - Data Sources"

echo "✅ Created Issue #4: Monitor Data Source"

# Issue #5: Monitors Data Source (List)
gh issue create --repo $REPO \
  --title "Implement Monitors Data Source (List)" \
  --label "data-source,feature,priority:medium" \
  --milestone "v0.3.0 - Data Sources" \
  --body "## Description

Implement \`data \"uptrace_monitors\"\` data source for listing all monitors in a project.

## Scope
- Create \`/internal/provider/monitors_data_source.go\`
- List all monitors for the configured project
- Optional filtering by type, state, name pattern
- Pagination support if API provides it
- Tests and documentation

## Usage Example
\`\`\`hcl
data \"uptrace_monitors\" \"all_metric\" {
  type = \"metric\"
}

output \"monitor_count\" {
  value = length(data.uptrace_monitors.all_metric.monitors)
}
\`\`\`

## Files to Create
- \`/internal/provider/monitors_data_source.go\`
- \`/internal/provider/monitors_data_source_test.go\`
- \`/examples/data-sources/uptrace_monitors/data-source.tf\`

## Acceptance Criteria
- [ ] List all monitors
- [ ] Filter by type (metric/error)
- [ ] Filter by state
- [ ] Full test coverage
- [ ] Generated documentation

## Related
Part of Sprint 2 - Data Sources"

echo "✅ Created Issue #5: Monitors Data Source (List)"

# Issue #6: Notification Channel Resource
gh issue create --repo $REPO \
  --title "Implement Notification Channel Resource" \
  --label "resource,feature,priority:high,needs-research" \
  --milestone "v0.4.0 - Core Resources" \
  --body "## Description

Implement \`uptrace_notification_channel\` resource for managing alert notification channels.

## Research Needed
- [ ] Review Uptrace API endpoints for notification channels
- [ ] Identify supported channel types (Slack, PagerDuty, Email, Webhook, etc.)
- [ ] Update OpenAPI spec with notification channel endpoints

## Scope
- Update \`/api/openapi.yaml\` with channel endpoints
- Regenerate client with \`task generate\`
- Create \`/internal/provider/notification_channel_resource.go\`
- Support multiple channel types (discriminated union)
- Full CRUD operations
- Tests and documentation

## Usage Example
\`\`\`hcl
resource \"uptrace_notification_channel\" \"slack\" {
  name = \"Engineering Alerts\"
  type = \"slack\"

  slack_config = {
    webhook_url = var.slack_webhook
    channel     = \"#alerts\"
  }
}
\`\`\`

## Files to Create
- \`/internal/provider/notification_channel_resource.go\`
- \`/internal/provider/notification_channel_models.go\`
- \`/internal/provider/notification_channel_resource_test.go\`
- \`/examples/resources/uptrace_notification_channel/*.tf\`

## Acceptance Criteria
- [ ] Support 5+ channel types (Slack, Email, PagerDuty, Webhook, Teams)
- [ ] Full CRUD operations tested
- [ ] Integration with monitor resource (channel_ids)
- [ ] Generated documentation

## Dependencies
- API endpoint research and OpenAPI spec update

## Related
Part of Sprint 3 - Core Resources"

echo "✅ Created Issue #6: Notification Channel Resource"

# Issue #7: Dashboard Resource
gh issue create --repo $REPO \
  --title "Implement Dashboard Resource" \
  --label "resource,feature,priority:high,needs-research" \
  --milestone "v0.4.0 - Core Resources" \
  --body "## Description

Implement \`uptrace_dashboard\` resource for managing Uptrace dashboards.

## Research Needed
- [ ] Review Uptrace dashboard API endpoints
- [ ] Understand dashboard schema (grid-based vs table-based)
- [ ] Identify chart types and configuration options

## Scope
- Update \`/api/openapi.yaml\` with dashboard endpoints
- Regenerate client
- Create dashboard resource implementation
- Support grid-based dashboards initially
- Chart configuration schema
- Full CRUD operations
- Tests and documentation

## Usage Example
\`\`\`hcl
resource \"uptrace_dashboard\" \"services\" {
  name        = \"Service Overview\"
  description = \"Overview of all services\"

  grid_items = [
    {
      chart_type = \"timeseries\"
      title      = \"Request Rate\"
      query      = \"sum(rate(requests_total[5m]))\"
      position   = { x = 0, y = 0, w = 12, h = 6 }
    }
  ]
}
\`\`\`

## Files to Create
- \`/internal/provider/dashboard_resource.go\`
- \`/internal/provider/dashboard_models.go\`
- \`/internal/provider/dashboard_resource_test.go\`
- \`/examples/resources/uptrace_dashboard/*.tf\`

## Acceptance Criteria
- [ ] Create/read/update/delete dashboards
- [ ] Support multiple chart types
- [ ] Grid-based layout configuration
- [ ] Full test coverage
- [ ] Generated documentation

## Dependencies
- API endpoint research and OpenAPI spec update

## Related
Part of Sprint 3 - Core Resources"

echo "✅ Created Issue #7: Dashboard Resource"

# Issue #8: Team Resource
gh issue create --repo $REPO \
  --title "Implement Team Resource (Enterprise)" \
  --label "resource,feature,priority:medium,enterprise,needs-research" \
  --milestone "v0.6.0 - Advanced Features" \
  --body "## Description

Implement \`uptrace_team\` resource for managing teams in Uptrace Enterprise.

## Research Needed
- [ ] Confirm teams API is available in Uptrace
- [ ] Understand team/organization hierarchy
- [ ] Identify required vs optional fields

## Scope
- Update OpenAPI spec with teams endpoints
- Create team resource implementation
- Member management (user assignments)
- Project access control
- Tests and documentation

## Usage Example
\`\`\`hcl
resource \"uptrace_team\" \"platform\" {
  name        = \"Platform Team\"
  description = \"Infrastructure and platform services\"

  members = [
    { user_id = 1, role = \"admin\" },
    { user_id = 2, role = \"member\" }
  ]
}
\`\`\`

## Files to Create
- \`/internal/provider/team_resource.go\`
- \`/internal/provider/team_models.go\`
- \`/internal/provider/team_resource_test.go\`
- \`/examples/resources/uptrace_team/*.tf\`

## Acceptance Criteria
- [ ] Create/read/update/delete teams
- [ ] Manage team members
- [ ] Control project access
- [ ] Test coverage
- [ ] Documentation

## Dependencies
- Enterprise Uptrace instance for testing
- API endpoint confirmation

## Related
Part of Sprint 5 - Advanced Features"

echo "✅ Created Issue #8: Team Resource"

# Issue #9: Generate Provider Documentation
gh issue create --repo $REPO \
  --title "Generate Provider Documentation" \
  --label "documentation,priority:high" \
  --milestone "v0.5.0 - Documentation" \
  --body "## Description

Generate complete provider documentation using \`tfplugindocs\` tool.

## Scope
- Run \`task docs\` to generate documentation
- Review and enhance resource/data source descriptions
- Add detailed attribute documentation
- Create comprehensive examples for all resources
- Document common patterns and best practices

## Files to Generate/Modify
- \`/docs/index.md\` - Provider configuration
- \`/docs/resources/*.md\` - Resource documentation
- \`/docs/data-sources/*.md\` - Data source documentation
- \`/docs/guides/*.md\` - Usage guides

## Acceptance Criteria
- [ ] All resources documented
- [ ] All attributes have descriptions
- [ ] Examples for every resource
- [ ] Provider configuration guide
- [ ] Ready for Terraform Registry (future)

## Dependencies
- All resources and data sources implemented

## Related
Part of Sprint 4 - Documentation"

echo "✅ Created Issue #9: Generate Provider Documentation"

# Issue #10: Comprehensive Examples
gh issue create --repo $REPO \
  --title "Add Comprehensive Examples" \
  --label "documentation,enhancement,priority:medium" \
  --milestone "v0.5.0 - Documentation" \
  --body "## Description

Create comprehensive real-world examples demonstrating provider capabilities.

## Scope
- Complete monitoring setup example
- Multi-channel notification example
- Dashboard with multiple metrics
- Import existing resources example
- Error handling patterns
- Best practices guide

## Files to Create
- \`/examples/complete-setup/\` - Full monitoring stack
- \`/examples/multi-channel-alerts/\` - Complex notification setup
- \`/examples/dashboard-examples/\` - Dashboard configurations
- \`/examples/import/\` - Import existing resources
- \`/docs/guides/getting-started.md\`
- \`/docs/guides/best-practices.md\`

## Acceptance Criteria
- [ ] 5+ complete working examples
- [ ] All examples tested against real Uptrace
- [ ] README in each example directory
- [ ] Best practices documented

## Dependencies
- Issue #9 (documentation generation)

## Related
Part of Sprint 4 - Documentation"

echo "✅ Created Issue #10: Comprehensive Examples"

# Issue #11: Alert Naming Templates
gh issue create --repo $REPO \
  --title "Add Alert Naming Templates" \
  --label "enhancement,priority:medium" \
  --milestone "v0.6.0 - Advanced Features" \
  --body "## Description

Add support for alert naming templates in monitor resource using Go templates.

## Scope
- Add \`alert_name_template\` field to monitor schema
- Support template variables: \`{{ .DisplayName }}\`, \`{{ .Attrs.service_name }}\`
- Validation of template syntax
- Tests for template rendering
- Documentation with examples

## Usage Example
\`\`\`hcl
resource \"uptrace_monitor\" \"cpu\" {
  name                 = \"High CPU\"
  alert_name_template  = \"High CPU on {{ .Attrs.service_name }}\"
  # ... other config
}
\`\`\`

## Files to Modify
- \`/api/openapi.yaml\`
- \`/internal/provider/monitor_resource.go\`
- \`/internal/provider/monitor_models.go\`
- Tests and examples

## Acceptance Criteria
- [ ] Template field added to schema
- [ ] Template variables supported
- [ ] Validation of template syntax
- [ ] Documentation with examples

## Related
Part of Sprint 5 - Advanced Features"

echo "✅ Created Issue #11: Alert Naming Templates"

# Issue #12: Channel Conditions
gh issue create --repo $REPO \
  --title "Add Channel Conditions (Expr Language)" \
  --label "enhancement,priority:medium" \
  --milestone "v0.6.0 - Advanced Features" \
  --body "## Description

Add support for channel conditions using Expr language for smart alert routing.

## Scope
- Add channel condition configuration to notification channels
- Support Expr language functions:
  - \`monitorName()\` - filter by monitor name
  - \`alertName()\` - filter by alert name
  - \`alertType()\` - filter by alert type
  - \`attr()\` - access alert attributes
  - \`hasAttr()\` - check attribute existence
- Validation of Expr expressions
- Tests and documentation

## Usage Example
\`\`\`hcl
resource \"uptrace_notification_channel\" \"critical_only\" {
  name = \"Critical Alerts\"
  type = \"pagerduty\"

  condition = \"attr('severity') == 'critical'\"

  pagerduty_config = {
    routing_key = var.pagerduty_key
  }
}
\`\`\`

## Acceptance Criteria
- [ ] Condition field in channel schema
- [ ] Expr expression validation
- [ ] Documentation with examples
- [ ] Test coverage

## Dependencies
- Issue #6 (Notification channel resource)

## Related
Part of Sprint 5 - Advanced Features"

echo "✅ Created Issue #12: Channel Conditions"

# Issue #13: Service Graph Configuration
gh issue create --repo $REPO \
  --title "Add Service Graph Configuration" \
  --label "resource,feature,priority:low,needs-research" \
  --body "## Description

Research and implement service graph configuration resource if supported by API.

## Research Needed
- [ ] Confirm API endpoints exist
- [ ] Understand service graph configuration options
- [ ] Identify use cases and requirements

## Scope
TBD - depends on research findings

## Related
Future enhancement - low priority"

echo "✅ Created Issue #13: Service Graph Configuration"

# Issue #14: Saved Query Resource
gh issue create --repo $REPO \
  --title "Add Query/Saved Query Resource" \
  --label "resource,feature,priority:low,needs-research" \
  --body "## Description

Implement resource for saving and managing UQL queries.

## Research Needed
- [ ] Confirm saved queries API exists
- [ ] Understand query storage and sharing
- [ ] Identify required fields

## Scope
TBD - depends on research findings

## Related
Future enhancement - low priority"

echo "✅ Created Issue #14: Saved Query Resource"

echo ""
echo "🎉 All done!"
echo ""
echo "Summary:"
echo "- ✅ Created all labels"
echo "- ✅ Created 5 milestones"
echo "- ✅ Created 14 issues"
echo ""
echo "Next steps:"
echo "1. Create GitHub Project: https://github.com/users/riccap/projects/new"
echo "2. Add issues to project board"
echo "3. Configure project views (Board, Sprint, Priority, Type)"
echo "4. Start Sprint 1 with issues #1, #2, #3"
echo ""
echo "View all issues: gh issue list --repo $REPO"
