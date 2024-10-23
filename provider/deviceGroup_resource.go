package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &deviceGroupResource{}
	_ resource.ResourceWithConfigure   = &deviceGroupResource{}
	_ resource.ResourceWithImportState = &deviceGroupResource{}
)

// deviceGroupResourceModel maps the resource schema data.
type deviceGroupResourceModel struct {
	Name           types.String `tfsdk:"name"`
	ID             types.String `tfsdk:"id"`
	Attributes     types.Map    `tfsdk:"attributes"`
	Profiles       types.Set    `tfsdk:"profiles"`
	CustomProfiles types.Set    `tfsdk:"customprofiles"`
}

// deviceGroupResource is a helper function to simplify the provider implementation.
func DeviceGroupResource() resource.Resource {
	return &deviceGroupResource{}
}

// deviceGroupResource is the resource implementation.
type deviceGroupResource struct {
	client *simplemdm.Client
}

// Configure adds the provider configured client to the resource.
func (r *deviceGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*simplemdm.Client)
}

// Metadata returns the resource type name.
func (r *deviceGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devicegroup"
}

// Schema defines the schema for the resource.
func (r *deviceGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Device Group resource can be used to manage Device Group. Can be used together with Custom Profile(s), Attribute(s), Assignment Group(s) or Device Group(s) and set addition details regarding Device Group.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Required. The name of the device group.",
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "ID of a Device Group in SimpleMDM",
			},
			"profiles": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Optional. List of Configuration Profiles assigned to this Device Group",
			},
			"customprofiles": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Optional. List of Custom Configuration Profiles assigned to this Device Group",
			},
			"attributes": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Optional. Map of Custom Configuration Profiles and values set for this Device Group",
			},
		},
	}
}

