package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &profileResource{}
	_ resource.ResourceWithConfigure   = &profileResource{}
	_ resource.ResourceWithImportState = &profileResource{}
)

type profileResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	AutoDeploy             types.Bool   `tfsdk:"auto_deploy"`
	InstallType            types.String `tfsdk:"install_type"`
	ReinstallAfterOSUpdate types.Bool   `tfsdk:"reinstall_after_os_update"`
	ProfileIdentifier      types.String `tfsdk:"profile_identifier"`
	UserScope              types.Bool   `tfsdk:"user_scope"`
	AttributeSupport       types.Bool   `tfsdk:"attribute_support"`
	EscapeAttributes       types.Bool   `tfsdk:"escape_attributes"`
	GroupCount             types.Int64  `tfsdk:"group_count"`
	DeviceCount            types.Int64  `tfsdk:"device_count"`
	ProfileSHA             types.String `tfsdk:"profile_sha"`
	Source                 types.String `tfsdk:"source"`
	CreatedAt              types.String `tfsdk:"created_at"`
	UpdatedAt              types.String `tfsdk:"updated_at"`
}

func ProfileResource() resource.Resource {
	return &profileResource{}
}

type profileResource struct {
	client *simplemdm.Client
}

func (r *profileResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*simplemdm.Client)
}

func (r *profileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}

func (r *profileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Profile resource exposes read-only information for existing configuration profiles in SimpleMDM.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of an existing profile in SimpleMDM.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the profile.",
			},
			"auto_deploy": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the profile is auto-deployed when assigned.",
			},
			"install_type": schema.StringAttribute{
				Computed:    true,
				Description: "The install type configured for the profile.",
			},
			"reinstall_after_os_update": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the profile reinstalls automatically after macOS updates.",
			},
			"profile_identifier": schema.StringAttribute{
				Computed:    true,
				Description: "The identifier contained within the profile payload.",
			},
			"user_scope": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates if the profile installs in the user scope.",
			},
			"attribute_support": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the profile supports attribute substitution.",
			},
			"escape_attributes": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether attribute values are escaped during installation.",
			},
			"group_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of device groups currently assigned to the profile.",
			},
			"device_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of devices currently assigned to the profile.",
			},
			"profile_sha": schema.StringAttribute{
				Computed:    true,
				Description: "SHA hash reported by SimpleMDM for the profile contents.",
			},
			"source": schema.StringAttribute{
				Computed:    true,
				Description: "Origin of the profile within SimpleMDM.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the profile was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the profile was last updated.",
			},
		},
	}
}

func (r *profileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *profileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan profileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.readProfile(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating SimpleMDM profile reference",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func (r *profileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state profileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.readProfile(state.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading SimpleMDM profile",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func (r *profileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan profileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model, err := r.readProfile(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error refreshing SimpleMDM profile",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}

func (r *profileResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}

func (r *profileResource) readProfile(profileID string) (*profileResourceModel, error) {
	profile, err := r.client.ProfileGet(profileID)
	if err != nil {
		return nil, err
	}

	model := &profileResourceModel{
		ID:                     types.StringValue(strconv.Itoa(profile.Data.ID)),
		Name:                   types.StringValue(profile.Data.Attributes.Name),
		AutoDeploy:             types.BoolValue(profile.Data.Attributes.AutoDeploy),
		InstallType:            types.StringValue(profile.Data.Attributes.InstallType),
		ReinstallAfterOSUpdate: types.BoolValue(profile.Data.Attributes.ReinstallAfterOsUpdate),
		ProfileIdentifier:      types.StringValue(profile.Data.Attributes.ProfileIdentifier),
		UserScope:              types.BoolValue(profile.Data.Attributes.UserScope),
		AttributeSupport:       types.BoolValue(profile.Data.Attributes.AttributeSupport),
		EscapeAttributes:       types.BoolValue(profile.Data.Attributes.EscapeAttributes),
		GroupCount:             types.Int64Value(int64(profile.Data.Attributes.GroupCount)),
		DeviceCount:            types.Int64Value(int64(profile.Data.Attributes.DeviceCount)),
		ProfileSHA:             types.StringValue(profile.Data.Attributes.ProfileSHA),
		Source:                 types.StringValue(profile.Data.Attributes.Source),
		CreatedAt:              types.StringValue(profile.Data.Attributes.CreatedAt),
		UpdatedAt:              types.StringValue(profile.Data.Attributes.UpdatedAt),
	}

	return model, nil
}
