package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/riccap/terraform-provider-uptrace/internal/client/generated"
)

// planToChannelInput converts a Terraform plan to an API NotificationChannelInput.
//
//nolint:gocritic // Plan passed by value to keep function signatures consistent
func planToChannelInput(ctx context.Context, plan NotificationChannelResourceModel, diags *diag.Diagnostics) generated.NotificationChannelInput {
	input := generated.NotificationChannelInput{
		Name: plan.Name.ValueString(),
		Type: generated.NotificationChannelInputType(plan.Type.ValueString()),
	}

	// Set optional condition
	if !plan.Condition.IsNull() && !plan.Condition.IsUnknown() {
		condition := plan.Condition.ValueString()
		input.Condition = &condition
	}

	// Convert priority list
	if !plan.Priority.IsNull() && !plan.Priority.IsUnknown() {
		var priorityList []string
		convertDiags := plan.Priority.ElementsAs(ctx, &priorityList, false)
		tflog.Debug(ctx, "Converting priority list", map[string]any{
			"is_null":     plan.Priority.IsNull(),
			"is_unknown":  plan.Priority.IsUnknown(),
			"has_error":   convertDiags.HasError(),
			"list_length": len(priorityList),
			"list":        priorityList,
		})
		diags.Append(convertDiags...)
		if !convertDiags.HasError() && len(priorityList) > 0 {
			// Convert []string to []generated.NotificationChannelInputPriority
			typedPriorities := make([]generated.NotificationChannelInputPriority, len(priorityList))
			for i, p := range priorityList {
				typedPriorities[i] = generated.NotificationChannelInputPriority(p)
			}
			input.Priority = &typedPriorities
			tflog.Debug(ctx, "Setting priority on channel input", map[string]any{
				"priority": priorityList,
			})
		} else {
			tflog.Warn(ctx, "Priority list not set", map[string]any{
				"has_error":   convertDiags.HasError(),
				"list_length": len(priorityList),
			})
		}
	} else {
		tflog.Debug(ctx, "Priority is null or unknown", map[string]any{
			"is_null":    plan.Priority.IsNull(),
			"is_unknown": plan.Priority.IsUnknown(),
		})
	}

	// Convert params map to interface{} for JSON marshaling
	if !plan.Params.IsNull() && !plan.Params.IsUnknown() {
		paramsMap := make(map[string]string)
		diags.Append(plan.Params.ElementsAs(ctx, &paramsMap, false)...)
		if diags.HasError() {
			return input
		}

		// Convert string map to interface{} map for API
		paramsInterface := make(map[string]interface{})
		for k, v := range paramsMap {
			// Try to parse as JSON for nested objects/numbers
			var jsonValue interface{}
			if err := json.Unmarshal([]byte(v), &jsonValue); err == nil {
				paramsInterface[k] = jsonValue
			} else {
				// Keep as string if not valid JSON
				paramsInterface[k] = v
			}
		}
		input.Params = paramsInterface
	}

	// Debug: log the marshaled input
	if jsonBytes, err := json.Marshal(input); err == nil {
		tflog.Debug(ctx, "Marshaled notification channel input", map[string]any{
			"json": string(jsonBytes),
		})
	}

	return input
}

// channelToState converts an API NotificationChannel to Terraform state.
func channelToState(ctx context.Context, channel *generated.NotificationChannel, state *NotificationChannelResourceModel, diags *diag.Diagnostics) {
	state.ID = types.StringValue(fmt.Sprintf("%d", channel.Id))
	state.Name = types.StringValue(channel.Name)
	state.Type = types.StringValue(string(channel.Type))

	// Set optional condition
	if channel.Condition != nil && *channel.Condition != "" {
		state.Condition = types.StringValue(*channel.Condition)
	} else {
		state.Condition = types.StringNull()
	}

	// Convert priority from API response
	if channel.Priority != nil && len(*channel.Priority) > 0 {
		priorityValues := make([]attr.Value, len(*channel.Priority))
		for i, p := range *channel.Priority {
			// Convert generated.NotificationChannelPriority to string
			priorityValues[i] = types.StringValue(string(p))
		}
		state.Priority = types.ListValueMust(types.StringType, priorityValues)
	} else {
		state.Priority = types.ListNull(types.StringType)
	}

	// Set status
	state.Status = types.StringValue(channel.Status)

	// Convert params from interface{} to string map
	if channel.Params != nil {
		paramsMap := make(map[string]string)
		for k, v := range channel.Params {
			// Marshal to JSON string for complex types
			switch val := v.(type) {
			case string:
				paramsMap[k] = val
			case float64, int, int64, bool:
				paramsMap[k] = fmt.Sprintf("%v", val)
			default:
				// For complex types, marshal to JSON
				jsonBytes, err := json.Marshal(val)
				if err != nil {
					diags.AddWarning(
						"Failed to marshal param",
						fmt.Sprintf("Could not marshal param %s: %s", k, err.Error()),
					)
					continue
				}
				paramsMap[k] = string(jsonBytes)
			}
		}

		var diag diag.Diagnostics
		state.Params, diag = types.MapValueFrom(ctx, types.StringType, paramsMap)
		diags.Append(diag...)
	} else {
		state.Params = types.MapNull(types.StringType)
	}

	// Note: Uptrace API doesn't return created_at/updated_at for notification channels
	// Keep these as null in state
	state.CreatedAt = types.StringNull()
	state.UpdatedAt = types.StringNull()
}
