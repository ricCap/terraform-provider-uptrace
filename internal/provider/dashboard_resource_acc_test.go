package provider

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	acceptancetests "github.com/riccap/tofu-uptrace-provider/internal/acceptance_tests"
)

func TestAccDashboardResource_Basic(t *testing.T) {
	resourceName := "uptrace_dashboard.test"
	dashboardName := acceptancetests.RandomTestName("tf-acc-dashboard")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptancetests.PreCheck(t) },
		ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDashboardDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDashboardResourceConfigBasic(dashboardName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDashboardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", dashboardName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "yaml"),
					resource.TestCheckResourceAttr(resourceName, "pinned", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update testing
			{
				Config: testAccDashboardResourceConfigUpdated(dashboardName + " Updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDashboardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", dashboardName+" Updated"),
				),
			},
		},
	})
}

func TestAccDashboardResource_ComplexYAML(t *testing.T) {
	resourceName := "uptrace_dashboard.test"
	dashboardName := acceptancetests.RandomTestName("tf-acc-dashboard-complex")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptancetests.PreCheck(t) },
		ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDashboardResourceConfigComplex(dashboardName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDashboardExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", dashboardName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "yaml"),
				),
			},
		},
	})
}

func TestAccDashboardResource_Disappears(t *testing.T) {
	resourceName := "uptrace_dashboard.test"
	dashboardName := acceptancetests.RandomTestName("tf-acc-dashboard-disappears")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptancetests.PreCheck(t) },
		ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDashboardResourceConfigBasic(dashboardName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDashboardExists(resourceName),
					testAccCheckDashboardDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Helper functions

func testAccCheckDashboardExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No dashboard ID is set")
		}

		dashboardID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Invalid dashboard ID %s: %w", rs.Primary.ID, err)
		}

		client := acceptancetests.GetTestClient()
		_, err = client.GetDashboard(context.Background(), dashboardID)
		if err != nil {
			return fmt.Errorf("Dashboard %s not found: %w", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccCheckDashboardDestroy(s *terraform.State) error {
	client := acceptancetests.GetTestClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "uptrace_dashboard" {
			continue
		}

		dashboardID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			continue // Skip invalid IDs
		}

		_, err = client.GetDashboard(context.Background(), dashboardID)
		if err == nil {
			return fmt.Errorf("Dashboard %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckDashboardDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		dashboardID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Invalid dashboard ID %s: %w", rs.Primary.ID, err)
		}

		client := acceptancetests.GetTestClient()
		err = client.DeleteDashboard(context.Background(), dashboardID)
		if err != nil {
			return fmt.Errorf("Error deleting dashboard %s: %w", rs.Primary.ID, err)
		}

		return nil
	}
}

// Test configurations

func testAccDashboardResourceConfigBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "uptrace_dashboard" "test" {
  yaml = <<-YAML
    name: %s
    gridRows:
      - items:
          - type: chart
            title: CPU Usage
            params:
              metrics:
                - name: system.cpu.utilization
                  alias: cpu
              query: avg(cpu)
              chartKind: line
              legend:
                show: true
  YAML
}
`, acceptancetests.GetTestProviderConfig(), name)
}

func testAccDashboardResourceConfigUpdated(name string) string {
	return fmt.Sprintf(`
%s

resource "uptrace_dashboard" "test" {
  yaml = <<-YAML
    name: %s
    gridRows:
      - items:
          - type: chart
            title: CPU and Memory Usage
            params:
              metrics:
                - name: system.cpu.utilization
                  alias: cpu
                - name: system.memory.usage
                  alias: mem
              query: avg(cpu), avg(mem)
              chartKind: line
              legend:
                show: true
  YAML
}
`, acceptancetests.GetTestProviderConfig(), name)
}

func testAccDashboardResourceConfigComplex(name string) string {
	return fmt.Sprintf(`
%s

resource "uptrace_dashboard" "test" {
  yaml = <<-YAML
    name: %s
    gridRows:
      - items:
          - type: chart
            title: Request Rate
            params:
              metrics:
                - name: http_requests_total
                  alias: requests
              query: rate(requests[5m])
              chartKind: line
              legend:
                show: true
          - type: table
            title: Service List
            params:
              metrics:
                - name: service_info
                  alias: services
              columns:
                - name: service_name
                  label: Service
      - items:
          - type: gauge
            title: Error Rate
            params:
              metrics:
                - name: http_errors_total
                  alias: errors
              valueMapping:
                - value: 0-5
                  color: green
                - value: 5-10
                  color: yellow
                - value: 10+
                  color: red
          - type: heatmap
            title: Response Time Distribution
            params:
              metrics:
                - name: http_response_time
                  alias: response_time
  YAML
}
`, acceptancetests.GetTestProviderConfig(), name)
}
