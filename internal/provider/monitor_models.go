package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/riccap/tofu-uptrace-provider/internal/client/generated"
)

// RepeatIntervalModel represents the repeat interval configuration.
type RepeatIntervalModel struct {
	Strategy types.String `tfsdk:"strategy"`
	Interval types.Int64  `tfsdk:"interval"`
}

// MonitorParamsModel represents monitor parameters.
type MonitorParamsModel struct {
	Metrics          types.List    `tfsdk:"metrics"`
	Query            types.String  `tfsdk:"query"`
	Column           types.String  `tfsdk:"column"`
	MinAllowedValue  types.Float64 `tfsdk:"min_allowed_value"`
	MaxAllowedValue  types.Float64 `tfsdk:"max_allowed_value"`
	GroupingInterval types.Float64 `tfsdk:"grouping_interval"`
	CheckNumPoint    types.Int64   `tfsdk:"check_num_point"`
	NullsMode        types.String  `tfsdk:"nulls_mode"`
	TimeOffset       types.Float64 `tfsdk:"time_offset"`
}

// MetricDefinitionModel represents a metric definition.
type MetricDefinitionModel struct {
	Name  types.String `tfsdk:"name"`
	Alias types.String `tfsdk:"alias"`
}

// planToMonitorInput converts a Terraform plan to an API MonitorInput.
func planToMonitorInput(ctx context.Context, plan MonitorResourceModel, diags *diag.Diagnostics) generated.MonitorInput {
	input := generated.MonitorInput{
		Name: plan.Name.ValueString(),
		Type: generated.MonitorType(plan.Type.ValueString()),
	}

	// Set optional notification fields
	if !plan.NotifyEveryoneByEmail.IsNull() {
		notifyEmail := plan.NotifyEveryoneByEmail.ValueBool()
		input.NotifyEveryoneByEmail = &notifyEmail
	}

	// Convert team IDs
	if !plan.TeamIDs.IsNull() && !plan.TeamIDs.IsUnknown() {
		var teamIDs []int64
		diags.Append(plan.TeamIDs.ElementsAs(ctx, &teamIDs, false)...)
		if !diags.HasError() {
			input.TeamIds = &teamIDs
		}
	}

	// Convert channel IDs
	if !plan.ChannelIDs.IsNull() && !plan.ChannelIDs.IsUnknown() {
		var channelIDs []int64
		diags.Append(plan.ChannelIDs.ElementsAs(ctx, &channelIDs, false)...)
		if !diags.HasError() {
			input.ChannelIds = &channelIDs
		}
	}

	// Convert repeat interval
	if !plan.RepeatInterval.IsNull() && !plan.RepeatInterval.IsUnknown() {
		var repeatInterval RepeatIntervalModel
		diags.Append(plan.RepeatInterval.As(ctx, &repeatInterval, basetypes.ObjectAsOptions{})...)
		if !diags.HasError() {
			ri := generated.RepeatInterval{}
			if !repeatInterval.Strategy.IsNull() {
				strategy := generated.RepeatIntervalStrategy(repeatInterval.Strategy.ValueString())
				ri.Strategy = &strategy
			}
			if !repeatInterval.Interval.IsNull() {
				interval := repeatInterval.Interval.ValueInt64()
				ri.Interval = &interval
			}
			input.RepeatInterval = &ri
		}
	}

	// Convert params based on monitor type
	if !plan.Params.IsNull() && !plan.Params.IsUnknown() {
		var params MonitorParamsModel
		diags.Append(plan.Params.As(ctx, &params, basetypes.ObjectAsOptions{})...)
		if !diags.HasError() {
			if plan.Type.ValueString() == "metric" {
				convertToMetricParams(ctx, params, &input.Params, diags)
			} else if plan.Type.ValueString() == "error" {
				convertToErrorParams(ctx, params, &input.Params, diags)
			}
		}
	}

	return input
}

// convertMetricsToAPI converts Terraform metric definitions to API format.
func convertMetricsToAPI(ctx context.Context, params MonitorParamsModel, diags *diag.Diagnostics) []generated.MetricDefinition {
	if params.Metrics.IsNull() || params.Metrics.IsUnknown() {
		return nil
	}

	var metrics []MetricDefinitionModel
	diags.Append(params.Metrics.ElementsAs(ctx, &metrics, false)...)
	if diags.HasError() {
		return nil
	}

	apiMetrics := make([]generated.MetricDefinition, len(metrics))
	for i, m := range metrics {
		apiMetrics[i] = generated.MetricDefinition{
			Name: m.Name.ValueString(),
		}
		if !m.Alias.IsNull() {
			alias := m.Alias.ValueString()
			apiMetrics[i].Alias = &alias
		}
	}
	return apiMetrics
}

