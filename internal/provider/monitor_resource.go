package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/riccap/tofu-uptrace-provider/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &MonitorResource{}
	_ resource.ResourceWithConfigure   = &MonitorResource{}
	_ resource.ResourceWithImportState = &MonitorResource{}
)

// NewMonitorResource is a helper function to create the resource.
func NewMonitorResource() resource.Resource {
	return &MonitorResource{}
}

// MonitorResource is the resource implementation.
type MonitorResource struct {
	client *client.Client
}

// MonitorResourceModel describes the resource data model.
type MonitorResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Type                  types.String `tfsdk:"type"`
	State                 types.String `tfsdk:"state"`
	NotifyEveryoneByEmail types.Bool   `tfsdk:"notify_everyone_by_email"`
	TeamIDs               types.List   `tfsdk:"team_ids"`
	ChannelIDs            types.List   `tfsdk:"channel_ids"`
	RepeatInterval        types.Object `tfsdk:"repeat_interval"`
	Params                types.Object `tfsdk:"params"`
	CreatedAt             types.String `tfsdk:"created_at"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *MonitorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

// Schema defines the schema for the resource.
func (r *MonitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Uptrace monitor for metrics or errors.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Monitor identifier.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Monitor name.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Monitor type. Must be 'metric' or 'error'.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("metric", "error"),
				},
			},
			"state": schema.StringAttribute{
				Description: "Current monitor state (open, firing, paused).",
				Computed:    true,
			},
			"notify_everyone_by_email": schema.BoolAttribute{
				Description: "Whether to notify all project members by email.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"team_ids": schema.ListAttribute{
				Description: "List of team IDs to notify.",
				ElementType: types.Int64Type,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.Int64Type, []attr.Value{})),
			},
			"channel_ids": schema.ListAttribute{
				Description: "List of notification channel IDs.",
				ElementType: types.Int64Type,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.Int64Type, []attr.Value{})),
			},
			"repeat_interval": schema.SingleNestedAttribute{
				Description: "Repeat interval configuration.",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"strategy": schema.StringAttribute{
						Description: "Repeat interval strategy. Must be 'default' or 'custom'.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("default"),
						Validators: []validator.String{
							stringvalidator.OneOf("default", "custom"),
						},
					},
					"interval": schema.Int64Attribute{
						Description: "Custom interval in seconds (only for custom strategy, minimum 60).",
						Optional:    true,
						Validators: []validator.Int64{
							int64validator.AtLeast(60),
						},
					},
				},
			},
			"params": schema.SingleNestedAttribute{
				Description: "Monitor parameters (metric or error specific).",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					// Metric monitor params
					"metrics": schema.ListNestedAttribute{
						Description: "List of metrics to monitor.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "Metric name.",
									Required:    true,
								},
								"alias": schema.StringAttribute{
									Description: "Optional alias for the metric.",
									Optional:    true,
								},
							},
						},
					},
					"query": schema.StringAttribute{
						Description: "UQL query for metric evaluation or error filtering.",
						Optional:    true,
						Computed:    true,
					},
					"column": schema.StringAttribute{
						Description: "Column name to evaluate (metric monitors only).",
						Optional:    true,
						Computed:    true,
					},
					"min_allowed_value": schema.Float64Attribute{
						Description: "Minimum allowed value for the metric.",
						Optional:    true,
					},
					"max_allowed_value": schema.Float64Attribute{
						Description: "Maximum allowed value for the metric.",
						Optional:    true,
					},
					"grouping_interval": schema.Float64Attribute{
						Description: "Grouping interval in milliseconds.",
						Optional:    true,
						Computed:    true,
					},
					"check_num_point": schema.Int64Attribute{
						Description: "Number of consecutive points that must breach threshold.",
						Optional:    true,
						Computed:    true,
					},
					"nulls_mode": schema.StringAttribute{
						Description: "How to handle null values: 'allow', 'forbid', or 'convert'.",
						Optional:    true,
						Computed:    true,
					},
					"time_offset": schema.Float64Attribute{
						Description: "Time offset in milliseconds.",
						Optional:    true,
						Computed:    true,
					},
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Monitor creation timestamp.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Monitor last update timestamp.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *MonitorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *MonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MonitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating monitor", map[string]any{"name": plan.Name.ValueString()})

	// Convert plan to API input
	input := planToMonitorInput(ctx, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create monitor via API
	monitor, err := r.client.CreateMonitor(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Monitor",
			fmt.Sprintf("Could not create monitor: %s", err.Error()),
		)
		return
	}

	// Convert API response to state
	monitorToState(ctx, monitor, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Created monitor", map[string]any{"id": plan.ID.ValueString()})
}

// Read refreshes the Terraform state with the latest data.
//
//nolint:gocritic // Request type defined by Terraform Plugin Framework interface
func (r *MonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Reading monitor", map[string]any{"id": state.ID.ValueString()})

	// Get monitor from API
	monitor, err := r.client.GetMonitor(ctx, state.ID.ValueString())
	if err != nil {
		// If the monitor doesn't exist (404), remove from state
		if isNotFoundError(err) {
			tflog.Warn(ctx, "Monitor not found, removing from state", map[string]any{"id": state.ID.ValueString()})
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Monitor",
			fmt.Sprintf("Could not read monitor ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to state
	monitorToState(ctx, monitor, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
//
//nolint:gocritic // Request type defined by Terraform Plugin Framework interface
func (r *MonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MonitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Updating monitor", map[string]any{"id": plan.ID.ValueString()})

	// Convert plan to API input
	input := planToMonitorInput(ctx, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update monitor via API
	monitor, err := r.client.UpdateMonitor(ctx, plan.ID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Monitor",
			fmt.Sprintf("Could not update monitor ID %s: %s", plan.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to state
	monitorToState(ctx, monitor, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Updated monitor", map[string]any{"id": plan.ID.ValueString()})
}

// Delete deletes the resource and removes the Terraform state on success.
//
//nolint:gocritic // Request type defined by Terraform Plugin Framework interface
func (r *MonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MonitorResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Deleting monitor", map[string]any{"id": state.ID.ValueString()})

	// Delete monitor via API
	err := r.client.DeleteMonitor(ctx, state.ID.ValueString())
	if err != nil {
		// If the monitor doesn't exist (404), treat as already deleted
		if isNotFoundError(err) {
			tflog.Warn(ctx, "Monitor already deleted", map[string]any{"id": state.ID.ValueString()})
			return
		}

		resp.Diagnostics.AddError(
			"Error Deleting Monitor",
			fmt.Sprintf("Could not delete monitor ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Info(ctx, "Deleted monitor", map[string]any{"id": state.ID.ValueString()})
}

// ImportState imports an existing resource into Terraform.
func (r *MonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
