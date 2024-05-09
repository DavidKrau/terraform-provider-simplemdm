package provider

import (
	"context"
	"strconv"

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
	_ resource.Resource                = &profileResource{}
	_ resource.ResourceWithConfigure   = &profileResource{}
	_ resource.ResourceWithImportState = &profileResource{}
)

// profileResourceModel maps the resource schema data.
type profileResourceModel struct {
	Name                   types.String `tfsdk:"name"`
	MobileConfig           types.String `tfsdk:"mobileconfig"`
	FileSHA                types.String `tfsdk:"filesha"`
	UserScope              types.Bool   `tfsdk:"userscope"`
	AttributeSupport       types.Bool   `tfsdk:"attributesupport"`
	EscapeAttributes       types.Bool   `tfsdk:"escapeattributes"`
	ReinstallAfterOSUpdate types.Bool   `tfsdk:"reinstallafterosupdate"`
	ID                     types.String `tfsdk:"id"`
}

// ProfileResource is a helper function to simplify the provider implementation.
func CustomProfileResource() resource.Resource {
	return &profileResource{}
}

// profileResource is the resource implementation.
type profileResource struct {
	client *simplemdm.Client
}

// Configure adds the provider configured client to the resource.
func (r *profileResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*simplemdm.Client)
}

// Metadata returns the resource type name.
func (r *profileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customprofile"
}

// Schema defines the schema for the resource.
func (r *profileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Custom Profile resource can be used to manage Custom Profile. Can be used together with Device(s), Assignment Group(s) or Device Group(s) and set addition details regarding Custom Profile.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Required. A name for the profile. Example \"My First profile by terraform\"",
			},
			"mobileconfig": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Required. The mobileconfig file. Example: \"./profiles/my_first_profile.mobileconfig\" ",
			},
			"filesha": schema.StringAttribute{
				Optional:    false,
				Required:    true,
				Description: "Required. The mobileconfig file. Example: ${filesha256(\"./profiles/my_first_profile.mobileconfig\")}",
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
				Description: "Optional. A boolean true or false. When enabled, SimpleMDM will process variables in the uploaded profile. Defaults to false",
			},
			"escapeattributes": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Optional. A boolean true or false. When enabled, SimpleMDM escape the values of the custom variables in the uploaded profile. Defaults to false",
			},
			"reinstallafterosupdate": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Optional. A boolean true or false. When enabled, SimpleMDM will re-install the profile automatically after macOS software updates are detected. Defaults to false",
			},
		},
	}
}

func (r *profileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource
func (r *profileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan profileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	Profile, err := r.client.CreateProfile(plan.Name.ValueString(), plan.MobileConfig.ValueString(), plan.UserScope.ValueBool(), plan.AttributeSupport.ValueBool(), plan.EscapeAttributes.ValueBool(), plan.ReinstallAfterOSUpdate.ValueBool())
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

func (r *profileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state profileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//retrive one needs to be implemented first
	// Get refreshed profile values from SimpleMDM
	profiles, err := r.client.GetAllProfiles(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM custom profile",
			"Could not read custom profles ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	profilefound := false
	for _, profile := range profiles.Data {
		if state.ID.ValueString() == strconv.Itoa(profile.ID) {
			state.Name = types.StringValue(profile.Attributes.Name)
			state.UserScope = types.BoolValue(profile.Attributes.UserScope)
			state.AttributeSupport = types.BoolValue(profile.Attributes.AttributeSupport)
			state.EscapeAttributes = types.BoolValue(profile.Attributes.EscapeAttributes)
			state.ReinstallAfterOSUpdate = types.BoolValue(profile.Attributes.ReinstallAfterOsUpdate)
			profilefound = true
			break
		}
	}

	if !profilefound {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM custom profile",
			"Could not read custom profles ID  from array "+state.ID.ValueString(),
		)
		return
	}

	sha, _, err := r.client.GetProfileSHA(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM custom profile",
			"Could not read custom profles ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	if state.FileSHA.ValueString() != "" {
		if sha != state.FileSHA.ValueString()[0:32] {
			//fmt.Println("SHA is same")
			state.FileSHA = types.StringValue(sha)
		}
	} else {
		state.FileSHA = types.StringValue(sha)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *profileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Retrieve values from plan
	var plan profileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	_, err := r.client.UpdateProfile(plan.Name.ValueString(), plan.MobileConfig.ValueString(), plan.UserScope.ValueBool(), plan.AttributeSupport.ValueBool(), plan.EscapeAttributes.ValueBool(), plan.ReinstallAfterOSUpdate.ValueBool(), plan.FileSHA.ValueString(), plan.ID.ValueString())
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

func (r *profileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state profileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteProfile(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SimpleMDM custom profile",
			"Could not delete custom profile, unexpected error: "+err.Error(),
		)
		return
	}
}
