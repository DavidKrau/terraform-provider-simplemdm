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
	_ datasource.DataSource              = &customProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &customProfileDataSource{}
)

// ProfileDataSourceModel maps the data source schema data.
type customProfileDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// ProfileDataSource is a helper function to simplify the provider implementation.
func CustomProfileDataSource() datasource.DataSource {
	return &customProfileDataSource{}
}

// profileDataSource is the data source implementation.
type customProfileDataSource struct {
	client *simplemdm.Client
}

// Metadata returns the data source type name.
func (d *customProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customprofile"
}

// Schema defines the schema for the data source.
func (d *customProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Custom Profile data source can be used together with Device(s), Assignment Group(s) or Device Group(s) to assign profiles to these objects.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the custom profile.",
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the custom profile.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *customProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state customProfileDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	profile, err := d.client.ProfileGet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM custom profile",
			"Could not read custom profile ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Set state
	state.Name = types.StringValue(profile.Data.Attributes.Name)
	state.ID = types.StringValue(strconv.Itoa(profile.Data.ID))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *customProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*simplemdm.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
