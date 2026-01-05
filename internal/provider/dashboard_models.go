package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/riccap/terraform-provider-uptrace/internal/client/generated"
)

// dashboardToState converts an API Dashboard to Terraform state.
func dashboardToState(
	_ context.Context,
	dashboard *generated.Dashboard,
	yamlContent string,
	state *DashboardResourceModel,
	_ *diag.Diagnostics,
) {
	state.ID = types.StringValue(fmt.Sprintf("%d", dashboard.Id))
	state.Name = types.StringValue(dashboard.Name)
	state.YAML = types.StringValue(yamlContent)

	// Set pinned field (default to false if not provided)
	if dashboard.Pinned != nil {
		state.Pinned = types.BoolValue(*dashboard.Pinned)
	} else {
		state.Pinned = types.BoolValue(false)
	}

	// Set timestamps if available
	if dashboard.CreatedAt != nil {
		state.CreatedAt = types.StringValue(fmt.Sprintf("%f", *dashboard.CreatedAt))
	} else {
		state.CreatedAt = types.StringNull()
	}

	if dashboard.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(fmt.Sprintf("%f", *dashboard.UpdatedAt))
	} else {
		state.UpdatedAt = types.StringNull()
	}
}

// parseDashboardID parses a dashboard ID string and adds diagnostics on error.
// Returns the parsed ID and a boolean indicating success.
func parseDashboardID(idStr string, diags *diag.Diagnostics) (int64, bool) {
	dashboardID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		diags.AddError(
			"Invalid Dashboard ID",
			fmt.Sprintf("Could not parse dashboard ID %s: %s", idStr, err.Error()),
		)
		return 0, false
	}
	return dashboardID, true
}
