package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	acceptancetests "github.com/riccap/terraform-provider-uptrace/internal/acceptance_tests"
)

func TestAccMonitorsDataSource_All(t *testing.T) {
	dataSourceName := "data.uptrace_monitors.test"
	name1 := acceptancetests.RandomTestName("tf-acc-ds-all-1")
	name2 := acceptancetests.RandomTestName("tf-acc-ds-all-2")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptancetests.PreCheck(t) },
		ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorsDataSourceConfigAll(name1, name2),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify we get at least 2 monitors (the ones we created)
					resource.TestMatchResourceAttr(dataSourceName, "monitors.#", regexp.MustCompile(`^[2-9]|[1-9]\d+$`)),
					// Verify monitors have expected attributes
					resource.TestCheckResourceAttrSet(dataSourceName, "monitors.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "monitors.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "monitors.0.type"),
					// Note: state field not checked as it may not be set immediately after monitor creation
				),
			},
		},
	})
}

func TestAccMonitorsDataSource_FilterByType(t *testing.T) {
	dataSourceName := "data.uptrace_monitors.test"
	metricName1 := acceptancetests.RandomTestName("tf-acc-ds-metric-1")
	metricName2 := acceptancetests.RandomTestName("tf-acc-ds-metric-2")
	errorName := acceptancetests.RandomTestName("tf-acc-ds-error")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptancetests.PreCheck(t) },
		ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorsDataSourceConfigFilterByType(metricName1, metricName2, errorName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify we get at least 2 metric monitors
					resource.TestMatchResourceAttr(dataSourceName, "monitors.#", regexp.MustCompile(`^[2-9]|[1-9]\d+$`)),
					// Verify all returned monitors are of type "metric"
					resource.TestCheckResourceAttr(dataSourceName, "monitors.0.type", "metric"),
					resource.TestCheckResourceAttr(dataSourceName, "monitors.1.type", "metric"),
				),
			},
		},
	})
}

func TestAccMonitorsDataSource_FilterByName(t *testing.T) {
	dataSourceName := "data.uptrace_monitors.test"
	cpuName := acceptancetests.RandomTestName("tf-acc-ds-CPU-monitor")
	memName := acceptancetests.RandomTestName("tf-acc-ds-memory-monitor")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptancetests.PreCheck(t) },
		ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorsDataSourceConfigFilterByName(cpuName, memName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify we get at least 1 monitor with "CPU" in the name
					resource.TestMatchResourceAttr(dataSourceName, "monitors.#", regexp.MustCompile(`^[1-9]\d*$`)),
					// Verify the returned monitor has "CPU" in the name
					resource.TestMatchResourceAttr(dataSourceName, "monitors.0.name", regexp.MustCompile("(?i)cpu")),
				),
			},
		},
	})
}

func TestAccMonitorsDataSource_EmptyResults(t *testing.T) {
	dataSourceName := "data.uptrace_monitors.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptancetests.PreCheck(t) },
		ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorsDataSourceConfigEmptyResults(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify we get 0 monitors
					resource.TestCheckResourceAttr(dataSourceName, "monitors.#", "0"),
				),
			},
		},
	})
}

func TestAccMonitorsDataSource_MultipleFilters(t *testing.T) {
	dataSourceName := "data.uptrace_monitors.test"
	metricCPUName := acceptancetests.RandomTestName("tf-acc-ds-CPU-metric")
	metricMemName := acceptancetests.RandomTestName("tf-acc-ds-memory-metric")
	errorName := acceptancetests.RandomTestName("tf-acc-ds-CPU-error")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptancetests.PreCheck(t) },
		ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorsDataSourceConfigMultipleFilters(metricCPUName, metricMemName, errorName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify we get at least 1 monitor (metric type with "CPU" in name)
					resource.TestMatchResourceAttr(dataSourceName, "monitors.#", regexp.MustCompile(`^[1-9]\d*$`)),
					// Verify the monitor is type metric
					resource.TestCheckResourceAttr(dataSourceName, "monitors.0.type", "metric"),
					// Verify the monitor has "CPU" in the name
					resource.TestMatchResourceAttr(dataSourceName, "monitors.0.name", regexp.MustCompile("(?i)cpu")),
				),
			},
		},
	})
}

