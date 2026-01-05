package provider

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	acceptancetests "github.com/riccap/tofu-uptrace-provider/internal/acceptance_tests"
)

// TestAccExampleCompleteSetup validates the complete-setup example
func TestAccExampleCompleteSetup(t *testing.T) {
	if testing.Short() {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptancetests.PreCheck(t) },
		ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getExampleConfig(t, "complete-setup"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify notification channels created
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.slack_critical", "id"),
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.slack_warnings", "id"),
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.webhook_pagerduty", "id"),
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.telegram_oncall", "id"),

					// Verify monitors created
					resource.TestCheckResourceAttrSet("uptrace_monitor.critical_error_rate", "id"),
					resource.TestCheckResourceAttrSet("uptrace_monitor.warning_4xx_errors", "id"),
					resource.TestCheckResourceAttrSet("uptrace_monitor.critical_database_errors", "id"),
					resource.TestCheckResourceAttrSet("uptrace_monitor.critical_api_latency", "id"),
					resource.TestCheckResourceAttrSet("uptrace_monitor.warning_db_slow_queries", "id"),
					resource.TestCheckResourceAttrSet("uptrace_monitor.warning_traffic_spike", "id"),
					resource.TestCheckResourceAttrSet("uptrace_monitor.critical_no_traffic", "id"),

					// Verify dashboards created
					resource.TestCheckResourceAttrSet("uptrace_dashboard.application_overview", "id"),
					resource.TestCheckResourceAttrSet("uptrace_dashboard.error_tracking", "id"),

					// Verify channel names
					resource.TestCheckResourceAttr("uptrace_notification_channel.slack_critical", "name", "Critical Alerts - Slack"),
					resource.TestCheckResourceAttr("uptrace_notification_channel.slack_warnings", "name", "Warning Alerts - Slack"),

					// Verify monitor names
					resource.TestCheckResourceAttr("uptrace_monitor.critical_error_rate", "name", "Critical: High Error Rate"),
					resource.TestCheckResourceAttr("uptrace_monitor.warning_4xx_errors", "name", "Warning: High 4xx Error Rate"),
				),
			},
		},
	})
}

// TestAccExampleMultiChannelAlerts validates the multi-channel-alerts example
func TestAccExampleMultiChannelAlerts(t *testing.T) {
	if testing.Short() {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptancetests.PreCheck(t) },
		ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getExampleConfig(t, "multi-channel-alerts"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify severity-based channels
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.sev1_pagerduty", "id"),
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.sev2_slack", "id"),
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.sev3_slack", "id"),

					// Verify team-based channels
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.team_backend", "id"),
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.team_frontend", "id"),
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.team_data", "id"),

					// Verify fan-out channels
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.fanout_slack", "id"),
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.fanout_telegram", "id"),
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.fanout_webhook", "id"),

					// Verify monitors with routing
					resource.TestCheckResourceAttrSet("uptrace_monitor.sev1_complete_outage", "id"),
					resource.TestCheckResourceAttrSet("uptrace_monitor.sev2_high_error_rate", "id"),
					resource.TestCheckResourceAttrSet("uptrace_monitor.sev3_elevated_latency", "id"),
					resource.TestCheckResourceAttrSet("uptrace_monitor.backend_api_errors", "id"),
					resource.TestCheckResourceAttrSet("uptrace_monitor.frontend_client_errors", "id"),
					resource.TestCheckResourceAttrSet("uptrace_monitor.data_pipeline_errors", "id"),
					resource.TestCheckResourceAttrSet("uptrace_monitor.critical_payment_failures", "id"),

					// Verify monitor names match routing patterns
					resource.TestCheckResourceAttr("uptrace_monitor.sev1_complete_outage", "name", "SEV1: Complete Service Outage"),
					resource.TestCheckResourceAttr("uptrace_monitor.sev2_high_error_rate", "name", "SEV2: High Error Rate"),
					resource.TestCheckResourceAttr("uptrace_monitor.sev3_elevated_latency", "name", "SEV3: Elevated API Latency"),
				),
			},
		},
	})
}