// convertToMetricParams converts params to MetricMonitorParams.
func convertToMetricParams(ctx context.Context, params MonitorParamsModel, result *generated.MonitorInput_Params, diags *diag.Diagnostics) {
	metricParams := generated.MetricMonitorParams{}

	// Convert metrics
	if apiMetrics := convertMetricsToAPI(ctx, params, diags); apiMetrics != nil {
		metricParams.Metrics = apiMetrics
	}

	if !params.Query.IsNull() {
		metricParams.Query = params.Query.ValueString()
	}

	if !params.Column.IsNull() {
		metricParams.Column = params.Column.ValueString()
	}

	if !params.MinAllowedValue.IsNull() {
		val := params.MinAllowedValue.ValueFloat64()
		metricParams.MinAllowedValue = &val
	}

	if !params.MaxAllowedValue.IsNull() {
		val := params.MaxAllowedValue.ValueFloat64()
		metricParams.MaxAllowedValue = &val
	}

	if !params.GroupingInterval.IsNull() {
		val := params.GroupingInterval.ValueFloat64()
		metricParams.GroupingInterval = &val
	}

	if !params.CheckNumPoint.IsNull() {
		val := int(params.CheckNumPoint.ValueInt64())
		metricParams.CheckNumPoint = &val
	}

	if !params.NullsMode.IsNull() {
		mode := generated.MetricMonitorParamsNullsMode(params.NullsMode.ValueString())
		metricParams.NullsMode = &mode
	}

	if !params.TimeOffset.IsNull() {
		val := params.TimeOffset.ValueFloat64()
		metricParams.TimeOffset = &val
	}

	if err := result.FromMetricMonitorParams(metricParams); err != nil {
		diags.AddError("Failed to convert metric params", err.Error())
	}
}

// convertToErrorParams converts params to ErrorMonitorParams.
func convertToErrorParams(ctx context.Context, params MonitorParamsModel, result *generated.MonitorInput_Params, diags *diag.Diagnostics) {
	errorParams := generated.ErrorMonitorParams{}

	// Convert metrics
	if apiMetrics := convertMetricsToAPI(ctx, params, diags); apiMetrics != nil {
		errorParams.Metrics = apiMetrics
	}

	if !params.Query.IsNull() {
		query := params.Query.ValueString()
		errorParams.Query = &query
	}

	if err := result.FromErrorMonitorParams(errorParams); err != nil {
		diags.AddError("Failed to convert error params", err.Error())
	}
}

// monitorToState converts an API Monitor to Terraform state.
func monitorToState(ctx context.Context, monitor *generated.Monitor, state *MonitorResourceModel, diags *diag.Diagnostics) {
	state.ID = types.StringValue(fmt.Sprintf("%d", monitor.Id))
	state.Name = types.StringValue(monitor.Name)
	state.Type = types.StringValue(string(monitor.Type))
	state.State = types.StringValue(string(monitor.State))

	// Convert optional fields
	if monitor.NotifyEveryoneByEmail != nil {
		state.NotifyEveryoneByEmail = types.BoolValue(*monitor.NotifyEveryoneByEmail)
	} else {
		state.NotifyEveryoneByEmail = types.BoolValue(false)
	}

	// Convert team IDs
	if monitor.TeamIds != nil && len(*monitor.TeamIds) > 0 {
		teamIDs := make([]attr.Value, len(*monitor.TeamIds))
		for i, id := range *monitor.TeamIds {
			teamIDs[i] = types.Int64Value(id)
		}
		state.TeamIDs = types.ListValueMust(types.Int64Type, teamIDs)
	} else {
		state.TeamIDs = types.ListValueMust(types.Int64Type, []attr.Value{})
	}

	// Convert channel IDs
	if monitor.ChannelIds != nil && len(*monitor.ChannelIds) > 0 {
		channelIDs := make([]attr.Value, len(*monitor.ChannelIds))
		for i, id := range *monitor.ChannelIds {
			channelIDs[i] = types.Int64Value(id)
		}
		state.ChannelIDs = types.ListValueMust(types.Int64Type, channelIDs)
	} else {
		state.ChannelIDs = types.ListValueMust(types.Int64Type, []attr.Value{})
	}

	// Convert repeat interval
	if monitor.RepeatInterval != nil {
		riAttrs := map[string]attr.Value{
			"strategy": types.StringNull(),
			"interval": types.Int64Null(),
		}
		if monitor.RepeatInterval.Strategy != nil {
			riAttrs["strategy"] = types.StringValue(string(*monitor.RepeatInterval.Strategy))
		}
		if monitor.RepeatInterval.Interval != nil {
			riAttrs["interval"] = types.Int64Value(*monitor.RepeatInterval.Interval)
		}
		state.RepeatInterval = types.ObjectValueMust(
			map[string]attr.Type{
				"strategy": types.StringType,
				"interval": types.Int64Type,
			},
			riAttrs,
		)
	}

	// Convert params - always present
	convertParamsToState(ctx, monitor, state, diags)

	// Convert timestamps (Unix milliseconds)
	if monitor.CreatedAt != nil {
		state.CreatedAt = types.StringValue(fmt.Sprintf("%.0f", *monitor.CreatedAt))
	}
	if monitor.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(fmt.Sprintf("%.0f", *monitor.UpdatedAt))
	}
}

