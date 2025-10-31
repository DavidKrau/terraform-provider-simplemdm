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
	_ datasource.DataSource              = &profilesDataSource{}
	_ datasource.DataSourceWithConfigure = &profilesDataSource{}
)

type profilesDataSource struct {
	client *simplemdm.Client
}

type profilesDataSourceModel struct {
	Profiles []profilesDataSourceProfileModel `tfsdk:"profiles"`
}

type profilesDataSourceProfileModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	ProfileIdentifier      types.String `tfsdk:"profile_identifier"`
	UserScope              types.Bool   `tfsdk:"user_scope"`
	AttributeSupport       types.Bool   `tfsdk:"attribute_support"`
	EscapeAttributes       types.Bool   `tfsdk:"escape_attributes"`
	ReinstallAfterOSUpdate types.Bool   `tfsdk:"reinstall_after_os_update"`
	GroupCount             types.Int64  `tfsdk:"group_count"`
	DeviceCount            types.Int64  `tfsdk:"device_count"`
}

func ProfilesDataSource() datasource.DataSource {
	return &profilesDataSource{}
}

func (d *profilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profiles"
}

func (d *profilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the collection of configuration profiles from your SimpleMDM account.",
		Blocks: map[string]schema.Block{
			"profiles": schema.ListNestedBlock{
				Description: "Collection of profile records returned by the API.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Profile identifier.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the profile.",
						},
						"profile_identifier": schema.StringAttribute{
							Computed:    true,
							Description: "The identifier contained within the profile payload.",
						},
						"user_scope": schema.BoolAttribute{
							Computed:    true,
							Description: "Indicates if the profile installs in the user scope.",
						},
						"attribute_support": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the profile supports attribute substitution.",
						},
						"escape_attributes": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether attribute values are escaped during installation.",
						},
						"reinstall_after_os_update": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the profile reinstalls automatically after macOS updates.",
						},
						"group_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of device groups currently assigned to the profile.",
						},
						"device_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of devices currently assigned to the profile.",
						},
					},
				},
			},
		},
	}
}

func (d *profilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config profilesDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	profiles, err := fetchAllProfiles(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list SimpleMDM profiles",
			err.Error(),
		)
		return
	}

	entries := make([]profilesDataSourceProfileModel, 0, len(profiles))
	for _, profile := range profiles {
		entry := profilesDataSourceProfileModel{
			ID:                     types.StringValue(strconv.Itoa(profile.ID)),
			Name:                   types.StringValue(profile.Attributes.Name),
			ProfileIdentifier:      types.StringValue(profile.Attributes.ProfileIdentifier),
			UserScope:              types.BoolValue(profile.Attributes.UserScope),
			AttributeSupport:       types.BoolValue(profile.Attributes.AttributeSupport),
			EscapeAttributes:       types.BoolValue(profile.Attributes.EscapeAttributes),
			ReinstallAfterOSUpdate: types.BoolValue(profile.Attributes.ReinstallAfterOSUpdate),
			GroupCount:             types.Int64Value(int64(profile.Attributes.GroupCount)),
			DeviceCount:            types.Int64Value(int64(profile.Attributes.DeviceCount)),
		}

		entries = append(entries, entry)
	}

	config.Profiles = entries

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func (d *profilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// fetchAllProfiles retrieves all profiles with pagination support
func fetchAllProfiles(ctx context.Context, client *simplemdm.Client) ([]profileDataList, error) {
	var allProfiles []profileDataList
	startingAfter := 0
	limit := 100

	for {
		url := fmt.Sprintf("https://%s/api/v1/profiles?limit=%d", client.HostName, limit)
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

		var response profilesAPIResponse
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

// profilesAPIResponse represents the paginated API response
type profilesAPIResponse struct {
	Data    []profileDataList `json:"data"`
	HasMore bool              `json:"has_more"`
}

// profileDataList represents a single profile in the list response
type profileDataList struct {
	ID         int                   `json:"id"`
	Type       string                `json:"type"`
	Attributes profileListAttributes `json:"attributes"`
}

type profileListAttributes struct {
	Name                   string `json:"name"`
	ProfileIdentifier      string `json:"profile_identifier"`
	UserScope              bool   `json:"user_scope"`
	AttributeSupport       bool   `json:"attribute_support"`
	EscapeAttributes       bool   `json:"escape_attributes"`
	ReinstallAfterOSUpdate bool   `json:"reinstall_after_os_update"`
	GroupCount             int    `json:"group_count"`
	DeviceCount            int    `json:"device_count"`
}