// TestAccExampleDashboards validates the dashboard-examples
func TestAccExampleDashboards(t *testing.T) {
	if testing.Short() {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptancetests.PreCheck(t) },
		ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getExampleConfig(t, "dashboard-examples"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify all 8 dashboards created
					resource.TestCheckResourceAttrSet("uptrace_dashboard.service_overview", "id"),
					resource.TestCheckResourceAttrSet("uptrace_dashboard.service_comparison", "id"),
					resource.TestCheckResourceAttrSet("uptrace_dashboard.red_metrics", "id"),
					resource.TestCheckResourceAttrSet("uptrace_dashboard.database_performance", "id"),
					resource.TestCheckResourceAttrSet("uptrace_dashboard.http_status_codes", "id"),
					resource.TestCheckResourceAttrSet("uptrace_dashboard.external_dependencies", "id"),
					resource.TestCheckResourceAttrSet("uptrace_dashboard.business_metrics", "id"),
					resource.TestCheckResourceAttrSet("uptrace_dashboard.microservices_map", "id"),

					// Verify dashboard names
					resource.TestCheckResourceAttr("uptrace_dashboard.service_overview", "name", "Service Overview - Simple"),
					resource.TestCheckResourceAttr("uptrace_dashboard.service_comparison", "name", "Multi-Service Comparison"),
					resource.TestCheckResourceAttr("uptrace_dashboard.red_metrics", "name", "RED Metrics Dashboard"),
					resource.TestCheckResourceAttr("uptrace_dashboard.database_performance", "name", "Database Performance"),
					resource.TestCheckResourceAttr("uptrace_dashboard.http_status_codes", "name", "HTTP Status Codes"),
					resource.TestCheckResourceAttr("uptrace_dashboard.external_dependencies", "name", "External Dependencies"),
					resource.TestCheckResourceAttr("uptrace_dashboard.business_metrics", "name", "Business Metrics"),
					resource.TestCheckResourceAttr("uptrace_dashboard.microservices_map", "name", "Microservices Health Map"),

					// Verify dashboards have YAML config
					resource.TestCheckResourceAttrSet("uptrace_dashboard.service_overview", "yaml_config"),
					resource.TestCheckResourceAttrSet("uptrace_dashboard.red_metrics", "yaml_config"),
				),
			},
		},
	})
}

// TestAccExampleResourceFiles validates individual resource examples
func TestAccExampleResourceFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
		return
	}

	// Test monitor resource example
	t.Run("MonitorResource", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acceptancetests.PreCheck(t) },
			ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: getResourceExampleConfig(t, "uptrace_monitor"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("uptrace_monitor.cpu_usage", "id"),
						resource.TestCheckResourceAttrSet("uptrace_monitor.error_rate", "id"),
					),
				},
			},
		})
	})

	// Test dashboard resource example
	t.Run("DashboardResource", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acceptancetests.PreCheck(t) },
			ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: getResourceExampleConfig(t, "uptrace_dashboard"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("uptrace_dashboard.example", "id"),
						resource.TestCheckResourceAttr("uptrace_dashboard.example", "name", "Example Dashboard"),
					),
				},
			},
		})
	})

	// Test notification channel resource example
	t.Run("NotificationChannelResource", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acceptancetests.PreCheck(t) },
			ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: getResourceExampleConfig(t, "uptrace_notification_channel"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("uptrace_notification_channel.slack", "id"),
						resource.TestCheckResourceAttr("uptrace_notification_channel.slack", "name", "Engineering Alerts"),
						resource.TestCheckResourceAttr("uptrace_notification_channel.slack", "type", "slack"),
					),
				},
			},
		})
	})
}

