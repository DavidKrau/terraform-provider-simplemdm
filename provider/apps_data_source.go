package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &appsDataSource{}
	_ datasource.DataSourceWithConfigure = &appsDataSource{}
)

type appsDataSource struct {
	client *simplemdm.Client
}

type appsDataSourceModel struct {
	IncludeShared types.Bool               `tfsdk:"include_shared"`
	Apps          []appsDataSourceAppModel `tfsdk:"apps"`
}

type appsDataSourceAppModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	AppStoreID           types.String `tfsdk:"app_store_id"`
	BundleID             types.String `tfsdk:"bundle_id"`
	AppType              types.String `tfsdk:"app_type"`
	Version              types.String `tfsdk:"version"`
	PlatformSupport      types.String `tfsdk:"platform_support"`
	ProcessingStatus     types.String `tfsdk:"processing_status"`
	InstallationChannels types.List   `tfsdk:"installation_channels"`
}

func AppsDataSource() datasource.DataSource {
	return &appsDataSource{}
}

func (d *appsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apps"
}

func (d *appsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the collection of apps from your SimpleMDM account, with optional filtering by shared catalog.",
		Attributes: map[string]schema.Attribute{
			"include_shared": schema.BoolAttribute{
				Optional:    true,
				Description: "Include apps from the SimpleMDM shared catalog. When set to true, the data source will query apps available in the shared catalog in addition to account-specific apps. Defaults to false.",
			},
		},
		Blocks: map[string]schema.Block{
			"apps": schema.ListNestedBlock{
				Description: "Collection of app records returned by the API.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "App identifier.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the app.",
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
							Description: "The catalog classification of the app (e.g., app store, enterprise, custom b2b).",
						},
						"version": schema.StringAttribute{
							Computed:    true,
							Description: "The latest version reported by SimpleMDM for the app.",
						},
						"platform_support": schema.StringAttribute{
							Computed:    true,
							Description: "The platform supported by the app (iOS, macOS, etc.).",
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
					},
				},
			},
		},
	}
}

func (d *appsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config appsDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apps, err := fetchAllApps(ctx, d.client, config.IncludeShared)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list SimpleMDM apps",
			err.Error(),
		)
		return
	}

	entries := make([]appsDataSourceAppModel, 0, len(apps))
	for _, app := range apps {
		installationChannels, diags := types.ListValueFrom(ctx, types.StringType, app.Attributes.InstallationChannels)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		entry := appsDataSourceAppModel{
			ID:                   types.StringValue(strconv.Itoa(app.ID)),
			Name:                 types.StringValue(app.Attributes.Name),
			AppStoreID:           stringValueOrNull(app.Attributes.AppStoreID),
			BundleID:             stringValueOrNull(app.Attributes.BundleIdentifier),
			AppType:              stringValueOrNull(app.Attributes.AppType),
			Version:              stringValueOrNull(app.Attributes.Version),
			PlatformSupport:      stringValueOrNull(app.Attributes.PlatformSupport),
			ProcessingStatus:     stringValueOrNull(app.Attributes.ProcessingStatus),
			InstallationChannels: installationChannels,
		}

		entries = append(entries, entry)
	}

	config.Apps = entries

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func (d *appsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// fetchAllApps retrieves all apps with pagination support
func fetchAllApps(ctx context.Context, client *simplemdm.Client, includeShared types.Bool) ([]appData, error) {
	var allApps []appData
	startingAfter := 0
	limit := 100

	for {
		url := fmt.Sprintf("https://%s/api/v1/apps?limit=%d", client.HostName, limit)
		if startingAfter > 0 {
			url += fmt.Sprintf("&starting_after=%d", startingAfter)
		}
		if !includeShared.IsNull() && includeShared.ValueBool() {
			url += "&include_shared=true"
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		body, err := client.RequestResponse200(req)
		if err != nil {
			return nil, err
		}

		var response appsAPIResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}

		allApps = append(allApps, response.Data...)

		if !response.HasMore {
			break
		}

		if len(response.Data) > 0 {
			startingAfter = response.Data[len(response.Data)-1].ID
		} else {
			break
		}
	}

	return allApps, nil
}

// appsAPIResponse represents the paginated API response for apps list
type appsAPIResponse struct {
	Data    []appData `json:"data"`
	HasMore bool      `json:"has_more"`
}

// appData represents a single app in the list response
type appData struct {
	ID         int           `json:"id"`
	Type       string        `json:"type"`
	Attributes appAttributes `json:"attributes"`
}

type appAttributes struct {
	Name                 string   `json:"name"`
	AppStoreID           string   `json:"itunes_store_id"`
	BundleIdentifier     string   `json:"bundle_identifier"`
	AppType              string   `json:"app_type"`
	Version              string   `json:"version"`
	PlatformSupport      string   `json:"platform_support"`
	ProcessingStatus     string   `json:"processing_status"`
	InstallationChannels []string `json:"installation_channels"`
}
