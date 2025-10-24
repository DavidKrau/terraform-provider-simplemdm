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
	ID                     types.String `tfsdk:"id"`
	Type                   types.String `tfsdk:"type"`
	Name                   types.String `tfsdk:"name"`
	AutoDeploy             types.Bool   `tfsdk:"auto_deploy"`
	InstallType            types.String `tfsdk:"install_type"`
	ReinstallAfterOSUpdate types.Bool   `tfsdk:"reinstall_after_os_update"`
	ProfileIdentifier      types.String `tfsdk:"profile_identifier"`
	UserScope              types.Bool   `tfsdk:"user_scope"`
	AttributeSupport       types.Bool   `tfsdk:"attribute_support"`
	EscapeAttributes       types.Bool   `tfsdk:"escape_attributes"`
	GroupCount             types.Int64  `tfsdk:"group_count"`
	DeviceCount            types.Int64  `tfsdk:"device_count"`
	GroupIDs               types.Set    `tfsdk:"group_ids"`
	ProfileSHA             types.String `tfsdk:"profile_sha"`
	Source                 types.String `tfsdk:"source"`
	CreatedAt              types.String `tfsdk:"created_at"`
	UpdatedAt              types.String `tfsdk:"updated_at"`
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
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Profile.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The profile payload type reported by SimpleMDM.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the Profile.",
			},
			"auto_deploy": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the profile is auto-deployed when assigned.",
			},
			"install_type": schema.StringAttribute{
				Computed:    true,
				Description: "The install type configured for the profile.",
			},
			"reinstall_after_os_update": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the profile reinstalls automatically after macOS updates.",
			},
			"profile_identifier": schema.StringAttribute{
				Computed:    true,
				Description: "The identifier contained within the profile payload.",
			},
			"user_scope": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates if the profile is installed in the user scope.",
			},
			"attribute_support": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the profile supports attribute substitution.",
			},
			"escape_attributes": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether attribute values are escaped during installation.",
			},
			"group_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of device groups currently assigned to the profile.",
			},
			"device_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of devices currently assigned to the profile.",
			},
			"group_ids": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "IDs of device or assignment groups currently assigned to the profile.",
			},
			"profile_sha": schema.StringAttribute{
				Computed:    true,
				Description: "SHA hash reported by SimpleMDM for the profile contents.",
			},
			"source": schema.StringAttribute{
				Computed:    true,
				Description: "Origin of the profile within SimpleMDM.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the profile was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the profile was last updated.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *profileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state profileDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	profile, err := fetchProfile(ctx, d.client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SimpleMDM profile",
			err.Error(),
		)
		return
	}

	groupIDs, err := convertGroupIDs(ctx, profile.Data.Relationships)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SimpleMDM profile",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state.ID = types.StringValue(strconv.Itoa(profile.Data.ID))
	if profile.Data.Type != "" {
		state.Type = types.StringValue(profile.Data.Type)
	} else {
		state.Type = types.StringNull()
	}
	state.Name = types.StringValue(profile.Data.Attributes.Name)
	state.AutoDeploy = types.BoolValue(profile.Data.Attributes.AutoDeploy)
	state.InstallType = types.StringValue(profile.Data.Attributes.InstallType)
	state.ReinstallAfterOSUpdate = types.BoolValue(profile.Data.Attributes.ReinstallAfterOSUpdate)
	state.ProfileIdentifier = types.StringValue(profile.Data.Attributes.ProfileIdentifier)
	state.UserScope = types.BoolValue(profile.Data.Attributes.UserScope)
	state.AttributeSupport = types.BoolValue(profile.Data.Attributes.AttributeSupport)
	state.EscapeAttributes = types.BoolValue(profile.Data.Attributes.EscapeAttributes)
	state.GroupCount = types.Int64Value(int64(profile.Data.Attributes.GroupCount))
	state.DeviceCount = types.Int64Value(int64(profile.Data.Attributes.DeviceCount))
	state.GroupIDs = groupIDs
	state.ProfileSHA = types.StringValue(profile.Data.Attributes.ProfileSHA)
	state.Source = types.StringValue(profile.Data.Attributes.Source)
	state.CreatedAt = types.StringValue(profile.Data.Attributes.CreatedAt)
	state.UpdatedAt = types.StringValue(profile.Data.Attributes.UpdatedAt)

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
