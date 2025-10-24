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
	_ datasource.DataSource              = &deviceProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceProfilesDataSource{}
)

type deviceProfilesDataSource struct {
	client *simplemdm.Client
}

type deviceProfilesDataSourceModel struct {
	DeviceID types.String                       `tfsdk:"device_id"`
	Profiles []deviceRelatedItemDataSourceModel `tfsdk:"profiles"`
}

type deviceRelatedItemDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Type           types.String `tfsdk:"type"`
	AttributesJSON types.String `tfsdk:"attributes_json"`
}

func DeviceProfilesDataSource() datasource.DataSource {
	return &deviceProfilesDataSource{}
}

func (d *deviceProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_profiles"
}

func (d *deviceProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of profiles directly assigned to a device.",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier of the device.",
			},
		},
		Blocks: map[string]schema.Block{
			"profiles": schema.ListNestedBlock{
				Description: "Collection of profiles applied directly to the device.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Profile identifier.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Profile resource type.",
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

func (d *deviceProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deviceProfilesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	profiles, err := simplemdmext.ListDeviceProfiles(ctx, d.client, state.DeviceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list device profiles",
			err.Error(),
		)
		return
	}

	converted := simplemdmext.ConvertRelatedItems(profiles.Data)
	items := make([]deviceRelatedItemDataSourceModel, 0, len(converted))
	for _, item := range converted {
		profile := deviceRelatedItemDataSourceModel{
			ID:   types.StringValue(item["id"]),
			Type: types.StringValue(item["type"]),
		}

		if raw := item["attributes"]; raw != "" {
			profile.AttributesJSON = types.StringValue(raw)
		} else {
			profile.AttributesJSON = types.StringNull()
		}

		items = append(items, profile)
	}

	state.Profiles = items

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *deviceProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
