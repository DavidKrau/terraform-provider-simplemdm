package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	_ validator.String                 = &managedConfigValueValidator{}
)

// managedConfigValueValidator validates that value format matches value_type
type managedConfigValueValidator struct{}

func (v managedConfigValueValidator) Description(ctx context.Context) string {
	return "Validates that the value format matches the specified value_type"
}

func (v managedConfigValueValidator) MarkdownDescription(ctx context.Context) string {
	return "Validates that the value format matches the specified value_type"
}

func (v managedConfigValueValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Get value_type from config
	var valueType types.String
	diags := req.Config.GetAttribute(ctx, path.Root("value_type"), &valueType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() || valueType.IsNull() || valueType.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	valueTypeStr := valueType.ValueString()

	// Validate based on type
	switch valueTypeStr {
	case "boolean":
		if value != "0" && value != "1" {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid boolean value",
				"Boolean value must be '0' or '1'",
			)
		}
	case "integer":
		if _, err := strconv.Atoi(value); err != nil {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid integer value",
				fmt.Sprintf("Value must be a valid integer: %v", err),
			)
		}
	case "float":
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid float value",
				fmt.Sprintf("Value must be a valid float: %v", err),
			)
		}
	case "integer array":
		parts := strings.Split(value, ",")
		for _, part := range parts {
			if _, err := strconv.Atoi(strings.TrimSpace(part)); err != nil {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid integer array value",
					fmt.Sprintf("Integer array must contain comma-separated integers (e.g., '1,452,-129'): %v", err),
				)
				return
			}
		}
	case "float array":
		parts := strings.Split(value, ",")
		for _, part := range parts {
			if _, err := strconv.ParseFloat(strings.TrimSpace(part), 64); err != nil {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid float array value",
					fmt.Sprintf("Float array must contain comma-separated floats (e.g., '0.123,923.1,42'): %v", err),
				)
				return
			}
		}
	case "string array":
		// String array format: strings in quotes separated by commas (e.g., "First","Second")
		if !strings.Contains(value, "\"") {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid string array value",
				"String array must contain quoted strings separated by commas (e.g., '\"First\",\"Second\"')",
			)
		}
	case "date":
		// Date format: Timestamp with timezone (e.g., 2017-01-01T12:31:15-07:00)
		if _, err := time.Parse(time.RFC3339, value); err != nil {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid date value",
				fmt.Sprintf("Date must be in RFC3339 format with timezone (e.g., '2017-01-01T12:31:15-07:00'): %v", err),
			)
		}
	case "string":
		// String values are always valid
	default:
		// Unknown value_type - let the API reject it
	}
}

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
		Description: "Managed Config resource manages managed app configurations for a specific app and automatically pushes updates after apply.\n\n" +
			"**Note:** The SimpleMDM API does not provide an endpoint to retrieve a single managed config by ID. " +
			"Read operations fetch all configs for the app and filter to find the specific config, which may impact performance for apps with many configurations. " +
			"Since there is no UPDATE endpoint, all configuration changes require resource replacement.",
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
				},
			},
			"value": schema.StringAttribute{
				Required:    true,
				Description: "Raw value supplied to the managed configuration. Format must match value_type.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					managedConfigValueValidator{},
				},
			},
			"value_type": schema.StringAttribute{
				Required:    true,
				Description: "Data type of value accepted by the app (boolean, date, float, float array, integer, integer array, string, string array).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
			"⚠️  IMPORTANT: Managed config push failed",
			fmt.Sprintf("The managed config was created successfully in SimpleMDM but was NOT pushed to devices. "+
				"Devices will not receive this configuration until you manually push updates or fix the error and re-apply. "+
				"App ID: %s, Error: %v", plan.AppID.ValueString(), err),
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
			"⚠️  IMPORTANT: Managed config push failed",
			fmt.Sprintf("The managed config was deleted from SimpleMDM but the deletion was NOT pushed to devices. "+
				"Devices may still have the old configuration until you manually push updates or fix the error and re-apply. "+
				"App ID: %s, Error: %v", appID, err),
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

	// Verify the config exists
	config, err := fetchManagedConfig(ctx, r.client, appID, configID)
	if err != nil {
		if errors.Is(err, errManagedConfigNotFound) {
			resp.Diagnostics.AddError(
				"Managed config not found",
				fmt.Sprintf("The managed config %s for app %s does not exist", configID, appID),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error reading managed config",
				err.Error(),
			)
		}
		return
	}

	// Set full state
	compositeID := fmt.Sprintf("%s:%s", appID, configID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), compositeID)...)      //nolint:errcheck
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("app_id"), appID)...)        //nolint:errcheck
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), config.Attributes.Key)...)           //nolint:errcheck
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("value"), config.Attributes.Value)...)       //nolint:errcheck
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("value_type"), config.Attributes.ValueType)...) //nolint:errcheck
}

func parseManagedConfigID(id string) (string, string, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("expected identifier in app_id:managed_config_id format")
	}

	return parts[0], parts[1], nil
}
