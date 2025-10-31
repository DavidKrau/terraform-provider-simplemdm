package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/DavidKrau/terraform-provider-simplemdm/internal/simplemdmext"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &deviceResource{}
	_ resource.ResourceWithConfigure   = &deviceResource{}
	_ resource.ResourceWithImportState = &deviceResource{}
)

// deviceGroupResourceModel maps the resource schema data.
type deviceResourceModel struct {
	Name           types.String `tfsdk:"name"`
	ID             types.String `tfsdk:"id"`
	Attributes     types.Map    `tfsdk:"attributes"`
	CustomProfiles types.Set    `tfsdk:"customprofiles"`
	Profiles       types.Set    `tfsdk:"profiles"`
	DeviceGroup    types.String `tfsdk:"devicegroup"`
	DeviceName     types.String `tfsdk:"devicename"`
	EnrollmentURL  types.String `tfsdk:"enrollmenturl"`
	Details        types.Map    `tfsdk:"details"`
}

// deviceGroupResource is a helper function to simplify the provider implementation.
func DeviceResource() resource.Resource {
	return &deviceResource{}
}

// deviceGroupResource is the resource implementation.
type deviceResource struct {
	client *simplemdm.Client
}

// Configure adds the provider configured client to the resource.
func (r *deviceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*simplemdm.Client)
}

// Metadata returns the resource type name.
func (r *deviceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

// Schema defines the schema for the resource.
func (r *deviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Device resource can be used to manage Device. Can be used together with Custom Profile(s), Attribute(s), Assignment Group(s) or Device Group(s) and set addition details regarding Device.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Required. The SimpleMDM name of the device.",
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The ID of the Device in SimpleMDM",
			},
			"profiles": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Optional. List of Configuration Profiles assigned to this Device",
			},
			"customprofiles": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Optional. List of Custom Configuration Profiles assigned to this Device",
			},
			"attributes": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "The name of the Assignment Group.",
			},
			"devicegroup": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "The ID of Device Group where device will be assigned.",
			},
			"devicename": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: "The Device name (localhost name) of the device.",
			},
			"enrollmenturl": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "SimpleMDM enrollment URL is generated when new device is created via API.",
			},
			"details": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "Full set of attributes returned by the SimpleMDM device record.",
			},
		},
	}
}

