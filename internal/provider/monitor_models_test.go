package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/riccap/tofu-uptrace-provider/internal/client/generated"
)

func TestPlanToMonitorInput_MetricMonitor(t *testing.T) {
	ctx := context.Background()

	// Create metric definition
	metricAttrs := map[string]attr.Value{
		"name":  types.StringValue("system.cpu.utilization"),
		"alias": types.StringValue("$cpu"),
	}
	metricObj := types.ObjectValueMust(
		map[string]attr.Type{
			"name":  types.StringType,
			"alias": types.StringType,
		},
		metricAttrs,
	)

	// Create params
	paramsAttrs := map[string]attr.Value{
		"metrics": types.ListValueMust(
			types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"alias": types.StringType,
			}},
			[]attr.Value{metricObj},
		),
		"query":              types.StringValue("avg($cpu) > 80"),
		"column":             types.StringValue("value"),
		"max_allowed_value":  types.Float64Value(80),
		"check_num_point":    types.Int64Value(2),
		"nulls_mode":         types.StringValue("allow"),
		"grouping_interval":  types.Float64Value(60000),
		"time_offset":        types.Float64Value(0),
		"min_allowed_value":  types.Float64Null(),
	}

	paramsType := map[string]attr.Type{
		"metrics": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"alias": types.StringType,
		}}},
		"query":              types.StringType,
		"column":             types.StringType,
		"max_allowed_value":  types.Float64Type,
		"min_allowed_value":  types.Float64Type,
		"grouping_interval":  types.Float64Type,
		"check_num_point":    types.Int64Type,
		"nulls_mode":         types.StringType,
		"time_offset":        types.Float64Type,
	}

	plan := MonitorResourceModel{
		Name:                  types.StringValue("Test Metric Monitor"),
		Type:                  types.StringValue("metric"),
		NotifyEveryoneByEmail: types.BoolValue(false),
		TeamIds:               types.ListValueMust(types.Int64Type, []attr.Value{}),
		ChannelIds:            types.ListValueMust(types.Int64Type, []attr.Value{}),
		Params:                types.ObjectValueMust(paramsType, paramsAttrs),
	}

	diags := diag.Diagnostics{}
	input := planToMonitorInput(ctx, plan, &diags)

	require.False(t, diags.HasError(), "Conversion should not produce errors")
	assert.Equal(t, "Test Metric Monitor", input.Name)
	assert.Equal(t, generated.MonitorTypeMetric, input.Type)
	assert.NotNil(t, input.NotifyEveryoneByEmail)
	assert.False(t, *input.NotifyEveryoneByEmail)

	// Verify metric params
	metricParams, err := input.Params.AsMetricMonitorParams()
	require.NoError(t, err)
	assert.Len(t, metricParams.Metrics, 1)
	assert.Equal(t, "system.cpu.utilization", metricParams.Metrics[0].Name)
	assert.Equal(t, "$cpu", *metricParams.Metrics[0].Alias)
	assert.Equal(t, "avg($cpu) > 80", metricParams.Query)
	assert.Equal(t, float64(80), *metricParams.MaxAllowedValue)
}

func TestPlanToMonitorInput_ErrorMonitor(t *testing.T) {
	ctx := context.Background()

	// Create metric definition
	metricAttrs := map[string]attr.Value{
		"name":  types.StringValue("uptrace_tracing_events"),
		"alias": types.StringValue("$logs"),
	}
	metricObj := types.ObjectValueMust(
		map[string]attr.Type{
			"name":  types.StringType,
			"alias": types.StringType,
		},
		metricAttrs,
	)

	// Create params for error monitor
	paramsAttrs := map[string]attr.Value{
		"metrics": types.ListValueMust(
			types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"alias": types.StringType,
			}},
			[]attr.Value{metricObj},
		),
		"query":              types.StringValue("sum($logs) | where span.event_name exists"),
		"column":             types.StringNull(),
		"max_allowed_value":  types.Float64Null(),
		"min_allowed_value":  types.Float64Null(),
		"grouping_interval":  types.Float64Null(),
		"check_num_point":    types.Int64Null(),
		"nulls_mode":         types.StringNull(),
		"time_offset":        types.Float64Null(),
	}

	paramsType := map[string]attr.Type{
		"metrics": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"alias": types.StringType,
		}}},
		"query":              types.StringType,
		"column":             types.StringType,
		"max_allowed_value":  types.Float64Type,
		"min_allowed_value":  types.Float64Type,
		"grouping_interval":  types.Float64Type,
		"check_num_point":    types.Int64Type,
		"nulls_mode":         types.StringType,
		"time_offset":        types.Float64Type,
	}

	plan := MonitorResourceModel{
		Name:                  types.StringValue("Test Error Monitor"),
		Type:                  types.StringValue("error"),
		NotifyEveryoneByEmail: types.BoolValue(false),
		TeamIds:               types.ListValueMust(types.Int64Type, []attr.Value{}),
		ChannelIds:            types.ListValueMust(types.Int64Type, []attr.Value{}),
		Params:                types.ObjectValueMust(paramsType, paramsAttrs),
	}

	diags := diag.Diagnostics{}
	input := planToMonitorInput(ctx, plan, &diags)

	require.False(t, diags.HasError(), "Conversion should not produce errors")
	assert.Equal(t, "Test Error Monitor", input.Name)
	assert.Equal(t, generated.MonitorTypeError, input.Type)

	// Verify error params
	errorParams, err := input.Params.AsErrorMonitorParams()
	require.NoError(t, err)
	assert.Len(t, errorParams.Metrics, 1)
	assert.Equal(t, "uptrace_tracing_events", errorParams.Metrics[0].Name)
	assert.Equal(t, "$logs", *errorParams.Metrics[0].Alias)
	assert.Equal(t, "sum($logs) | where span.event_name exists", *errorParams.Query)
}