func (r *deviceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to state
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource
func (r *deviceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan deviceGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError(
		"Resource can not be created!",
		"Device groups currently do not support creation via API request, if you wish to create new group  "+
			"go to website and create group and use import. Name of the group also can not be managed via provider, "+
			"same as deletion of the group can not be done via terraform. This will be implemented properly once API will have correct endpoints.",
	)
	if resp.Diagnostics.HasError() {
		return
	}
	// // Generate API request body from plan
	// deviceGroup, err := r.client.CreateDeviceGroup(plan.Name.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error creating device group",
	// 		"Could not create device group, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }

	// plan.ID = types.StringValue(strconv.Itoa(deviceGroup.Data.ID))

	// //setting attributes
	// for attribute, value := range plan.Attributes.Elements() {
	// 	err := r.client.SetAttributeForDeviceGroupAttribute(plan.ID.ValueString(), attribute, strings.Replace(value.String(), "\"", "", 2))
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"Error updating device group attribute",
	// 			"Could not set attribute value for device group, unexpected error: "+err.Error(),
	// 		)
	// 		return
	// 	}
	// }

	// // Assign all profiles in plan
	// for _, profileId := range plan.Profiles.Elements() {
	// 	err := r.client.AssignToAssignmentGroup(plan.ID.ValueString(), strings.Replace(profileId.String(), "\"", "", 2), "profiles")
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"Error updating device group profile assignment",
	// 			"Could not update device group profile assignment, unexpected error: "+err.Error(),
	// 		)
	// 		return
	// 	}
	// }
	// // Map response body to schema and populate Computed attribute values

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *deviceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state deviceGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get device group value from SimpleMDM
	devicegroup, err := r.client.DeviceGroupGet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM device group",
			"Could not read SimpleMDM device group "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	//load attributes for given group
	attributes, err := r.client.AttributeGetAttributesForDeviceGroup(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM device group attributes",
			"Could not read SimpleMDM device group attributes"+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	//adding attributes to the map
	attributePresent := false
	attributesElements := map[string]attr.Value{}
	for _, attribute := range attributes.Data {
		if attribute.Attributes.Source == "group" {
			attributesElements[attribute.ID] = types.StringValue(attribute.Attributes.Value)
			attributePresent = true
		}
	}
	if attributePresent {
		attributesSetValue, _ := types.MapValue(types.StringType, attributesElements)
		state.Attributes = attributesSetValue
	} else {
		attributesSetValue := types.MapNull(types.StringType)
		state.Attributes = attributesSetValue
	}

	// Overwrite items with refreshed state
	state.Name = types.StringValue(devicegroup.Data.Attributes.Name)

	// Load all profiles in SimpleMDM
	profiles, err := r.client.ProfileGetAll()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM profiles",
			"Could not read SimpleMDM profiles: "+err.Error(),
		)
		return
	}
	// //read all profiles and put them to slice
	profilesPresent := false
	profilesElements := []attr.Value{}
	customProfilesPresent := false
	customProfilesElements := []attr.Value{}

	for _, profile := range profiles.Data { //<<edit here
		for _, group := range profile.Relationships.DeviceGroups.Data {
			if strconv.Itoa(group.ID) == state.ID.ValueString() {
				if profile.Type == "custom_configuration_profile" {
					customProfilesElements = append(customProfilesElements, types.StringValue(strconv.Itoa(profile.ID)))
					customProfilesPresent = true
				} else {
					profilesElements = append(profilesElements, types.StringValue(strconv.Itoa(profile.ID)))
					profilesPresent = true
				}
			}
		}

	}

	//if there are profile or custom profiles return them to state
	if profilesPresent {
		profilesSetValue, _ := types.SetValue(types.StringType, profilesElements)
		state.Profiles = profilesSetValue
	} else {
		profilesSetValue := types.SetNull(types.StringType)
		state.Profiles = profilesSetValue
	}

	if customProfilesPresent {
		customProfilesSetValue, _ := types.SetValue(types.StringType, customProfilesElements)
		state.CustomProfiles = customProfilesSetValue
	} else {
		customProfilesSetValue := types.SetNull(types.StringType)
		state.CustomProfiles = customProfilesSetValue
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *deviceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Retrieve values from plan
	var plan, state deviceGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	// err := r.client.UpdateDeviceGroup(plan.ID.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error updating device group",
	// 		"Could not update device group, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }

	resp.Diagnostics.AddWarning(
		"Name can not be changed via terraform",
		"Device groups currently do not support change of the name via API request, in case you wish to change "+
			"name of the device group please do it via website.",
	)

	//comparing planed attributes and their values to attributes in SimpleMDM
	for planAttribute, planValue := range plan.Attributes.Elements() {
		found := false
		for stateAttribute, stateValue := range state.Attributes.Elements() {
			if planAttribute == stateAttribute {
				found = true
				if planValue != stateValue {
					err := r.client.AttributeSetAttributeForDeviceGroup(plan.ID.ValueString(), planAttribute, strings.Replace(planValue.String(), "\"", "", 2))
					if err != nil {
						resp.Diagnostics.AddError(
							"Error updating SimpleMDM device group attributes value",
							"Could not update SimpleMDM device group attributes value "+plan.ID.ValueString()+": "+err.Error(),
						)
						return
					}
				}
				break
			}
		}
		if !found {
			err := r.client.AttributeSetAttributeForDeviceGroup(plan.ID.ValueString(), planAttribute, strings.Replace(planValue.String(), "\"", "", 2))
			if err != nil {
				resp.Diagnostics.AddError(
					"Error updating SimpleMDM device group attributes value",
					"Could not update SimpleMDM device group attributes value "+plan.ID.ValueString()+": "+err.Error(),
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
			err := r.client.AttributeSetAttributeForDeviceGroup(plan.ID.ValueString(), stateAttribute, "")
			if err != nil {
				resp.Diagnostics.AddError(
					"Error updating SimpleMDM device group attributes value",
					"Could not update SimpleMDM device group attributes value "+plan.ID.ValueString()+": "+err.Error(),
				)
				return
			}
		}
	}

	//Handling assigned profiles
	//reading assigned profiles from simpleMDM
	stateProfiles := []string{}
	for _, profileId := range state.Profiles.Elements() {
		stateProfiles = append(stateProfiles, strings.Replace(profileId.String(), "\"", "", 2))
	}

	//reading configured profiles from TF file
	planProfiles := []string{}
	for _, profileId := range plan.Profiles.Elements() {
		planProfiles = append(planProfiles, strings.Replace(profileId.String(), "\"", "", 2))
	}

	// // creating diff
	profilesToAdd, profilesToRemove := diffFunction(stateProfiles, planProfiles)

	// //adding profiles
	for _, profileId := range profilesToAdd {
		err := r.client.ProfileAssignToGroup(profileId, plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating device group profile assignment",
				"Could not update device group profile assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	//removing profiles
	for _, profileId := range profilesToRemove {
		err := r.client.ProfileUnAssignToGroup(profileId, plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating device group profile assignment",
				"Could not update device group profile assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	//Handling assigned custom profiles
	stateCustomProfiles := []string{}
	for _, profileId := range state.CustomProfiles.Elements() {
		stateCustomProfiles = append(stateCustomProfiles, strings.Replace(profileId.String(), "\"", "", 2))
	}

	//reading configured profiles from TF file
	planCustomProfiles := []string{}
	for _, profileId := range plan.CustomProfiles.Elements() {
		planCustomProfiles = append(planCustomProfiles, strings.Replace(profileId.String(), "\"", "", 2))
	}

	// // creating diff
	customProfilesToAdd, customProfilesToRemove := diffFunction(stateCustomProfiles, planCustomProfiles)

	// //adding profiles
	for _, profileId := range customProfilesToAdd {
		err := r.client.CustomProfileAssignToDeviceGroup(profileId, plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating device group profile assignment",
				"Could not update device group profile assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	//removing profiles
	for _, profileId := range customProfilesToRemove {
		err := r.client.CustomProfileUnassignFromDeviceGroup(profileId, plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating device group profile assignment",
				"Could not update device group profile assignment, unexpected error: "+err.Error(),
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

func (r *deviceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state deviceGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddWarning(
		"Device group can not be deleted",
		"Applying this resource destruction will only remove the resource from the Terraform state "+
			"and will not call the deletion API due to API limitations. Manually use the web interface to fully destroy this resource.",
	)

	// // Delete existing group
	// err := r.client.DeleteDeviceGroup(state.ID.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Deleting SimpleMDM device groups",
	// 		"Could not delte device group, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }
}