// testAccMonitorsDataSourceConfigAll generates a config that lists all monitors.
func testAccMonitorsDataSourceConfigAll(name1, name2 string) string {
	return fmt.Sprintf(`
%s

resource "uptrace_monitor" "metric1" {
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

resource "uptrace_monitor" "metric2" {
  name = "%s"
  type = "metric"

  params = {
    metrics = [{
      name  = "system.memory.usage"
      alias = "$mem"
    }]
    query             = "avg($mem) > 90"
    max_allowed_value = 90
  }
}

data "uptrace_monitors" "test" {
  depends_on = [
    uptrace_monitor.metric1,
    uptrace_monitor.metric2,
  ]
}
`, acceptancetests.GetTestProviderConfig(), name1, name2)
}

// testAccMonitorsDataSourceConfigFilterByType generates a config that filters by type.
func testAccMonitorsDataSourceConfigFilterByType(metricName1, metricName2, errorName string) string {
	return fmt.Sprintf(`
%s

resource "uptrace_monitor" "metric1" {
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

resource "uptrace_monitor" "metric2" {
  name = "%s"
  type = "metric"

  params = {
    metrics = [{
      name  = "system.memory.usage"
      alias = "$mem"
    }]
    query             = "avg($mem) > 90"
    max_allowed_value = 90
  }
}

resource "uptrace_monitor" "error1" {
  name = "%s"
  type = "error"

  params = {
    metrics = [{
      name  = "uptrace_tracing_events"
      alias = "$logs"
    }]
    query = "sum($logs) | where span.event_name exists"
  }
}

data "uptrace_monitors" "test" {
  type = "metric"

  depends_on = [
    uptrace_monitor.metric1,
    uptrace_monitor.metric2,
    uptrace_monitor.error1,
  ]
}
`, acceptancetests.GetTestProviderConfig(), metricName1, metricName2, errorName)
}

// testAccMonitorsDataSourceConfigFilterByName generates a config that filters by name.
func testAccMonitorsDataSourceConfigFilterByName(cpuName, memName string) string {
	return fmt.Sprintf(`
%s

resource "uptrace_monitor" "cpu" {
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

resource "uptrace_monitor" "mem" {
  name = "%s"
  type = "metric"

  params = {
    metrics = [{
      name  = "system.memory.usage"
      alias = "$mem"
    }]
    query             = "avg($mem) > 90"
    max_allowed_value = 90
  }
}

data "uptrace_monitors" "test" {
  name = "CPU"

  depends_on = [
    uptrace_monitor.cpu,
    uptrace_monitor.mem,
  ]
}
`, acceptancetests.GetTestProviderConfig(), cpuName, memName)
}

// testAccMonitorsDataSourceConfigEmptyResults generates a config that returns no results.
func testAccMonitorsDataSourceConfigEmptyResults() string {
	return fmt.Sprintf(`
%s

data "uptrace_monitors" "test" {
  name = "this-monitor-name-should-never-exist-in-any-test-run-12345"
}
`, acceptancetests.GetTestProviderConfig())
}

// testAccMonitorsDataSourceConfigMultipleFilters generates a config with multiple filters.
func testAccMonitorsDataSourceConfigMultipleFilters(metricCPUName, metricMemName, errorName string) string {
	return fmt.Sprintf(`
%s

resource "uptrace_monitor" "metric_cpu" {
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

resource "uptrace_monitor" "metric_mem" {
  name = "%s"
  type = "metric"

  params = {
    metrics = [{
      name  = "system.memory.usage"
      alias = "$mem"
    }]
    query             = "avg($mem) > 90"
    max_allowed_value = 90
  }
}

resource "uptrace_monitor" "error_cpu" {
  name = "%s"
  type = "error"

  params = {
    metrics = [{
      name  = "uptrace_tracing_events"
      alias = "$logs"
    }]
    query = "sum($logs) | where span.event_name exists"
  }
}

data "uptrace_monitors" "test" {
  type = "metric"
  name = "CPU"

  depends_on = [
    uptrace_monitor.metric_cpu,
    uptrace_monitor.metric_mem,
    uptrace_monitor.error_cpu,
  ]
}
`, acceptancetests.GetTestProviderConfig(), metricCPUName, metricMemName, errorName)
}
