package provider

import (
	"context"
	"fmt"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &scriptJobDataSource{}
	_ datasource.DataSourceWithConfigure = &scriptJobDataSource{}
)

type scriptJobDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	JobName              types.String `tfsdk:"job_name"`
	JobIdentifier        types.String `tfsdk:"job_identifier"`
	Status               types.String `tfsdk:"status"`
	PendingCount         types.Int64  `tfsdk:"pending_count"`
	SuccessCount         types.Int64  `tfsdk:"success_count"`
	ErroredCount         types.Int64  `tfsdk:"errored_count"`
	ScriptName           types.String `tfsdk:"script_name"`
	CustomAttribute      types.String `tfsdk:"custom_attribute"`
	CustomAttributeRegex types.String `tfsdk:"custom_attribute_regex"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
	CreatedBy            types.String `tfsdk:"created_by"`
	VariableSupport      types.Bool   `tfsdk:"variable_support"`
	Content              types.String `tfsdk:"content"`
	Devices              types.List   `tfsdk:"devices"`
}

func ScriptJobDataSource() datasource.DataSource {
	return &scriptJobDataSource{}
}

type scriptJobDataSource struct {
	client *simplemdm.Client
}

func (d *scriptJobDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scriptjob"
}

func (d *scriptJobDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Script Job data source provides status information for an existing script job execution.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the script job.",
			},
			"job_name": schema.StringAttribute{
				Computed:    true,
				Description: "Human friendly name of the job.",
			},
			"job_identifier": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier reported by the SimpleMDM API for the job.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Current execution status of the job.",
			},
			"pending_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of devices that have not yet reported a result.",
			},
			"success_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of devices that completed successfully.",
			},
			"errored_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of devices that failed to execute the script.",
			},
			"script_name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the script that was executed.",
			},
			"custom_attribute": schema.StringAttribute{
				Computed:    true,
				Description: "Custom attribute that stores the job output, when configured.",
			},
			"custom_attribute_regex": schema.StringAttribute{
				Computed:    true,
				Description: "Regular expression used to filter the custom attribute output.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Creation timestamp returned by the API.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Last update timestamp returned by the API.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "User or API key that created the job.",
			},
			"variable_support": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether the script supports variables.",
			},
			"content": schema.StringAttribute{
				Computed:    true,
				Description: "Script contents that were executed by the job.",
			},
			"devices": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Execution results for each targeted device.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Device identifier.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Execution status reported for the device.",
						},
						"status_code": schema.StringAttribute{
							Computed:    true,
							Description: "Optional status code returned by the device.",
						},
						"response": schema.StringAttribute{
							Computed:    true,
							Description: "Output returned by the device, when available.",
						},
					},
				},
			},
		},
	}
}

func (d *scriptJobDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *scriptJobDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state scriptJobDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	details, err := fetchScriptJobDetails(ctx, d.client, state.ID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Unable to Read SimpleMDM script job",
				fmt.Sprintf("Script job %s was not found", state.ID.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to Read SimpleMDM script job",
				err.Error(),
			)
		}
		return
	}

	state.ID = types.StringValue(details.ID)
	state.JobName = stringValueOrNull(details.JobName)
	state.JobIdentifier = stringValueOrNull(details.JobIdentifier)
	state.Status = stringValueOrNull(details.Status)
	state.PendingCount = types.Int64Value(details.PendingCount)
	state.SuccessCount = types.Int64Value(details.SuccessCount)
	state.ErroredCount = types.Int64Value(details.ErroredCount)
	state.ScriptName = stringValueOrNull(details.ScriptName)

	if details.CustomAttribute != "" {
		state.CustomAttribute = types.StringValue(details.CustomAttribute)
	} else {
		state.CustomAttribute = types.StringNull()
	}

	if details.CustomAttributeRegex != "" {
		state.CustomAttributeRegex = types.StringValue(details.CustomAttributeRegex)
	} else {
		state.CustomAttributeRegex = types.StringNull()
	}

	state.CreatedAt = stringValueOrNull(details.CreatedAt)
	state.UpdatedAt = stringValueOrNull(details.UpdatedAt)
	state.CreatedBy = stringValueOrNull(details.CreatedBy)
	state.VariableSupport = types.BoolValue(details.VariableSupport)
	state.Content = stringValueOrNull(details.Content)

	devices, diags := scriptJobDevicesListValue(ctx, details.Devices)
	resp.Diagnostics.Append(diags...)
	state.Devices = devices

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
