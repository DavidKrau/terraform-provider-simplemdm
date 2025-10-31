package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
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
	Type                   types.String `tfsdk:"type"`
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
	GroupIDs               types.Set    `tfsdk:"group_ids"`
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
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The profile payload type reported by SimpleMDM.",
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
			"group_ids": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "IDs of device or assignment groups currently assigned to the profile.",
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

	model, err := r.readProfile(ctx, plan.ID.ValueString())
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

	model, err := r.readProfile(ctx, state.ID.ValueString())
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

	model, err := r.readProfile(ctx, plan.ID.ValueString())
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

func (r *profileResource) readProfile(ctx context.Context, profileID string) (*profileResourceModel, error) {
	profile, err := fetchProfile(ctx, r.client, profileID)
	if err != nil {
		return nil, err
	}

	groupIDs, err := convertGroupIDs(ctx, profile.Data.Relationships)
	if err != nil {
		return nil, err
	}

	// Handle optional/computed string fields - use null when empty
	typeValue := types.StringNull()
	if profile.Data.Type != "" {
		typeValue = types.StringValue(profile.Data.Type)
	}

	installTypeValue := types.StringNull()
	if profile.Data.Attributes.InstallType != "" {
		installTypeValue = types.StringValue(profile.Data.Attributes.InstallType)
	}

	sourceValue := types.StringNull()
	if profile.Data.Attributes.Source != "" {
		sourceValue = types.StringValue(profile.Data.Attributes.Source)
	}

	createdAtValue := types.StringNull()
	if profile.Data.Attributes.CreatedAt != "" {
		createdAtValue = types.StringValue(profile.Data.Attributes.CreatedAt)
	}

	updatedAtValue := types.StringNull()
	if profile.Data.Attributes.UpdatedAt != "" {
		updatedAtValue = types.StringValue(profile.Data.Attributes.UpdatedAt)
	}

	model := &profileResourceModel{
		ID:                     types.StringValue(strconv.Itoa(profile.Data.ID)),
		Type:                   typeValue,
		Name:                   types.StringValue(profile.Data.Attributes.Name),
		AutoDeploy:             types.BoolValue(profile.Data.Attributes.AutoDeploy),
		InstallType:            installTypeValue,
		ReinstallAfterOSUpdate: types.BoolValue(profile.Data.Attributes.ReinstallAfterOSUpdate),
		ProfileIdentifier:      types.StringValue(profile.Data.Attributes.ProfileIdentifier),
		UserScope:              types.BoolValue(profile.Data.Attributes.UserScope),
		AttributeSupport:       types.BoolValue(profile.Data.Attributes.AttributeSupport),
		EscapeAttributes:       types.BoolValue(profile.Data.Attributes.EscapeAttributes),
		GroupCount:             types.Int64Value(int64(profile.Data.Attributes.GroupCount)),
		DeviceCount:            types.Int64Value(int64(profile.Data.Attributes.DeviceCount)),
		GroupIDs:               groupIDs,
		ProfileSHA:             types.StringValue(profile.Data.Attributes.ProfileSHA),
		Source:                 sourceValue,
		CreatedAt:              createdAtValue,
		UpdatedAt:              updatedAtValue,
	}

	return model, nil
}

type profileAPIResponse struct {
	Data struct {
		Type          string               `json:"type"`
		ID            int                  `json:"id"`
		Attributes    profileAttributes    `json:"attributes"`
		Relationships profileRelationships `json:"relationships"`
	} `json:"data"`
}

type profileAttributes struct {
	Name                   string `json:"name"`
	AutoDeploy             bool   `json:"auto_deploy"`
	InstallType            string `json:"install_type"`
	ReinstallAfterOSUpdate bool   `json:"reinstall_after_os_update"`
	ProfileIdentifier      string `json:"profile_identifier"`
	UserScope              bool   `json:"user_scope"`
	AttributeSupport       bool   `json:"attribute_support"`
	EscapeAttributes       bool   `json:"escape_attributes"`
	GroupCount             int    `json:"group_count"`
	DeviceCount            int    `json:"device_count"`
	ProfileSHA             string `json:"profile_sha"`
	Source                 string `json:"source"`
	CreatedAt              string `json:"created_at"`
	UpdatedAt              string `json:"updated_at"`
}

type profileRelationships struct {
	DeviceGroups relationshipCollection `json:"device_groups"`
	Groups       relationshipCollection `json:"groups"`
}

type relationshipCollection struct {
	Data []relationshipReference `json:"data"`
}

type relationshipReference struct {
	ID int `json:"id"`
}

func fetchProfile(ctx context.Context, client *simplemdm.Client, profileID string) (*profileAPIResponse, error) {
	url := fmt.Sprintf("https://%s/api/v1/profiles/%s", client.HostName, profileID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := client.RequestResponse200(req)
	if err != nil {
		return nil, err
	}

	var profile profileAPIResponse
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func convertGroupIDs(ctx context.Context, relationships profileRelationships) (types.Set, error) {
	unique := make(map[int]struct{})
	for _, item := range relationships.DeviceGroups.Data {
		unique[item.ID] = struct{}{}
	}
	for _, item := range relationships.Groups.Data {
		unique[item.ID] = struct{}{}
	}

	if len(unique) == 0 {
		return types.SetNull(types.StringType), nil
	}

	ids := make([]string, 0, len(unique))
	for id := range unique {
		ids = append(ids, strconv.Itoa(id))
	}
	sort.Strings(ids)

	value, diags := types.SetValueFrom(ctx, types.StringType, ids)
	if diags.HasError() {
		return types.SetNull(types.StringType), fmt.Errorf("unable to convert profile group IDs: %s", diags)
	}

	return value, nil
}
