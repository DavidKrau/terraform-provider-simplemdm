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
	_ datasource.DataSource              = &deviceGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceGroupDataSource{}
)

// deviceGroupDataSourceModel maps the data source schema data.
type deviceGroupDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// DeviceGroupDataSource is a helper function to simplify the provider implementation.
func DeviceGroupDataSource() datasource.DataSource {
	return &deviceGroupDataSource{}
}

// deviceGroupDataSource is the data source implementation.
type deviceGroupDataSource struct {
	client *simplemdm.Client
}

// Metadata returns the data source type name.
func (d *deviceGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devicegroup"
}

// Schema defines the schema for the data source.
func (d *deviceGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "⚠️ DEPRECATED: Device Groups have been superseded by Assignment Groups in SimpleMDM. " +
			"Please use the simplemdm_assignmentgroup data source instead. " +
			"This data source is maintained for backward compatibility only. " +
			"Device Group data source can be used together with Assignment Group(s) to assign group(s) to these objects.",
		DeprecationMessage: "Device Groups are deprecated. Use simplemdm_assignmentgroup instead.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Device Group name.",
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of a Device Group in SimpleMDM.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *deviceGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deviceGroupDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	devicegroup, err := d.client.DeviceGroupGet(state.ID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"SimpleMDM device group not found",
				fmt.Sprintf("The device group with ID %s was not found. "+
					"Note: Only legacy device group IDs from migrated groups are supported. "+
					"Use simplemdm_assignmentgroup for current group functionality.",
					state.ID.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to Read SimpleMDM device group",
				"Could not read device group "+state.ID.ValueString()+": "+err.Error()+". "+
					"Ensure you are using a valid legacy device group ID.",
			)
		}
		return
	}

	// Map response body to model
	state.Name = types.StringValue(devicegroup.Data.Attributes.Name)
	state.ID = types.StringValue(strconv.Itoa(devicegroup.Data.ID))

	// Set state

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *deviceGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
