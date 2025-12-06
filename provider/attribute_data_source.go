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
	ID           types.String `tfsdk:"id"`
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
		Description: "Attribute data source can be used together with Device(s) or Device Group(s) to set values or in lifecycle management.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for the custom attribute (same as name).",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name (and ID) of the Attribute.",
			},
			"default_value": schema.StringAttribute{
				Computed:    true,
				Description: "Default (global) value of the Attribute.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *attributeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state attributeDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	attribute, err := d.client.AttributeGet(state.Name.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"SimpleMDM attribute not found",
				fmt.Sprintf("The attribute with name %s was not found. It may have been deleted.", state.Name.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to Read SimpleMDM attribute",
				err.Error(),
			)
		}
		return
	}

	// Map response body to model
	state.ID = types.StringValue(attribute.Data.Attributes.Name)
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
