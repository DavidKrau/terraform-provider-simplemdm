package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/DavidKrau/terraform-provider-simplemdm/internal/simplemdmext"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &deviceProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceProfilesDataSource{}
)

type deviceProfilesDataSource struct {
	client *simplemdm.Client
}

type deviceProfilesDataSourceModel struct {
	DeviceID types.String                      `tfsdk:"device_id"`
	Profiles []deviceProfileDataSourceModel    `tfsdk:"profiles"`
}

type deviceProfileDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	Type             types.String `tfsdk:"type"`
	Name             types.String `tfsdk:"name"`
	ProfileIdentifier types.String `tfsdk:"profile_identifier"`
	UserScope        types.Bool   `tfsdk:"user_scope"`
	AttributeSupport types.Bool   `tfsdk:"attribute_support"`
}

type deviceRelatedItemDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	Type           types.String `tfsdk:"type"`
	AttributesJSON types.String `tfsdk:"attributes_json"`
}

func DeviceProfilesDataSource() datasource.DataSource {
	return &deviceProfilesDataSource{}
}

func (d *deviceProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_profiles"
}

func (d *deviceProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of profiles directly assigned to a device. Note: Profiles assigned through groups are not included.",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier of the device.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d+$`),
						"device_id must be a numeric string",
					),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"profiles": schema.ListNestedBlock{
				Description: "Collection of profiles applied directly to the device (excludes profiles assigned through groups).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Profile identifier.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Profile resource type.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Profile name.",
						},
						"profile_identifier": schema.StringAttribute{
							Computed:    true,
							Description: "Profile identifier string from the configuration profile.",
						},
						"user_scope": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the profile is user-scoped.",
						},
						"attribute_support": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the profile supports custom attributes.",
						},
					},
				},
			},
		},
	}
}

func (d *deviceProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deviceProfilesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	profiles, err := simplemdmext.ListDeviceProfiles(ctx, d.client, state.DeviceID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Device not found",
				fmt.Sprintf("The device with ID %s was not found. It may have been deleted.", state.DeviceID.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to list device profiles",
				fmt.Sprintf("Failed to retrieve profiles for device %s: %s", state.DeviceID.ValueString(), err.Error()),
			)
		}
		return
	}

	items := make([]deviceProfileDataSourceModel, 0, len(profiles.Data))
	for _, item := range profiles.Data {
		profile := deviceProfileDataSourceModel{
			ID:   types.StringValue(item.ID.String()),
			Type: types.StringValue(item.Type),
		}

		// Parse attributes from the map
		if name, ok := item.Attributes["name"].(string); ok {
			profile.Name = types.StringValue(name)
		} else {
			profile.Name = types.StringNull()
		}

		if profileIdentifier, ok := item.Attributes["profile_identifier"].(string); ok {
			profile.ProfileIdentifier = types.StringValue(profileIdentifier)
		} else {
			profile.ProfileIdentifier = types.StringNull()
		}

		if userScope, ok := item.Attributes["user_scope"].(bool); ok {
			profile.UserScope = types.BoolValue(userScope)
		} else {
			profile.UserScope = types.BoolNull()
		}

		if attributeSupport, ok := item.Attributes["attribute_support"].(bool); ok {
			profile.AttributeSupport = types.BoolValue(attributeSupport)
		} else {
			profile.AttributeSupport = types.BoolNull()
		}

		items = append(items, profile)
	}

	state.Profiles = items

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *deviceProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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