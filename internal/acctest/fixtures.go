package acctest

import (
	"github.com/riccap/tofu-uptrace-provider/internal/client/generated"
)

// GetMetricMonitorInput returns a basic metric monitor input for testing.
func GetMetricMonitorInput(name string) generated.MonitorInput {
	notifyEmail := false
	maxValue := float64(80)
	checkNumPoint := int(2)
	nullsMode := generated.MetricMonitorParamsNullsModeAllow

	metricName := "system.cpu.utilization"
	alias := "$cpu"

	var params generated.MonitorInput_Params
	//nolint:errcheck // Test fixture: error handling not required
	_ = params.FromMetricMonitorParams(generated.MetricMonitorParams{
		Metrics: []generated.MetricDefinition{
			{
				Name:  metricName,
				Alias: &alias,
			},
		},
		Query:           "avg($cpu) > 80",
		Column:          "value",
		MaxAllowedValue: &maxValue,
		CheckNumPoint:   &checkNumPoint,
		NullsMode:       &nullsMode,
	})

	return generated.MonitorInput{
		Name:                  name,
		Type:                  generated.MonitorTypeMetric,
		NotifyEveryoneByEmail: &notifyEmail,
		TeamIds:               &[]int64{},
		ChannelIds:            &[]int64{},
		Params:                params,
	}
}

// GetErrorMonitorInput returns a basic error monitor input for testing.
func GetErrorMonitorInput(name string) generated.MonitorInput {
	notifyEmail := false
	query := "sum($logs) | where span.event_name exists"

	metricName := "uptrace_tracing_events"
	alias := "$logs"

	var params generated.MonitorInput_Params
	//nolint:errcheck // Test fixture: error handling not required
	_ = params.FromErrorMonitorParams(generated.ErrorMonitorParams{
		Metrics: []generated.MetricDefinition{
			{
				Name:  metricName,
				Alias: &alias,
			},
		},
		Query: &query,
	})

	return generated.MonitorInput{
		Name:                  name,
		Type:                  generated.MonitorTypeError,
		NotifyEveryoneByEmail: &notifyEmail,
		TeamIds:               &[]int64{},
		ChannelIds:            &[]int64{},
		Params:                params,
	}
}

// GetMetricMonitorWithAllFields returns a fully populated metric monitor.
func GetMetricMonitorWithAllFields(name string) generated.MonitorInput {
	notifyEmail := true
	teamIDs := []int64{1}
	channelIDs := []int64{10}
	minValue := float64(0)
	maxValue := float64(90)
	groupingInterval := float64(60000)
	checkNumPoint := int(3)
	nullsMode := generated.MetricMonitorParamsNullsModeForbid
	timeOffset := float64(0)

	metricName := "system.cpu.utilization"
	alias := "$cpu"

	strategy := generated.RepeatIntervalStrategyDefault

	var params generated.MonitorInput_Params
	//nolint:errcheck // Test fixture: error handling not required
	_ = params.FromMetricMonitorParams(generated.MetricMonitorParams{
		Metrics: []generated.MetricDefinition{
			{
				Name:  metricName,
				Alias: &alias,
			},
		},
		Query:            "avg($cpu) > 90",
		Column:           "value",
		MinAllowedValue:  &minValue,
		MaxAllowedValue:  &maxValue,
		GroupingInterval: &groupingInterval,
		CheckNumPoint:    &checkNumPoint,
		NullsMode:        &nullsMode,
		TimeOffset:       &timeOffset,
	})

	return generated.MonitorInput{
		Name:                  name,
		Type:                  generated.MonitorTypeMetric,
		NotifyEveryoneByEmail: &notifyEmail,
		TeamIds:               &teamIDs,
		ChannelIds:            &channelIDs,
		RepeatInterval: &generated.RepeatInterval{
			Strategy: &strategy,
		},
		Params: params,
	}
}
