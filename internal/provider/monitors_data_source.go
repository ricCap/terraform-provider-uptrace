package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/riccap/terraform-provider-uptrace/internal/client"
	"github.com/riccap/terraform-provider-uptrace/internal/client/generated"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &MonitorsDataSource{}
	_ datasource.DataSourceWithConfigure = &MonitorsDataSource{}
)

// NewMonitorsDataSource is a helper function to create the data source.
func NewMonitorsDataSource() datasource.DataSource {
	return &MonitorsDataSource{}
}

// MonitorsDataSource is the data source implementation.
type MonitorsDataSource struct {
	client *client.Client
}

// MonitorsDataSourceModel describes the data source data model.
type MonitorsDataSourceModel struct {
	Type     types.String `tfsdk:"type"`
	State    types.String `tfsdk:"state"`
	Name     types.String `tfsdk:"name"`
	Monitors types.List   `tfsdk:"monitors"`
}

// MonitorModel describes an individual monitor in the list.
type MonitorModel struct {
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

// Metadata returns the data source type name.
func (d *MonitorsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitors"
}

// Schema defines the schema for the data source.
func (d *MonitorsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of Uptrace monitors with optional filtering.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "Filter monitors by type (metric or error).",
				Optional:    true,
			},
			"state": schema.StringAttribute{
				Description: "Filter monitors by state (open, firing, paused).",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Filter monitors by name (case-insensitive substring match).",
				Optional:    true,
			},
			//nolint:dupl // Schema duplication with monitor_data_source acceptable - different data sources
			"monitors": schema.ListNestedAttribute{
				Description: "List of monitors matching the filter criteria.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Monitor identifier.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Monitor name.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Monitor type (metric or error).",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "Current monitor state (open, firing, paused).",
							Computed:    true,
						},
						"notify_everyone_by_email": schema.BoolAttribute{
							Description: "Whether to notify all project members by email.",
							Computed:    true,
						},
						"team_ids": schema.ListAttribute{
							Description: "List of team IDs to notify.",
							ElementType: types.Int64Type,
							Computed:    true,
						},
						"channel_ids": schema.ListAttribute{
							Description: "List of notification channel IDs.",
							ElementType: types.Int64Type,
							Computed:    true,
						},
						"repeat_interval": schema.SingleNestedAttribute{
							Description: "Repeat interval configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"strategy": schema.StringAttribute{
									Description: "Repeat interval strategy (default or custom).",
									Computed:    true,
								},
								"interval": schema.Int64Attribute{
									Description: "Custom interval in seconds.",
									Computed:    true,
								},
							},
						},
						"params": schema.SingleNestedAttribute{
							Description: "Monitor parameters (metric or error specific).",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"metrics": schema.ListNestedAttribute{
									Description: "List of metrics to monitor.",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"name": schema.StringAttribute{
												Description: "Metric name.",
												Computed:    true,
											},
											"alias": schema.StringAttribute{
												Description: "Optional alias for the metric.",
												Computed:    true,
											},
										},
									},
								},
								"query": schema.StringAttribute{
									Description: "UQL query for metric evaluation or error filtering.",
									Computed:    true,
								},
								"column": schema.StringAttribute{
									Description: "Column name to evaluate (metric monitors only).",
									Computed:    true,
								},
								"min_allowed_value": schema.Float64Attribute{
									Description: "Minimum allowed value for the metric.",
									Computed:    true,
								},
								"max_allowed_value": schema.Float64Attribute{
									Description: "Maximum allowed value for the metric.",
									Computed:    true,
								},
								"grouping_interval": schema.Float64Attribute{
									Description: "Grouping interval in milliseconds.",
									Computed:    true,
								},
								"check_num_point": schema.Int64Attribute{
									Description: "Number of consecutive points that must breach threshold.",
									Computed:    true,
								},
								"nulls_mode": schema.StringAttribute{
									Description: "How to handle null values: allow, forbid, or convert.",
									Computed:    true,
								},
								"time_offset": schema.Float64Attribute{
									Description: "Time offset in milliseconds.",
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
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *MonitorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	uptraceClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = uptraceClient
}

