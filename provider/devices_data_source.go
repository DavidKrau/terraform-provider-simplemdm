package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/DavidKrau/terraform-provider-simplemdm/internal/simplemdmext"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &devicesDataSource{}
	_ datasource.DataSourceWithConfigure = &devicesDataSource{}
)

type devicesDataSource struct {
	client *simplemdm.Client
}

type devicesDataSourceModel struct {
	Search                        types.String                   `tfsdk:"search"`
	IncludeAwaitingEnrollment     types.Bool                     `tfsdk:"include_awaiting_enrollment"`
	IncludeSecretCustomAttributes types.Bool                     `tfsdk:"include_secret_custom_attributes"`
	Devices                       []devicesDataSourceDeviceModel `tfsdk:"devices"`
}

type devicesDataSourceDeviceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	DeviceName    types.String `tfsdk:"device_name"`
	Status        types.String `tfsdk:"status"`
	DeviceGroupID types.String `tfsdk:"device_group_id"`
	Details       types.Map    `tfsdk:"details"`
}

func DevicesDataSource() datasource.DataSource {
	return &devicesDataSource{}
}

func (d *devicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices"
}

func (d *devicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the collection of devices exposed by the SimpleMDM API.",
		Attributes: map[string]schema.Attribute{
			"search": schema.StringAttribute{
				Optional:    true,
				Description: "Optional filter to restrict the device results by identifier or name.",
			},
			"include_awaiting_enrollment": schema.BoolAttribute{
				Optional:    true,
				Description: "Include devices that are still awaiting enrollment in the response.",
			},
			"include_secret_custom_attributes": schema.BoolAttribute{
				Optional:    true,
				Description: "Request secret custom attribute values from the API.",
			},
		},
		Blocks: map[string]schema.Block{
			"devices": schema.ListNestedBlock{
				Description: "Collection of device records returned by the API.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Device identifier.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "SimpleMDM display name for the device.",
						},
						"device_name": schema.StringAttribute{
							Computed:    true,
							Description: "Hostname reported by the device.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Current enrollment status reported by SimpleMDM.",
						},
						"device_group_id": schema.StringAttribute{
							Computed:    true,
							Description: "Device group identifier for the device.",
						},
						"details": schema.MapAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Flattened device attribute payload for additional inspection.",
						},
					},
				},
			},
		},
	}
}

func (d *devicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config devicesDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	devices, err := simplemdmext.ListDevices(ctx, d.client, config.Search.ValueString(), config.IncludeAwaitingEnrollment.ValueBool(), config.IncludeSecretCustomAttributes.ValueBool())
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"SimpleMDM devices not found",
				"No devices were found matching the search criteria.",
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to list SimpleMDM devices",
				err.Error(),
			)
		}
		return
	}

	entries := make([]devicesDataSourceDeviceModel, 0, len(devices))
	for _, item := range devices {
		attributes := simplemdmext.FlattenAttributes(item.Attributes)
		detailsValue, detailsDiags := types.MapValueFrom(ctx, types.StringType, attributes)
		resp.Diagnostics.Append(detailsDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		entry := devicesDataSourceDeviceModel{
			ID:            types.StringValue(strconv.Itoa(item.ID)),
			Name:          types.StringNull(),
			DeviceName:    types.StringNull(),
			Status:        types.StringNull(),
			DeviceGroupID: types.StringNull(),
			Details:       detailsValue,
		}

		if name, ok := attributes["name"]; ok && name != "" {
			entry.Name = types.StringValue(name)
		}

		if deviceName, ok := attributes["device_name"]; ok && deviceName != "" {
			entry.DeviceName = types.StringValue(deviceName)
		}

		if status, ok := attributes["status"]; ok && status != "" {
			entry.Status = types.StringValue(status)
		}

		if groupID := item.Relationships.DeviceGroup.Data.ID; groupID != 0 {
			entry.DeviceGroupID = types.StringValue(strconv.Itoa(groupID))
		}

		entries = append(entries, entry)
	}

	config.Devices = entries

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func (d *devicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
