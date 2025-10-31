package provider

import (
	"context"
	"fmt"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/DavidKrau/terraform-provider-simplemdm/internal/simplemdmext"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &deviceInstalledAppsDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceInstalledAppsDataSource{}
)

type deviceInstalledAppsDataSource struct {
	client *simplemdm.Client
}

type deviceInstalledAppsDataSourceModel struct {
	DeviceID      types.String                       `tfsdk:"device_id"`
	InstalledApps []deviceRelatedItemDataSourceModel `tfsdk:"installed_apps"`
}

func DeviceInstalledAppsDataSource() datasource.DataSource {
	return &deviceInstalledAppsDataSource{}
}

func (d *deviceInstalledAppsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_installed_apps"
}

func (d *deviceInstalledAppsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves installed applications for a device.",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier of the device.",
			},
		},
		Blocks: map[string]schema.Block{
			"installed_apps": schema.ListNestedBlock{
				Description: "Collection of installed applications on the device.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Installed app identifier.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Installed app resource type.",
						},
						"attributes_json": schema.StringAttribute{
							Computed:    true,
							Description: "Raw attributes payload returned by the API in JSON format.",
						},
					},
				},
			},
		},
	}
}

func (d *deviceInstalledAppsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deviceInstalledAppsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apps, err := simplemdmext.ListDeviceInstalledApps(ctx, d.client, state.DeviceID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Device not found",
				fmt.Sprintf("The device with ID %s was not found. It may have been deleted.", state.DeviceID.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to list installed apps",
				err.Error(),
			)
		}
		return
	}

	converted := simplemdmext.ConvertRelatedItems(apps.Data)
	items := make([]deviceRelatedItemDataSourceModel, 0, len(converted))
	for _, item := range converted {
		app := deviceRelatedItemDataSourceModel{
			ID:   types.StringValue(item["id"]),
			Type: types.StringValue(item["type"]),
		}

		if raw := item["attributes"]; raw != "" {
			app.AttributesJSON = types.StringValue(raw)
		} else {
			app.AttributesJSON = types.StringNull()
		}

		items = append(items, app)
	}

	state.InstalledApps = items

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *deviceInstalledAppsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
