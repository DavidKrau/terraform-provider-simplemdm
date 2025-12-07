package provider

import (
	"context"
	"fmt"
	"strings"

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
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	MobileConfig           types.String `tfsdk:"mobileconfig"`
	UserScope              types.Bool   `tfsdk:"user_scope"`
	AttributeSupport       types.Bool   `tfsdk:"attribute_support"`
	EscapeAttributes       types.Bool   `tfsdk:"escape_attributes"`
	ReinstallAfterOSUpdate types.Bool   `tfsdk:"reinstall_after_os_update"`
	ProfileIdentifier      types.String `tfsdk:"profile_identifier"`
	GroupCount             types.Int64  `tfsdk:"group_count"`
	DeviceCount            types.Int64  `tfsdk:"device_count"`
	ProfileSHA             types.String `tfsdk:"profile_sha"`
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
			"mobileconfig": schema.StringAttribute{
				Computed:    true,
				Description: "Contents of the downloaded custom configuration profile.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the custom profile.",
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the custom profile.",
			},
			"user_scope": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the profile deploys as a user profile for macOS devices.",
			},
			"attribute_support": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether variable substitution is enabled for the profile.",
			},
			"escape_attributes": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether custom attribute values are escaped when substituted into the profile.",
			},
			"reinstall_after_os_update": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the profile reinstalls automatically after macOS updates.",
			},
			"profile_identifier": schema.StringAttribute{
				Computed:    true,
				Description: "Profile identifier assigned by SimpleMDM.",
			},
			"group_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of device groups currently assigned to this profile.",
			},
			"device_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of devices currently assigned to this profile.",
			},
			"profile_sha": schema.StringAttribute{
				Computed:    true,
				Description: "SHA-256 checksum reported by SimpleMDM for the profile payload.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *customProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state customProfileDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	// NOTE: CustomProfileGet uses GET /api/v1/custom_configuration_profiles/{id}
	// This endpoint is not documented in the API specification but is functional.
	profile, err := d.client.CustomProfileGet(state.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.Diagnostics.AddError(
				"Error reading SimpleMDM custom profile",
				"Custom profile with ID "+state.ID.ValueString()+" was not found.",
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading SimpleMDM custom profile",
			"Could not read custom profile ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(profile.Data.Attributes.Name)
	state.UserScope = types.BoolValue(profile.Data.Attributes.UserScope)
	state.AttributeSupport = types.BoolValue(profile.Data.Attributes.AttributeSupport)
	state.EscapeAttributes = types.BoolValue(profile.Data.Attributes.EscapeAttributes)
	state.ReinstallAfterOSUpdate = types.BoolValue(profile.Data.Attributes.ReinstallAfterOsUpdate)
	state.ProfileIdentifier = stringValueOrNull(profile.Data.Attributes.ProfileIdentifier)
	state.GroupCount = types.Int64Value(int64(profile.Data.Attributes.GroupCount))
	state.DeviceCount = types.Int64Value(int64(profile.Data.Attributes.DeviceCount))

	// NOTE: CustomProfileSHA downloads the profile using GET /api/v1/custom_configuration_profiles/{id}/download
	// and computes the SHA-256 checksum locally. The 'profile_sha' field is not returned by the API.
	sha, body, err := d.client.CustomProfileSHA(state.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.Diagnostics.AddError(
				"Error reading SimpleMDM custom profile",
				"Custom profile payload for ID "+state.ID.ValueString()+" was not found.",
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading SimpleMDM custom profile",
			"Could not download custom profile ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.MobileConfig = types.StringValue(body)
	state.ProfileSHA = stringValueOrNull(sha)

	// Set state

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
