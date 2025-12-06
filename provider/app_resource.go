package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &appResource{}
	_ resource.ResourceWithConfigure   = &appResource{}
	_ resource.ResourceWithImportState = &appResource{}
)

// appResourceModel maps the resource schema data.
type appResourceModel struct {
	Name                 types.String `tfsdk:"name"`
	ID                   types.String `tfsdk:"id"`
	AppStoreId           types.String `tfsdk:"app_store_id"`
	BundleId             types.String `tfsdk:"bundle_id"`
	BinaryFile           types.String `tfsdk:"binary_file"`
	DeployTo             types.String `tfsdk:"deploy_to"`
	Status               types.String `tfsdk:"status"`
	AppType              types.String `tfsdk:"app_type"`
	Version              types.String `tfsdk:"version"`
	PlatformSupport      types.String `tfsdk:"platform_support"`
	ProcessingStatus     types.String `tfsdk:"processing_status"`
	InstallationChannels types.List   `tfsdk:"installation_channels"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

func AppResource() resource.Resource {
	return &appResource{}
}

// appResource is the resource implementation.
type appResource struct {
	client *simplemdm.Client
}

// Configure adds the provider configured client to the resource.
func (r *appResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*simplemdm.Client)
}

// Metadata returns the resource type name.
func (r *appResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

// Schema defines the schema for the resource.
func (r *appResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "App resource can be used to manage Apps.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name for this app in SimpleMDM. For App Store apps, this is computed from the store. For binary uploads, you may optionally specify a name, otherwise it's extracted from the binary.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_store_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The Apple App Store ID. Required when adding App Store apps via store ID. Use either this, bundle_id, or binary_file. Example: '1090161858'",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("app_store_id"),
						path.MatchRoot("bundle_id"),
						path.MatchRoot("binary_file"),
					),
				},
			},
			"bundle_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The bundle identifier of the Apple App Store app. Required when adding App Store apps via bundle ID. Use either this, app_store_id, or binary_file. Example: com.myCompany.MyApp1",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("app_store_id"),
						path.MatchRoot("bundle_id"),
						path.MatchRoot("binary_file"),
					),
				},
			},
			"binary_file": schema.StringAttribute{
				Optional:    true,
				Description: "Path to app binary (ipa or pkg) to upload. Required when managing enterprise, custom B2B, or macOS package apps. Use either this, app_store_id, or bundle_id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("app_store_id"),
						path.MatchRoot("bundle_id"),
						path.MatchRoot("binary_file"),
					),
				},
			},
			"deploy_to": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Deploy the app after upload. Values: 'none' (default), 'outdated' (devices with older version), 'all' (all devices). Note: Only applies during updates; create operations require subsequent update. API does not return this value, so state shows last configured value.",
				Default:     stringdefault.StaticString("none"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("none", "outdated", "all"),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The current deployment status of the app. Note: This is a write-only parameter; API does not return this value.",
			},
			"app_type": schema.StringAttribute{
				Computed:    true,
				Description: "The catalog classification of the app, for example app store, enterprise, or custom b2b.",
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "The latest version reported by SimpleMDM for the app.",
			},
			"platform_support": schema.StringAttribute{
				Computed:    true,
				Description: "The platform supported by the app, such as iOS or macOS.",
			},
			"processing_status": schema.StringAttribute{
				Computed:    true,
				Description: "The current processing status of the app binary within SimpleMDM.",
			},
			"installation_channels": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The deployment channels supported by the app.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the app was added to SimpleMDM.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the app was last updated in SimpleMDM.",
			},
		},
	}
}

type appAPIResponse struct {
	Data struct {
		ID         int `json:"id"`
		Attributes struct {
			Name                 string   `json:"name"`
			BundleIdentifier     string   `json:"bundle_identifier"`
			AppType              string   `json:"app_type"`
			ITunesStoreID        int      `json:"itunes_store_id"`
			InstallationChannels []string `json:"installation_channels"`
			PlatformSupport      string   `json:"platform_support"`
			ProcessingStatus     string   `json:"processing_status"`
			Version              string   `json:"version"`
			CreatedAt            string   `json:"created_at"`
			UpdatedAt            string   `json:"updated_at"`
		} `json:"attributes"`
	} `json:"data"`
}

func newAppResourceModelFromAPI(ctx context.Context, app *appAPIResponse) (appResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Set computed fields using null for empty values to maintain proper Terraform semantics
	model := appResourceModel{
		ID:         types.StringValue(strconv.Itoa(app.Data.ID)),
		BinaryFile: types.StringNull(),
	}

	// Handle optional/computed string fields - use null when empty
	if app.Data.Attributes.Name != "" {
		model.Name = types.StringValue(app.Data.Attributes.Name)
	} else {
		model.Name = types.StringNull()
	}

	if app.Data.Attributes.BundleIdentifier != "" {
		model.BundleId = types.StringValue(app.Data.Attributes.BundleIdentifier)
	} else {
		model.BundleId = types.StringNull()
	}

	if app.Data.Attributes.AppType != "" {
		model.AppType = types.StringValue(app.Data.Attributes.AppType)
	} else {
		model.AppType = types.StringNull()
	}

	if app.Data.Attributes.Version != "" {
		model.Version = types.StringValue(app.Data.Attributes.Version)
	} else {
		model.Version = types.StringNull()
	}

	if app.Data.Attributes.PlatformSupport != "" {
		model.PlatformSupport = types.StringValue(app.Data.Attributes.PlatformSupport)
	} else {
		model.PlatformSupport = types.StringNull()
	}

	if app.Data.Attributes.ProcessingStatus != "" {
		model.ProcessingStatus = types.StringValue(app.Data.Attributes.ProcessingStatus)
	} else {
		model.ProcessingStatus = types.StringNull()
	}

	if app.Data.Attributes.CreatedAt != "" {
		model.CreatedAt = types.StringValue(app.Data.Attributes.CreatedAt)
	} else {
		model.CreatedAt = types.StringNull()
	}

	if app.Data.Attributes.UpdatedAt != "" {
		model.UpdatedAt = types.StringValue(app.Data.Attributes.UpdatedAt)
	} else {
		model.UpdatedAt = types.StringNull()
	}

	// Handle AppStoreId
	if app.Data.Attributes.ITunesStoreID != 0 {
		model.AppStoreId = types.StringValue(strconv.Itoa(app.Data.Attributes.ITunesStoreID))
	} else {
		model.AppStoreId = types.StringNull()
	}

	// DeployTo and Status are write-only parameters that API does not return
	// Preserve from plan/state as they are not provided by API
	model.DeployTo = types.StringValue("none")
	model.Status = types.StringNull()

	// Validate bundle_identifier is present (should always be returned by API)
	if app.Data.Attributes.BundleIdentifier == "" {
		diags.AddWarning(
			"Missing bundle_identifier",
			fmt.Sprintf("App %d returned by API without bundle_identifier", app.Data.ID),
		)
	}

	// Validate app_type is a known value
	knownAppTypes := map[string]bool{
		"app store":  true,
		"enterprise": true,
		"custom b2b": true,
	}
	if app.Data.Attributes.AppType != "" && !knownAppTypes[app.Data.Attributes.AppType] {
		diags.AddWarning(
			"Unknown app type",
			fmt.Sprintf("App returned unexpected app_type: %s", app.Data.Attributes.AppType),
		)
	}

	// Validate platform_support is a known value
	knownPlatforms := map[string]bool{
		"iOS":   true,
		"macOS": true,
	}
	if app.Data.Attributes.PlatformSupport != "" && !knownPlatforms[app.Data.Attributes.PlatformSupport] {
		diags.AddWarning(
			"Unknown platform",
			fmt.Sprintf("App returned unexpected platform_support: %s", app.Data.Attributes.PlatformSupport),
		)
	}

	// Handle installation channels
	if len(app.Data.Attributes.InstallationChannels) > 0 {
		listValue, listDiags := types.ListValueFrom(ctx, types.StringType, app.Data.Attributes.InstallationChannels)
		diags.Append(listDiags...)
		if !listDiags.HasError() {
			model.InstallationChannels = listValue
		}
	} else {
		model.InstallationChannels = types.ListNull(types.StringType)
	}

	return model, diags
}

func fetchApp(ctx context.Context, client *simplemdm.Client, appID string) (*appAPIResponse, error) {
	url := fmt.Sprintf("https://%s/api/v1/apps/%s", client.HostName, appID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := client.RequestResponse200(req)
	if err != nil {
		return nil, err
	}

	app := appAPIResponse{}
	if err := json.Unmarshal(body, &app); err != nil {
		return nil, err
	}

	return &app, nil
}

// isNotFoundError checks if an error is a 404 not found error
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found")
}

// waitForProcessingComplete polls until app processing is complete
func waitForProcessingComplete(ctx context.Context, client *simplemdm.Client, appID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if time.Now().After(deadline) {
				return fmt.Errorf("timeout waiting for app processing to complete after %v", timeout)
			}

			app, err := fetchApp(ctx, client, appID)
			if err != nil {
				return fmt.Errorf("error checking processing status: %w", err)
			}

			if app.Data.Attributes.ProcessingStatus == "processed" {
				return nil
			}

			// If processing failed, return error
			if app.Data.Attributes.ProcessingStatus == "failed" {
				return fmt.Errorf("app processing failed")
			}
		}
	}
}

func (r *appResource) appCreateWithBinary(ctx context.Context, binaryPath, name string) (_ *simplemdm.SimplemdmDefaultStruct, err error) {
	file, err := os.Open(binaryPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open app binary %q: %w", binaryPath, err)
	}
	defer func() {
		if cerr := file.Close(); err == nil && cerr != nil {
			err = fmt.Errorf("unable to close app binary %q: %w", binaryPath, cerr)
		}
	}()

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	part, err := writer.CreateFormFile("binary", filepath.Base(binaryPath))
	if err != nil {
		return nil, fmt.Errorf("unable to create app binary form data: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("unable to read app binary %q: %w", binaryPath, err)
	}

	if name != "" {
		if err := writer.WriteField("name", name); err != nil {
			return nil, fmt.Errorf("unable to encode app name: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("unable to finalize app upload payload: %w", err)
	}

	url := fmt.Sprintf("https://%s/api/v1/apps", r.client.HostName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	body, err := r.client.RequestResponse201(req)
	if err != nil {
		return nil, err
	}

	app := simplemdm.SimplemdmDefaultStruct{}
	if err := json.Unmarshal(body, &app); err != nil {
		return nil, err
	}

	return &app, nil
}

func (r *appResource) appUpdateWithBinary(ctx context.Context, appID, binaryPath, name, deployTo string) (err error) {
	file, err := os.Open(binaryPath)
	if err != nil {
		return fmt.Errorf("unable to open app binary %q: %w", binaryPath, err)
	}
	defer func() {
		if cerr := file.Close(); err == nil && cerr != nil {
			err = fmt.Errorf("unable to close app binary %q: %w", binaryPath, cerr)
		}
	}()

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	part, err := writer.CreateFormFile("binary", filepath.Base(binaryPath))
	if err != nil {
		return fmt.Errorf("unable to create app binary form data: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("unable to read app binary %q: %w", binaryPath, err)
	}

	if name != "" {
		if err := writer.WriteField("name", name); err != nil {
			return fmt.Errorf("unable to encode app name: %w", err)
		}
	}

	if deployTo != "" {
		if err := writer.WriteField("deploy_to", deployTo); err != nil {
			return fmt.Errorf("unable to encode deploy_to value: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("unable to finalize app update payload: %w", err)
	}

	url := fmt.Sprintf("https://%s/api/v1/apps/%s", r.client.HostName, appID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	_, err = r.client.RequestResponse200(req)
	if err != nil {
		return err
	}

	return nil
}

func (r *appResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource
func (r *appResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan appResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var appStoreId, bundleId, name, binaryPath string
	if !plan.AppStoreId.IsNull() {
		appStoreId = plan.AppStoreId.ValueString()
	}
	if !plan.BundleId.IsNull() {
		bundleId = plan.BundleId.ValueString()
	}
	if !plan.Name.IsNull() {
		name = plan.Name.ValueString()
	}
	if !plan.BinaryFile.IsNull() {
		binaryPath = plan.BinaryFile.ValueString()
	}

	// Generate API request body from plan
	var app *simplemdm.SimplemdmDefaultStruct
	var err error

	switch {
	case binaryPath != "":
		app, err = r.appCreateWithBinary(ctx, binaryPath, name)
	default:
		app, err = r.client.AppCreate(
			appStoreId,
			bundleId,
			name,
		)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating app",
			"Could not create app, unexpected error: "+err.Error(),
		)
		return
	}

	appID := strconv.Itoa(app.Data.ID)

	// If deploy_to specified at creation time, SimpleMDM requires a follow-up update.
	if !plan.DeployTo.IsNull() && plan.DeployTo.ValueString() != "" && plan.DeployTo.ValueString() != "none" {
		_, err = r.client.AppUpdate(appID, name, plan.DeployTo.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating deploy_to",
				"Failed to configure deploy_to during app creation: "+err.Error(),
			)
			return
		}
	}

	apiApp, err := fetchApp(ctx, r.client, appID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading created app",
			"Could not read newly created app "+appID+": "+err.Error(),
		)
		return
	}

	newState, diagsFromAPI := newAppResourceModelFromAPI(ctx, apiApp)
	resp.Diagnostics.Append(diagsFromAPI...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState.Name.IsNull() && !plan.Name.IsNull() {
		newState.Name = plan.Name
	}
	if newState.AppStoreId.IsNull() && !plan.AppStoreId.IsNull() {
		newState.AppStoreId = plan.AppStoreId
	}
	if newState.BundleId.IsNull() && !plan.BundleId.IsNull() {
		newState.BundleId = plan.BundleId
	}
	if (newState.DeployTo.IsNull() || newState.DeployTo.ValueString() == "") && !plan.DeployTo.IsNull() {
		newState.DeployTo = plan.DeployTo
	}
	if !plan.BinaryFile.IsNull() {
		newState.BinaryFile = plan.BinaryFile
	}

	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
}

func (r *appResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state appResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing app
	err := r.client.AppDelete(state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SimpleMDM app",
			"Could not delete app, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *appResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state appResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := fetchApp(ctx, r.client, state.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM App",
			"Could not read SimpleMDM App "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	newState, diagsFromAPI := newAppResourceModelFromAPI(ctx, app)
	resp.Diagnostics.Append(diagsFromAPI...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.BinaryFile.IsNull() {
		newState.BinaryFile = state.BinaryFile
	}

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
}

func (r *appResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state appResourceModel
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

	appID := plan.ID.ValueString()

	name := ""
	if !plan.Name.IsNull() {
		name = plan.Name.ValueString()
	}

	deployTo := ""
	if !plan.DeployTo.IsNull() {
		deployTo = plan.DeployTo.ValueString()
	}

	if !plan.BinaryFile.IsNull() && plan.BinaryFile.ValueString() != "" {
		err := r.appUpdateWithBinary(ctx, appID, plan.BinaryFile.ValueString(), name, deployTo)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating app",
				"Failed to upload new binary: "+err.Error(),
			)
			return
		}
	} else {
		_, err := r.client.AppUpdate(
			appID,
			name,
			deployTo,
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating app",
				"Failed to update app: "+err.Error(),
			)
			return
		}
	}

	apiApp, err := fetchApp(ctx, r.client, appID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading updated app",
			"Failed to refresh app state: "+err.Error(),
		)
		return
	}

	newState, diagsFromAPI := newAppResourceModelFromAPI(ctx, apiApp)
	resp.Diagnostics.Append(diagsFromAPI...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState.AppStoreId.IsNull() && !state.AppStoreId.IsNull() {
		newState.AppStoreId = state.AppStoreId
	}
	if newState.BundleId.IsNull() && !state.BundleId.IsNull() {
		newState.BundleId = state.BundleId
	}
	if newState.Name.IsNull() && !plan.Name.IsNull() {
		newState.Name = plan.Name
	}
	if (newState.DeployTo.IsNull() || newState.DeployTo.ValueString() == "") && !plan.DeployTo.IsNull() {
		newState.DeployTo = plan.DeployTo
	}
	if !plan.BinaryFile.IsNull() {
		newState.BinaryFile = plan.BinaryFile
	}

	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
}
