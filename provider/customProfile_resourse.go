package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &customProfileResource{}
	_ resource.ResourceWithConfigure      = &customProfileResource{}
	_ resource.ResourceWithImportState    = &customProfileResource{}
	_ resource.ResourceWithValidateConfig = &customProfileResource{}
)

// profileResourceModel maps the resource schema data.
type customProfileResourceModel struct {
	Name                           types.String `tfsdk:"name"`
	MobileConfig                   types.String `tfsdk:"mobileconfig"`
	ID                             types.String `tfsdk:"id"`
	UserScope                      types.Bool   `tfsdk:"userscope"`
	AttributeSupport               types.Bool   `tfsdk:"attributesupport"`
	EscapeAttributes               types.Bool   `tfsdk:"escapeattributes"`
	ReinstallAfterOSUpdate         types.Bool   `tfsdk:"reinstallafterosupdate"`
	Declarative                    types.Bool   `tfsdk:"declarative"`
	AutoRenewScepBasedCertificates types.Bool   `tfsdk:"auto_renew_scep_based_certificates"`
	AllowedPlatforms               types.Set    `tfsdk:"allowed_platforms"`
	MinimumMacosVersion            types.String `tfsdk:"minimum_macos_version"`
	MaximumMacosVersion            types.String `tfsdk:"maximum_macos_version"`
	AllowedMacosArchitecture       types.String `tfsdk:"allowed_macos_architecture"`
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
			"userscope": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(true),
				Computed:    true,
				Description: "Optional. A boolean true or false. If false, deploy as a device profile instead of a user profile for macOS devices. Defaults to true.",
			},
			"attributesupport": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Optional. A boolean true or false. When enabled, SimpleMDM will process variables in the uploaded profile. Defaults to false.",
			},
			"escapeattributes": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Optional. A boolean true or false. When enabled, SimpleMDM escape the values of the custom variables in the uploaded profile. Defaults to false.",
			},
			"reinstallafterosupdate": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Optional. A boolean true or false. When enabled, SimpleMDM will re-install the profile automatically after macOS software updates are detected. Defaults to false.",
			},
			"declarative": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Optional. A boolean true or false. When enabled, this profile will be installed using Declarative Management on any device that has Declarative Management enabled. Defaults to false.",
			},
			"auto_renew_scep_based_certificates": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Optional. A boolean true or false. When enabled, SimpleMDM will automatically re-issue SCEP based certificates in the profile before they expire. Cannot be enabled when declarative is enabled. Defaults to false.",
			},
			"allowed_platforms": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Default: setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("macos"),
					types.StringValue("ios"),
					types.StringValue("ipados"),
					types.StringValue("tvos"),
					types.StringValue("visionos"),
				})),
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf("macos", "ios", "ipados", "tvos", "visionos")),
				},
				Computed:    true,
				Description: "Optional. An array of operating systems the profile is allowed to be installed on. Valid values: macos, ios, ipados, tvos, visionos. Defaults to all platforms.",
			},
			"minimum_macos_version": schema.StringAttribute{
				Optional:    true,
				Description: "Optional. The minimum macOS version (e.g. 14.0) the profile is allowed to be installed on. Returns an error if the version is not recognized.",
			},
			"maximum_macos_version": schema.StringAttribute{
				Optional:    true,
				Description: "Optional. The maximum macOS version (e.g. 15.0) the profile is allowed to be installed on. Must be greater than or equal to minimum_macos_version. Returns an error if the version is not recognized.",
			},
			"allowed_macos_architecture": schema.StringAttribute{
				Optional: true,
				Default:  stringdefault.StaticString("any"),
				Validators: []validator.String{
					stringvalidator.OneOf("any", "x86", "arm"),
				},
				Computed:    true,
				Description: "Optional. Restricts the profile to a Mac architecture. Valid values: any, x86 (Intel), arm (Apple Silicon). Defaults to any.",
			},
		},
	}
}

func (r *customProfileResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data customProfileResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if both declarative and auto_renew_scep_based_certificates are true
	if !data.Declarative.IsNull() && !data.Declarative.IsUnknown() && data.Declarative.ValueBool() &&
		!data.AutoRenewScepBasedCertificates.IsNull() && !data.AutoRenewScepBasedCertificates.IsUnknown() && data.AutoRenewScepBasedCertificates.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			path.Root("auto_renew_scep_based_certificates"),
			"Conflicting Attribute Configuration",
			"auto_renew_scep_based_certificates cannot be enabled when declarative is enabled",
		)
	}
}

