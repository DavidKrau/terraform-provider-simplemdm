package provider

import (
	"context"
	"os"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &simplemdmProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &simplemdmProvider{
			version: version,
		}
	}
}

// simplemdmProvider is the provider implementation.
type simplemdmProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// simplemdmProviderModel maps provider schema data to a Go type.
type simplemdmProviderModel struct {
	Host   types.String `tfsdk:"host"`
	APIKey types.String `tfsdk:"apikey"`
}

// Metadata returns the provider type name.
func (p *simplemdmProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "simplemdm"
}

// Schema defines the provider-level schema for configuration data.
func (p *simplemdmProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "SimpleMDM terraform provider developed by FreeNow.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "API host for you instance, can be set as environment variable SIMPLEMDM_HOST, if not set it will default to a.simplemdm.com",
			},
			"apikey": schema.StringAttribute{
				Optional:    true,
				Description: "API key for you instance, can be set as environment variable SIMPLEMDM_APIKEY",
			},
		},
	}
}

func (p *simplemdmProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring simplemdm client")

	//Retrieve provider data from configuration
	var config simplemdmProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown SimpleMDM Host",
			"The provider cannot create the simplemdm API client as there is an unknown configuration value for the SimpleMDM host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SIMPLEMDM_HOST environment variable.",
		)
	}

	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown SimpleMDM Host",
			"The provider cannot create the simplemdm API client as there is an unknown configuration value for the SimpleMDM host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SIMPLEMDM_APIKEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	host := os.Getenv("SIMPLEMDM_HOST")
	apikey := os.Getenv("SIMPLEMDM_APIKEY")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.APIKey.IsNull() {
		apikey = config.APIKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if host == "" {
		host = "a.simplemdm.com"
	}

	if apikey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Missing SimpleMDM Host",
			"The provider cannot create the SimpleMDM API client as there is a missing or empty value for the SimpleMDM host. "+
				"Set the host value in the configuration or use the SIMPLEMDM_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "simplemdm_host", host)
	ctx = tflog.SetField(ctx, "simplemdm_apikey", apikey)

	tflog.Debug(ctx, "Creating SimpleMDM client")

	apiClient := simplemdm.NewClient(host, apikey)

	// Make the SimpleMDM client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient

	tflog.Info(ctx, "Configured SimpleMDM client", map[string]any{"success": true})

}

// DataSources defines the data sources implemented in the provider.
func (p *simplemdmProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		AppDataSource, AttributeDataSource, DeviceGroupDataSource, CustomProfileDataSource, ProfileDataSource, DeviceDataSource, ScriptDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *simplemdmProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		CustomProfileResource, AttributeResource, AssignmentGroupResource, DeviceGroupResource, DeviceResource,
	}
}
