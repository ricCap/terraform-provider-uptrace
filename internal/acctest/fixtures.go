package acctest

import (
	"github.com/riccap/tofu-uptrace-provider/internal/client/generated"
)

// GetMetricMonitorInput returns a basic metric monitor input for testing
func GetMetricMonitorInput(name string) generated.MonitorInput {
	notifyEmail := false
	maxValue := float64(80)
	checkNumPoint := int64(2)
	nullsMode := "allow"

	metricName := "system.cpu.utilization"
	alias := "$cpu"

	return generated.MonitorInput{
		Name: name,
		Type: generated.MonitorTypeMetric,
		NotifyEveryoneByEmail: &notifyEmail,
		TeamIds:               &[]int64{},
		ChannelIds:            &[]int64{},
		Params: generated.MonitorInput_Params{
			union: generated.MetricMonitorParams{
				Metrics: []generated.MetricDefinition{
					{
						Name:  metricName,
						Alias: &alias,
					},
				},
				Query:            "avg($cpu) > 80",
				MaxAllowedValue:  &maxValue,
				CheckNumPoint:    &checkNumPoint,
				NullsMode:        &nullsMode,
			},
		},
	}
}

// GetErrorMonitorInput returns a basic error monitor input for testing
func GetErrorMonitorInput(name string) generated.MonitorInput {
	notifyEmail := false
	query := "sum($logs) | where span.event_name exists"

	metricName := "uptrace_tracing_events"
	alias := "$logs"

	return generated.MonitorInput{
		Name: name,
		Type: generated.MonitorTypeError,
		NotifyEveryoneByEmail: &notifyEmail,
		TeamIds:               &[]int64{},
		ChannelIds:            &[]int64{},
		Params: generated.MonitorInput_Params{
			union: generated.ErrorMonitorParams{
				Metrics: []generated.MetricDefinition{
					{
						Name:  metricName,
						Alias: &alias,
					},
				},
				Query: &query,
			},
		},
	}
}

// GetMetricMonitorWithAllFields returns a fully populated metric monitor
func GetMetricMonitorWithAllFields(name string) generated.MonitorInput {
	notifyEmail := true
	teamIds := []int64{1}
	channelIds := []int64{10}
	minValue := float64(0)
	maxValue := float64(90)
	groupingInterval := float64(60000)
	checkNumPoint := int64(3)
	nullsMode := "forbid"
	timeOffset := float64(0)
	column := "value"

	metricName := "system.cpu.utilization"
	alias := "$cpu"

	strategy := "default"

	return generated.MonitorInput{
		Name: name,
		Type: generated.MonitorTypeMetric,
		NotifyEveryoneByEmail: &notifyEmail,
		TeamIds:               &teamIds,
		ChannelIds:            &channelIds,
		RepeatInterval: &generated.RepeatInterval{
			Strategy: &strategy,
		},
		Params: generated.MonitorInput_Params{
			union: generated.MetricMonitorParams{
				Metrics: []generated.MetricDefinition{
					{
						Name:  metricName,
						Alias: &alias,
					},
				},
				Query:             "avg($cpu) > 90",
				Column:            &column,
				MinAllowedValue:   &minValue,
				MaxAllowedValue:   &maxValue,
				GroupingInterval:  &groupingInterval,
				CheckNumPoint:     &checkNumPoint,
				NullsMode:         &nullsMode,
				TimeOffset:        &timeOffset,
			},
		},
	}
}
