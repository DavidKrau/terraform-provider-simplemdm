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
	Name             types.String `tfsdk:"name"`
	AutoDeploy       types.Bool   `tfsdk:"auto_deploy"`
	ID               types.String `tfsdk:"id"`
	Apps             []appModel   `tfsdk:"apps"`
	AppsUpdate       types.Bool   `tfsdk:"apps_update"`
	AppsPush         types.Bool   `tfsdk:"apps_push"`
	Profiles         types.Set    `tfsdk:"profiles"`
	ProfilesSync     types.Bool   `tfsdk:"profiles_sync"`
	Devices          types.Set    `tfsdk:"devices"`
	Attributes       types.Map    `tfsdk:"attributes"`
	Priority         types.String `tfsdk:"priority"`
	AppTrackLocation types.Bool   `tfsdk:"app_track_location"`
}

type appModel struct {
	AppID          types.String `tfsdk:"app_id"`
	DeploymnetType types.String `tfsdk:"deployment_type"`
	InstallType    types.String `tfsdk:"install_type"`
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
		Description: "Assignment Group resource is used to manage group, you can assign App(s), Profile(s), Custom Profile(s), Custom Declaration(s), Device(s), Device Group(s) and set addition details regarding Group.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "The name of the Group.",
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "ID of the Group in SimpleMDM",
			},
			"auto_deploy": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Optional. Whether the Apps should be automatically pushed to device(s) when they join this Group. Defaults to true",
			},
			"app_track_location": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Optional. If true, it tracks the location of IOS device when the SimpleMDM mobile app is installed. Defaults to true.",
			},
			"apps": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"app_id": schema.StringAttribute{
							Required:    true,
							Description: "ID of the Application in SimpleMDM",
						},
						"deployment_type": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Optional. Type of assignment group. Must be one of standard (for MDM app/media deployments) or munki for Munki app deployments. Defaults to standard",
							Default:     stringdefault.StaticString("standard"),
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("standard", "munki"),
							},
						},
						"install_type": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Optional. The install type for munki assignment groups. Must be one of managed, self_serve, default_installs or managed_updates. This setting has no effect for non-munki (standard) assignment groups. Defaults to managed.",
							Default:     stringdefault.StaticString("managed"),
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Validators: []validator.String{
								stringvalidator.OneOf("managed", "self_serve", "default_installs", "managed_updates"),
							},
						},
					},
				},
				Description: "Optional. List of Apps assigned to this group",
			},
			"apps_update": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Optional. Updates associated apps on associated devices. A munki catalog refresh or MDM install command will be sent to all associated devices. Defaults to true",
			},
			"apps_push": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Optional. Installs associated apps to associated devices. A munki catalog refresh or MDM install command will be sent to all associated devices. Defaults to true.",
			},
			"profiles": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Optional. List of Configuration Profiles (Custom or predefined Profiles and Custom Declarations) assigned to this group",
			},
			"profiles_sync": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Optional. Set true if you would like to send Sync Profiles command after Group creation or changes. Defaults to true.",
			},
			"devices": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Optional. List of Devices assigned to this Group",
			},
			"attributes": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Optional. Map of Attributes and values set for this Group",
			},
			"priority": schema.StringAttribute{
				Optional:    true,
				Description: "Optional. The priority (0 to 20) of the assignment group. Default to 0",
				Default:     stringdefault.StaticString("0"),
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"),
				},
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
	assignmentgroup, err := r.client.AssignmentGroupCreate(plan.Name.ValueString(), plan.AutoDeploy.ValueBool(), plan.Priority.ValueString(), plan.AppTrackLocation.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating assignment group",
			"Could not create assignment group, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(assignmentgroup.Data.ID))

	//setting attributes
	for attribute, value := range plan.Attributes.Elements() {
		err := r.client.AttributeSetAttributeForDeviceGroup(plan.ID.ValueString(), attribute, strings.Replace(value.String(), "\"", "", 2))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating device group attribute",
				"Could not set attribute value for device group, unexpected error: "+err.Error(),
			)
			return
		}
	}

	for _, app := range plan.Apps {
		err := r.client.AssignmentGroupAssignApp(plan.ID.ValueString(), app.AppID.ValueString(), app.DeploymnetType.ValueString(), app.InstallType.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating device group apps",
				"Could not assing app to device group, unexpected error: "+err.Error(),
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
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM assignment group",
			"Could not read assignment group ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	//load attributes for given group
	attributes, err := r.client.AttributeGetAttributesForGroup(state.ID.ValueString())
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

	// Load all profiles in SimpleMDM
	profiles, err := r.client.ProfileGetAll()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM profiles",
			"Could not read SimpleMDM profiles: "+err.Error(),
		)
		return
	}

	//read apps and add them to state
	if len(assignmentGroup.Data.Relationships.Apps.Data) >= 1 {
		state.Apps = []appModel{}
		for _, app := range assignmentGroup.Data.Relationships.Apps.Data {
			state.Apps = append(state.Apps, appModel{
				AppID:          types.StringValue(strconv.Itoa(app.ID)),
				DeploymnetType: types.StringValue(app.DeploymnetType),
				InstallType:    types.StringValue(app.InstallType),
			})
		}
	} else {
		state.Apps = nil
	}
	//read all profiles and put them to slice
	profilesPresent := false
	profilesElements := []attr.Value{}

	for _, profile := range profiles.Data {
		for _, group := range profile.Relationships.DeviceGroups.Groups.Data {
			if strconv.Itoa(group.ID) == state.ID.ValueString() {
				profilesElements = append(profilesElements, types.StringValue(strconv.Itoa(profile.ID)))
				profilesPresent = true

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

	// Overwrite items with refreshed state
	state.Name = types.StringValue(assignmentGroup.Data.Attributes.Name)
	state.AutoDeploy = types.BoolValue(assignmentGroup.Data.Attributes.AutoDeploy)
	state.AppTrackLocation = types.BoolValue(assignmentGroup.Data.Attributes.AppTrackLocation)
	state.Priority = types.StringValue(strconv.Itoa(assignmentGroup.Data.Attributes.Priority))

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

	// compare app in state, add missing apps and remove and re-add apps with missmatch in config of the app
	for _, planApp := range plan.Apps {
		found := false
		update := false
		for _, stateApp := range state.Apps {
			if stateApp.AppID == planApp.AppID {
				found = true
				if stateApp.DeploymnetType != planApp.DeploymnetType {
					update = true
				}
				if stateApp.InstallType != planApp.InstallType {
					update = true
				}
			}
		}
		// app needs update remove it first, and later add it again
		if update {
			err := r.client.AssignmentGroupUnAssignApp(plan.ID.ValueString(), planApp.AppID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error updating device group apps",
					"Could not un-assing app from device group, unexpected error: "+err.Error(),
				)
				return
			}
			found = false
		}
		if !found {
			err := r.client.AssignmentGroupAssignApp(plan.ID.ValueString(), planApp.AppID.ValueString(), planApp.DeploymnetType.ValueString(), planApp.InstallType.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error updating device group apps",
					"Could not assing app to device group, unexpected error: "+err.Error(),
				)
				return
			}
		}

	}

	//remove any left over apps in state
	for _, stateApp := range state.Apps {
		found := false
		for _, planApp := range plan.Apps {
			if stateApp.AppID == planApp.AppID {
				found = true
			}
		}
		if !found {
			err := r.client.AssignmentGroupUnAssignApp(plan.ID.ValueString(), stateApp.AppID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error updating device group apps",
					"Could not un-assing app from device group, unexpected error: "+err.Error(),
				)
				return
			}
		}
	}

	// Generate API request body from plan
	err := r.client.AssignmentGroupUpdate(plan.Name.ValueString(), plan.AutoDeploy.ValueBool(), plan.ID.ValueString(), plan.AppTrackLocation.ValueBool(), plan.Priority.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating assignment group",
			"Could not update assignment group, unexpected error: "+err.Error(),
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

	//handling assigned devices
	//reading currently assigned devices
	stateDevices := []string{}
	for _, device := range state.Devices.Elements() {
		stateDevices = append(stateDevices, strings.Replace(device.String(), "\"", "", 2))
	}
	//reading configured apps in TF file
	planDevices := []string{}
	for _, device := range plan.Devices.Elements() {
		planDevices = append(planDevices, strings.Replace(device.String(), "\"", "", 2))
	}
	//creating diff
	devicesToAdd, devicesToRemove := diffFunction(stateDevices, planDevices)

	//devices to add
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

	//devices to remove
	for _, deviceId := range devicesToRemove {
		err := r.client.AssignmentGroupUnAssignObject(plan.ID.ValueString(), deviceId, "devices")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group device assignment",
				"Could not update assignment group device assignment, unexpected error: "+err.Error(),
			)
			return
		}
	}

	if plan.AppsUpdate.ValueBool() {
		err := r.client.AssignmentGroupUpdateInstalledApps(plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment App update failed",
				"Could not update assignment App update failed, unexpected error: "+err.Error(),
			)
			return
		}
	}

	if plan.AppsPush.ValueBool() {
		err := r.client.AssignmentGroupPushApps(plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment App push failed",
				"Could not update assignment App push failed, unexpected error: "+err.Error(),
			)
			return
		}
	}

	if plan.ProfilesSync.ValueBool() {
		err := r.client.AssignmentGroupSyncProfiles(plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating assignment group profile sync",
				"Could not update assignment group profile sync, unexpected error: "+err.Error(),
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
