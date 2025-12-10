package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &deviceCommandResource{}
	_ resource.ResourceWithConfigure   = &deviceCommandResource{}
	_ resource.ResourceWithImportState = &deviceCommandResource{}
)

type deviceCommandResource struct {
	client *simplemdm.Client
}

type deviceCommandResourceModel struct {
	ID         types.String `tfsdk:"id"`
	DeviceID   types.String `tfsdk:"device_id"`
	Command    types.String `tfsdk:"command"`
	Parameters types.Map    `tfsdk:"parameters"`
	StatusCode types.Int64  `tfsdk:"status_code"`
	Response   types.String `tfsdk:"response"`
}

type deviceCommandSpec struct {
	method         string
	pathTemplate   string
	expectedStatus int
}

const (
	deviceCommandEndpointFormat   = "https://%s/api/v1/devices/%s/%s"
	contentTypeFormURLEncoded     = "application/x-www-form-urlencoded"
	deviceCommandIDFormatTemplate = "%s:%s:%d"
)

// deviceCommandCatalog maps command names to their API specifications
// Command name mappings (provider name -> API endpoint):
//   enable_bluetooth       -> bluetooth (POST)
//   disable_bluetooth      -> bluetooth (DELETE)
//   enable_remote_desktop  -> remote_desktop (POST)
//   disable_remote_desktop -> remote_desktop (DELETE)
//   rotate_filevault_recovery_key -> rotate_filevault_key
var deviceCommandCatalog = map[string]deviceCommandSpec{
	"push_assigned_apps":            {method: http.MethodPost, pathTemplate: "push_apps", expectedStatus: http.StatusAccepted},
	"refresh":                       {method: http.MethodPost, pathTemplate: "refresh", expectedStatus: http.StatusAccepted},
	"restart":                       {method: http.MethodPost, pathTemplate: "restart", expectedStatus: http.StatusAccepted},
	"shutdown":                      {method: http.MethodPost, pathTemplate: "shutdown", expectedStatus: http.StatusAccepted},
	"lock":                          {method: http.MethodPost, pathTemplate: "lock", expectedStatus: http.StatusAccepted},
	"clear_passcode":                {method: http.MethodPost, pathTemplate: "clear_passcode", expectedStatus: http.StatusAccepted},
	"clear_firmware_password":       {method: http.MethodPost, pathTemplate: "clear_firmware_password", expectedStatus: http.StatusAccepted},
	"rotate_firmware_password":      {method: http.MethodPost, pathTemplate: "rotate_firmware_password", expectedStatus: http.StatusAccepted},
	"clear_recovery_lock_password":  {method: http.MethodPost, pathTemplate: "clear_recovery_lock_password", expectedStatus: http.StatusAccepted},
	"clear_restrictions_password":   {method: http.MethodPost, pathTemplate: "clear_restrictions_password", expectedStatus: http.StatusAccepted},
	"rotate_recovery_lock_password": {method: http.MethodPost, pathTemplate: "rotate_recovery_lock_password", expectedStatus: http.StatusAccepted},
	"rotate_filevault_recovery_key": {method: http.MethodPost, pathTemplate: "rotate_filevault_key", expectedStatus: http.StatusAccepted},
	"set_admin_password":            {method: http.MethodPost, pathTemplate: "set_admin_password", expectedStatus: http.StatusAccepted},
	"rotate_admin_password":         {method: http.MethodPost, pathTemplate: "rotate_admin_password", expectedStatus: http.StatusAccepted},
	"wipe":                          {method: http.MethodPost, pathTemplate: "wipe", expectedStatus: http.StatusAccepted},
	"update_os":                     {method: http.MethodPost, pathTemplate: "update_os", expectedStatus: http.StatusAccepted},
	"enable_remote_desktop":         {method: http.MethodPost, pathTemplate: "remote_desktop", expectedStatus: http.StatusAccepted},
	"disable_remote_desktop":        {method: http.MethodDelete, pathTemplate: "remote_desktop", expectedStatus: http.StatusAccepted},
	"enable_bluetooth":              {method: http.MethodPost, pathTemplate: "bluetooth", expectedStatus: http.StatusAccepted},
	"disable_bluetooth":             {method: http.MethodDelete, pathTemplate: "bluetooth", expectedStatus: http.StatusAccepted},
	"set_time_zone":                 {method: http.MethodPost, pathTemplate: "set_time_zone", expectedStatus: http.StatusNoContent},
	"unenroll":                      {method: http.MethodPost, pathTemplate: "unenroll", expectedStatus: http.StatusAccepted},
	"delete_user":                   {method: http.MethodDelete, pathTemplate: "users/{user_id}", expectedStatus: http.StatusAccepted},
	// Lost Mode commands (BUG-DC-001)
	"enable_lost_mode":          {method: http.MethodPost, pathTemplate: "lost_mode", expectedStatus: http.StatusAccepted},
	"disable_lost_mode":         {method: http.MethodDelete, pathTemplate: "lost_mode", expectedStatus: http.StatusAccepted},
	"lost_mode_play_sound":      {method: http.MethodPost, pathTemplate: "lost_mode/play_sound", expectedStatus: http.StatusAccepted},
	"lost_mode_update_location": {method: http.MethodPost, pathTemplate: "lost_mode/update_location", expectedStatus: http.StatusAccepted},
}

