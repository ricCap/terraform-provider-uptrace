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
    grid_rows:
      - items:
          - title: CPU Usage
            metrics:
              - system.cpu.utilization as $cpu
            query:
              - avg($cpu)
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
    grid_rows:
      - items:
          - title: CPU and Memory Usage
            metrics:
              - system.cpu.utilization as $cpu
              - system.memory.usage as $mem
            query:
              - avg($cpu)
              - avg($mem)
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
    grid_rows:
      - title: Performance
        items:
          - title: Request Rate
            metrics:
              - http_requests_total as $requests
            query:
              - per_min(sum($requests))
          - title: Error Rate
            metrics:
              - http_errors_total as $errors
            query:
              - per_min(sum($errors))
      - title: Resources
        items:
          - title: CPU Usage
            metrics:
              - system.cpu.utilization as $cpu
            query:
              - avg($cpu)
          - title: Memory Usage
            metrics:
              - system.memory.usage as $mem
            query:
              - avg($mem)
  YAML
}
`, acceptancetests.GetTestProviderConfig(), name)
}
