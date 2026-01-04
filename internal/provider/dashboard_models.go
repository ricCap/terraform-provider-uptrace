package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/riccap/tofu-uptrace-provider/internal/client/generated"
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

// isNotFoundError checks if an error is a "not found" (404) error.
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "not found")
}
