package provider

import (
	"context"
	"fmt"
	"strconv"

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
	_ resource.Resource                = &NotificationChannelResource{}
	_ resource.ResourceWithConfigure   = &NotificationChannelResource{}
	_ resource.ResourceWithImportState = &NotificationChannelResource{}
)

// NewNotificationChannelResource is a helper function to create the resource.
func NewNotificationChannelResource() resource.Resource {
	return &NotificationChannelResource{}
}

// NotificationChannelResource is the resource implementation.
type NotificationChannelResource struct {
	client *client.Client
}

// NotificationChannelResourceModel describes the resource data model.
type NotificationChannelResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Type      types.String `tfsdk:"type"`
	Condition types.String `tfsdk:"condition"`
	Params    types.Map    `tfsdk:"params"`
	Status    types.String `tfsdk:"status"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *NotificationChannelResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_channel"
}

// Schema defines the schema for the resource.
func (r *NotificationChannelResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Uptrace notification channel for alert notifications.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Notification channel identifier.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Channel name.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Channel type. Supported values: slack, webhook, telegram, mattermost.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"condition": schema.StringAttribute{
				Description: "Optional condition expression to filter notifications.",
				Optional:    true,
			},
			"params": schema.MapAttribute{
				Description: "Channel-specific configuration parameters. Structure varies by channel type.",
				ElementType: types.StringType,
				Required:    true,
				Sensitive:   true,
			},
			"status": schema.StringAttribute{
				Description: "Channel delivery status (computed).",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Channel creation timestamp (computed).",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Channel last update timestamp (computed).",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *NotificationChannelResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *NotificationChannelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NotificationChannelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating notification channel", map[string]any{"name": plan.Name.ValueString()})

	// Convert plan to API input
	input := planToChannelInput(ctx, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create channel via API
	channel, err := r.client.CreateNotificationChannel(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Notification Channel",
			fmt.Sprintf("Could not create notification channel: %s", err.Error()),
		)
		return
	}

	// Preserve params from plan (sensitive field that API may not return)
	plannedParams := plan.Params

	// Convert API response to state
	channelToState(ctx, channel, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Restore params from plan (sensitive values not returned by API)
	plan.Params = plannedParams

	// Save state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Successfully created notification channel", map[string]any{"id": plan.ID.ValueString()})
}

// Read refreshes the Terraform state with the latest data.
//
//nolint:gocritic // Request type defined by Terraform Plugin Framework interface
func (r *NotificationChannelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NotificationChannelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Reading notification channel", map[string]any{"id": state.ID.ValueString()})

	// Parse channel ID
	channelID, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Channel ID",
			fmt.Sprintf("Could not parse channel ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Get channel from API
	channel, err := r.client.GetNotificationChannel(ctx, channelID)
	if err != nil {
		if isNotFoundError(err) {
			tflog.Warn(ctx, "Notification channel not found, removing from state", map[string]any{"id": state.ID.ValueString()})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Notification Channel",
			fmt.Sprintf("Could not read notification channel ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Preserve params from prior state (sensitive field that API may not return)
	priorParams := state.Params

	// Convert API response to state
	channelToState(ctx, channel, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Restore params from prior state (sensitive values not returned by API)
	state.Params = priorParams

	// Save state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "Successfully read notification channel", map[string]any{"id": state.ID.ValueString()})
}

// Update updates the resource and sets the updated Terraform state on success.
//
//nolint:gocritic // Request type defined by Terraform Plugin Framework interface
func (r *NotificationChannelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NotificationChannelResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Updating notification channel", map[string]any{"id": plan.ID.ValueString()})

	// Parse channel ID
	channelID, err := strconv.ParseInt(plan.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Channel ID",
			fmt.Sprintf("Could not parse channel ID %s: %s", plan.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert plan to API input
	input := planToChannelInput(ctx, plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update channel via API
	channel, err := r.client.UpdateNotificationChannel(ctx, channelID, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Notification Channel",
			fmt.Sprintf("Could not update notification channel ID %s: %s", plan.ID.ValueString(), err.Error()),
		)
		return
	}

	// Preserve params from plan (sensitive field that API may not return)
	plannedParams := plan.Params

	// Convert API response to state
	channelToState(ctx, channel, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Restore params from plan (sensitive values not returned by API)
	plan.Params = plannedParams

	// Save state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Successfully updated notification channel", map[string]any{"id": plan.ID.ValueString()})
}

// Delete deletes the resource and removes the Terraform state on success.
//
//nolint:gocritic // Request type defined by Terraform Plugin Framework interface
func (r *NotificationChannelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NotificationChannelResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Deleting notification channel", map[string]any{"id": state.ID.ValueString()})

	// Parse channel ID
	channelID, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Channel ID",
			fmt.Sprintf("Could not parse channel ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	// Delete channel via API
	err = r.client.DeleteNotificationChannel(ctx, channelID)
	if err != nil {
		if isNotFoundError(err) {
			tflog.Warn(ctx, "Notification channel already deleted", map[string]any{"id": state.ID.ValueString()})
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Notification Channel",
			fmt.Sprintf("Could not delete notification channel ID %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	tflog.Info(ctx, "Successfully deleted notification channel", map[string]any{"id": state.ID.ValueString()})
}

// ImportState imports the resource state.
func (r *NotificationChannelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
