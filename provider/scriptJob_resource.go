package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &scriptJobResource{}
	_ resource.ResourceWithConfigure = &scriptJobResource{}
)

// scriptJobsResourceModel maps the resource schema data.
type scriptJobResourceModel struct {
	ScriptId             types.String `tfsdk:"script_id"`
	DeviceIds            types.Set    `tfsdk:"device_ids"`
	GroupIds             types.Set    `tfsdk:"group_ids"`
	AssignmentGroupIds   types.Set    `tfsdk:"assignment_group_ids"`
	CustomAttribute      types.String `tfsdk:"custom_attribute"`
	CustomAttributeRegex types.String `tfsdk:"custom_attribute_regex"`
	ID                   types.String `tfsdk:"id"`
	JobName              types.String `tfsdk:"job_name"`
	JobIdentifier        types.String `tfsdk:"job_identifier"`
	Status               types.String `tfsdk:"status"`
	PendingCount         types.Int64  `tfsdk:"pending_count"`
	SuccessCount         types.Int64  `tfsdk:"success_count"`
	ErroredCount         types.Int64  `tfsdk:"errored_count"`
	ScriptName           types.String `tfsdk:"script_name"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
	CreatedBy            types.String `tfsdk:"created_by"`
	VariableSupport      types.Bool   `tfsdk:"variable_support"`
	Content              types.String `tfsdk:"content"`
	Devices              types.List   `tfsdk:"devices"`
}

// scriptJobResource is a helper function to simplify the provider implementation.
func ScriptJobResource() resource.Resource {
	return &scriptJobResource{}
}

// scriptJobResource is the resource implementation.
type scriptJobResource struct {
	client *simplemdm.Client
}

// Configure adds the provider configured client to the resource.
func (r *scriptJobResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*simplemdm.Client)
}

// Metadata returns the resource type name.
func (r *scriptJobResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scriptjob"
}

// Schema defines the schema for the resource.
func (r *scriptJobResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Script resource can be used to manage Scripts Jobs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "ID of a Script Job in SimpleMDM",
			},
			"script_id": schema.StringAttribute{
				Required:    true,
				Description: "Required. The ID of the script to be run on the devices",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device_ids": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A comma separated list of device IDs to run the script on. At least one of `device_ids`, `group_ids`, or `assignment_group_ids` must be provided.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Set{
					setvalidator.AtLeastOneOf(
						path.MatchRoot("device_ids"),
						path.MatchRoot("group_ids"),
						path.MatchRoot("assignment_group_ids"),
					),
				},
			},
			"group_ids": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A comma separated list of group IDs to run the script on. All macOS devices from these groups will be included. At least one of `device_ids`, `group_ids`, or `assignment_group_ids` must be provided.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Set{
					setvalidator.AtLeastOneOf(
						path.MatchRoot("device_ids"),
						path.MatchRoot("group_ids"),
						path.MatchRoot("assignment_group_ids"),
					),
				},
			},
			"assignment_group_ids": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A comma separated list of assignment group IDs to run the script on. All macOS devices from these assignment groups will be included At least one of `device_ids`, `group_ids`, or `assignment_group_ids` must be provided.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Set{
					setvalidator.AtLeastOneOf(
						path.MatchRoot("device_ids"),
						path.MatchRoot("group_ids"),
						path.MatchRoot("assignment_group_ids"),
					),
				},
			},
			"custom_attribute": schema.StringAttribute{
				Optional:    true,
				Description: "Optional. If provided the output from the script will be stored in this custom attribute on each device.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"custom_attribute_regex": schema.StringAttribute{
				Optional:    true,
				Description: "Optional. Used to sanitize the output from the script before storing it in the custom attribute. Can be left empty but \\n is recommended.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"job_name": schema.StringAttribute{
				Computed:    true,
				Description: "Human friendly name of the job.",
			},
			"job_identifier": schema.StringAttribute{
				Computed:    true,
				Description: "Short identifier string for the job (maps to API's 'job_id' field, different from the numeric ID).",
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

// Create a new resource
func (r *scriptJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan scriptJobResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceIDs, err := convertSetToSlice(ctx, plan.DeviceIds)
	if err != nil {
		resp.Diagnostics.AddError("Device IDs Conversion Error", "Failed to convert device IDs: "+err.Error())
		return
	}

	groupIDs, err := convertSetToSlice(ctx, plan.GroupIds)
	if err != nil {
		resp.Diagnostics.AddError("Group IDs Conversion Error", "Failed to convert group IDs: "+err.Error())
		return
	}

	assignmentGroupIDs, err := convertSetToSlice(ctx, plan.AssignmentGroupIds)
	if err != nil {
		resp.Diagnostics.AddError("Assignment Group IDs Conversion Error", "Failed to convert assignment group IDs: "+err.Error())
		return
	}

	var customAttribute string
	if !plan.CustomAttribute.IsNull() && !plan.CustomAttribute.IsUnknown() {
		customAttribute = plan.CustomAttribute.ValueString()
	}

	customAttributeRegex := ""
	if !plan.CustomAttributeRegex.IsNull() && !plan.CustomAttributeRegex.IsUnknown() {
		customAttributeRegex = plan.CustomAttributeRegex.ValueString()
	}

	// Note: Validation for at least one target being provided is already handled by schema validators
	// at lines 101-106, 115-121, and 130-136, so no additional validation is needed here.

	scriptJob, err := r.client.ScriptJobCreate(
		plan.ScriptId.ValueString(),
		deviceIDs,
		groupIDs,
		assignmentGroupIDs,
		customAttribute,
		customAttributeRegex,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating script job",
			"Could not create script job, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(scriptJob.Data.ID))

	plan.DeviceIds, diags = types.SetValueFrom(ctx, types.StringType, deviceIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.GroupIds, diags = types.SetValueFrom(ctx, types.StringType, groupIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.AssignmentGroupIds, diags = types.SetValueFrom(ctx, types.StringType, assignmentGroupIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	details, err := fetchScriptJobDetails(ctx, r.client, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving created script job",
			"Could not retrieve script job details: "+err.Error(),
		)
		return
	}

	applyScriptJobDetailsToResourceModel(ctx, details, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *scriptJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state scriptJobResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError(
			"Client not configured",
			"The SimpleMDM client was not configured.",
		)
		return
	}

	if err := r.client.ScriptCancelJob(state.ID.ValueString()); err != nil && !isNotFoundError(err) {
		resp.Diagnostics.AddError(
			"Error cancelling script job",
			fmt.Sprintf("Could not cancel SimpleMDM Script Job %s. Note: Jobs can only be canceled before devices receive the command. Error: %s",
				state.ID.ValueString(), err.Error()),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *scriptJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state scriptJobResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	details, err := fetchScriptJobDetails(ctx, r.client, state.ID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM Script Job",
			"Could not read SimpleMDM Script Job "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	applyScriptJobDetailsToResourceModel(ctx, details, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *scriptJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Force the recreation by showing an appropriate error
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Updating this resource is not supported. Please destroy and recreate the resource.",
	)
}

func (r *scriptJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import using the script job ID
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func applyScriptJobDetailsToResourceModel(ctx context.Context, details *scriptJobDetailsData, model *scriptJobResourceModel, diagnostics *diag.Diagnostics) {
	if details == nil || model == nil {
		return
	}

	model.ID = types.StringValue(details.ID)
	model.JobName = stringValueOrNull(details.JobName)
	model.JobIdentifier = stringValueOrNull(details.JobIdentifier)
	model.Status = stringValueOrNull(details.Status)
	model.PendingCount = types.Int64Value(details.PendingCount)
	model.SuccessCount = types.Int64Value(details.SuccessCount)
	model.ErroredCount = types.Int64Value(details.ErroredCount)
	model.ScriptName = stringValueOrNull(details.ScriptName)
	model.CreatedAt = stringValueOrNull(details.CreatedAt)
	model.UpdatedAt = stringValueOrNull(details.UpdatedAt)
	model.CreatedBy = stringValueOrNull(details.CreatedBy)
	model.VariableSupport = types.BoolValue(details.VariableSupport)
	model.Content = stringValueOrNull(details.Content)

	if details.CustomAttribute != "" {
		model.CustomAttribute = types.StringValue(details.CustomAttribute)
	} else {
		model.CustomAttribute = types.StringNull()
	}

	if details.CustomAttributeRegex != "" {
		model.CustomAttributeRegex = types.StringValue(details.CustomAttributeRegex)
	} else {
		model.CustomAttributeRegex = types.StringNull()
	}

	devices, diags := scriptJobDevicesListValue(ctx, details.Devices)
	if diagnostics != nil {
		diagnostics.Append(diags...)
	}
	model.Devices = devices
}

func convertSetToSlice(ctx context.Context, set types.Set) ([]string, error) {
	if set.IsNull() || set.IsUnknown() {
		return nil, nil
	}

	var result []string
	diags := set.ElementsAs(ctx, &result, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert set to slice: %s", diags.Errors()[0].Detail())
	}
	return result, nil
}
