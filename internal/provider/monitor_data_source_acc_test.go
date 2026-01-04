package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/riccap/tofu-uptrace-provider/internal/acctest"
)

func TestAccMonitorDataSource_Basic(t *testing.T) {
	resourceName := "uptrace_monitor.test"
	dataSourceName := "data.uptrace_monitor.test"
	monitorName := acctest.RandomTestName("tf-acc-ds-basic")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorDataSourceConfigBasic(monitorName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify data source attributes match resource
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "type", resourceName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "state", resourceName, "state"),
					resource.TestCheckResourceAttrPair(dataSourceName, "params.metrics.0.name", resourceName, "params.metrics.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "params.query", resourceName, "params.query"),

					// Verify computed fields are set
					resource.TestCheckResourceAttrSet(dataSourceName, "created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "updated_at"),
				),
			},
		},
	})
}

func TestAccMonitorDataSource_AllFields(t *testing.T) {
	resourceName := "uptrace_monitor.test"
	dataSourceName := "data.uptrace_monitor.test"
	monitorName := acctest.RandomTestName("tf-acc-ds-full")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorDataSourceConfigAllFields(monitorName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify all main attributes match
					resource.TestCheckResourceAttrPair(dataSourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "type", resourceName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "state", resourceName, "state"),

					// Verify notification settings
					resource.TestCheckResourceAttrPair(dataSourceName, "notify_everyone_by_email", resourceName, "notify_everyone_by_email"),

					// Verify params
					resource.TestCheckResourceAttrPair(dataSourceName, "params.metrics.#", resourceName, "params.metrics.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "params.metrics.0.name", resourceName, "params.metrics.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "params.metrics.0.alias", resourceName, "params.metrics.0.alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "params.query", resourceName, "params.query"),
					resource.TestCheckResourceAttrPair(dataSourceName, "params.max_allowed_value", resourceName, "params.max_allowed_value"),
					resource.TestCheckResourceAttrPair(dataSourceName, "params.check_num_point", resourceName, "params.check_num_point"),

					// Verify repeat interval
					resource.TestCheckResourceAttrPair(dataSourceName, "repeat_interval.strategy", resourceName, "repeat_interval.strategy"),

					// Verify timestamps
					resource.TestCheckResourceAttrSet(dataSourceName, "created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "updated_at"),
				),
			},
		},
	})
}

func TestAccMonitorDataSource_NotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccMonitorDataSourceConfigNotFound(),
				ExpectError: regexp.MustCompile("(Error Reading Monitor|not found)"),
			},
		},
	})
}

// testAccMonitorDataSourceConfigBasic generates a basic test configuration.
func testAccMonitorDataSourceConfigBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "uptrace_monitor" "test" {
  name = "%s"
  type = "metric"

  params = {
    metrics = [{
      name  = "system.cpu.utilization"
      alias = "$cpu"
    }]
    query             = "avg($cpu) > 80"
    max_allowed_value = 80
  }
}

data "uptrace_monitor" "test" {
  id = uptrace_monitor.test.id
}
`, acctest.GetTestProviderConfig(), name)
}

// testAccMonitorDataSourceConfigAllFields generates a full configuration with all fields.
func testAccMonitorDataSourceConfigAllFields(name string) string {
	return fmt.Sprintf(`
%s

resource "uptrace_monitor" "test" {
  name                     = "%s"
  type                     = "metric"
  notify_everyone_by_email = false

  repeat_interval = {
    strategy = "default"
  }

  params = {
    metrics = [{
      name  = "system.cpu.utilization"
      alias = "$cpu"
    }]
    query             = "avg($cpu) > 90"
    max_allowed_value = 90
    check_num_point   = 3
  }
}

data "uptrace_monitor" "test" {
  id = uptrace_monitor.test.id
}
`, acctest.GetTestProviderConfig(), name)
}

// testAccMonitorDataSourceConfigNotFound generates a configuration with non-existent ID.
func testAccMonitorDataSourceConfigNotFound() string {
	return fmt.Sprintf(`
%s

data "uptrace_monitor" "test" {
  id = "99999999"
}
`, acctest.GetTestProviderConfig())
}
