package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/riccap/tofu-uptrace-provider/internal/client"
)

// Ensure the implementation satisfies the provider.Provider interface.
var _ provider.Provider = &UptraceProvider{}

// UptraceProvider defines the provider implementation.
type UptraceProvider struct {
	version string
}

// UptraceProviderModel describes the provider data model.
type UptraceProviderModel struct {
	Endpoint  types.String `tfsdk:"endpoint"`
	Token     types.String `tfsdk:"token"`
	ProjectID types.Int64  `tfsdk:"project_id"`
}

// Metadata returns the provider type name.
func (p *UptraceProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "uptrace"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *UptraceProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Uptrace API to manage monitors and dashboards.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "The Uptrace API endpoint. May also be provided via UPTRACE_ENDPOINT environment variable.",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "The authentication token for Uptrace API. May also be provided via UPTRACE_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"project_id": schema.Int64Attribute{
				Description: "The default project ID for Uptrace operations. May also be provided via UPTRACE_PROJECT_ID environment variable.",
				Optional:    true,
			},
		},
	}
}

// Configure prepares an API client for data sources and resources.
func (p *UptraceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Uptrace client")

	var config UptraceProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get endpoint from config or environment
	endpoint := os.Getenv("UPTRACE_ENDPOINT")
	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Uptrace API Endpoint",
			"The provider cannot create the Uptrace API client as there is a missing or empty value for the Uptrace API endpoint. "+
				"Set the endpoint value in the configuration or use the UPTRACE_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	// Get token from config or environment
	token := os.Getenv("UPTRACE_TOKEN")
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Uptrace API Token",
			"The provider cannot create the Uptrace API client as there is a missing or empty value for the Uptrace API token. "+
				"Set the token value in the configuration or use the UPTRACE_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	// Get project ID from config or environment
	var projectID int64
	if !config.ProjectID.IsNull() {
		projectID = config.ProjectID.ValueInt64()
	} else if envProjectID := os.Getenv("UPTRACE_PROJECT_ID"); envProjectID != "" {
		var err error
		projectID, err = strconv.ParseInt(envProjectID, 10, 64)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("project_id"),
				"Invalid Uptrace Project ID",
				fmt.Sprintf("The UPTRACE_PROJECT_ID environment variable value %q is not a valid integer: %s", envProjectID, err),
			)
		}
	}

	if projectID <= 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("project_id"),
			"Missing Uptrace Project ID",
			"The provider cannot create the Uptrace API client as there is a missing or invalid value for the Uptrace project ID. "+
				"Set the project_id value in the configuration or use the UPTRACE_PROJECT_ID environment variable. "+
				"The project ID must be greater than 0.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create client
	uptraceClient, err := client.New(client.Config{
		Endpoint:  endpoint,
		Token:     token,
		ProjectID: projectID,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Uptrace API Client",
			"An unexpected error occurred when creating the Uptrace API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Uptrace Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = uptraceClient
	resp.ResourceData = uptraceClient

	tflog.Info(ctx, "Configured Uptrace client", map[string]any{
		"endpoint":   endpoint,
		"project_id": projectID,
	})
}

// DataSources defines the data sources implemented in the provider.
func (p *UptraceProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMonitorDataSource,
		NewMonitorsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *UptraceProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewMonitorResource,
		NewDashboardResource,
		NewNotificationChannelResource,
	}
}

// New returns a new provider instance.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &UptraceProvider{
			version: version,
		}
	}
}