func DeviceCommandResource() resource.Resource {
	return &deviceCommandResource{}
}

func (r *deviceCommandResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_command"
}

func (r *deviceCommandResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Executes management commands against a SimpleMDM device. Commands are executed immediately during resource creation and cannot be reversed by removing the resource. This is a fire-and-forget operation - removing the resource from Terraform state does not undo the command.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Internal identifier for the executed command.",
			},
			"device_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Identifier of the target device.",
			},
			"command": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: `Command to execute. Supported commands:

Basic Commands:
  - push_assigned_apps: Push all assigned apps to device
  - refresh: Request device inventory refresh (may be rate limited with HTTP 429)
  - restart: Restart device (params: rebuild_kernel_cache, notify_user as strings 'true'/'false')
  - shutdown: Shut down device
  - unenroll: Unenroll device from MDM

Security Commands:
  - lock: Lock device (params: message, phone_number, pin [required for macOS])
  - clear_passcode: Clear device passcode
  - wipe: Erase device (params: pin [required for macOS without T2 chip])
  
Lost Mode Commands:
  - enable_lost_mode: Enable lost mode (params: message, phone_number, footnote)
  - disable_lost_mode: Disable lost mode
  - lost_mode_play_sound: Play sound on device in lost mode
  - lost_mode_update_location: Update device location in lost mode

Password Management:
  - clear_firmware_password: Clear firmware password (macOS)
  - rotate_firmware_password: Rotate firmware password (macOS)
  - clear_recovery_lock_password: Clear recovery lock password
  - rotate_recovery_lock_password: Rotate recovery lock password
  - rotate_filevault_recovery_key: Rotate FileVault key (macOS)
  - set_admin_password: Set admin password (params: new_password [required])
  - rotate_admin_password: Rotate admin password
  - clear_restrictions_password: Clear restrictions password (iOS)

System Configuration:
  - update_os: Update OS (params: os_update_mode [required for macOS], version_type [optional])
  - set_time_zone: Set time zone (params: time_zone [required])
  - enable_remote_desktop: Enable Remote Desktop (macOS 10.14.4+)
  - disable_remote_desktop: Disable Remote Desktop
  - enable_bluetooth: Enable Bluetooth
  - disable_bluetooth: Disable Bluetooth
  - delete_user: Delete user account (params: user_id [required, used in path])`,
			},
			"parameters": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
				Description: "Optional parameters to pass to the API call. Boolean parameters should be specified as 'true' or 'false' strings. Required parameters vary by command - see command description for details.",
			},
			"status_code": schema.Int64Attribute{
				Computed:    true,
				Description: "HTTP status code returned by the SimpleMDM API. Note: A 202 (Accepted) status indicates the command was queued, not that it completed successfully on the device. A 204 (No Content) indicates successful completion of commands like set_time_zone.",
			},
			"response": schema.StringAttribute{
				Computed:    true,
				Description: "Raw response payload, if any, returned by the API.",
			},
		},
	}
}

