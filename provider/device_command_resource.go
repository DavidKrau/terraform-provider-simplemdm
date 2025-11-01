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
}

func DeviceCommandResource() resource.Resource {
	return &deviceCommandResource{}
}

func (r *deviceCommandResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_command"
}

func (r *deviceCommandResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Executes management commands against a SimpleMDM device.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Internal identifier for the executed command.",
			},
			"device_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier of the target device.",
			},
			"command": schema.StringAttribute{
				Required:    true,
				Description: "Command to execute. Supported values include push_assigned_apps, refresh, restart, shutdown, lock, clear_passcode, clear_firmware_password, rotate_firmware_password, clear_recovery_lock_password, clear_restrictions_password, rotate_recovery_lock_password, rotate_filevault_recovery_key, set_admin_password, rotate_admin_password, wipe, update_os, enable_remote_desktop, disable_remote_desktop, enable_bluetooth, disable_bluetooth, set_time_zone, unenroll, delete_user.",
			},
			"parameters": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Optional parameters to pass to the API call.",
			},
			"status_code": schema.Int64Attribute{
				Computed:    true,
				Description: "HTTP status code returned by the SimpleMDM API.",
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

	body, err := r.executeCommand(reqObj, spec.expectedStatus)
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

func (r *deviceCommandResource) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *deviceCommandResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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
