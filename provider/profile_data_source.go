package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
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
	ReinstallAfterOSUpdate types.Bool   `tfsdk:"reinstall_after_os_update"`
	ProfileIdentifier      types.String `tfsdk:"profile_identifier"`
	UserScope              types.Bool   `tfsdk:"user_scope"`
	GroupCount             types.Int64  `tfsdk:"group_count"`
	DeviceCount            types.Int64  `tfsdk:"device_count"`
	GroupIDs               types.Set    `tfsdk:"group_ids"`
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
		Description: "Profile data source can be used to reference existing configuration profiles in SimpleMDM. Profiles represent both built-in and custom configuration profiles.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Profile.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The profile payload type (e.g., 'apn', 'email', 'app_restrictions', 'custom_configuration_profile').",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the Profile.",
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
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"SimpleMDM profile not found",
				fmt.Sprintf("The profile with ID %s was not found. It may have been deleted.", state.ID.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to Read SimpleMDM profile",
				err.Error(),
			)
		}
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

	// Map response body to model with consistent null handling
	state.ID = types.StringValue(strconv.Itoa(profile.Data.ID))
	
	if profile.Data.Type != "" {
		state.Type = types.StringValue(profile.Data.Type)
	} else {
		state.Type = types.StringNull()
	}
	
	state.Name = types.StringValue(profile.Data.Attributes.Name)
	state.ReinstallAfterOSUpdate = types.BoolValue(profile.Data.Attributes.ReinstallAfterOSUpdate)
	state.ProfileIdentifier = types.StringValue(profile.Data.Attributes.ProfileIdentifier)
	state.UserScope = types.BoolValue(profile.Data.Attributes.UserScope)
	state.GroupCount = types.Int64Value(int64(profile.Data.Attributes.GroupCount))
	state.DeviceCount = types.Int64Value(int64(profile.Data.Attributes.DeviceCount))
	state.GroupIDs = groupIDs

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

// Helper types and functions for profile API interaction

type profileAPIResponse struct {
	Data struct {
		Type          string               `json:"type"`
		ID            int                  `json:"id"`
		Attributes    profileAttributes    `json:"attributes"`
		Relationships profileRelationships `json:"relationships"`
	} `json:"data"`
}

type profileAttributes struct {
	Name                   string `json:"name"`
	ReinstallAfterOSUpdate bool   `json:"reinstall_after_os_update"`
	ProfileIdentifier      string `json:"profile_identifier"`
	UserScope              bool   `json:"user_scope"`
	GroupCount             int    `json:"group_count"`
	DeviceCount            int    `json:"device_count"`
}

type profileRelationships struct {
	DeviceGroups relationshipCollection `json:"device_groups"`
	Groups       relationshipCollection `json:"groups"`
}

type relationshipCollection struct {
	Data []relationshipReference `json:"data"`
}

type relationshipReference struct {
	ID int `json:"id"`
}

func fetchProfile(ctx context.Context, client *simplemdm.Client, profileID string) (*profileAPIResponse, error) {
	url := fmt.Sprintf("https://%s/api/v1/profiles/%s", client.HostName, profileID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := client.RequestResponse200(req)
	if err != nil {
		return nil, err
	}

	var profile profileAPIResponse
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func convertGroupIDs(ctx context.Context, relationships profileRelationships) (types.Set, error) {
	unique := make(map[int]struct{})
	for _, item := range relationships.DeviceGroups.Data {
		unique[item.ID] = struct{}{}
	}
	for _, item := range relationships.Groups.Data {
		unique[item.ID] = struct{}{}
	}

	if len(unique) == 0 {
		return types.SetNull(types.StringType), nil
	}

	ids := make([]string, 0, len(unique))
	for id := range unique {
		ids = append(ids, strconv.Itoa(id))
	}
	sort.Strings(ids)

	value, diags := types.SetValueFrom(ctx, types.StringType, ids)
	if diags.HasError() {
		return types.SetNull(types.StringType), fmt.Errorf("unable to convert profile group IDs: %s", diags)
	}

	return value, nil
}