func (r *deviceCommandResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*simplemdm.Client)
	if !ok {
		return
	}

	r.client = client
}

func (r *deviceCommandResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Unable to create device command because the client was not configured")
		return
	}

	plan, ok := r.readCreatePlan(ctx, req, resp)
	if !ok {
		return
	}

	commandKey := plan.Command.ValueString()
	spec, ok := resolveDeviceCommandSpec(commandKey, resp)
	if !ok {
		return
	}

	params, ok := decodeCommandParameters(ctx, plan.Parameters, resp)
	if !ok {
		return
	}

	// BUG-DC-002: Validate required parameters
	if err := validateRequiredParameters(commandKey, params); err != nil {
		resp.Diagnostics.AddError("Missing required parameter", err.Error())
		return
	}

	pathFragment, consumedKeys, err := expandCommandPath(spec.pathTemplate, params)
	if err != nil {
		resp.Diagnostics.AddError("Invalid command parameters", err.Error())
		return
	}

	removeConsumedParameters(params, consumedKeys)

	reqObj, err := r.buildCommandRequest(ctx, spec.method, plan.DeviceID.ValueString(), pathFragment, params)
	if err != nil {
		resp.Diagnostics.AddError("Error creating request", err.Error())
		return
	}

	// BUG-DC-004: Handle rate limiting for refresh command
	body, err := r.executeCommandWithRetry(ctx, reqObj, spec.expectedStatus, commandKey)
	if err != nil {
		resp.Diagnostics.AddError("Error executing device command", err.Error())
		return
	}

	r.updateCreateState(ctx, plan, commandKey, spec.expectedStatus, body, resp)
}

func (r *deviceCommandResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state deviceCommandResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Commands are fire-and-forget; retain state as-is.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *deviceCommandResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Device commands cannot be updated", "Remove the resource and recreate it to issue another command.")
}

// Delete implements resource.Resource.
// Device commands are fire-and-forget operations that cannot be reversed.
// Removing the resource from Terraform state does not undo the command.
// To reverse a command's effects, create a new device_command resource
// with the opposite command (e.g., disable_bluetooth after enable_bluetooth).
func (r *deviceCommandResource) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// No-op: commands cannot be undone
}

// ImportState implements resource.ResourceWithImportState.
// Device commands cannot be meaningfully imported because they are fire-and-forget
// operations that execute once and cannot be read back from the API.
func (r *deviceCommandResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.AddError(
		"Import not supported",
		"Device commands are fire-and-forget operations and cannot be imported. "+
			"Commands must be defined in Terraform configuration to be executed.",
	)
}

func expandCommandPath(template string, params map[string]string) (string, []string, error) {
	expanded := template
	consumed := make([]string, 0)
	if strings.Contains(template, "{") {
		for key, value := range params {
			placeholder := fmt.Sprintf("{%s}", key)
			if strings.Contains(expanded, placeholder) {
				expanded = strings.ReplaceAll(expanded, placeholder, url.PathEscape(value))
				consumed = append(consumed, key)
			}
		}

		if strings.Contains(expanded, "{") {
			return "", nil, fmt.Errorf("missing required parameter for path template %q", template)
		}
	}

	return expanded, consumed, nil
}

func (r *deviceCommandResource) readCreatePlan(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) (*deviceCommandResourceModel, bool) {
	var plan deviceCommandResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return nil, false
	}

	return &plan, true
}

func resolveDeviceCommandSpec(commandKey string, resp *resource.CreateResponse) (deviceCommandSpec, bool) {
	spec, ok := deviceCommandCatalog[commandKey]
	if !ok {
		resp.Diagnostics.AddError(
			"Unsupported device command",
			fmt.Sprintf("Command %q is not currently supported by the provider", commandKey),
		)
		return deviceCommandSpec{}, false
	}

	return spec, true
}