func (r *deviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to state
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource
func (r *deviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan deviceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	device, err := r.client.DeviceCreate(plan.Name.ValueString(), plan.DeviceGroup.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating device",
			"Could not create device, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(device.Data.ID))
	plan.EnrollmentURL = types.StringValue(device.Data.Attributes.EnrollmentURL)

	//setting attributes
	for attribute, value := range plan.Attributes.Elements() {
		err := r.client.AttributeSetAttributeForDevice(plan.ID.ValueString(), attribute, value.(types.String).ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating device attribute",
				"Could not set attribute value for device, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Assign all custom profiles in plan
	for _, profileId := range plan.CustomProfiles.Elements() {
		err := r.client.CustomProfileAssignToDevice(profileId.(types.String).ValueString(), plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating device profile assignment",
				"Could not update device profile assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Assign all custom profiles in plan
	for _, profileId := range plan.Profiles.Elements() {
		err := r.client.ProfileAssignToDevice(profileId.(types.String).ValueString(), plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating device profile assignment",
				"Could not update device profile assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Refresh state from API to populate computed attributes and relationships
	apiDevice, err := simplemdmext.GetDevice(ctx, r.client, plan.ID.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM device",
			"Could not read SimpleMDM device "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(r.assignAPIValues(ctx, apiDevice, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *deviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state deviceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiDevice, err := simplemdmext.GetDevice(ctx, r.client, state.ID.ValueString(), true)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM device",
			"Could not read SimpleMDM device "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(r.assignAPIValues(ctx, apiDevice, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *deviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Retrieve values from plan
	var plan, state deviceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	_, err := r.client.DeviceUpdate(plan.ID.ValueString(), plan.Name.ValueString(), plan.DeviceName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating device",
			"Could not update device, unexpected error: "+err.Error(),
		)
		return
	}

	//assign device to correct group
	err2 := r.client.DeviceGroupAssignDevice(plan.ID.ValueString(), plan.DeviceGroup.ValueString())
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error updating device",
			"Could not update device, unexpected error: "+err2.Error(),
		)
		return
	}
	//comparing planed attributes and their values to attributes in SimpleMDM
	for planAttribute, planValue := range plan.Attributes.Elements() {
		found := false
		for stateAttribute, stateValue := range state.Attributes.Elements() {
			if planAttribute == stateAttribute {
				found = true
				if planValue != stateValue {
					err := r.client.AttributeSetAttributeForDevice(plan.ID.ValueString(), planAttribute, planValue.(types.String).ValueString())
					if err != nil {
						resp.Diagnostics.AddError(
							"Error updating SimpleMDM device attributes value",
							"Could not update SimpleMDM device attributes value "+plan.ID.ValueString()+": "+err.Error(),
						)
						return
					}
				}
				break
			}
		}
		if !found {
			err := r.client.AttributeSetAttributeForDevice(plan.ID.ValueString(), planAttribute, planValue.(types.String).ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error updating SimpleMDM device attributes value",
					"Could not update SimpleMDM device attributes value "+plan.ID.ValueString()+": "+err.Error(),
				)
				return
			}
		}
	}

	//comparing attributes from SimpleMDM to the plan to find attributes set manually in MDM
	for stateAttribute := range state.Attributes.Elements() {
		found := false
		for planAttribute := range plan.Attributes.Elements() {
			if stateAttribute == planAttribute {
				found = true
				break
			}
		}
		if !found {
			err := r.client.AttributeSetAttributeForDevice(plan.ID.ValueString(), stateAttribute, "")
			if err != nil {
				resp.Diagnostics.AddError(
					"Error updating SimpleMDM device attributes value",
					"Could not update SimpleMDM device attributes value "+plan.ID.ValueString()+": "+err.Error(),
				)
				return
			}
		}
	}

	//Handling assigned profiles
	//reading assigned profiles from simpleMDM
	stateProfiles := []string{}
	for _, profileId := range state.Profiles.Elements() {
		stateProfiles = append(stateProfiles, profileId.(types.String).ValueString())
	}

	//reading configured profiles from TF file
	planProfiles := []string{}
	for _, profileId := range plan.Profiles.Elements() {
		planProfiles = append(planProfiles, profileId.(types.String).ValueString())
	}

	// // creating diff
	profilesToAdd, profilesToRemove := diffFunction(stateProfiles, planProfiles)

	// //adding profiles
	for _, profileId := range profilesToAdd {
		err := r.client.ProfileAssignToDevice(profileId, plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating device custom profile assignment",
				"Could not update device custom profile assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	//removing profiles
	for _, profileId := range profilesToRemove {
		err := r.client.ProfileUnAssignToDevice(profileId, plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating device custom profile assignment",
				"Could not update device custom profile assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	//Handling assigned custom prfiles profiles
	//reading assigned profiles from simpleMDM
	stateCustomProfiles := []string{}
	for _, profileId := range state.CustomProfiles.Elements() {
		stateCustomProfiles = append(stateCustomProfiles, profileId.(types.String).ValueString())
	}

	//reading configured profiles from TF file
	planCustomProfiles := []string{}
	for _, profileId := range plan.CustomProfiles.Elements() {
		planCustomProfiles = append(planCustomProfiles, profileId.(types.String).ValueString())
	}

	// // creating diff
	customProfilesToAdd, customProfilesToRemove := diffFunction(stateCustomProfiles, planCustomProfiles)

	// //adding profiles
	for _, profileId := range customProfilesToAdd {
		err := r.client.CustomProfileAssignToDevice(profileId, plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating device custom profile assignment",
				"Could not update device custom profile assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	//removing profiles
	for _, profileId := range customProfilesToRemove {
		err := r.client.CustomProfileUnAssignToDevice(profileId, plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating device custom profile assignment",
				"Could not update device custom profile assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *deviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state deviceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing device
	err := r.client.DeviceDelete(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SimpleMDM device",
			"Could not delte device, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *deviceResource) assignAPIValues(ctx context.Context, apiDevice *simplemdmext.DeviceResponse, model *deviceResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	flatAttributes := simplemdmext.FlattenAttributes(apiDevice.Data.Attributes)
	detailsValue, detailsDiags := types.MapValueFrom(ctx, types.StringType, flatAttributes)
	diags.Append(detailsDiags...)
	if !detailsValue.IsNull() {
		model.Details = detailsValue
	} else {
		model.Details = types.MapNull(types.StringType)
	}

	if id := apiDevice.Data.ID; id != 0 {
		model.ID = types.StringValue(strconv.Itoa(id))
	}

	if name, ok := flatAttributes["name"]; ok && name != "" {
		model.Name = types.StringValue(name)
	}

	if deviceName, ok := flatAttributes["device_name"]; ok && deviceName != "" {
		model.DeviceName = types.StringValue(deviceName)
	} else if model.DeviceName.IsNull() {
		model.DeviceName = types.StringNull()
	}

	enrollmentURL := flatAttributes["enrollment_url"]
	if enrollmentURL != "" && enrollmentURL != "null" {
		model.EnrollmentURL = types.StringValue(enrollmentURL)
	} else {
		model.EnrollmentURL = types.StringNull()
	}

	if groupID := apiDevice.Data.Relationships.DeviceGroup.Data.ID; groupID != 0 {
		model.DeviceGroup = types.StringValue(strconv.Itoa(groupID))
	}

	attributeValues := map[string]attr.Value{}
	for _, attribute := range apiDevice.Data.Relationships.CustomAttributeValues.Data {
		if attribute.Attributes.Value != "" {
			attributeValues[attribute.ID] = types.StringValue(attribute.Attributes.Value)
		}
	}

	if len(attributeValues) > 0 {
		attributesMap, attrDiags := types.MapValue(types.StringType, attributeValues)
		diags.Append(attrDiags...)
		model.Attributes = attributesMap
	} else {
		model.Attributes = types.MapNull(types.StringType)
	}

	return diags
}
