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
	_ datasource.DataSource              = &deviceUsersDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceUsersDataSource{}
)

type deviceUsersDataSource struct {
	client *simplemdm.Client
}

type deviceUsersDataSourceModel struct {
	DeviceID types.String                       `tfsdk:"device_id"`
	Users    []deviceRelatedItemDataSourceModel `tfsdk:"users"`
}

func DeviceUsersDataSource() datasource.DataSource {
	return &deviceUsersDataSource{}
}

func (d *deviceUsersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_users"
}

func (d *deviceUsersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves user accounts currently configured on a macOS device.",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier of the device.",
			},
		},
		Blocks: map[string]schema.Block{
			"users": schema.ListNestedBlock{
				Description: "Collection of user accounts present on the device.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "User identifier.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "User resource type.",
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

func (d *deviceUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deviceUsersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	users, err := simplemdmext.ListDeviceUsers(ctx, d.client, state.DeviceID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Device not found",
				fmt.Sprintf("The device with ID %s was not found. It may have been deleted.", state.DeviceID.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to list device users",
				err.Error(),
			)
		}
		return
	}

	converted := simplemdmext.ConvertRelatedItems(users.Data)
	items := make([]deviceRelatedItemDataSourceModel, 0, len(converted))
	for _, item := range converted {
		user := deviceRelatedItemDataSourceModel{
			ID:   types.StringValue(item["id"]),
			Type: types.StringValue(item["type"]),
		}

		if raw := item["attributes"]; raw != "" {
			user.AttributesJSON = types.StringValue(raw)
		} else {
			user.AttributesJSON = types.StringNull()
		}

		items = append(items, user)
	}

	state.Users = items

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *deviceUsersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
