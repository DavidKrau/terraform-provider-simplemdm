package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &attributeResource{}
	_ resource.ResourceWithConfigure   = &attributeResource{}
	_ resource.ResourceWithImportState = &attributeResource{}
)

// scriptJobsResourceModel maps the resource schema data.
type scriptJobResourceModel struct {
	ScriptId             types.String `tfsdk:"script_id"`
	DeviceIds            types.Set    `tfsdk:"device_ids"`
	AssignmentGroupIds   types.Set    `tfsdk:"assignment_group_ids"`
	CustomAttribute      types.String `tfsdk:"custom_attribute"`
	CustomAttributeRegex types.String `tfsdk:"custom_attribute_regex"`
	ID                   types.String `tfsdk:"id"`
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
				Required:    true,
				ElementType: types.StringType,
				Description: "A comma separated list of device IDs to run the script on. At least one of `device_ids`, `group_ids`, or `assignment_group_ids` must be provided.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"assignment_group_ids": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "A comma separated list of assignment group IDs to run the script on. All macOS devices from these assignment groups will be included At least one of `device_ids`, `group_ids`, or `assignment_group_ids` must be provided.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
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
				Description: "Optional. Used to sanitize the output from the script before storing it in the custom attribute. Can be left empty but \n is recommended.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// func (r *scriptJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }

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

	assignmentGroupIDs, err := convertSetToSlice(ctx, plan.AssignmentGroupIds)
	if err != nil {
		resp.Diagnostics.AddError("Assignment Group IDs Conversion Error", "Failed to convert assignment group IDs: "+err.Error())
		return
	}

	var customAttribute string
	if !plan.CustomAttribute.IsNull() {
		customAttribute = plan.CustomAttribute.String()
	}

	customAttributeRegex := ""
	if !plan.CustomAttributeRegex.IsNull() {
		customAttributeRegex = plan.CustomAttributeRegex.ValueString()
	}

	// Generate API request body from plan
	scriptJob, err := r.client.ScriptJobCreate(
		plan.ScriptId.ValueString(),
		deviceIDs,
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

	// Preserve input fields in state (since they are not returned by the API)
	plan.DeviceIds = types.SetValueMust(types.StringType, stringSliceToAttrValues(deviceIDs))
	plan.AssignmentGroupIds = types.SetValueMust(types.StringType, stringSliceToAttrValues(assignmentGroupIDs))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *scriptJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Delete doesn't have sense here, beause a job could be canceled but not deleted.
	// For now, the schedule of a job is not perimted by the API so do nothing here.
	var state scriptJobResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"Deleting this resource is not supported.",
	)

}

func (r *scriptJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state scriptJobResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call API to get the script job
	scriptJob, err := r.client.ScriptJobGet(state.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM Script Job",
			"Could not read SimpleMDM Script Job "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	// resp.Diagnostics.AddError(
	// 	"Error Reading SimpleMDM Script Job",
	// 	"Could not read SimpleMDM Script Job "+state.ID.ValueString(),
	// )
	// Update fields returned by the API
	//state.ScriptId = types.StringValue(scriptJob.Data.Attributes.ScriptName)
	state.ID = types.StringValue(strconv.Itoa(scriptJob.Data.ID))
	if scriptJob.Data.Attributes.CustomAttributeRegex != "" {
		state.CustomAttribute = types.StringValue(scriptJob.Data.Relationships.CustomAttribute.Data.ID)
	}
	if scriptJob.Data.Attributes.CustomAttributeRegex != "" {
		state.CustomAttributeRegex = types.StringValue(scriptJob.Data.Attributes.CustomAttributeRegex)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *scriptJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Force the recreation by seeing an appropriate error
	resp.Diagnostics.AddError(
		"Update Not Supported",
		"Updating this resource is not supported. Please destroy and recreate the resource.",
	)
}

func convertSetToSlice(ctx context.Context, set types.Set) ([]string, error) {
	if set.IsNull() || set.IsUnknown() {
		return nil, nil
	}

	var result []string
	set.ElementsAs(ctx, &result, false)
	return result, nil
}

func stringSliceToAttrValues(slice []string) []attr.Value {
	values := make([]attr.Value, len(slice))
	for i, s := range slice {
		values[i] = types.StringValue(s)
	}
	return values
}
