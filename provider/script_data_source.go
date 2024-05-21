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
	_ datasource.DataSource              = &scriptDataSource{}
	_ datasource.DataSourceWithConfigure = &scriptDataSource{}
)

// appDataSourceModel maps the data source schema data.
type scriptDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// appDataSource is a helper function to simplify the provider implementation.
func ScriptDataSource() datasource.DataSource {
	return &scriptDataSource{}
}

// appDataSource is the data source implementation.
type scriptDataSource struct {
	client *simplemdm.Client
}

// Metadata returns the data source type name.
func (d *scriptDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_script"
}

// Schema defines the schema for the data source.
func (d *scriptDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Script data source can be used together with Script Jobs(s).",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the Script.",
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Script.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *scriptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state scriptDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	app, err := d.client.ScriptGet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SimpleMDM Script",
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
func (d *scriptDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
