package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &managedConfigResource{}
	_ resource.ResourceWithConfigure   = &managedConfigResource{}
	_ resource.ResourceWithImportState = &managedConfigResource{}
)

type managedConfigModel struct {
	ID        types.String `tfsdk:"id"`
	AppID     types.String `tfsdk:"app_id"`
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	ValueType types.String `tfsdk:"value_type"`
}

type managedConfigResource struct {
	client *simplemdm.Client
}

func ManagedConfigResource() resource.Resource {
	return &managedConfigResource{}
}

func (r *managedConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_config"
}

func (r *managedConfigResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Managed Config resource manages managed app configurations for a specific app and automatically pushes updates after apply.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Composite identifier of the managed app configuration in the format &lt;app_id&gt;:&lt;managed_config_id&gt;.",
			},
			"app_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the app that owns the managed configuration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key": schema.StringAttribute{
				Required:    true,
				Description: "Configuration key as defined by the managed app schema.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"value": schema.StringAttribute{
				Required:    true,
				Description: "Raw value supplied to the managed configuration. Format must match value_type.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"value_type": schema.StringAttribute{
				Required:    true,
				Description: "Data type of value accepted by the app (boolean, date, float, float array, integer, integer array, string, string array).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"boolean",
						"date",
						"float",
						"float array",
						"integer",
						"integer array",
						"string",
						"string array",
					),
				},
			},
		},
	}
}

func (r *managedConfigResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*simplemdm.Client)
	if ok {
		r.client = client
	}
}

func (r *managedConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan managedConfigModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := createManagedConfig(ctx, r.client, plan.AppID.ValueString(), plan.Key.ValueString(), plan.Value.ValueString(), plan.ValueType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating managed app configuration",
			err.Error(),
		)
		return
	}

	if err := pushManagedConfigUpdates(ctx, r.client, plan.AppID.ValueString()); err != nil {
		resp.Diagnostics.AddWarning(
			"Managed config push failed",
			fmt.Sprintf("Failed to push managed config updates for app %s: %v", plan.AppID.ValueString(), err),
		)
	}

	compositeID := fmt.Sprintf("%s:%d", plan.AppID.ValueString(), created.ID)
	plan.ID = types.StringValue(compositeID)
	plan.Key = types.StringValue(created.Attributes.Key)
	plan.Value = types.StringValue(created.Attributes.Value)
	plan.ValueType = types.StringValue(created.Attributes.ValueType)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *managedConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state managedConfigModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	appID, configID, err := parseManagedConfigID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid managed config identifier", err.Error())
		return
	}

	config, err := fetchManagedConfig(ctx, r.client, appID, configID)
	if err != nil {
		if errors.Is(err, errManagedConfigNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading managed app configuration",
			err.Error(),
		)
		return
	}

	state.AppID = types.StringValue(appID)
	state.Key = types.StringValue(config.Attributes.Key)
	state.Value = types.StringValue(config.Attributes.Value)
	state.ValueType = types.StringValue(config.Attributes.ValueType)
	state.ID = types.StringValue(fmt.Sprintf("%s:%d", appID, config.ID))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *managedConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Managed config update not supported",
		"Managed app configurations must be replaced to change values.",
	)
}

func (r *managedConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state managedConfigModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	appID, configID, err := parseManagedConfigID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid managed config identifier", err.Error())
		return
	}

	if err := deleteManagedConfig(ctx, r.client, appID, configID); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting managed app configuration",
			err.Error(),
		)
		return
	}

	if err := pushManagedConfigUpdates(ctx, r.client, appID); err != nil {
		resp.Diagnostics.AddWarning(
			"Managed config push failed",
			fmt.Sprintf("Failed to push managed config updates for app %s: %v", appID, err),
		)
	}
}

func (r *managedConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	appID, configID, err := parseManagedConfigID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected import identifier format",
			"Expected app_id:managed_config_id",
		)
		return
	}

	compositeID := fmt.Sprintf("%s:%s", appID, configID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), compositeID)...) //nolint:errcheck
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("app_id"), appID)...)   //nolint:errcheck
}

func parseManagedConfigID(id string) (string, string, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("expected identifier in app_id:managed_config_id format")
	}

	return parts[0], parts[1], nil
}