// getExampleConfig reads and prepares an example configuration for testing
func getExampleConfig(t *testing.T, exampleName string) string {
	t.Helper()

	// Get the repo root
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %s", err)
	}

	// Navigate to repo root (from internal/provider to root)
	repoRoot := filepath.Join(cwd, "..", "..")
	examplePath := filepath.Join(repoRoot, "examples", exampleName, "main.tf")

	// Read the example file
	content, err := os.ReadFile(examplePath)
	if err != nil {
		t.Fatalf("Failed to read example %s: %s", exampleName, err)
	}

	// Add provider configuration with test credentials
	// Note: terraform and provider blocks are stripped from examples
	// The test framework provides the provider via ProtoV6ProviderFactories
	config := ""

	// Add test variable values based on example
	switch exampleName {
	case "complete-setup":
		config += `
variable "uptrace_endpoint" { default = "http://localhost:14318/internal/v1" }
variable "uptrace_token" { default = "user1_secret_token" }
variable "uptrace_project_id" { default = 1 }
variable "slack_critical_webhook" { default = "https://hooks.slack.com/test/critical" }
variable "slack_warnings_webhook" { default = "https://hooks.slack.com/test/warnings" }
variable "pagerduty_webhook_url" { default = "https://events.pagerduty.com/test" }
variable "telegram_bot_token" { default = "123456:test_token" }
variable "telegram_chat_id" { default = "-1001234567890" }
`
	case "multi-channel-alerts":
		config += `
variable "uptrace_endpoint" { default = "http://localhost:14318/internal/v1" }
variable "uptrace_token" { default = "user1_secret_token" }
variable "uptrace_project_id" { default = 1 }
variable "pagerduty_url" { default = "https://events.pagerduty.com/test/sev1" }
variable "slack_oncall_webhook" { default = "https://hooks.slack.com/test/oncall" }
variable "slack_engineering_webhook" { default = "https://hooks.slack.com/test/eng" }
variable "slack_backend_webhook" { default = "https://hooks.slack.com/test/backend" }
variable "slack_frontend_webhook" { default = "https://hooks.slack.com/test/frontend" }
variable "slack_data_webhook" { default = "https://hooks.slack.com/test/data" }
variable "slack_critical_webhook" { default = "https://hooks.slack.com/test/critical" }
variable "telegram_bot_token" { default = "123456:test_token" }
variable "telegram_chat_id" { default = "-1001234567890" }
variable "custom_webhook_url" { default = "https://example.com/webhook" }
variable "slack_production_webhook" { default = "https://hooks.slack.com/test/prod" }
variable "slack_staging_webhook" { default = "https://hooks.slack.com/test/staging" }
variable "slack_business_hours_webhook" { default = "https://hooks.slack.com/test/biz" }
variable "pagerduty_after_hours_url" { default = "https://events.pagerduty.com/test/after" }
`
	case "dashboard-examples":
		config += `
variable "uptrace_endpoint" { default = "http://localhost:14318/internal/v1" }
variable "uptrace_token" { default = "user1_secret_token" }
variable "uptrace_project_id" { default = 1 }
`
	}

	// Strip the terraform and provider blocks from the example file
	// (we provide our own above with test values)
	exampleContent := string(content)

	// Simple approach: skip until we find the first resource
	// This works because all examples start with terraform{} then provider{}
	// Find first "resource " or "data " and take everything from there
	lines := ""
	inResource := false
	for _, line := range splitLines(exampleContent) {
		if !inResource && (contains(line, "resource \"") || contains(line, "data \"") || contains(line, "# ===")) {
			inResource = true
		}
		if inResource {
			lines += line + "\n"
		}
	}

	return config + lines
}

// getResourceExampleConfig reads individual resource example files
func getResourceExampleConfig(t *testing.T, resourceName string) string {
	t.Helper()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %s", err)
	}

	repoRoot := filepath.Join(cwd, "..", "..")
	examplePath := filepath.Join(repoRoot, "examples", "resources", resourceName, "resource.tf")

	content, err := os.ReadFile(examplePath)
	if err != nil {
		t.Fatalf("Failed to read resource example %s: %s", resourceName, err)
	}

	// Test variables for resource examples
	config := `
variable "slack_webhook_url" { default = "https://hooks.slack.com/test" }
variable "telegram_bot_token" { default = "123456:test_token" }
variable "mattermost_webhook_url" { default = "https://mattermost.example.com/test" }

`

	return config + string(content)
}

// Helper functions
func splitLines(s string) []string {
	result := []string{}
	line := ""
	for _, c := range s {
		if c == '\n' {
			result = append(result, line)
			line = ""
		} else {
			line += string(c)
		}
	}
	if line != "" {
		result = append(result, line)
	}
	return result
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
