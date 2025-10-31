package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &appDataSource{}
	_ datasource.DataSourceWithConfigure = &appDataSource{}
)

// appDataSourceModel maps the data source schema data.
type appDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	AppStoreId           types.String `tfsdk:"app_store_id"`
	BundleId             types.String `tfsdk:"bundle_id"`
	DeployTo             types.String `tfsdk:"deploy_to"`
	Status               types.String `tfsdk:"status"`
	AppType              types.String `tfsdk:"app_type"`
	Version              types.String `tfsdk:"version"`
	PlatformSupport      types.String `tfsdk:"platform_support"`
	ProcessingStatus     types.String `tfsdk:"processing_status"`
	InstallationChannels types.List   `tfsdk:"installation_channels"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
	IncludeShared        types.Bool   `tfsdk:"include_shared"`
}

// appDataSource is a helper function to simplify the provider implementation.
func AppDataSource() datasource.DataSource {
	return &appDataSource{}
}

// appDataSource is the data source implementation.
type appDataSource struct {
	client *simplemdm.Client
}

// Metadata returns the data source type name.
func (d *appDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

// Schema defines the schema for the data source.
func (d *appDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "App data source can be used together with Assignment Group(s) to assign App(s) to the group(s).",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the attribute.",
			},
			"app_store_id": schema.StringAttribute{
				Computed:    true,
				Description: "The Apple App Store ID associated with the app.",
			},
			"bundle_id": schema.StringAttribute{
				Computed:    true,
				Description: "The bundle identifier of the app.",
			},
			"app_type": schema.StringAttribute{
				Computed:    true,
				Description: "The catalog classification of the app, for example app store, enterprise, or custom b2b.",
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "The latest version reported by SimpleMDM for the app.",
			},
			"platform_support": schema.StringAttribute{
				Computed:    true,
				Description: "The platform supported by the app, such as iOS or macOS.",
			},
			"processing_status": schema.StringAttribute{
				Computed:    true,
				Description: "The current processing status of the app binary within SimpleMDM.",
			},
			"installation_channels": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The deployment channels supported by the app.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the app was added to SimpleMDM.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the app was last updated in SimpleMDM.",
			},
			"deploy_to": schema.StringAttribute{
				Computed:    true,
				Description: "Where the app is deployed (none, outdated, or all).",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The current deployment status of the app.",
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the attribute.",
			},
			"include_shared": schema.BoolAttribute{
				Optional:    true,
				Description: "Include apps from the SimpleMDM shared catalog. When set to true, the data source will query apps available in the shared catalog in addition to account-specific apps. Defaults to false.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *appDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state appDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	app, err := fetchAppWithParams(ctx, d.client, state.ID.ValueString(), state.IncludeShared)
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"SimpleMDM app not found",
				fmt.Sprintf("The app with ID %s was not found. It may have been deleted.", state.ID.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to Read SimpleMDM app",
				err.Error(),
			)
		}
		return
	}

	// Map response body to model
	resourceModel, diags := newAppResourceModelFromAPI(ctx, app)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = resourceModel.ID
	state.Name = resourceModel.Name
	state.AppStoreId = resourceModel.AppStoreId
	state.BundleId = resourceModel.BundleId
	state.DeployTo = resourceModel.DeployTo
	state.Status = resourceModel.Status
	state.AppType = resourceModel.AppType
	state.Version = resourceModel.Version
	state.PlatformSupport = resourceModel.PlatformSupport
	state.ProcessingStatus = resourceModel.ProcessingStatus
	state.InstallationChannels = resourceModel.InstallationChannels
	state.CreatedAt = resourceModel.CreatedAt
	state.UpdatedAt = resourceModel.UpdatedAt

// fetchAppWithParams fetches an app with optional query parameters
func fetchAppWithParams(ctx context.Context, client *simplemdm.Client, appID string, includeShared types.Bool) (*appAPIResponse, error) {
	url := fmt.Sprintf("https://%s/api/v1/apps/%s", client.HostName, appID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Add query parameters if include_shared is explicitly set to true
	if !includeShared.IsNull() && includeShared.ValueBool() {
		q := req.URL.Query()
		q.Add("include_shared", "true")
		req.URL.RawQuery = q.Encode()
	}

	body, err := client.RequestResponse200(req)
	if err != nil {
		return nil, err
	}

	app := appAPIResponse{}
	if err := json.Unmarshal(body, &app); err != nil {
		return nil, err
	}

	return &app, nil
}

	// Set state

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *appDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*simplemdm.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
