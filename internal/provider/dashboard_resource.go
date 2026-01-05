package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/riccap/terraform-provider-uptrace/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &DashboardResource{}
	_ resource.ResourceWithConfigure   = &DashboardResource{}
	_ resource.ResourceWithImportState = &DashboardResource{}
)

// NewDashboardResource is a helper function to create the resource.
func NewDashboardResource() resource.Resource {
	return &DashboardResource{}
}

// DashboardResource is the resource implementation.
type DashboardResource struct {
	client *client.Client
}

// DashboardResourceModel describes the resource data model.
type DashboardResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	YAML      types.String `tfsdk:"yaml"`
	Pinned    types.Bool   `tfsdk:"pinned"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *DashboardResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard"
}

// Schema defines the schema for the resource.
func (r *DashboardResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Uptrace dashboard using YAML definition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Dashboard identifier.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Dashboard name (extracted from YAML or API response).",
				Computed:    true,
			},
			"yaml": schema.StringAttribute{
				Description: "Dashboard YAML definition. Supports all dashboard features including grid layout, charts, tables, heatmaps, and gauges.",
				Required:    true,
			},
			"pinned": schema.BoolAttribute{
				Description: "Whether the dashboard is pinned.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Dashboard creation timestamp.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Dashboard last update timestamp.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *DashboardResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	uptraceClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = uptraceClient
}

// Create creates the resource and sets the initial Terraform state.
//
//nolint:gocritic // Request type defined by Terraform Plugin Framework interface
func (r *DashboardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DashboardResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating dashboard", map[string]any{"yaml_length": len(plan.YAML.ValueString())})

	// Create dashboard via API
	dashboard, err := r.client.CreateDashboardFromYAML(ctx, plan.YAML.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Dashboard",
			fmt.Sprintf("Could not create dashboard: %s", err.Error()),
		)
		return
	}

	// Convert API response to state
	dashboardToState(ctx, dashboard, plan.YAML.ValueString(), &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Successfully created dashboard", map[string]any{"id": plan.ID.ValueString()})
}

// Read refreshes the Terraform state with the latest data.
//
//nolint:gocritic // Request type defined by Terraform Plugin Framework interface
func (r *DashboardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DashboardResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Reading dashboard", map[string]any{"id": state.ID.ValueString()})

	// Parse dashboard ID
	dashboardID, ok := parseDashboardID(state.ID.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	// Get dashboard from API
	dashboard, err := r.client.GetDashboard(ctx, dashboardID)
	if err != nil {
		if isNotFoundError(err) {
			tflog.Warn(ctx, "Dashboard not found, removing from state", map[string]any{"id": state.ID.ValueString()})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Dashboard",
			fmt.Sprintf("Could not read dashboard ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Only fetch YAML from API during import (when state.YAML is empty)
	// Otherwise, preserve the user's original YAML to avoid drift from API-added defaults
	var yamlContent string
	if state.YAML.ValueString() == "" {
		tflog.Debug(ctx, "Fetching dashboard YAML (import scenario)", map[string]any{"id": state.ID.ValueString()})
		fetchedYAML, err := r.client.GetDashboardYAML(ctx, dashboardID)
		if err != nil {
			tflog.Error(ctx, "Could not fetch dashboard YAML during import", map[string]any{
				"id":    state.ID.ValueString(),
				"error": err.Error(),
			})
			resp.Diagnostics.AddError(
				"Error Fetching Dashboard YAML",
				fmt.Sprintf("Could not fetch YAML for dashboard ID %s during import: %s", state.ID.ValueString(), err.Error()),
			)
			return
		}
		yamlContent = fetchedYAML
	} else {
		// Preserve user's original YAML from state
		yamlContent = state.YAML.ValueString()
	}

	// Convert API response to state
	dashboardToState(ctx, dashboard, yamlContent, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Successfully read dashboard", map[string]any{"id": state.ID.ValueString()})
}

// Update updates the resource and sets the updated Terraform state on success.
//
//nolint:gocritic // Request type defined by Terraform Plugin Framework interface
func (r *DashboardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DashboardResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Updating dashboard", map[string]any{"id": plan.ID.ValueString()})

	// Parse dashboard ID
	dashboardID, ok := parseDashboardID(plan.ID.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	// Update dashboard via API
	dashboard, err := r.client.UpdateDashboardFromYAML(ctx, dashboardID, plan.YAML.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Dashboard",
			fmt.Sprintf("Could not update dashboard ID %s: %s", plan.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to state
	dashboardToState(ctx, dashboard, plan.YAML.ValueString(), &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Successfully updated dashboard", map[string]any{"id": plan.ID.ValueString()})
}

// Delete deletes the resource and removes the Terraform state on success.
//
//nolint:gocritic // Request type defined by Terraform Plugin Framework interface
func (r *DashboardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DashboardResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Deleting dashboard", map[string]any{"id": state.ID.ValueString()})

	// Parse dashboard ID
	dashboardID, ok := parseDashboardID(state.ID.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	// Delete dashboard via API
	err := r.client.DeleteDashboard(ctx, dashboardID)
	if err != nil {
		if isNotFoundError(err) {
			tflog.Warn(ctx, "Dashboard already deleted", map[string]any{"id": state.ID.ValueString()})
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Dashboard",
			fmt.Sprintf("Could not delete dashboard ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Info(ctx, "Successfully deleted dashboard", map[string]any{"id": state.ID.ValueString()})
}

// ImportState imports the resource state.
func (r *DashboardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
