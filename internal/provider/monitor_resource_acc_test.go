package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/riccap/tofu-uptrace-provider/internal/acctest"
)

func init() {
	// Initialize provider factories to avoid import cycle
	acctest.TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"uptrace": providerserver.NewProtocol6WithError(New("test")()),
	}
}

func TestAccMonitorResource_MetricBasic(t *testing.T) {
	resourceName := "uptrace_monitor.test"
	monitorName := acctest.RandomTestName("tf-acc-metric")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMonitorDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMonitorResourceConfig_MetricBasic(monitorName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", monitorName),
					resource.TestCheckResourceAttr(resourceName, "type", "metric"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "notify_everyone_by_email", "false"),
					resource.TestCheckResourceAttr(resourceName, "params.metrics.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "params.metrics.0.name", "system.cpu.utilization"),
					resource.TestCheckResourceAttr(resourceName, "params.query", "avg($cpu) > 80"),
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
				Config: testAccMonitorResourceConfig_MetricBasic(monitorName + "-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", monitorName+"-updated"),
				),
			},
		},
	})
}

func TestAccMonitorResource_ErrorBasic(t *testing.T) {
	resourceName := "uptrace_monitor.test"
	monitorName := acctest.RandomTestName("tf-acc-error")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_ErrorBasic(monitorName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", monitorName),
					resource.TestCheckResourceAttr(resourceName, "type", "error"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "params.metrics.#", "1"),
				),
			},
		},
	})
}

func TestAccMonitorResource_Disappears(t *testing.T) {
	resourceName := "uptrace_monitor.test"
	monitorName := acctest.RandomTestName("tf-acc-disappears")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorResourceConfig_MetricBasic(monitorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMonitorExists(resourceName),
					testAccCheckMonitorDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Helper functions

func testAccCheckMonitorExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No monitor ID is set")
		}

		client := acctest.GetTestClient()
		_, err := client.GetMonitor(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Monitor %s not found: %w", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccCheckMonitorDestroy(s *terraform.State) error {
	client := acctest.GetTestClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "uptrace_monitor" {
			continue
		}

		_, err := client.GetMonitor(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Monitor %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMonitorDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		client := acctest.GetTestClient()
		return client.DeleteMonitor(context.Background(), rs.Primary.ID)
	}
}

// Configuration helpers

func testAccMonitorResourceConfig_MetricBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "uptrace_monitor" "test" {
  name = "%s"
  type = "metric"

  notify_everyone_by_email = false

  params = {
    metrics = [
      {
        name  = "system.cpu.utilization"
        alias = "$cpu"
      }
    ]
    query             = "avg($cpu) > 80"
    max_allowed_value = 80
    check_num_point   = 2
  }
}
`, acctest.GetTestProviderConfig(), name)
}

func testAccMonitorResourceConfig_ErrorBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "uptrace_monitor" "test" {
  name = "%s"
  type = "error"

  notify_everyone_by_email = false

  params = {
    metrics = [
      {
        name  = "uptrace_tracing_events"
        alias = "$logs"
      }
    ]
    query = "sum($logs) | where span.event_name exists"
  }
}
`, acctest.GetTestProviderConfig(), name)
}
