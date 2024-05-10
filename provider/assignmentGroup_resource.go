package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &assignment_groupResource{}
	_ resource.ResourceWithConfigure   = &assignment_groupResource{}
	_ resource.ResourceWithImportState = &assignment_groupResource{}
)

// assignment_groupResourceModel maps the resource schema data.
type assignment_groupResourceModel struct {
	Name         types.String `tfsdk:"name"`
	AutoDeploy   types.Bool   `tfsdk:"auto_deploy"`
	GroupType    types.String `tfsdk:"group_type"`
	InstallType  types.String `tfsdk:"install_type"`
	ID           types.String `tfsdk:"id"`
	Apps         types.Set    `tfsdk:"apps"`
	AppsUpdate   types.Bool   `tfsdk:"apps_update"`
	AppsPush     types.Bool   `tfsdk:"apps_push"`
	Profiles     types.Set    `tfsdk:"profiles"`
	ProfilesSync types.Bool   `tfsdk:"profiles_sync"`
	Groups       types.Set    `tfsdk:"groups"`
	Devices      types.Set    `tfsdk:"devices"`
}

// AssignmentGroupResource is a helper function to simplify the provider implementation.
func AssignmentGroupResource() resource.Resource {
	return &assignment_groupResource{}
}

// assignment_groupResource is the resource implementation.
type assignment_groupResource struct {
	client *simplemdm.Client
}

// Configure adds the provider configured client to the resource.
func (r *assignment_groupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*simplemdm.Client)
}

// Metadata returns the resource type name.
func (r *assignment_groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_assignmentgroup"
}

// Schema defines the schema for the resource.
func (r *assignment_groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Assignment Group resource is used to manage group, you can assign App(s), Custom Profile(s), Device(s), Device Group(s) and set addition details regarding Assignemtn Group.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "The name of the Assignment Group.",
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "ID of the Assignment Group in SimpleMDM",
			},
			"auto_deploy": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Optional. Whether the Apps should be automatically pushed to device(s) when they join this Assignment Group. Defaults to true",
			},
			"group_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("standard"),
				Validators: []validator.String{
					// Validate string value must be "standard" or "munki"
					stringvalidator.OneOf("standard", "munki"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Optional. Type of assignment group. Must be one of standard (for MDM app/media deployments) or munki for Munki app deployments. Defaults to standard.",
			},
			"install_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("managed"),
				Validators: []validator.String{
					// Validate string value must be "managed", "self_serve" or "munki"
					stringvalidator.OneOf([]string{"managed", "self_serve", "managed_updates", "default_installs"}...),
				},
				Description: "Optional. The install type for munki assignment groups. Must be one of managed, self_serve, managed_updates or default_installs. This setting has no effect for non-munki (standard) assignment groups. Defaults to managed.",
			},
			"apps": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Optional. List of Apps assigned to this assignment group",
			},
			"apps_update": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Optional. Set true if you would like to send update Apps command after assignment group creation or changes. Defaults to false.",
			},
			"apps_push": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Optional. Set true if you would like to send push Apps command after assignment group creation or changes. Defaults to false.",
			},
			"profiles": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Optional. List of Configuration Profiles (both Custom and predefined Profiles) assigned to this assignment group",
			},
			"profiles_sync": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Optional. Set true if you would like to send Sync Profiles command after Assignment Group creation or changes. Defaults to false.",
			},
			"groups": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Optional. List of Device Groups assigned to this Assignment Group",
			},
			"devices": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Optional. List of Devices assigned to this Assignment Group",
			},
		},
	}
}

