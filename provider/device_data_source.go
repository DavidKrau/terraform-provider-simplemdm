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

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &deviceDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceDataSource{}
)

// deviceDataSourceModel maps the data source schema data.
type deviceDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	DeviceName    types.String `tfsdk:"devicename"`
	DeviceGroup   types.String `tfsdk:"devicegroup"`
	EnrollmentURL types.String `tfsdk:"enrollmenturl"`
	Details       types.Map    `tfsdk:"details"`
}

// deviceDataSource is a helper function to simplify the provider implementation.
func DeviceDataSource() datasource.DataSource {
	return &deviceDataSource{}
}

// deviceDataSource is the data source implementation.
type deviceDataSource struct {
	client *simplemdm.Client
}

// Metadata returns the data source type name.
func (d *deviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

// Schema defines the schema for the data source.
func (d *deviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Device data source can be used together Assignment Group(s) to assign device to these objects.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The SimpleMDM name of the device.",
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the device.",
			},
			"devicename": schema.StringAttribute{
				Computed:    true,
				Description: "The hostname reported by the device.",
			},
			"devicegroup": schema.StringAttribute{
				Computed:    true,
				Description: "Device group identifier for the device.",
			},
			"enrollmenturl": schema.StringAttribute{
				Computed:    true,
				Description: "Enrollment URL generated for the device, when available.",
			},
			"details": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Full set of attributes returned by the SimpleMDM API for the device.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *deviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deviceDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	device, err := simplemdmext.GetDevice(ctx, d.client, state.ID.ValueString(), true)
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"SimpleMDM device not found",
				fmt.Sprintf("The device with ID %s was not found. It may have been deleted.", state.ID.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to Read SimpleMDM device",
				err.Error(),
			)
		}
		return
	}

	// Map response body to model
	state.ID = types.StringValue(strconv.Itoa(device.Data.ID))
	flatAttributes := simplemdmext.FlattenAttributes(device.Data.Attributes)
	if name, ok := flatAttributes["name"]; ok && name != "" {
		state.Name = types.StringValue(name)
	}
	if deviceName, ok := flatAttributes["device_name"]; ok && deviceName != "" {
		state.DeviceName = types.StringValue(deviceName)
	} else {
		state.DeviceName = types.StringNull()
	}
	if enrollmentURL, ok := flatAttributes["enrollment_url"]; ok && enrollmentURL != "" && enrollmentURL != "null" {
		state.EnrollmentURL = types.StringValue(enrollmentURL)
	} else {
		state.EnrollmentURL = types.StringNull()
	}
	if groupID := device.Data.Relationships.DeviceGroup.Data.ID; groupID != 0 {
		state.DeviceGroup = types.StringValue(strconv.Itoa(groupID))
	} else {
		state.DeviceGroup = types.StringNull()
	}

	detailsValue, detailsDiags := types.MapValueFrom(ctx, types.StringType, flatAttributes)
	resp.Diagnostics.Append(detailsDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Details = detailsValue

	// Set state

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *deviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
