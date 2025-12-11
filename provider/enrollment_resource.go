package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	AssignmentGroupID types.String `tfsdk:"assignment_group_id"`
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
		Description: "Enrollment resource manages SimpleMDM enrollment links. There are two types: " +
			"One-time enrollments (with a URL) can be used once by a single device and support sending invitations. " +
			"Account driven enrollments (URL is null) can be used multiple times but do not support invitations. " +
			"Note: You must specify either device_group_id (for legacy device groups) or assignment_group_id (for modern assignment groups).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Numeric identifier of the enrollment.",
			},
			"device_group_id": schema.StringAttribute{
				Optional:    true,
				Description: "Identifier of the legacy device group (deprecated). For accounts using the New Groups Experience, use assignment_group_id instead. Either device_group_id or assignment_group_id must be specified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("assignment_group_id")),
					stringvalidator.AtLeastOneOf(path.MatchRoot("assignment_group_id")),
				},
			},
			"assignment_group_id": schema.StringAttribute{
				Optional:    true,
				Description: "Identifier of the assignment group to associate with this enrollment. For accounts using the New Groups Experience. Either device_group_id or assignment_group_id must be specified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("device_group_id")),
					stringvalidator.AtLeastOneOf(path.MatchRoot("device_group_id")),
				},
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Description: "Enrollment URL returned by SimpleMDM for one-time enrollments. Will be null for account driven enrollments.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_enrollment": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "When true, creates a user enrollment instead of a device enrollment. This setting cannot be changed after creation and will force a new resource if modified.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"welcome_screen": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Controls whether the welcome screen is shown to end users during enrollment. This setting cannot be changed after creation and will force a new resource if modified.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"authentication": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Requires authentication before enrollment can proceed. This setting cannot be changed after creation and will force a new resource if modified.",
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
				Description: "True when the enrollment is account driven (URL is null). Account driven enrollments do not support sending invitations.",
			},
			"invitation_contact": schema.StringAttribute{
				Optional:    true,
				Description: "Email address or phone number to send an enrollment invitation to after creation or when updated. Phone numbers should be prefixed with + for international numbers. Note: This is write-only - the API does not return invitation history, so Terraform cannot detect if invitations were sent outside of Terraform. Only works for one-time enrollments (not account driven).",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(3),
				},
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
		DeviceGroupID:     plan.DeviceGroupID.ValueString(),
		AssignmentGroupID: plan.AssignmentGroupID.ValueString(),
		UserEnrollment:    boolPointer(plan.UserEnrollment),
		WelcomeScreen:     boolPointer(plan.WelcomeScreen),
		Authentication:    boolPointer(plan.Authentication),
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
		// Check if this is an account driven enrollment before sending invitation
		if state.AccountDriven.ValueBool() {
			resp.Diagnostics.AddError(
				"Cannot send invitation for account driven enrollment",
				"Account driven enrollments (where URL is null) do not support sending invitations. Only one-time enrollments with a URL can have invitations sent.",
			)
			return
		}
		
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
		if isNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
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
		// Check if this is an account driven enrollment before sending invitation
		if state.AccountDriven.ValueBool() {
			resp.Diagnostics.AddError(
				"Cannot send invitation for account driven enrollment",
				"Account driven enrollments (where URL is null) do not support sending invitations. Only one-time enrollments with a URL can have invitations sent.",
			)
			return
		}
		
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

	if flat.AssignmentGroupID != nil {
		model.AssignmentGroupID = types.StringValue(strconv.Itoa(*flat.AssignmentGroupID))
	} else {
		model.AssignmentGroupID = types.StringNull()
	}

	// More defensive URL null handling
	if flat.URL == nil {
		model.URL = types.StringNull()
		model.AccountDriven = types.BoolValue(true)
	} else if *flat.URL == "" {
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
