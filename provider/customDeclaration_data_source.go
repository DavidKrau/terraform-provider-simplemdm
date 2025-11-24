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
	_ datasource.DataSource              = &customDeclarationDataSource{}
	_ datasource.DataSourceWithConfigure = &customDeclarationDataSource{}
)

// customDeclarationModel maps the data source schema data.
type customDeclarationDataSourceModel struct {
	Name types.String `tfsdk:"name"`
	ID   types.String `tfsdk:"id"`
}

// AttributeDataSource is a helper function to simplify the provider implementation.
func CustomDeclarationDataSource() datasource.DataSource {
	return &customDeclarationDataSource{}
}

// AttributeDataSource is the data source implementation.
type customDeclarationDataSource struct {
	client *simplemdm.Client
}

// Metadata returns the data source type name.
func (d *customDeclarationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customdeclaration"
}

// Schema defines the schema for the data source.
func (d *customDeclarationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Custom Declaration data source can be used together with Device(s) or Device Group(s) to set values or in lifecycle management.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the custom Declaration.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Default (global) value of the Attribute.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *customDeclarationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state customDeclarationDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	declaration, err := d.client.ProfileGet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SimpleMDM custom declaration",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state.Name = types.StringValue(declaration.Data.Attributes.Name)
	state.ID = types.StringValue(strconv.Itoa(declaration.Data.ID))

	// Set state

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *customDeclarationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
