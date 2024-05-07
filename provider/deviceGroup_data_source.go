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
	_ datasource.DataSource              = &deviceGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceGroupDataSource{}
)

// coffeesDataSourceModel maps the data source schema data.
type deviceGroupDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func DeviceGroupDataSource() datasource.DataSource {
	return &deviceGroupDataSource{}
}

// coffeesDataSource is the data source implementation.
type deviceGroupDataSource struct {
	client *simplemdm.Client
}

// Metadata returns the data source type name.
func (d *deviceGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devicegroup"
}

// Schema defines the schema for the data source.
func (d *deviceGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Device group name.",
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of a Device Group in SimpleMDM.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *deviceGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deviceGroupDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	app, err := d.client.GetDeviceGroup(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SimpleMDM device group",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state.Name = types.StringValue(app.Data.Attributes.Name)
	state.ID = types.StringValue(strconv.Itoa(app.Data.ID))

	// Set state

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *deviceGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