// Import function
func (r *assignment_groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource
func (r *assignment_groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan assignment_groupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	assignmentgroup, err := r.client.AssignmentGroupCreate(plan.Name.ValueString(), plan.AutoDeploy.ValueBool(), plan.GroupType.ValueString(), plan.InstallType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating assignment group",
			"Could not create assignment group, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(assignmentgroup.Data.ID))

	// Assign all apps in plan
	for _, appId := range plan.Apps.Elements() {
		err := r.client.AssignmentGroupAssignObject(plan.ID.ValueString(), strings.Replace(appId.String(), "\"", "", 2), "apps")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group app assignment",
				"Could not update assignment group app assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Assign all profiles in plan
	for _, profileId := range plan.Profiles.Elements() {
		err := r.client.AssignmentGroupAssignObject(plan.ID.ValueString(), strings.Replace(profileId.String(), "\"", "", 2), "profiles")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group profile assignment",
				"Could not update assignment group profile assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	//assign all groups in plan
	for _, groupId := range plan.Groups.Elements() {
		err := r.client.AssignmentGroupAssignObject(plan.ID.ValueString(), strings.Replace(groupId.String(), "\"", "", 2), "device_groups")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group device group assignment",
				"Could not update assignment group device group, unexpected error: "+err.Error(),
			)
			return
		}
	}

	//assign all devices in plan
	for _, deviceId := range plan.Devices.Elements() {
		err := r.client.AssignmentGroupAssignObject(plan.ID.ValueString(), strings.Replace(deviceId.String(), "\"", "", 2), "devices")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group device group assignment",
				"Could not update assignment group device group, unexpected error: "+err.Error(),
			)
			return
		}
	}

	if plan.AppsUpdate.ValueBool() {
		err := r.client.AssignmentGroupUpdateInstalledApps(plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error when sending command to Update Apps, deleting group to prevent issus next run.",
				"Could not send Apps Update command, unexpected error: "+err.Error(),
			)
			err := r.client.AssignmentGroupDelete(plan.ID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Deleting SimpleMDM assignment group",
					"Could not delete assignment group, unexpected error: "+err.Error(),
				)
				return
			}
			return
		}
	}

	if plan.AppsPush.ValueBool() {
		err := r.client.AssignmentGroupPushApps(plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error when sending command to Push Apps, deleting group to prevent issus next run.",
				"Could not send Push Apps command, unexpected error: "+err.Error(),
			)
			err := r.client.AssignmentGroupDelete(plan.ID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Deleting SimpleMDM assignment group",
					"Could not delete assignment group, unexpected error: "+err.Error(),
				)
				return
			}
			return
		}
	}

	if plan.ProfilesSync.ValueBool() {
		err := r.client.AssignmentGroupSyncProfiles(plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error when sending command to Sync Profiles, deleting group to prevent issus next run.",
				"Could not send Sync Profiles command, unexpected error: "+err.Error(),
			)
			err := r.client.AssignmentGroupDelete(plan.ID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error Deleting SimpleMDM assignment group",
					"Could not delete assignment group, unexpected error: "+err.Error(),
				)
				return
			}
			return
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read group data
func (r *assignment_groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state assignment_groupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed assignment group values from SimpleMDM
	assignmentGroup, err := r.client.AssignmentGroupGet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM assignment group",
			"Could not read assignment group ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	//read apps and put them to slice
	appsPresent := false
	appsElements := []attr.Value{}
	for _, appAssigned := range assignmentGroup.Data.Relationships.Apps.Data {
		appsElements = append(appsElements, types.StringValue(strconv.Itoa(appAssigned.ID)))
		appsPresent = true
	}
	//if there are apps return them to state
	if appsPresent {
		appsSetValue, _ := types.SetValue(types.StringType, appsElements)
		state.Apps = appsSetValue
	} else {
		appsSetValue := types.SetNull(types.StringType)
		state.Apps = appsSetValue
	}

	//read all groups and put them to slice
	groupsPresent := false
	groupsElements := []attr.Value{}
	for _, groupAssigned := range assignmentGroup.Data.Relationships.DeviceGroups.Data {
		groupsElements = append(groupsElements, types.StringValue(strconv.Itoa(groupAssigned.ID)))
		groupsPresent = true
	}
	//if there are groups return them to state
	if groupsPresent {
		groupsSetValue, _ := types.SetValue(types.StringType, groupsElements)
		state.Groups = groupsSetValue
	} else {
		groupsSetValue := types.SetNull(types.StringType)
		state.Groups = groupsSetValue
	}

	//read all devices and put them to slice
	devicesPresent := false
	devicesElements := []attr.Value{}
	for _, deviceAssigned := range assignmentGroup.Data.Relationships.Devices.Data {
		devicesElements = append(devicesElements, types.StringValue(strconv.Itoa(deviceAssigned.ID)))
		devicesPresent = true
	}
	//if there are groups return them to state
	if devicesPresent {
		devicesSetValue, _ := types.SetValue(types.StringType, devicesElements)
		state.Devices = devicesSetValue
	} else {
		devicesSetValue := types.SetNull(types.StringType)
		state.Devices = devicesSetValue
	}

	resp.Diagnostics.AddWarning(
		"Notice about profiles:",
		"API limitations is curretly not allowing terraform provider to get state of the profiles assigned to assigment group."+
			" This is not issue as long as you are using only terraform provider to manage profiles for assigment group."+
			" This will be implemented properly once API will have correct responses and we will be able to load profiles for assignment group via API.",
	)

	// //read all profiles and put them to slice
	// profilesPresent := false
	// profilesElements := []attr.Value{}
	// for _, profileAssigned := range assignmentGroup.Data.Relationships.DeviceGroups.Data { //<<edit here
	// 	profilesElements = append(profilesElements, types.StringValue(strconv.Itoa(profileAssigned.ID)))
	// 	profilesPresent = true
	// }
	// //if there are groups return them to state
	// if profilesPresent {
	// 	profilesSetValue, _ := types.SetValue(types.StringType, profilesElements)
	// 	state.Profiles = profilesSetValue
	// }else {
	// 	profilesSetValue := types.SetNull(types.StringType)
	// 	state.Profiles = profilesSetValue
	// }

	// Overwrite items with refreshed state
	state.Name = types.StringValue(assignmentGroup.Data.Attributes.Name)
	state.AutoDeploy = types.BoolValue(assignmentGroup.Data.Attributes.AutoDeploy)
	state.GroupType = types.StringValue(assignmentGroup.Data.Attributes.Type)
	if assignmentGroup.Data.Attributes.Type == "munki" {
		state.InstallType = types.StringValue(assignmentGroup.Data.Attributes.InstallType)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// update group
func (r *assignment_groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Retrieve values from plan
	var plan, state assignment_groupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	err := r.client.AssignmentGroupUpdate(plan.Name.ValueString(), plan.AutoDeploy.ValueBool(), plan.ID.ValueString(), plan.GroupType.ValueString(), plan.InstallType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating assignment group",
			"Could not update assignment group, unexpected error: "+err.Error(),
		)
		return
	}

	//Handling assigned apps
	//reading assigned apps from simpleMDM
	stateApps := []string{}
	for _, appID := range state.Apps.Elements() {
		stateApps = append(stateApps, strings.Replace(appID.String(), "\"", "", 2))
	}

	//reading configured apps from TF file
	planApps := []string{}
	for _, appID := range plan.Apps.Elements() {
		planApps = append(planApps, strings.Replace(appID.String(), "\"", "", 2))
	}

	// creating diff
	appsToAdd, appsToRemove := diffFunction(stateApps, planApps)

	//adding apps
	for _, appId := range appsToAdd {
		err := r.client.AssignmentGroupAssignObject(plan.ID.ValueString(), appId, "apps")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group app assignment",
				"Could not update assignment group app assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	//removing apps
	for _, appId := range appsToRemove {
		err := r.client.AssignmentGroupUnAssignObject(plan.ID.ValueString(), appId, "apps")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group app assignment",
				"Could not update assignment group app assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	//Handling assigned profiles
	//reading assigned profiles from simpleMDM
	stateProfiles := []string{}
	for _, profileId := range state.Profiles.Elements() { //<< edit here
		stateProfiles = append(stateProfiles, strings.Replace(profileId.String(), "\"", "", 2))
	}

	//reading configured profiles from TF file
	planProfiles := []string{}
	for _, profileId := range plan.Profiles.Elements() {
		planProfiles = append(planProfiles, strings.Replace(profileId.String(), "\"", "", 2))
	}

	// creating diff
	profilesToAdd, profilesToRemove := diffFunction(stateProfiles, planProfiles)

	//adding profiles
	for _, profileId := range profilesToAdd {
		err := r.client.AssignmentGroupAssignObject(plan.ID.ValueString(), profileId, "profiles")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group profile assignment",
				"Could not update assignment group profile assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	//removing profiles
	for _, profileId := range profilesToRemove {
		err := r.client.AssignmentGroupUnAssignObject(plan.ID.ValueString(), profileId, "profiles")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group app assignment",
				"Could not update assignment group app assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	//handling assigned groups
	// reading currently assigned apps
	stateGroups := []string{}
	for _, groupId := range state.Groups.Elements() {
		stateGroups = append(stateGroups, strings.Replace(groupId.String(), "\"", "", 2))
	}
	//reading configured apps in TF file
	planGroups := []string{}
	for _, groupId := range plan.Groups.Elements() {
		planGroups = append(planGroups, strings.Replace(groupId.String(), "\"", "", 2))
	}
	//creating diff
	groupsToAdd, groupsToRemove := diffFunction(stateGroups, planGroups)

	//groups to add
	for _, groupId := range groupsToAdd {
		err := r.client.AssignmentGroupAssignObject(plan.ID.ValueString(), groupId, "device_groups")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group device group assignment",
				"Could not update assignment group device group, unexpected error: "+err.Error(),
			)
			return
		}
	}

	//groups to remove
	for _, groupId := range groupsToRemove {
		err := r.client.AssignmentGroupUnAssignObject(plan.ID.ValueString(), groupId, "device_groups")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group device group assignment",
				"Could not update assignment group device group, unexpected error: "+err.Error(),
			)
			return
		}
	}

	//handling assigned devices
	//reading currently assigned devices
	stateDevices := []string{}
	for _, device := range state.Devices.Elements() {
		stateDevices = append(stateDevices, strings.Replace(device.String(), "\"", "", 2))
	}
	//reading configured apps in TF file
	planDevices := []string{}
	for _, device := range plan.Groups.Elements() {
		planDevices = append(planDevices, strings.Replace(device.String(), "\"", "", 2))
	}
	//creating diff
	devicesToAdd, devicesToRemove := diffFunction(stateDevices, planDevices)

	//groups to add
	for _, deviceId := range devicesToAdd {
		err := r.client.AssignmentGroupAssignObject(plan.ID.ValueString(), deviceId, "devices")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group device group assignment",
				"Could not update assignment group device group, unexpected error: "+err.Error(),
			)
			return
		}
	}

	//groups to remove
	for _, deviceId := range devicesToRemove {
		err := r.client.AssignmentGroupUnAssignObject(plan.ID.ValueString(), deviceId, "devices")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group device group assignment",
				"Could not update assignment group device group, unexpected error: "+err.Error(),
			)
			return
		}
	}

	if plan.AppsUpdate.ValueBool() {
		err := r.client.AssignmentGroupUpdateInstalledApps(plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group profile assignment",
				"Could not update assignment group profile assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	if plan.AppsPush.ValueBool() {
		err := r.client.AssignmentGroupPushApps(plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group profile assignment",
				"Could not update assignment group profile assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	if plan.ProfilesSync.ValueBool() {
		err := r.client.AssignmentGroupSyncProfiles(plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group profile assignment",
				"Could not update assignment group profile assignment, unexpected error: "+err.Error(),
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

// Delete group
func (r *assignment_groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state assignment_groupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing assignment group
	err := r.client.AssignmentGroupDelete(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SimpleMDM assignment group",
			"Could not delete assignment group, unexpected error: "+err.Error(),
		)
		return
	}
}

// helper function to get diff between two groups
func diffFunction(state []string, plan []string) (add []string, remove []string) {
	IDsToAdd := []string{}
	IDsToRemove := []string{}
	for _, planObject := range plan {
		ispresent := false
		for _, stateObject := range state {
			if planObject == stateObject {
				ispresent = true
				break
			}
		}

		if !ispresent {
			IDsToAdd = append(IDsToAdd, planObject)
		}
	}

	for _, stateObject := range state {
		ispresent := false
		for _, planObject := range plan {
			if stateObject == planObject {
				ispresent = true
				break
			}
		}
		if !ispresent {
			IDsToRemove = append(IDsToRemove, stateObject)
		}
	}
	return IDsToAdd, IDsToRemove
}
