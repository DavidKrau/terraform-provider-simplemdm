package provider

import (
	"context"
	"fmt"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &managedConfigDataSource{}
	_ datasource.DataSourceWithConfigure = &managedConfigDataSource{}
)

type managedConfigDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	AppID     types.String `tfsdk:"app_id"`
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	ValueType types.String `tfsdk:"value_type"`
}

type managedConfigDataSource struct {
	client *simplemdm.Client
}

func ManagedConfigDataSource() datasource.DataSource {
	return &managedConfigDataSource{}
}

func (d *managedConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_config"
}

func (d *managedConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Managed Config data source retrieves the value of a managed app configuration for a given app.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Managed config identifier in the format &lt;app_id&gt;:&lt;managed_config_id&gt;.",
			},
			"app_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the app that owns the managed configuration.",
			},
			"key": schema.StringAttribute{
				Computed:    true,
				Description: "Configuration key returned by the SimpleMDM API.",
			},
			"value": schema.StringAttribute{
				Computed:    true,
				Description: "Raw value returned by the SimpleMDM API.",
			},
			"value_type": schema.StringAttribute{
				Computed:    true,
				Description: "Data type that SimpleMDM reports for the value.",
			},
		},
	}
}

func (d *managedConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state managedConfigDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	appID, configID, err := parseManagedConfigID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid managed config identifier", err.Error())
		return
	}

	config, err := fetchManagedConfig(ctx, d.client, appID, configID)
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"SimpleMDM managed config not found",
				fmt.Sprintf("The managed config %s for app %s was not found. It may have been deleted.", configID, appID),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to read managed app configuration",
				fmt.Sprintf("failed to fetch managed config %s for app %s: %v", configID, appID, err),
			)
		}
		return
	}

	state.AppID = types.StringValue(appID)
	state.Key = types.StringValue(config.Attributes.Key)
	state.Value = types.StringValue(config.Attributes.Value)
	state.ValueType = types.StringValue(config.Attributes.ValueType)
	state.ID = types.StringValue(fmt.Sprintf("%s:%d", appID, config.ID))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *managedConfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