func decodeCommandParameters(ctx context.Context, parameters types.Map, resp *resource.CreateResponse) (map[string]string, bool) {
	if parameters.IsNull() || parameters.IsUnknown() {
		return map[string]string{}, true
	}

	result := make(map[string]string)
	diags := parameters.ElementsAs(ctx, &result, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return nil, false
	}

	for key, value := range result {
		if strings.TrimSpace(key) == "" || value == "" {
			delete(result, key)
		}
	}

	return result, true
}

func removeConsumedParameters(params map[string]string, keys []string) {
	for _, key := range keys {
		delete(params, key)
	}
}

func (r *deviceCommandResource) buildCommandRequest(ctx context.Context, method, deviceID, pathFragment string, params map[string]string) (*http.Request, error) {
	endpoint := fmt.Sprintf(deviceCommandEndpointFormat, r.client.HostName, deviceID, pathFragment)

	bodyReader, hasBody := prepareCommandBody(method, params)

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bodyReader)
	if err != nil {
		return nil, err
	}

	if hasBody {
		req.Header.Set("Content-Type", contentTypeFormURLEncoded)
	}

	return req, nil
}

func prepareCommandBody(method string, params map[string]string) (*strings.Reader, bool) {
	if method != http.MethodPost || len(params) == 0 {
		return nil, false
	}

	requestBody := url.Values{}
	for key, value := range params {
		requestBody.Set(key, value)
	}

	return strings.NewReader(requestBody.Encode()), true
}

// validateRequiredParameters validates that required parameters are present for commands that need them.
// BUG-DC-002: Add validation for required parameters
func validateRequiredParameters(command string, params map[string]string) error {
	requiredParams := map[string][]string{
		"set_admin_password": {"new_password"},
		"set_time_zone":      {"time_zone"},
		"delete_user":        {"user_id"},
		// Note: lock (macOS), wipe (macOS without T2), and update_os (macOS) have platform-specific
		// required parameters that cannot be validated without device information
	}

	if required, ok := requiredParams[command]; ok {
		for _, param := range required {
			if _, exists := params[param]; !exists {
				return fmt.Errorf("command %q requires parameter %q", command, param)
			}
		}
	}
	return nil
}

func (r *deviceCommandResource) executeCommand(req *http.Request, expectedStatus int) ([]byte, error) {
	switch expectedStatus {
	case http.StatusAccepted:
		return r.client.RequestResponse202(req)
	case http.StatusNoContent:
		return r.client.RequestResponse204(req)
	default:
		return r.client.RequestResponse200(req)
	}
}

// executeCommandWithRetry executes a command with retry logic for rate-limited requests.
// BUG-DC-004: Add rate limit handling for refresh command
func (r *deviceCommandResource) executeCommandWithRetry(ctx context.Context, req *http.Request, expectedStatus int, command string) ([]byte, error) {
	maxRetries := 3
	retryDelay := 10 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		body, err := r.executeCommand(req, expectedStatus)
		if err == nil {
			return body, nil
		}

		// Check if this is a rate limit error (429) and we should retry
		// This is particularly important for the refresh command
		if command == "refresh" && strings.Contains(err.Error(), "429") && attempt < maxRetries-1 {
			time.Sleep(retryDelay)
			retryDelay *= 2 // Exponential backoff
			continue
		}

		return nil, err
	}

	return nil, fmt.Errorf("max retries exceeded for command %q", command)
}

func (r *deviceCommandResource) updateCreateState(ctx context.Context, plan *deviceCommandResourceModel, commandKey string, expectedStatus int, body []byte, resp *resource.CreateResponse) {
	plan.StatusCode = types.Int64Value(int64(expectedStatus))
	if len(body) > 0 {
		plan.Response = types.StringValue(string(body))
	} else {
		plan.Response = types.StringNull()
	}

	plan.ID = types.StringValue(fmt.Sprintf(deviceCommandIDFormatTemplate, plan.DeviceID.ValueString(), commandKey, time.Now().UTC().Unix()))

	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}
