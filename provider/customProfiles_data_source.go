package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &customProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &customProfilesDataSource{}
)

type customProfilesDataSource struct {
	client *simplemdm.Client
}

type customProfilesDataSourceModel struct {
	CustomProfiles []customProfilesDataSourceProfileModel `tfsdk:"custom_profiles"`
}

type customProfilesDataSourceProfileModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	UserScope              types.Bool   `tfsdk:"user_scope"`
	AttributeSupport       types.Bool   `tfsdk:"attribute_support"`
	EscapeAttributes       types.Bool   `tfsdk:"escape_attributes"`
	ReinstallAfterOSUpdate types.Bool   `tfsdk:"reinstall_after_os_update"`
	ProfileIdentifier      types.String `tfsdk:"profile_identifier"`
	GroupCount             types.Int64  `tfsdk:"group_count"`
	DeviceCount            types.Int64  `tfsdk:"device_count"`
}

func CustomProfilesDataSource() datasource.DataSource {
	return &customProfilesDataSource{}
}

func (d *customProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customprofiles"
}

func (d *customProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the collection of custom configuration profiles from your SimpleMDM account.",
		Blocks: map[string]schema.Block{
			"custom_profiles": schema.ListNestedBlock{
				Description: "Collection of custom profile records returned by the API.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Custom profile identifier.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the custom profile.",
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
					},
				},
			},
		},
	}
}

func (d *customProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config customProfilesDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	customProfiles, err := fetchAllCustomProfiles(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list SimpleMDM custom profiles",
			err.Error(),
		)
		return
	}

	entries := make([]customProfilesDataSourceProfileModel, 0, len(customProfiles))
	for _, profile := range customProfiles {
		entry := customProfilesDataSourceProfileModel{
			ID:                     types.StringValue(strconv.Itoa(profile.ID)),
			Name:                   types.StringValue(profile.Attributes.Name),
			UserScope:              types.BoolValue(profile.Attributes.UserScope),
			AttributeSupport:       types.BoolValue(profile.Attributes.AttributeSupport),
			EscapeAttributes:       types.BoolValue(profile.Attributes.EscapeAttributes),
			ReinstallAfterOSUpdate: types.BoolValue(profile.Attributes.ReinstallAfterOsUpdate),
			ProfileIdentifier:      stringValueOrNull(profile.Attributes.ProfileIdentifier),
			GroupCount:             types.Int64Value(int64(profile.Attributes.GroupCount)),
			DeviceCount:            types.Int64Value(int64(profile.Attributes.DeviceCount)),
		}

		entries = append(entries, entry)
	}

	config.CustomProfiles = entries

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func (d *customProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// fetchAllCustomProfiles retrieves all custom profiles with pagination support
func fetchAllCustomProfiles(ctx context.Context, client *simplemdm.Client) ([]customProfileData, error) {
	var allProfiles []customProfileData
	startingAfter := 0
	limit := 100

	for {
		// Use correct API endpoint: /api/v1/custom_configuration_profiles
		url := fmt.Sprintf("https://%s/api/v1/custom_configuration_profiles?limit=%d", client.HostName, limit)
		if startingAfter > 0 {
			url += fmt.Sprintf("&starting_after=%d", startingAfter)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		body, err := client.RequestResponse200(req)
		if err != nil {
			return nil, err
		}

		var response customProfilesAPIResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}

		allProfiles = append(allProfiles, response.Data...)

		if !response.HasMore {
			break
		}

		if len(response.Data) > 0 {
			startingAfter = response.Data[len(response.Data)-1].ID
		} else {
			break
		}
	}

	return allProfiles, nil
}

// customProfilesAPIResponse represents the paginated API response for custom profiles list
type customProfilesAPIResponse struct {
	Data    []customProfileData `json:"data"`
	HasMore bool                `json:"has_more"`
}

// customProfileData represents a single custom profile in the list response
type customProfileData struct {
	ID         int                     `json:"id"`
	Type       string                  `json:"type"`
	Attributes customProfileAttributes `json:"attributes"`
}

type customProfileAttributes struct {
	Name                   string `json:"name"`
	UserScope              bool   `json:"user_scope"`
	AttributeSupport       bool   `json:"attribute_support"`
	EscapeAttributes       bool   `json:"escape_attributes"`
	ReinstallAfterOsUpdate bool   `json:"reinstall_after_os_update"`
	ProfileIdentifier      string `json:"profile_identifier"`
	GroupCount             int    `json:"group_count"`
	DeviceCount            int    `json:"device_count"`
}
