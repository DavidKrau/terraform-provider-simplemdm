package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &deviceDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceDataSource{}
)

// deviceDataSourceModel maps the data source schema data.
type deviceDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// deviceDataSource is a helper function to simplify the provider implementation.
func DeviceDataSource() datasource.DataSource {
	return &deviceDataSource{}
}

// deviceDataSource is the data source implementation.
type deviceDataSource struct {
	client *simplemdm.Client
}

// Metadata returns the data source type name.
func (d *deviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

// Schema defines the schema for the data source.
func (d *deviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Device data source can be used together Assignment Group(s) to assign device to these objects.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The SimpleMDM name of the device.",
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the device.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *deviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deviceDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	device, err := d.client.DeviceGet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SimpleMDM device",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state.Name = types.StringValue(device.Data.Attributes.Name)
	state.ID = types.StringValue(strconv.Itoa(device.Data.ID))

	// Set state

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *deviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
