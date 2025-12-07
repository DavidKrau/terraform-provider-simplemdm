package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &customProfileResource{}
	_ resource.ResourceWithConfigure   = &customProfileResource{}
	_ resource.ResourceWithImportState = &customProfileResource{}
)

// profileResourceModel maps the resource schema data.
type customProfileResourceModel struct {
	Name                   types.String `tfsdk:"name"`
	MobileConfig           types.String `tfsdk:"mobileconfig"`
	UserScope              types.Bool   `tfsdk:"user_scope"`
	AttributeSupport       types.Bool   `tfsdk:"attribute_support"`
	EscapeAttributes       types.Bool   `tfsdk:"escape_attributes"`
	ReinstallAfterOSUpdate types.Bool   `tfsdk:"reinstall_after_os_update"`
	ProfileIdentifier      types.String `tfsdk:"profile_identifier"`
	GroupCount             types.Int64  `tfsdk:"group_count"`
	DeviceCount            types.Int64  `tfsdk:"device_count"`
	ProfileSHA             types.String `tfsdk:"profile_sha"`
	ID                     types.String `tfsdk:"id"`
}

// ProfileResource is a helper function to simplify the provider implementation.
func CustomProfileResource() resource.Resource {
	return &customProfileResource{}
}

// profileResource is the resource implementation.
type customProfileResource struct {
	client *simplemdm.Client
}

// Configure adds the provider configured client to the resource.
func (r *customProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*simplemdm.Client)
}

// Metadata returns the resource type name.
func (r *customProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customprofile"
}

// Schema defines the schema for the resource.
func (r *customProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Custom Profile resource can be used to manage Custom Profile. Can be used together with Device(s), Assignment Group(s) or Device Group(s) and set addition details regarding Custom Profile.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Required. A name for the profile. Example: \"My First profile by terraform\"",
			},
			"mobileconfig": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Required. Can be string or you can use function 'file' or 'templatefile' to load string from file (see examples folder). Example: mobileconfig = file(\"./profiles/profile.mobileconfig\") or mobileconfig = <<-EOT PROFILE STRING EOT",
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "ID of a Custom Configuration Profile in SimpleMDM",
			},
			"user_scope": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(true),
				Computed:    true,
				Description: "Optional. A boolean true or false. If false, deploy as a device profile instead of a user profile for macOS devices. Defaults to true.",
			},
			"attribute_support": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Optional. A boolean true or false. When enabled, SimpleMDM will process variables in the uploaded profile. Defaults to false",
			},
			"escape_attributes": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Optional. A boolean true or false. When enabled, SimpleMDM escape the values of the custom variables in the uploaded profile. Defaults to false",
			},
			"reinstall_after_os_update": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Optional. A boolean true or false. When enabled, SimpleMDM will re-install the profile automatically after macOS software updates are detected. Defaults to false",
			},
			"profile_identifier": schema.StringAttribute{
				Computed:    true,
				Description: "Read-only profile identifier assigned by SimpleMDM.",
			},
			"group_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of device groups assigned to this custom configuration profile.",
			},
			"device_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of devices assigned to this custom configuration profile.",
			},
			"profile_sha": schema.StringAttribute{
				Computed:    true,
				Description: "SHA-256 checksum reported by SimpleMDM for the current mobileconfig payload.",
			},
		},
	}
}

func (r *customProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource
func (r *customProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan customProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	Profile, err := r.client.CustomProfileCreate(plan.Name.ValueString(), plan.MobileConfig.ValueString(), plan.UserScope.ValueBool(), plan.AttributeSupport.ValueBool(), plan.EscapeAttributes.ValueBool(), plan.ReinstallAfterOSUpdate.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating SimpleMDM custom profile",
			"Could not create custom profile: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(strconv.Itoa(Profile.Data.ID))
	assignCustomProfileAttributes(&plan, Profile.Data.Attributes)

	sha, body, err := r.client.CustomProfileSHA(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading SimpleMDM custom profile",
			"Could not download custom profile ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.MobileConfig = types.StringValue(body)
	plan.ProfileSHA = stringValueOrNull(sha)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *customProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state customProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile, err := r.client.CustomProfileGet(state.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading SimpleMDM custom profile",
			"Could not read custom profile ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	assignCustomProfileAttributes(&state, profile.Data.Attributes)
	state.ID = types.StringValue(strconv.Itoa(profile.Data.ID))

	sha, body, err := r.client.CustomProfileSHA(state.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading SimpleMDM custom profile",
			"Could not download custom profile ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.MobileConfig = types.StringValue(body)
	state.ProfileSHA = stringValueOrNull(sha)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *customProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Retrieve values from plan
	var plan customProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	_, err := r.client.CustomProfileUpdate(plan.Name.ValueString(), plan.MobileConfig.ValueString(), plan.UserScope.ValueBool(), plan.AttributeSupport.ValueBool(), plan.EscapeAttributes.ValueBool(), plan.ReinstallAfterOSUpdate.ValueBool(), "", plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating SimpleMDM custom profile",
			"Could not update custom profile ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	profile, err := r.client.CustomProfileGet(plan.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading SimpleMDM custom profile",
			"Could not read custom profile ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	assignCustomProfileAttributes(&plan, profile.Data.Attributes)

	sha, body, err := r.client.CustomProfileSHA(plan.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading SimpleMDM custom profile",
			"Could not download custom profile ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.MobileConfig = types.StringValue(body)
	plan.ProfileSHA = stringValueOrNull(sha)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *customProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state customProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing custom profile
	err := r.client.CustomProfileDelete(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting SimpleMDM custom profile",
			"Could not delete custom profile ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}

func assignCustomProfileAttributes(model *customProfileResourceModel, attributes simplemdm.Attributes) {
	model.Name = types.StringValue(attributes.Name)
	model.UserScope = types.BoolValue(attributes.UserScope)
	model.AttributeSupport = types.BoolValue(attributes.AttributeSupport)
	model.EscapeAttributes = types.BoolValue(attributes.EscapeAttributes)
	model.ReinstallAfterOSUpdate = types.BoolValue(attributes.ReinstallAfterOsUpdate)
	model.ProfileIdentifier = stringValueOrNull(attributes.ProfileIdentifier)
	model.GroupCount = types.Int64Value(int64(attributes.GroupCount))
	model.DeviceCount = types.Int64Value(int64(attributes.DeviceCount))
	if attributes.ProfileSHA != "" {
		model.ProfileSHA = types.StringValue(attributes.ProfileSHA)
	}
}

func stringValueOrNull(value string) types.String {
	if value == "" {
		return types.StringNull()
	}

	return types.StringValue(value)
}
