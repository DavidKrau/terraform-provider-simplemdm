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
	_ datasource.DataSource              = &deviceGroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceGroupsDataSource{}
)

type deviceGroupsDataSource struct {
	client *simplemdm.Client
}

type deviceGroupsDataSourceModel struct {
	DeviceGroups []deviceGroupsDataSourceGroupModel `tfsdk:"device_groups"`
}

type deviceGroupsDataSourceGroupModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func DeviceGroupsDataSource() datasource.DataSource {
	return &deviceGroupsDataSource{}
}

func (d *deviceGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devicegroups"
}

func (d *deviceGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "⚠️ DEPRECATED: Device Groups have been superseded by Assignment Groups in SimpleMDM. " +
			"Please use the simplemdm_assignmentgroups data source instead. " +
			"This data source is maintained for backward compatibility only. " +
			"Fetches the collection of device groups from your SimpleMDM account.",
		DeprecationMessage: "Device Groups are deprecated. Use simplemdm_assignmentgroups instead.",
		Blocks: map[string]schema.Block{
			"device_groups": schema.ListNestedBlock{
				Description: "Collection of device group records returned by the API.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Device group identifier.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Device group name.",
						},
					},
				},
			},
		},
	}
}

func (d *deviceGroupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config deviceGroupsDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceGroups, err := fetchAllDeviceGroups(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list SimpleMDM device groups",
			err.Error(),
		)
		return
	}

	entries := make([]deviceGroupsDataSourceGroupModel, 0, len(deviceGroups))
	for _, group := range deviceGroups {
		entry := deviceGroupsDataSourceGroupModel{
			ID:   types.StringValue(strconv.Itoa(group.ID)),
			Name: types.StringValue(group.Attributes.Name),
		}

		entries = append(entries, entry)
	}

	config.DeviceGroups = entries

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func (d *deviceGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// fetchAllDeviceGroups retrieves all device groups with pagination support
func fetchAllDeviceGroups(ctx context.Context, client *simplemdm.Client) ([]deviceGroupData, error) {
	var allGroups []deviceGroupData
	startingAfter := 0
	limit := 100

	for {
		url := fmt.Sprintf("https://%s/api/v1/device_groups?limit=%d", client.HostName, limit)
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

		var response deviceGroupsAPIResponse
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

// deviceGroupsAPIResponse represents the paginated API response
type deviceGroupsAPIResponse struct {
	Data    []deviceGroupData `json:"data"`
	HasMore bool              `json:"has_more"`
}

// deviceGroupData represents a single device group in the list response
type deviceGroupData struct {
	ID         int                   `json:"id"`
	Type       string                `json:"type"`
	Attributes deviceGroupAttributes `json:"attributes"`
}

type deviceGroupAttributes struct {
	Name string `json:"name"`
}