func (r *customProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource
func (r *customProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan customProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert AllowedPlatforms from types.Set to []string for API call
	var allowedPlatforms []string
	diags = plan.AllowedPlatforms.ElementsAs(ctx, &allowedPlatforms, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle optional version fields
	var minVersion, maxVersion *string
	if !plan.MinimumMacosVersion.IsNull() && !plan.MinimumMacosVersion.IsUnknown() {
		val := plan.MinimumMacosVersion.ValueString()
		minVersion = &val
	}
	if !plan.MaximumMacosVersion.IsNull() && !plan.MaximumMacosVersion.IsUnknown() {
		val := plan.MaximumMacosVersion.ValueString()
		maxVersion = &val
	}

	// Generate API request body from plan
	Profile, err := r.client.CustomProfileCreate(
		plan.Name.ValueString(),
		plan.MobileConfig.ValueString(),
		plan.UserScope.ValueBool(),
		plan.AttributeSupport.ValueBool(),
		plan.EscapeAttributes.ValueBool(),
		plan.ReinstallAfterOSUpdate.ValueBool(),
		plan.Declarative.ValueBool(),
		plan.AutoRenewScepBasedCertificates.ValueBool(),
		allowedPlatforms,
		minVersion,
		maxVersion,
		plan.AllowedMacosArchitecture.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating profile",
			"Could not create profile, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(strconv.Itoa(Profile.Data.ID))

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

	// Get refreshed profile values from SimpleMDM
	profile, err := r.client.ProfileGet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM custom profile",
			"Could not read custom profile ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(profile.Data.Attributes.Name)
	state.UserScope = types.BoolValue(profile.Data.Attributes.UserScope)
	state.AttributeSupport = types.BoolValue(profile.Data.Attributes.AttributeSupport)
	state.EscapeAttributes = types.BoolValue(profile.Data.Attributes.EscapeAttributes)
	state.ReinstallAfterOSUpdate = types.BoolValue(profile.Data.Attributes.ReinstallAfterOsUpdate)
	state.Declarative = types.BoolValue(profile.Data.Attributes.Declarative)
	state.AutoRenewScepBasedCertificates = types.BoolValue(profile.Data.Attributes.AutoRenewScepBasedCertificates)

	// Convert allowed_platforms from API response to types.Set
	// API returns ["macOS", "iOS", "iPadOS", "tvOS"] - convert to lowercase
	platformValues := make([]attr.Value, 0, len(profile.Data.Attributes.AllowedPlatforms))
	for _, platform := range profile.Data.Attributes.AllowedPlatforms {
		// Convert to lowercase to match schema expectations
		platformValues = append(platformValues, types.StringValue(strings.ToLower(platform)))
	}
	state.AllowedPlatforms = types.SetValueMust(types.StringType, platformValues)

	// Handle optional string fields - API returns null for unset values
	if profile.Data.Attributes.MinimumMacosVersion != nil && *profile.Data.Attributes.MinimumMacosVersion != "" {
		state.MinimumMacosVersion = types.StringValue(*profile.Data.Attributes.MinimumMacosVersion)
	} else {
		state.MinimumMacosVersion = types.StringNull()
	}

	if profile.Data.Attributes.MaximumMacosVersion != nil && *profile.Data.Attributes.MaximumMacosVersion != "" {
		state.MaximumMacosVersion = types.StringValue(*profile.Data.Attributes.MaximumMacosVersion)
	} else {
		state.MaximumMacosVersion = types.StringNull()
	}

	state.AllowedMacosArchitecture = types.StringValue(profile.Data.Attributes.AllowedMacosArchitecture)

	body, err := r.client.CustomProfileDownload(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM custom profile",
			"Could not read custom profles ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.MobileConfig = types.StringValue(body)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *customProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan customProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert AllowedPlatforms from types.Set to []string for API call
	var allowedPlatforms []string
	diags = plan.AllowedPlatforms.ElementsAs(ctx, &allowedPlatforms, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle optional version fields
	var minVersion, maxVersion *string
	if !plan.MinimumMacosVersion.IsNull() && !plan.MinimumMacosVersion.IsUnknown() {
		val := plan.MinimumMacosVersion.ValueString()
		minVersion = &val
	}
	if !plan.MaximumMacosVersion.IsNull() && !plan.MaximumMacosVersion.IsUnknown() {
		val := plan.MaximumMacosVersion.ValueString()
		maxVersion = &val
	}

	// Generate API request body from plan
	_, err := r.client.CustomProfileUpdate(
		plan.Name.ValueString(),
		plan.MobileConfig.ValueString(),
		plan.UserScope.ValueBool(),
		plan.AttributeSupport.ValueBool(),
		plan.EscapeAttributes.ValueBool(),
		plan.ReinstallAfterOSUpdate.ValueBool(),
		plan.Declarative.ValueBool(),
		plan.AutoRenewScepBasedCertificates.ValueBool(),
		allowedPlatforms,
		minVersion,
		maxVersion,
		plan.AllowedMacosArchitecture.ValueString(),
		plan.ID.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating profile",
			"Could not update profile, unexpected error: "+err.Error(),
		)
		return
	}

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
			"Error Deleting SimpleMDM custom profile",
			"Could not delete custom profile, unexpected error: "+err.Error(),
		)
		return
	}
}
