package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &enrollmentDataSource{}
	_ datasource.DataSourceWithConfigure = &enrollmentDataSource{}
)

type enrollmentDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	URL               types.String `tfsdk:"url"`
	UserEnrollment    types.Bool   `tfsdk:"user_enrollment"`
	WelcomeScreen     types.Bool   `tfsdk:"welcome_screen"`
	Authentication    types.Bool   `tfsdk:"authentication"`
	DeviceGroupID     types.String `tfsdk:"device_group_id"`
	AssignmentGroupID types.String `tfsdk:"assignment_group_id"`
	DeviceID          types.String `tfsdk:"device_id"`
	AccountDriven     types.Bool   `tfsdk:"account_driven"`
}

func EnrollmentDataSource() datasource.DataSource {
	return &enrollmentDataSource{}
}

type enrollmentDataSource struct {
	client *simplemdm.Client
}

func (d *enrollmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enrollment"
}

func (d *enrollmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*simplemdm.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *enrollmentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Enrollment data source retrieves details for an existing SimpleMDM enrollment. " +
			"One-time enrollments have a URL and can be used once. Account driven enrollments have a null URL and can be used multiple times.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier of the enrollment to retrieve.",
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Description: "Enrollment URL returned by SimpleMDM for one-time enrollments. Will be null for account driven enrollments.",
			},
			"user_enrollment": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether the enrollment is a user enrollment.",
			},
			"welcome_screen": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether the welcome screen is displayed during enrollment.",
			},
			"authentication": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates whether authentication is required before enrollment.",
			},
			"device_group_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the legacy device group associated with the enrollment (deprecated).",
			},
			"assignment_group_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the assignment group associated with the enrollment (for New Groups Experience).",
			},
			"device_id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the device that used the enrollment, when applicable.",
			},
			"account_driven": schema.BoolAttribute{
				Computed:    true,
				Description: "True when the enrollment is account driven (URL is null). Account driven enrollments do not support invitations.",
			},
		},
	}
}

func (d *enrollmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state enrollmentDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	enrollment, err := fetchEnrollment(ctx, d.client, state.ID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"SimpleMDM enrollment not found",
				fmt.Sprintf("The enrollment with ID %s was not found. It may have been deleted.", state.ID.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to read enrollment",
				err.Error(),
			)
		}
		return
	}

	flat := flattenEnrollment(enrollment)

	// More defensive URL null handling
	if flat.URL == nil {
		state.URL = types.StringNull()
		state.AccountDriven = types.BoolValue(true)
	} else if *flat.URL == "" {
		state.URL = types.StringNull()
		state.AccountDriven = types.BoolValue(true)
	} else {
		state.URL = types.StringValue(*flat.URL)
		state.AccountDriven = types.BoolValue(false)
	}

	state.UserEnrollment = types.BoolValue(flat.UserEnrollment)
	state.WelcomeScreen = types.BoolValue(flat.WelcomeScreen)
	state.Authentication = types.BoolValue(flat.Authentication)

	if flat.DeviceGroupID != nil {
		state.DeviceGroupID = types.StringValue(strconv.Itoa(*flat.DeviceGroupID))
	} else {
		state.DeviceGroupID = types.StringNull()
	}

	if flat.AssignmentGroupID != nil {
		state.AssignmentGroupID = types.StringValue(strconv.Itoa(*flat.AssignmentGroupID))
	} else {
		state.AssignmentGroupID = types.StringNull()
	}

	if flat.DeviceID != nil {
		state.DeviceID = types.StringValue(strconv.Itoa(*flat.DeviceID))
	} else {
		state.DeviceID = types.StringNull()
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
