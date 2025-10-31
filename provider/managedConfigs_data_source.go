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

var (
	_ datasource.DataSource              = &managedConfigsDataSource{}
	_ datasource.DataSourceWithConfigure = &managedConfigsDataSource{}
)

type managedConfigsDataSource struct {
	client *simplemdm.Client
}

type managedConfigsDataSourceModel struct {
	AppID          types.String                        `tfsdk:"app_id"`
	ManagedConfigs []managedConfigsItemDataSourceModel `tfsdk:"managed_configs"`
}

type managedConfigsItemDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	ValueType types.String `tfsdk:"value_type"`
}

func ManagedConfigsDataSource() datasource.DataSource {
	return &managedConfigsDataSource{}
}

func (d *managedConfigsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_configs"
}

func (d *managedConfigsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves all managed app configurations for a specific app.",
		Attributes: map[string]schema.Attribute{
			"app_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the app to list managed configurations for.",
			},
		},
		Blocks: map[string]schema.Block{
			"managed_configs": schema.ListNestedBlock{
				Description: "Collection of managed app configurations.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Managed config identifier.",
						},
						"key": schema.StringAttribute{
							Computed:    true,
							Description: "Configuration key.",
						},
						"value": schema.StringAttribute{
							Computed:    true,
							Description: "Configuration value.",
						},
						"value_type": schema.StringAttribute{
							Computed:    true,
							Description: "Data type of the value.",
						},
					},
				},
			},
		},
	}
}

func (d *managedConfigsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state managedConfigsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	configs, err := fetchAllManagedConfigs(ctx, d.client, state.AppID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"App not found",
				fmt.Sprintf("The app with ID %s was not found. It may have been deleted.", state.AppID.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to list managed configs",
				fmt.Sprintf("Failed to fetch managed configs for app %s: %v", state.AppID.ValueString(), err),
			)
		}
		return
	}

	items := make([]managedConfigsItemDataSourceModel, 0, len(configs))
	for _, config := range configs {
		item := managedConfigsItemDataSourceModel{
			ID:        types.StringValue(fmt.Sprintf("%d", config.ID)),
			Key:       types.StringValue(config.Attributes.Key),
			Value:     types.StringValue(config.Attributes.Value),
			ValueType: types.StringValue(config.Attributes.ValueType),
		}
		items = append(items, item)
	}

	state.ManagedConfigs = items

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *managedConfigsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*simplemdm.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *simplemdm.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// fetchAllManagedConfigs retrieves all managed configs for an app
func fetchAllManagedConfigs(ctx context.Context, client *simplemdm.Client, appID string) ([]managedConfigAPIResource, error) {
	if client == nil {
		return nil, fmt.Errorf("simplemdm client is not configured")
	}

	url := fmt.Sprintf("https://%s/api/v1/apps/%s/managed_configs", client.HostName, appID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := client.RequestResponse200(req)
	if err != nil {
		return nil, err
	}

	var response managedConfigListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}
