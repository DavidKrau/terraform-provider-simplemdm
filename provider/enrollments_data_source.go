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
	_ datasource.DataSource              = &enrollmentsDataSource{}
	_ datasource.DataSourceWithConfigure = &enrollmentsDataSource{}
)

type enrollmentsDataSource struct {
	client *simplemdm.Client
}

type enrollmentsDataSourceModel struct {
	Enrollments []enrollmentsDataSourceEnrollmentModel `tfsdk:"enrollments"`
}

type enrollmentsDataSourceEnrollmentModel struct {
	ID             types.String `tfsdk:"id"`
	DeviceGroupID  types.String `tfsdk:"device_group_id"`
	URL            types.String `tfsdk:"url"`
	UserEnrollment types.Bool   `tfsdk:"user_enrollment"`
	WelcomeScreen  types.Bool   `tfsdk:"welcome_screen"`
	Authentication types.Bool   `tfsdk:"authentication"`
	DeviceID       types.String `tfsdk:"device_id"`
	AccountDriven  types.Bool   `tfsdk:"account_driven"`
}

func EnrollmentsDataSource() datasource.DataSource {
	return &enrollmentsDataSource{}
}

func (d *enrollmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_enrollments"
}

func (d *enrollmentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the collection of enrollments from your SimpleMDM account.",
		Blocks: map[string]schema.Block{
			"enrollments": schema.ListNestedBlock{
				Description: "Collection of enrollment records returned by the API.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Enrollment identifier.",
						},
						"device_group_id": schema.StringAttribute{
							Computed:    true,
							Description: "Device group associated with the enrollment.",
						},
						"url": schema.StringAttribute{
							Computed:    true,
							Description: "Enrollment URL for one-time enrollments.",
						},
						"user_enrollment": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this is a user enrollment.",
						},
						"welcome_screen": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether welcome screen is shown.",
						},
						"authentication": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether authentication is required.",
						},
						"device_id": schema.StringAttribute{
							Computed:    true,
							Description: "Device that used this enrollment link.",
						},
						"account_driven": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether enrollment is account driven.",
						},
					},
				},
			},
		},
	}
}

func (d *enrollmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config enrollmentsDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	enrollments, err := fetchAllEnrollments(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list SimpleMDM enrollments",
			err.Error(),
		)
		return
	}

	entries := make([]enrollmentsDataSourceEnrollmentModel, 0, len(enrollments))
	for _, enrollment := range enrollments {
		flat := flattenEnrollment(&enrollment)

		entry := enrollmentsDataSourceEnrollmentModel{
			ID:             types.StringValue(strconv.Itoa(flat.ID)),
			UserEnrollment: types.BoolValue(flat.UserEnrollment),
			WelcomeScreen:  types.BoolValue(flat.WelcomeScreen),
			Authentication: types.BoolValue(flat.Authentication),
		}

		if flat.DeviceGroupID != nil {
			entry.DeviceGroupID = types.StringValue(strconv.Itoa(*flat.DeviceGroupID))
		} else {
			entry.DeviceGroupID = types.StringNull()
		}

		if flat.URL == nil || *flat.URL == "" {
			entry.URL = types.StringNull()
			entry.AccountDriven = types.BoolValue(true)
		} else {
			entry.URL = types.StringValue(*flat.URL)
			entry.AccountDriven = types.BoolValue(false)
		}

		if flat.DeviceID != nil {
			entry.DeviceID = types.StringValue(strconv.Itoa(*flat.DeviceID))
		} else {
			entry.DeviceID = types.StringNull()
		}

		entries = append(entries, entry)
	}

	config.Enrollments = entries

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func (d *enrollmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// fetchAllEnrollments retrieves all enrollments with pagination support
func fetchAllEnrollments(ctx context.Context, client *simplemdm.Client) ([]enrollmentResponse, error) {
	return listEnrollments(ctx, client, 0)
}
