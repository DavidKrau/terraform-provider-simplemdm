package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type customDeclarationDeviceAssignmentResource struct {
	client *simplemdm.Client
}

type customDeclarationDeviceAssignmentModel struct {
	ID                  types.String `tfsdk:"id"`
	CustomDeclarationID types.String `tfsdk:"custom_declaration_id"`
	DeviceID            types.String `tfsdk:"device_id"`
}

var (
	_ resource.Resource                = &customDeclarationDeviceAssignmentResource{}
	_ resource.ResourceWithConfigure   = &customDeclarationDeviceAssignmentResource{}
	_ resource.ResourceWithImportState = &customDeclarationDeviceAssignmentResource{}
)

func CustomDeclarationDeviceAssignmentResource() resource.Resource {
	return &customDeclarationDeviceAssignmentResource{}
}

func (r *customDeclarationDeviceAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customdeclaration_device_assignment"
}

func (r *customDeclarationDeviceAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the assignment of a custom declaration to a SimpleMDM device.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_declaration_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier of the custom declaration to assign.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"device_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier of the device that should receive the custom declaration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *customDeclarationDeviceAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*simplemdm.Client)
}

func (r *customDeclarationDeviceAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan customDeclarationDeviceAssignmentModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("https://%s/api/v1/custom_declarations/%s/devices/%s", r.client.HostName, plan.CustomDeclarationID.ValueString(), plan.DeviceID.ValueString())
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error creating SimpleMDM custom declaration assignment request", err.Error())
		return
	}

	if _, err := r.client.RequestResponse204or409(httpReq); err != nil {
		resp.Diagnostics.AddError("Error assigning custom declaration to device", err.Error())
		return
	}

	plan.ID = types.StringValue(buildCustomDeclarationAssignmentID(plan.CustomDeclarationID.ValueString(), plan.DeviceID.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *customDeclarationDeviceAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Cannot update custom declaration assignments",
		"Updates are not supported. Remove and recreate the assignment to target a different device or declaration.",
	)
}

func (r *customDeclarationDeviceAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state customDeclarationDeviceAssignmentModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("https://%s/api/v1/devices/%s", r.client.HostName, state.DeviceID.ValueString())
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error creating SimpleMDM device request", err.Error())
		return
	}

	body, err := r.client.RequestResponse200(httpReq)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading SimpleMDM device assignments", err.Error())
		return
	}

	assigned, err := deviceHasCustomDeclarationAssignment(body, state.CustomDeclarationID.ValueString(), state.DeviceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error parsing SimpleMDM device relationships", err.Error())
		return
	}

	if !assigned {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(buildCustomDeclarationAssignmentID(state.CustomDeclarationID.ValueString(), state.DeviceID.ValueString()))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *customDeclarationDeviceAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state customDeclarationDeviceAssignmentModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("https://%s/api/v1/custom_declarations/%s/devices/%s", r.client.HostName, state.CustomDeclarationID.ValueString(), state.DeviceID.ValueString())
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error creating SimpleMDM custom declaration assignment request", err.Error())
		return
	}

	if _, err := r.client.RequestResponse204or409(httpReq); err != nil {
		if strings.Contains(err.Error(), "404") {
			return
		}

		resp.Diagnostics.AddError("Error removing custom declaration assignment", err.Error())
	}
}

func (r *customDeclarationDeviceAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Support both : and | as separators for backward compatibility
	var parts []string
	var declarationID, deviceID string

	if strings.Contains(req.ID, "|") {
		parts = strings.Split(req.ID, "|")
	} else {
		parts = strings.Split(req.ID, ":")
	}

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected import identifier format",
			"Expected custom_declaration_id:device_id or custom_declaration_id|device_id",
		)
		return
	}

	declarationID = parts[0]
	deviceID = parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("custom_declaration_id"), declarationID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("device_id"), deviceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...) //nolint:errcheck
}

func deviceHasCustomDeclarationAssignment(body []byte, customDeclarationID string, deviceID string) (bool, error) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return false, fmt.Errorf("error parsing device payload for device %s: %w", deviceID, err)
	}

	data, ok := payload["data"].(map[string]any)
	if !ok {
		return false, fmt.Errorf("unexpected device payload structure for device %s: missing data node", deviceID)
	}

	relationships, ok := data["relationships"].(map[string]any)
	if !ok {
		return false, nil
	}

	rel, ok := relationships["custom_declarations"].(map[string]any)
	if !ok {
		return false, nil
	}

	assignments, ok := rel["data"].([]any)
	if !ok {
		return false, nil
	}

	for _, entry := range assignments {
		relEntry, ok := entry.(map[string]any)
		if !ok {
			continue
		}

		idValue, ok := relEntry["id"]
		if !ok {
			continue
		}

		if fmt.Sprint(idValue) == customDeclarationID {
			return true, nil
		}
	}

	return false, nil
}

func buildCustomDeclarationAssignmentID(customDeclarationID, deviceID string) string {
	// Use | separator to avoid conflicts with IDs that contain colons
	// This addresses BUG-CD-012
	if strings.Contains(customDeclarationID, ":") || strings.Contains(deviceID, ":") {
		return fmt.Sprintf("%s|%s", customDeclarationID, deviceID)
	}
	return fmt.Sprintf("%s:%s", customDeclarationID, deviceID)
}
