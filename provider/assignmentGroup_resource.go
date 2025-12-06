package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	Name                types.String `tfsdk:"name"`
	AutoDeploy          types.Bool   `tfsdk:"auto_deploy"`
	GroupType           types.String `tfsdk:"group_type"`
	InstallType         types.String `tfsdk:"install_type"`
	Priority            types.Int64  `tfsdk:"priority"`
	AppTrackLocation    types.Bool   `tfsdk:"app_track_location"`
	ID                  types.String `tfsdk:"id"`
	Apps                types.Set    `tfsdk:"apps"`
	AppsUpdate          types.Bool   `tfsdk:"apps_update"`
	AppsPush            types.Bool   `tfsdk:"apps_push"`
	Profiles            types.Set    `tfsdk:"profiles"`
	ProfilesSync        types.Bool   `tfsdk:"profiles_sync"`
	Groups              types.Set    `tfsdk:"groups"`
	Devices             types.Set    `tfsdk:"devices"`
	DevicesRemoveOthers types.Bool   `tfsdk:"devices_remove_others"`
	CreatedAt           types.String `tfsdk:"created_at"`
	UpdatedAt           types.String `tfsdk:"updated_at"`
	DeviceCount         types.Int64  `tfsdk:"device_count"`
	GroupCount          types.Int64  `tfsdk:"group_count"`
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
		Description: "Assignment Group resource is used to manage groups. You can assign apps, custom profiles, devices, and device groups, and configure additional assignment group settings.",
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
				Description: "Optional. Type of assignment group. Must be one of standard (for MDM app/media deployments) or munki for Munki app deployments. Defaults to standard. " +
					"⚠️ DEPRECATED: This field is deprecated by the SimpleMDM API and may be ignored for accounts using the New Groups Experience.",
			},
			"install_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					// Validate string value must be "managed", "self_serve" or "munki"
					stringvalidator.OneOf([]string{"managed", "self_serve", "managed_updates", "default_installs"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Optional. The install type for munki assignment groups. Must be one of managed, self_serve, managed_updates or default_installs. This setting has no effect for non-munki (standard) assignment groups. Defaults to managed for munki groups. " +
					"⚠️ DEPRECATED: The SimpleMDM API recommends setting install_type per-app using the Assign App endpoint instead of at the group level.",
			},
			"priority": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(0, 999),
				},
				Description: "Optional. Sets the priority order in which assignment groups are evaluated when devices are part of multiple groups. Lower numbers are evaluated first. Valid range: 0-999. If not set, SimpleMDM assigns a default priority.",
			},
			"app_track_location": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Optional. Controls whether the SimpleMDM app tracks device location when installed.",
			},
			"apps": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Optional. List of Apps assigned to this assignment group",
			},
			"apps_update": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Optional. Triggers 'Update Apps' command during apply. This sends an MDM install command to all associated devices for apps with available updates. Set to true when you want to push app updates. This is a one-time action on each apply where it's true. Difference from apps_push: update only installs if newer version available.",
			},
			"apps_push": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Optional. Triggers 'Push Apps' command during apply. This sends an MDM install command to all associated devices for all assigned apps, regardless of current version. Set to true when you want to reinstall or push apps. This is a one-time action on each apply where it's true. Difference from apps_update: push installs all apps.",
			},
			"profiles": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Optional. List of Configuration Profiles (both Custom and predefined Profiles) assigned to this assignment group",
			},
			"profiles_sync": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Optional. Triggers 'Sync Profiles' command during apply. This pushes all assigned profiles to devices in the assignment group. Set to true after profile changes to sync. ⚠️ Rate limited to 1 request per 30 seconds - wait between applies if true. This is a one-time action on each apply where it's true.",
			},
			"groups": schema.SetAttribute{
				ElementType:        types.StringType,
				Optional:           true,
				Computed:           true,
				DeprecationMessage: "The device_groups assignment API is deprecated by SimpleMDM. This only works with legacy_device_group_id from migrated groups. For accounts using the New Groups Experience, use device assignments instead.",
				Description:        "Optional. List of Device Groups assigned to this Assignment Group. ⚠️ DEPRECATED: This uses a deprecated API that only works with legacy_device_group_id from previously migrated groups.",
			},
			"devices": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Optional. List of Devices assigned to this Assignment Group",
			},
			"devices_remove_others": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Optional. When true, devices assigned through Terraform will be removed from other assignment groups before being added to this one.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the assignment group was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the assignment group was last updated.",
			},
			"device_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of devices currently assigned to the assignment group.",
			},
			"group_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of device groups currently assigned to the assignment group.",
			},
		},
	}
}

