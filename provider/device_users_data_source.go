package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/DavidKrau/terraform-provider-simplemdm/internal/simplemdmext"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &deviceUsersDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceUsersDataSource{}
)

type deviceUsersDataSource struct {
	client *simplemdm.Client
}

type deviceUsersDataSourceModel struct {
	DeviceID types.String                   `tfsdk:"device_id"`
	Users    []deviceUserDataSourceModel    `tfsdk:"users"`
}

type deviceUserDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Type          types.String `tfsdk:"type"`
	Username      types.String `tfsdk:"username"`
	FullName      types.String `tfsdk:"full_name"`
	UID           types.Int64  `tfsdk:"uid"`
	UserGUID      types.String `tfsdk:"user_guid"`
	DataQuota     types.Int64  `tfsdk:"data_quota"`
	DataUsed      types.Int64  `tfsdk:"data_used"`
	DataToSync    types.Bool   `tfsdk:"data_to_sync"`
	SecureToken   types.Bool   `tfsdk:"secure_token"`
	LoggedIn      types.Bool   `tfsdk:"logged_in"`
	MobileAccount types.Bool   `tfsdk:"mobile_account"`
}

func DeviceUsersDataSource() datasource.DataSource {
	return &deviceUsersDataSource{}
}

func (d *deviceUsersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_users"
}

func (d *deviceUsersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves user accounts currently configured on a macOS device. This endpoint only works for macOS devices and will fail for iOS/tvOS devices.",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier of the device. Note: This endpoint only supports macOS devices.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d+$`),
						"device_id must be a numeric string",
					),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"users": schema.ListNestedBlock{
				Description: "Collection of user accounts present on the device.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "User identifier.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "User resource type.",
						},
						"username": schema.StringAttribute{
							Computed:    true,
							Description: "Username of the account.",
						},
						"full_name": schema.StringAttribute{
							Computed:    true,
							Description: "Full name of the user.",
						},
						"uid": schema.Int64Attribute{
							Computed:    true,
							Description: "User ID (UID) of the account.",
						},
						"user_guid": schema.StringAttribute{
							Computed:    true,
							Description: "User GUID of the account.",
						},
						"data_quota": schema.Int64Attribute{
							Computed:    true,
							Description: "Data quota for the user in bytes.",
						},
						"data_used": schema.Int64Attribute{
							Computed:    true,
							Description: "Data used by the user in bytes.",
						},
						"data_to_sync": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the user has data to sync.",
						},
						"secure_token": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the user has a secure token.",
						},
						"logged_in": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the user is currently logged in.",
						},
						"mobile_account": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the user is a mobile account.",
						},
					},
				},
			},
		},
	}
}

func (d *deviceUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deviceUsersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	users, err := simplemdmext.ListDeviceUsers(ctx, d.client, state.DeviceID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Device not found",
				fmt.Sprintf("The device with ID %s was not found. It may have been deleted.", state.DeviceID.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to list device users",
				fmt.Sprintf("Failed to retrieve users for device %s. Note: This endpoint only works for macOS devices. Error: %s", state.DeviceID.ValueString(), err.Error()),
			)
		}
		return
	}

	items := make([]deviceUserDataSourceModel, 0, len(users.Data))
	for _, item := range users.Data {
		user := deviceUserDataSourceModel{
			ID:   types.StringValue(item.ID.String()),
			Type: types.StringValue(item.Type),
		}

		// Parse attributes from the map
		if username, ok := item.Attributes["username"].(string); ok {
			user.Username = types.StringValue(username)
		} else {
			user.Username = types.StringNull()
		}

		if fullName, ok := item.Attributes["full_name"].(string); ok {
			user.FullName = types.StringValue(fullName)
		} else {
			user.FullName = types.StringNull()
		}

		if uid, ok := item.Attributes["uid"].(float64); ok {
			user.UID = types.Int64Value(int64(uid))
		} else {
			user.UID = types.Int64Null()
		}

		if userGUID, ok := item.Attributes["user_guid"].(string); ok {
			user.UserGUID = types.StringValue(userGUID)
		} else {
			user.UserGUID = types.StringNull()
		}

		if dataQuota, ok := item.Attributes["data_quota"].(float64); ok {
			user.DataQuota = types.Int64Value(int64(dataQuota))
		} else {
			user.DataQuota = types.Int64Null()
		}

		if dataUsed, ok := item.Attributes["data_used"].(float64); ok {
			user.DataUsed = types.Int64Value(int64(dataUsed))
		} else {
			user.DataUsed = types.Int64Null()
		}

		if dataToSync, ok := item.Attributes["data_to_sync"].(bool); ok {
			user.DataToSync = types.BoolValue(dataToSync)
		} else {
			user.DataToSync = types.BoolNull()
		}

		if secureToken, ok := item.Attributes["secure_token"].(bool); ok {
			user.SecureToken = types.BoolValue(secureToken)
		} else {
			user.SecureToken = types.BoolNull()
		}

		if loggedIn, ok := item.Attributes["logged_in"].(bool); ok {
			user.LoggedIn = types.BoolValue(loggedIn)
		} else {
			user.LoggedIn = types.BoolNull()
		}

		if mobileAccount, ok := item.Attributes["mobile_account"].(bool); ok {
			user.MobileAccount = types.BoolValue(mobileAccount)
		} else {
			user.MobileAccount = types.BoolNull()
		}

		items = append(items, user)
	}

	state.Users = items

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *deviceUsersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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