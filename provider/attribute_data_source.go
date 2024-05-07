package provider

import (
	"context"
	"fmt"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &attributeDataSource{}
	_ datasource.DataSourceWithConfigure = &attributeDataSource{}
)

// attributeDataSourceModel maps the data source schema data.
type attributeDataSourceModel struct {
	DefaultValue types.String `tfsdk:"default_value"`
	Name         types.String `tfsdk:"name"`
}

// AttributeDataSource is a helper function to simplify the provider implementation.
func AttributeDataSource() datasource.DataSource {
	return &attributeDataSource{}
}

// AttributeDataSource is the data source implementation.
type attributeDataSource struct {
	client *simplemdm.Client
}

// Metadata returns the data source type name.
func (d *attributeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_attribute"
}

// Schema defines the schema for the data source.
func (d *attributeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed: true,
			},
			"default_value": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *attributeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state attributeDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	attribute, err := d.client.GetAttribute(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SimpleMDM attribute",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state.Name = types.StringValue(attribute.Data.Attributes.Name)
	state.DefaultValue = types.StringValue(attribute.Data.Attributes.DefaultValue)

	// Set state

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *attributeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