// syncProfilesWithRetry handles profile sync with rate limit retry logic
func (r *assignment_groupResource) syncProfilesWithRetry(ctx context.Context, groupID string) error {
	maxRetries := 3
	
	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := r.client.AssignmentGroupSyncProfiles(groupID)
		if err == nil {
			return nil
		}
		
		// Check for rate limit (429 status or rate limit in error message)
		if strings.Contains(err.Error(), "429") || strings.Contains(strings.ToLower(err.Error()), "rate limit") {
			if attempt < maxRetries {
				// Wait 30 seconds before retry
				select {
				case <-ctx.Done():
					return fmt.Errorf("operation cancelled: %w", ctx.Err())
				case <-time.After(30 * time.Second):
					continue
				}
			}
			return fmt.Errorf("profile sync rate limited after %d attempts. Please wait 30 seconds between sync operations", maxRetries+1)
		}
		
		// Non-rate-limit error, don't retry
		return err
	}
	
	return fmt.Errorf("profile sync failed after %d attempts", maxRetries+1)
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

	assignmentgroup, err := createAssignmentGroup(ctx, r.client, assignmentGroupUpsertRequest{
		Name:             plan.Name.ValueString(),
		AutoDeploy:       boolPointerFromType(plan.AutoDeploy),
		GroupType:        stringPointerFromType(plan.GroupType),
		InstallType:      stringPointerFromType(plan.InstallType),
		Priority:         int64PointerFromType(plan.Priority),
		AppTrackLocation: boolPointerFromType(plan.AppTrackLocation),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating assignment group",
			"Could not create assignment group, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(assignmentgroup.Data.ID))

	// Assign all apps in plan
	if err := assignObjectsToGroup(ctx, r.client, plan.ID.ValueString(), plan.Apps, "apps", false); err != nil {
		resp.Diagnostics.AddError(
			"Error assigning apps to assignment group",
			"Could not assign apps to assignment group, unexpected error: "+err.Error(),
		)
		return
	}

	// Assign all profiles in plan
	if err := assignObjectsToGroup(ctx, r.client, plan.ID.ValueString(), plan.Profiles, "profiles", false); err != nil {
		resp.Diagnostics.AddError(
			"Error assigning profiles to assignment group",
			"Could not assign profiles to assignment group, unexpected error: "+err.Error(),
		)
		return
	}

	// Assign all groups in plan
	if err := assignObjectsToGroup(ctx, r.client, plan.ID.ValueString(), plan.Groups, "device_groups", false); err != nil {
		resp.Diagnostics.AddError(
			"Error assigning device groups to assignment group",
			"Could not assign device groups to assignment group, unexpected error: "+err.Error(),
		)
		return
	}

	// Assign all devices in plan
	if err := assignObjectsToGroup(ctx, r.client, plan.ID.ValueString(), plan.Devices, "devices", boolValueOrDefault(plan.DevicesRemoveOthers, false)); err != nil {
		resp.Diagnostics.AddError(
			"Error assigning devices to assignment group",
			"Could not assign devices to assignment group, unexpected error: "+err.Error(),
		)
		return
	}

	if plan.AppsUpdate.ValueBool() {
		err := r.client.AssignmentGroupUpdateInstalledApps(plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Failed to send Update Apps command",
				fmt.Sprintf("Assignment group created successfully, but update apps command failed: %s. You may need to trigger manually.", err.Error()),
			)
		}
	}

	if plan.AppsPush.ValueBool() {
		err := r.client.AssignmentGroupPushApps(plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Failed to send Push Apps command",
				fmt.Sprintf("Assignment group created successfully, but push apps command failed: %s. You may need to trigger manually.", err.Error()),
			)
		}
	}

	if plan.ProfilesSync.ValueBool() {
		err := r.syncProfilesWithRetry(ctx, plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Failed to sync profiles",
				fmt.Sprintf("Assignment group created successfully, but profile sync failed: %s. Profiles may need manual sync.", err.Error()),
			)
		}
	}

	fetched, err := fetchAssignmentGroup(ctx, r.client, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error refreshing assignment group state",
			"Could not read assignment group after creation: "+err.Error(),
		)
		return
	}

	// Save the planned relationship values before applying API response
	// This handles eventual consistency where API may not immediately return assigned items
	plannedApps := plan.Apps
	plannedProfiles := plan.Profiles
	plannedGroups := plan.Groups
	plannedDevices := plan.Devices

	// Check what the API actually returned before applying to model
	apiReturnedApps := len(fetched.Data.Relationships.Apps.Data) > 0
	apiReturnedProfiles := len(fetched.Data.Relationships.Profiles.Data) > 0
	apiReturnedGroups := len(fetched.Data.Relationships.DeviceGroups.Data) > 0
	apiReturnedDevices := len(fetched.Data.Relationships.Devices.Data) > 0

	applyAssignmentGroupResponseToResourceModel(&plan, fetched)

	// Preserve planned relationships if API hasn't returned them yet (eventual consistency)
	preservePlannedRelationships(&plan, plannedApps, plannedProfiles, plannedGroups, plannedDevices,
		apiReturnedApps, apiReturnedProfiles, apiReturnedGroups, apiReturnedDevices)

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
	assignmentGroup, err := fetchAssignmentGroup(ctx, r.client, state.ID.ValueString())
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

	applyAssignmentGroupResponseToResourceModel(&state, assignmentGroup)

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

	err := updateAssignmentGroup(ctx, r.client, plan.ID.ValueString(), assignmentGroupUpsertRequest{
		Name:             plan.Name.ValueString(),
		AutoDeploy:       boolPointerFromType(plan.AutoDeploy),
		GroupType:        stringPointerFromType(plan.GroupType),
		InstallType:      stringPointerFromType(plan.InstallType),
		Priority:         int64PointerFromType(plan.Priority),
		AppTrackLocation: boolPointerFromType(plan.AppTrackLocation),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating assignment group",
			"Could not update assignment group, unexpected error: "+err.Error(),
		)
		return
	}

	// Update all assigned apps
	if err := updateAssignmentGroupObjects(ctx, r.client, plan.ID.ValueString(), state.Apps, plan.Apps, "apps", false); err != nil {
		resp.Diagnostics.AddError(
			"Error updating assignment group app assignments",
			"Could not update assignment group app assignments, unexpected error: "+err.Error(),
		)
		return
	}

	// Update all assigned profiles
	if err := updateAssignmentGroupObjects(ctx, r.client, plan.ID.ValueString(), state.Profiles, plan.Profiles, "profiles", false); err != nil {
		resp.Diagnostics.AddError(
			"Error updating assignment group profile assignments",
			"Could not update assignment group profile assignments, unexpected error: "+err.Error(),
		)
		return
	}

	// Update all assigned groups
	if err := updateAssignmentGroupObjects(ctx, r.client, plan.ID.ValueString(), state.Groups, plan.Groups, "device_groups", false); err != nil {
		resp.Diagnostics.AddError(
			"Error updating assignment group device group assignments",
			"Could not update assignment group device group assignments, unexpected error: "+err.Error(),
		)
		return
	}

	// Update all assigned devices
	if err := updateAssignmentGroupObjects(ctx, r.client, plan.ID.ValueString(), state.Devices, plan.Devices, "devices", boolValueOrDefault(plan.DevicesRemoveOthers, false)); err != nil {
		resp.Diagnostics.AddError(
			"Error updating assignment group device assignments",
			"Could not update assignment group device assignments, unexpected error: "+err.Error(),
		)
		return
	}

	if plan.AppsUpdate.ValueBool() {
		err := r.client.AssignmentGroupUpdateInstalledApps(plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Failed to send Update Apps command",
				fmt.Sprintf("Assignment group updated successfully, but update apps command failed: %s. You may need to trigger manually.", err.Error()),
			)
		}
	}

	if plan.AppsPush.ValueBool() {
		err := r.client.AssignmentGroupPushApps(plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Failed to send Push Apps command",
				fmt.Sprintf("Assignment group updated successfully, but push apps command failed: %s. You may need to trigger manually.", err.Error()),
			)
		}
	}

	if plan.ProfilesSync.ValueBool() {
		err := r.syncProfilesWithRetry(ctx, plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Failed to sync profiles",
				fmt.Sprintf("Assignment group updated successfully, but profile sync failed: %s. Profiles may need manual sync.", err.Error()),
			)
		}
	}

	fetched, err := fetchAssignmentGroup(ctx, r.client, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error refreshing assignment group state",
			"Could not read assignment group after update: "+err.Error(),
		)
		return
	}

	// Save the planned relationship values before applying API response
	// This handles eventual consistency where API may not immediately return assigned items
	plannedApps := plan.Apps
	plannedProfiles := plan.Profiles
	plannedGroups := plan.Groups
	plannedDevices := plan.Devices

	// Check what the API actually returned before applying to model
	apiReturnedApps := len(fetched.Data.Relationships.Apps.Data) > 0
	apiReturnedProfiles := len(fetched.Data.Relationships.Profiles.Data) > 0
	apiReturnedGroups := len(fetched.Data.Relationships.DeviceGroups.Data) > 0
	apiReturnedDevices := len(fetched.Data.Relationships.Devices.Data) > 0

	applyAssignmentGroupResponseToResourceModel(&plan, fetched)

	// Preserve planned relationships if API hasn't returned them yet (eventual consistency)
	preservePlannedRelationships(&plan, plannedApps, plannedProfiles, plannedGroups, plannedDevices,
		apiReturnedApps, apiReturnedProfiles, apiReturnedGroups, apiReturnedDevices)

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