// convertParamsToState converts API params to Terraform state.
func convertParamsToState(_ context.Context, monitor *generated.Monitor, state *MonitorResourceModel, diags *diag.Diagnostics) {
	paramsAttrs := map[string]attr.Value{
		"metrics":           types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType, "alias": types.StringType}}),
		"query":             types.StringNull(),
		"column":            types.StringNull(),
		"min_allowed_value": types.Float64Null(),
		"max_allowed_value": types.Float64Null(),
		"grouping_interval": types.Float64Null(),
		"check_num_point":   types.Int64Null(),
		"nulls_mode":        types.StringNull(),
		"time_offset":       types.Float64Null(),
	}

	if monitor.Type == generated.MonitorTypeMetric {
		metricParams, err := monitor.Params.AsMetricMonitorParams()
		if err == nil {
			convertMetricParamsToAttrs(metricParams, paramsAttrs)
		} else {
			diags.AddWarning("Failed to parse metric params", err.Error())
		}
	} else if monitor.Type == generated.MonitorTypeError {
		errorParams, err := monitor.Params.AsErrorMonitorParams()
		if err == nil {
			convertErrorParamsToAttrs(errorParams, paramsAttrs)
		} else {
			diags.AddWarning("Failed to parse error params", err.Error())
		}
	}

	state.Params = types.ObjectValueMust(
		map[string]attr.Type{
			"metrics":           types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType, "alias": types.StringType}}},
			"query":             types.StringType,
			"column":            types.StringType,
			"min_allowed_value": types.Float64Type,
			"max_allowed_value": types.Float64Type,
			"grouping_interval": types.Float64Type,
			"check_num_point":   types.Int64Type,
			"nulls_mode":        types.StringType,
			"time_offset":       types.Float64Type,
		},
		paramsAttrs,
	)
}

// convertMetricParamsToAttrs converts metric params to attribute map.
func convertMetricParamsToAttrs(params generated.MetricMonitorParams, attrs map[string]attr.Value) {
	metricAttrTypes := map[string]attr.Type{"name": types.StringType, "alias": types.StringType}

	if len(params.Metrics) > 0 {
		metrics := make([]attr.Value, len(params.Metrics))
		for i, m := range params.Metrics {
			metricAttrs := map[string]attr.Value{
				"name":  types.StringValue(m.Name),
				"alias": types.StringNull(),
			}
			if m.Alias != nil {
				metricAttrs["alias"] = types.StringValue(*m.Alias)
			}
			metrics[i] = types.ObjectValueMust(metricAttrTypes, metricAttrs)
		}
		attrs["metrics"] = types.ListValueMust(types.ObjectType{AttrTypes: metricAttrTypes}, metrics)
	}

	if params.Query != "" {
		attrs["query"] = types.StringValue(params.Query)
	}
	if params.Column != "" {
		attrs["column"] = types.StringValue(params.Column)
	}
	if params.MinAllowedValue != nil {
		attrs["min_allowed_value"] = types.Float64Value(*params.MinAllowedValue)
	}
	if params.MaxAllowedValue != nil {
		attrs["max_allowed_value"] = types.Float64Value(*params.MaxAllowedValue)
	}
	if params.GroupingInterval != nil {
		attrs["grouping_interval"] = types.Float64Value(*params.GroupingInterval)
	}
	if params.CheckNumPoint != nil {
		attrs["check_num_point"] = types.Int64Value(int64(*params.CheckNumPoint))
	}
	if params.NullsMode != nil {
		attrs["nulls_mode"] = types.StringValue(string(*params.NullsMode))
	}
	if params.TimeOffset != nil {
		attrs["time_offset"] = types.Float64Value(*params.TimeOffset)
	}
}

// convertErrorParamsToAttrs converts error params to attribute map.
func convertErrorParamsToAttrs(params generated.ErrorMonitorParams, attrs map[string]attr.Value) {
	metricAttrTypes := map[string]attr.Type{"name": types.StringType, "alias": types.StringType}

	if len(params.Metrics) > 0 {
		metrics := make([]attr.Value, len(params.Metrics))
		for i, m := range params.Metrics {
			metricAttrs := map[string]attr.Value{
				"name":  types.StringValue(m.Name),
				"alias": types.StringNull(),
			}
			if m.Alias != nil {
				metricAttrs["alias"] = types.StringValue(*m.Alias)
			}
			metrics[i] = types.ObjectValueMust(metricAttrTypes, metricAttrs)
		}
		attrs["metrics"] = types.ListValueMust(types.ObjectType{AttrTypes: metricAttrTypes}, metrics)
	}

	if params.Query != nil {
		attrs["query"] = types.StringValue(*params.Query)
	}
}
