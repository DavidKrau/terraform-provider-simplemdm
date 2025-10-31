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
	_ datasource.DataSource              = &assignmentGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &assignmentGroupsDataSource{}
)

type assignmentGroupsDataSource struct {
	client *simplemdm.Client
}

type assignmentGroupsDataSourceModel struct {
	AssignmentGroups []assignmentGroupsDataSourceGroupModel `tfsdk:"assignment_groups"`
}

type assignmentGroupsDataSourceGroupModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	AutoDeploy       types.Bool   `tfsdk:"auto_deploy"`
	GroupType        types.String `tfsdk:"group_type"`
	Priority         types.Int64  `tfsdk:"priority"`
	AppTrackLocation types.Bool   `tfsdk:"app_track_location"`
	DeviceCount      types.Int64  `tfsdk:"device_count"`
	GroupCount       types.Int64  `tfsdk:"group_count"`
}

func AssignmentGroupsDataSource() datasource.DataSource {
	return &assignmentGroupsDataSource{}
}

func (d *assignmentGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_assignmentgroups"
}

func (d *assignmentGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the collection of assignment groups from your SimpleMDM account.",
		Blocks: map[string]schema.Block{
			"assignment_groups": schema.ListNestedBlock{
				Description: "Collection of assignment group records returned by the API.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Assignment group identifier.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the assignment group.",
						},
						"auto_deploy": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the assignment group automatically deploys apps.",
						},
						"group_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of assignment group (standard or munki).",
						},
						"priority": schema.Int64Attribute{
							Computed:    true,
							Description: "The priority of the assignment group.",
						},
						"app_track_location": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the SimpleMDM app tracks device location when installed for this assignment group.",
						},
						"device_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of devices currently assigned to the assignment group.",
						},
						"group_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of device groups currently assigned to the assignment group.",
						},
					},
				},
			},
		},
	}
}

func (d *assignmentGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config assignmentGroupsDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	groups, err := fetchAllAssignmentGroups(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list SimpleMDM assignment groups",
			err.Error(),
		)
		return
	}

	entries := make([]assignmentGroupsDataSourceGroupModel, 0, len(groups))
	for _, group := range groups {
		entry := assignmentGroupsDataSourceGroupModel{
			ID:               types.StringValue(strconv.Itoa(group.ID)),
			Name:             types.StringValue(group.Attributes.Name),
			AutoDeploy:       types.BoolValue(group.Attributes.AutoDeploy),
			GroupType:        types.StringValue(group.Attributes.GroupType),
			Priority:         types.Int64Value(int64(group.Attributes.Priority)),
			AppTrackLocation: types.BoolValue(group.Attributes.AppTrackLocation),
			DeviceCount:      types.Int64Value(int64(group.Attributes.DeviceCount)),
			GroupCount:       types.Int64Value(int64(group.Attributes.GroupCount)),
		}

		entries = append(entries, entry)
	}

	config.AssignmentGroups = entries

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func (d *assignmentGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// fetchAllAssignmentGroups retrieves all assignment groups with pagination support
func fetchAllAssignmentGroups(ctx context.Context, client *simplemdm.Client) ([]assignmentGroupData, error) {
	var allGroups []assignmentGroupData
	startingAfter := 0
	limit := 100

	for {
		url := fmt.Sprintf("https://%s/api/v1/assignment_groups?limit=%d", client.HostName, limit)
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

		var response assignmentGroupsAPIResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}

		allGroups = append(allGroups, response.Data...)

		if !response.HasMore {
			break
		}

		if len(response.Data) > 0 {
			startingAfter = response.Data[len(response.Data)-1].ID
		} else {
			break
		}
	}

	return allGroups, nil
}

// assignmentGroupsAPIResponse represents the paginated API response for assignment groups list
type assignmentGroupsAPIResponse struct {
	Data    []assignmentGroupData `json:"data"`
	HasMore bool                  `json:"has_more"`
}

// assignmentGroupData represents a single assignment group in the list response
type assignmentGroupData struct {
	ID         int                           `json:"id"`
	Type       string                        `json:"type"`
	Attributes assignmentGroupDataAttributes `json:"attributes"`
}

type assignmentGroupDataAttributes struct {
	Name             string `json:"name"`
	AutoDeploy       bool   `json:"auto_deploy"`
	GroupType        string `json:"group_type"`
	Priority         int    `json:"priority"`
	AppTrackLocation bool   `json:"app_track_location"`
	DeviceCount      int    `json:"device_count"`
	GroupCount       int    `json:"group_count"`
}
