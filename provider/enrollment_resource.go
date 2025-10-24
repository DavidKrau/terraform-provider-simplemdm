package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &enrollmentResource{}
	_ resource.ResourceWithConfigure   = &enrollmentResource{}
	_ resource.ResourceWithImportState = &enrollmentResource{}
)

type enrollmentResourceModel struct {
	ID                types.String `tfsdk:"id"`
	DeviceGroupID     types.String `tfsdk:"device_group_id"`
	URL               types.String `tfsdk:"url"`
	UserEnrollment    types.Bool   `tfsdk:"user_enrollment"`
	WelcomeScreen     types.Bool   `tfsdk:"welcome_screen"`
	Authentication    types.Bool   `tfsdk:"authentication"`
	DeviceID          types.String `tfsdk:"device_id"`
	AccountDriven     types.Bool   `tfsdk:"account_driven"`
	InvitationContact types.String `tfsdk:"invitation_contact"`
}

func EnrollmentResource() resource.Resource {
	return &enrollmentResource{}
}

type enrollmentResource struct {
	client *simplemdm.Client
}

func (r *enrollmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enrollment"
}

func (r *enrollmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*simplemdm.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *enrollmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Enrollment resource manages SimpleMDM enrollment links used for one-time and account driven enrollments.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Numeric identifier of the enrollment.",
			},
			"device_group_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier of the device group associated with the enrollment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Description: "Enrollment URL returned by SimpleMDM for one-time enrollments.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_enrollment": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "When true, creates a user enrollment instead of a device enrollment.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"welcome_screen": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Controls whether the welcome screen is shown to end users during enrollment.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"authentication": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Requires authentication before enrollment can proceed.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"device_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the device that used this enrollment link, when applicable.",
			},
			"account_driven": schema.BoolAttribute{
				Computed:    true,
				Description: "True when the enrollment is account driven (URL is null).",
			},
			"invitation_contact": schema.StringAttribute{
				Optional:    true,
				Description: "Email address or phone number to send an enrollment invitation to after creation or when updated.",
			},
		},
	}
}

func (r *enrollmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan enrollmentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := enrollmentUpsertRequest{
		DeviceGroupID:  plan.DeviceGroupID.ValueString(),
		UserEnrollment: boolPointer(plan.UserEnrollment),
		WelcomeScreen:  boolPointer(plan.WelcomeScreen),
		Authentication: boolPointer(plan.Authentication),
	}

	created, err := createEnrollment(ctx, r.client, payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating enrollment",
			err.Error(),
		)
		return
	}

	flat := flattenEnrollment(created)
	state := enrollmentResourceModel{}
	applyEnrollmentFlat(&state, flat)

	if !plan.InvitationContact.IsNull() && !plan.InvitationContact.IsUnknown() {
		if err := sendEnrollmentInvitation(ctx, r.client, state.ID.ValueString(), plan.InvitationContact.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Error sending enrollment invitation",
				err.Error(),
			)
			return
		}
		state.InvitationContact = plan.InvitationContact
	} else {
		state.InvitationContact = types.StringNull()
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *enrollmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state enrollmentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	enrollment, err := fetchEnrollment(ctx, r.client, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading enrollment",
			err.Error(),
		)
		return
	}

	flat := flattenEnrollment(enrollment)
	updated := enrollmentResourceModel{}
	applyEnrollmentFlat(&updated, flat)

	// Preserve invitation contact stored in state since API does not expose it.
	updated.InvitationContact = state.InvitationContact

	diags = resp.State.Set(ctx, &updated)
	resp.Diagnostics.Append(diags...)
}

func (r *enrollmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan enrollmentResourceModel
	var state enrollmentResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.InvitationContact.IsNull() || plan.InvitationContact.IsUnknown() {
		state.InvitationContact = types.StringNull()
	} else if plan.InvitationContact.ValueString() != state.InvitationContact.ValueString() {
		if err := sendEnrollmentInvitation(ctx, r.client, state.ID.ValueString(), plan.InvitationContact.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Error sending enrollment invitation",
				err.Error(),
			)
			return
		}
		state.InvitationContact = plan.InvitationContact
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *enrollmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state enrollmentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := deleteEnrollment(ctx, r.client, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting enrollment",
			err.Error(),
		)
		return
	}
}

func (r *enrollmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func applyEnrollmentFlat(model *enrollmentResourceModel, flat enrollmentFlat) {
	model.ID = types.StringValue(strconv.Itoa(flat.ID))

	if flat.DeviceGroupID != nil {
		model.DeviceGroupID = types.StringValue(strconv.Itoa(*flat.DeviceGroupID))
	} else {
		model.DeviceGroupID = types.StringNull()
	}

	if flat.URL == nil || *flat.URL == "" {
		model.URL = types.StringNull()
		model.AccountDriven = types.BoolValue(true)
	} else {
		model.URL = types.StringValue(*flat.URL)
		model.AccountDriven = types.BoolValue(false)
	}

	model.UserEnrollment = types.BoolValue(flat.UserEnrollment)
	model.WelcomeScreen = types.BoolValue(flat.WelcomeScreen)
	model.Authentication = types.BoolValue(flat.Authentication)

	if flat.DeviceID != nil {
		model.DeviceID = types.StringValue(strconv.Itoa(*flat.DeviceID))
	} else {
		model.DeviceID = types.StringNull()
	}
}

func boolPointer(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	result := value.ValueBool()
	return &result
}