// Read refreshes the Terraform state with the latest data.
//
//nolint:gocritic // Request type defined by Terraform Plugin Framework interface
func (d *MonitorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config MonitorsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Reading monitors data source", map[string]any{
		"type_filter":  config.Type.ValueString(),
		"state_filter": config.State.ValueString(),
		"name_filter":  config.Name.ValueString(),
	})

	// Fetch all monitors from API
	monitors, err := d.client.ListMonitors(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Monitors",
			fmt.Sprintf("Could not list monitors: %s", err.Error()),
		)
		return
	}

	// Apply filters
	filtered := filterMonitors(monitors, config.Type, config.State, config.Name)

	tflog.Debug(ctx, "Filtered monitors", map[string]any{
		"total_count":    len(monitors),
		"filtered_count": len(filtered),
	})

	// Convert to Terraform state
	monitorModels := make([]MonitorModel, 0, len(filtered))
	for i := range filtered {
		model := convertMonitorToModel(ctx, &filtered[i], &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		monitorModels = append(monitorModels, model)
	}

	// Convert to types.List
	monitorList, diags := convertMonitorModelsToList(ctx, monitorModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Monitors = monitorList

	// Save state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	tflog.Info(ctx, "Successfully read monitors data source", map[string]any{"count": len(monitorModels)})
}

// filterMonitors applies filters to the monitor list.
func filterMonitors(monitors []generated.Monitor, typeFilter, stateFilter, nameFilter types.String) []generated.Monitor {
	var filtered []generated.Monitor

	//nolint:gocritic // Large struct copy acceptable for filtering logic
	for _, monitor := range monitors {
		// Skip if type filter doesn't match
		if !typeFilter.IsNull() && !typeFilter.IsUnknown() {
			if string(monitor.Type) != typeFilter.ValueString() {
				continue
			}
		}

		// Skip if state filter doesn't match
		if !stateFilter.IsNull() && !stateFilter.IsUnknown() {
			if string(monitor.State) != stateFilter.ValueString() {
				continue
			}
		}

		// Skip if name filter doesn't match
		if !nameFilter.IsNull() && !nameFilter.IsUnknown() {
			if !strings.Contains(
				strings.ToLower(monitor.Name),
				strings.ToLower(nameFilter.ValueString()),
			) {
				continue
			}
		}

		filtered = append(filtered, monitor)
	}

	return filtered
}

// convertMonitorToModel converts an API Monitor to MonitorModel.
func convertMonitorToModel(ctx context.Context, monitor *generated.Monitor, diags *diag.Diagnostics) MonitorModel {
	// Create a temporary MonitorResourceModel to reuse conversion logic
	var tempModel MonitorResourceModel
	monitorToState(ctx, monitor, &tempModel, diags)

	// Convert to MonitorModel
	//nolint:gosimple // MonitorModel and MonitorResourceModel are different types, struct literal is required
	model := MonitorModel{
		ID:                    tempModel.ID,
		Name:                  tempModel.Name,
		Type:                  tempModel.Type,
		State:                 tempModel.State,
		NotifyEveryoneByEmail: tempModel.NotifyEveryoneByEmail,
		TeamIDs:               tempModel.TeamIDs,
		ChannelIDs:            tempModel.ChannelIDs,
		RepeatInterval:        tempModel.RepeatInterval,
		Params:                tempModel.Params,
		CreatedAt:             tempModel.CreatedAt,
		UpdatedAt:             tempModel.UpdatedAt,
	}

	return model
}

// convertMonitorModelsToList converts a slice of MonitorModel to types.List.
func convertMonitorModelsToList(_ context.Context, models []MonitorModel) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Define the attribute types for a monitor
	monitorAttrTypes := map[string]attr.Type{
		"id":                       types.StringType,
		"name":                     types.StringType,
		"type":                     types.StringType,
		"state":                    types.StringType,
		"notify_everyone_by_email": types.BoolType,
		"team_ids":                 types.ListType{ElemType: types.Int64Type},
		"channel_ids":              types.ListType{ElemType: types.Int64Type},
		"repeat_interval": types.ObjectType{AttrTypes: map[string]attr.Type{
			"strategy": types.StringType,
			"interval": types.Int64Type,
		}},
		"params": types.ObjectType{AttrTypes: map[string]attr.Type{
			"metrics": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"alias": types.StringType,
			}}},
			"query":             types.StringType,
			"column":            types.StringType,
			"min_allowed_value": types.Float64Type,
			"max_allowed_value": types.Float64Type,
			"grouping_interval": types.Float64Type,
			"check_num_point":   types.Int64Type,
			"nulls_mode":        types.StringType,
			"time_offset":       types.Float64Type,
		}},
		"created_at": types.StringType,
		"updated_at": types.StringType,
	}

	// Convert each model to an object value
	monitorObjects := make([]attr.Value, 0, len(models))
	//nolint:gocritic // Large struct copy acceptable for conversion to Terraform types
	for _, model := range models {
		monitorAttrs := map[string]attr.Value{
			"id":                       model.ID,
			"name":                     model.Name,
			"type":                     model.Type,
			"state":                    model.State,
			"notify_everyone_by_email": model.NotifyEveryoneByEmail,
			"team_ids":                 model.TeamIDs,
			"channel_ids":              model.ChannelIDs,
			"repeat_interval":          model.RepeatInterval,
			"params":                   model.Params,
			"created_at":               model.CreatedAt,
			"updated_at":               model.UpdatedAt,
		}

		objVal, objDiags := types.ObjectValue(monitorAttrTypes, monitorAttrs)
		diags.Append(objDiags...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: monitorAttrTypes}), diags
		}

		monitorObjects = append(monitorObjects, objVal)
	}

	// Create list from objects
	listVal, listDiags := types.ListValue(types.ObjectType{AttrTypes: monitorAttrTypes}, monitorObjects)
	diags.Append(listDiags...)

	return listVal, diags
}
