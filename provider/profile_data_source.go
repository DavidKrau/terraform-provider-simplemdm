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
	_ datasource.DataSource              = &profileDataSource{}
	_ datasource.DataSourceWithConfigure = &profileDataSource{}
)

// ProfileDataSourceModel maps the data source schema data.
type profileDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// ProfileDataSource is a helper function to simplify the provider implementation.
func ProfileDataSource() datasource.DataSource {
	return &profileDataSource{}
}

// profileDataSource is the data source implementation.
type profileDataSource struct {
	client *simplemdm.Client
}

// Metadata returns the data source type name.
func (d *profileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}

// Schema defines the schema for the data source.
func (d *profileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Profile data source can be used together with Device(s), Assignment Group(s) or Device Group(s) to assign profiles to these objects.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the Profile.",
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Profile.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *profileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state profileDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	profiles, err := d.client.CustomProfileGetAll()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SimpleMDM profile",
			err.Error(),
		)
		return
	}

	profilefound := false
	for _, profile := range profiles.Data {
		if state.ID.ValueString() == strconv.Itoa(profile.ID) {
			state.Name = types.StringValue(profile.Attributes.Name)
			profilefound = true
			break
		}
	}

	if !profilefound {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM profile",
			"Could not read profile ID %s from array:"+state.ID.ValueString(),
		)
		return
	}

	// Set state

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *profileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