func TestMonitorToState_MetricMonitor(t *testing.T) {
	ctx := context.Background()

	// Create API monitor response
	notifyEmail := false
	maxValue := float64(90)
	checkNumPoint := int(3)
	nullsMode := generated.MetricMonitorParamsNullsModeAllow
	groupingInterval := float64(60000)
	timeOffset := float64(0)

	alias := "$cpu"

	var params generated.Monitor_Params
	_ = params.FromMetricMonitorParams(generated.MetricMonitorParams{
		Metrics: []generated.MetricDefinition{
			{
				Name:  "system.cpu.utilization",
				Alias: &alias,
			},
		},
		Query:            "avg($cpu) > 90",
		Column:           "value",
		MaxAllowedValue:  &maxValue,
		CheckNumPoint:    &checkNumPoint,
		NullsMode:        &nullsMode,
		GroupingInterval: &groupingInterval,
		TimeOffset:       &timeOffset,
	})

	monitor := &generated.Monitor{
		Id:                    123,
		Name:                  "Test Metric Monitor",
		State:                 generated.MonitorStateOpen,
		NotifyEveryoneByEmail: &notifyEmail,
		TeamIds:               &[]int64{},
		ChannelIds:            &[]int64{},
		Params:                params,
		CreatedAt:             func() *float64 { v := float64(1704067200000); return &v }(),
		UpdatedAt:             func() *float64 { v := float64(1704067200000); return &v }(),
	}

	state := &MonitorResourceModel{}
	diags := diag.Diagnostics{}

	monitorToState(ctx, monitor, state, &diags)

	require.False(t, diags.HasError())
	assert.Equal(t, "123", state.ID.ValueString())
	assert.Equal(t, "Test Metric Monitor", state.Name.ValueString())
	assert.Equal(t, "open", state.State.ValueString())
	assert.False(t, state.NotifyEveryoneByEmail.ValueBool())
	assert.Equal(t, "1704067200000", state.CreatedAt.ValueString())
}

func TestMonitorToState_ErrorMonitor(t *testing.T) {
	ctx := context.Background()

	// Create API monitor response
	notifyEmail := false
	query := "sum($logs) | where error exists"
	alias := "$logs"

	var params generated.Monitor_Params
	_ = params.FromErrorMonitorParams(generated.ErrorMonitorParams{
		Metrics: []generated.MetricDefinition{
			{
				Name:  "uptrace_tracing_events",
				Alias: &alias,
			},
		},
		Query: &query,
	})

	monitor := &generated.Monitor{
		Id:                    456,
		Name:                  "Test Error Monitor",
		Type:                  generated.MonitorTypeError,
		State:                 generated.MonitorStateFiring,
		NotifyEveryoneByEmail: &notifyEmail,
		TeamIds:               &[]int64{},
		ChannelIds:            &[]int64{},
		Params:                params,
		CreatedAt:             func() *float64 { v := float64(1704067200000); return &v }(),
		UpdatedAt:             func() *float64 { v := float64(1704067200000); return &v }(),
	}

	state := &MonitorResourceModel{}
	diags := diag.Diagnostics{}

	monitorToState(ctx, monitor, state, &diags)

	require.False(t, diags.HasError())
	assert.Equal(t, "456", state.ID.ValueString())
	assert.Equal(t, "Test Error Monitor", state.Name.ValueString())
	assert.Equal(t, "error", state.Type.ValueString())
	assert.Equal(t, "firing", state.State.ValueString())
	assert.False(t, state.NotifyEveryoneByEmail.ValueBool())
}
