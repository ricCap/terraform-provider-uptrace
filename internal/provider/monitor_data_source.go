package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/riccap/tofu-uptrace-provider/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &MonitorDataSource{}
	_ datasource.DataSourceWithConfigure = &MonitorDataSource{}
)

// NewMonitorDataSource is a helper function to create the data source.
func NewMonitorDataSource() datasource.DataSource {
	return &MonitorDataSource{}
}

// MonitorDataSource is the data source implementation.
type MonitorDataSource struct {
	client *client.Client
}

// Metadata returns the data source type name.
func (d *MonitorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

// Schema defines the schema for the data source.
func (d *MonitorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches an existing Uptrace monitor by ID.",
		//nolint:dupl // Schema duplication with monitors_data_source acceptable - different data sources
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Monitor identifier.",
				Required:    true,
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
						Description: "Custom interval in seconds (minimum 60).",
						Computed:    true,
					},
				},
			},
			"params": schema.SingleNestedAttribute{
				Description: "Monitor parameters (metric or error specific).",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					// Metric monitor params
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
	}
}

// Configure adds the provider configured client to the data source.
func (d *MonitorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *MonitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config MonitorResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Reading monitor data source", map[string]any{"id": config.ID.ValueString()})

	// Get monitor from API
	monitor, err := d.client.GetMonitor(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Monitor",
			fmt.Sprintf("Could not read monitor ID %s: %s", config.ID.ValueString(), err.Error()),
		)
		return
	}

	// Convert API response to state using existing helper
	monitorToState(ctx, monitor, &config, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	tflog.Info(ctx, "Successfully read monitor data source", map[string]any{"id": config.ID.ValueString()})
}
