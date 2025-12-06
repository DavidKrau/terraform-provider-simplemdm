package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type customDeclarationResource struct {
	client *simplemdm.Client
}

type customDeclarationResourceModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	DeclarationType        types.String `tfsdk:"declaration_type"`
	Payload                types.String `tfsdk:"payload"`
	UserScope              types.Bool   `tfsdk:"user_scope"`
	AttributeSupport       types.Bool   `tfsdk:"attribute_support"`
	EscapeAttributes       types.Bool   `tfsdk:"escape_attributes"`
	ActivationPredicate    types.String `tfsdk:"activation_predicate"`
	ReinstallAfterOsUpdate types.Bool   `tfsdk:"reinstall_after_os_update"`
	ProfileIdentifier      types.String `tfsdk:"profile_identifier"`
	GroupCount             types.Int64  `tfsdk:"group_count"`
	DeviceCount            types.Int64  `tfsdk:"device_count"`
}

type customDeclarationAttributes struct {
	Name                   string          `json:"name"`
	DeclarationType        string          `json:"declaration_type"`
	Payload                json.RawMessage `json:"payload"`
	UserScope              *bool           `json:"user_scope"`
	AttributeSupport       *bool           `json:"attribute_support"`
	EscapeAttributes       *bool           `json:"escape_attributes"`
	ActivationPredicate    string          `json:"activation_predicate"`
	ReinstallAfterOsUpdate *bool           `json:"reinstall_after_os_update"`
	ProfileIdentifier      string          `json:"profile_identifier"`
	GroupCount             *int64          `json:"group_count"`
	DeviceCount            *int64          `json:"device_count"`
}

type customDeclarationResponse struct {
	Data struct {
		ID         string                      `json:"id"`
		Attributes customDeclarationAttributes `json:"attributes"`
	} `json:"data"`
}

// customDeclarationPayload represents the multipart form data for Create/Update
type customDeclarationPayload struct {
	Name                   string `json:"name"`
	DeclarationType        string `json:"declaration_type"`
	Payload                []byte `json:"payload"`
	UserScope              *bool  `json:"user_scope,omitempty"`
	AttributeSupport       *bool  `json:"attribute_support,omitempty"`
	EscapeAttributes       *bool  `json:"escape_attributes,omitempty"`
	ActivationPredicate    string `json:"activation_predicate,omitempty"`
	ReinstallAfterOsUpdate *bool  `json:"reinstall_after_os_update,omitempty"`
}

func CustomDeclarationResource() resource.Resource {
	return &customDeclarationResource{}
}

var _ resource.Resource = &customDeclarationResource{}
var _ resource.ResourceWithConfigure = &customDeclarationResource{}
var _ resource.ResourceWithImportState = &customDeclarationResource{}

func (r *customDeclarationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customdeclaration"
}

func (r *customDeclarationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Custom Declaration resource manages Declarative Device Management custom declarations in SimpleMDM.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "A name for the custom declaration.",
			},
			"declaration_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of declaration being defined (e.g., com.apple.configuration.management.status-subscriptions).",
			},
			"payload": schema.StringAttribute{
				Required:    true,
				Description: "The JSON payload for the declaration.",
			},
			"user_scope": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the declaration is scoped to users (true) or devices (false). Defaults to true.",
				Default:     booldefault.StaticBool(true),
			},
			"attribute_support": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable variable expansion when processing the declaration payload. Defaults to false.",
				Default:     booldefault.StaticBool(false),
			},
			"escape_attributes": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Escape the values of custom variables within the payload before delivery. Defaults to false.",
				Default:     booldefault.StaticBool(false),
			},
			"activation_predicate": schema.StringAttribute{
				Optional:    true,
				Description: "Predicate format string that controls when the declaration activates on a device.",
			},
			"reinstall_after_os_update": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to reinstall the declaration after macOS updates. Defaults to false.",
				Default:     booldefault.StaticBool(false),
			},
			"profile_identifier": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier assigned by SimpleMDM for tracking the declaration profile.",
			},
			"group_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of device groups currently assigned to the declaration.",
			},
			"device_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of devices currently assigned to the declaration.",
			},
		},
	}
}

func (r *customDeclarationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*simplemdm.Client)
}

func (r *customDeclarationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan customDeclarationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diags := buildCustomDeclarationPayload(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build multipart form data
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add required fields
	if err := writer.WriteField("name", payload.Name); err != nil {
		resp.Diagnostics.AddError("Error building multipart request", err.Error())
		return
	}

	if err := writer.WriteField("declaration_type", payload.DeclarationType); err != nil {
		resp.Diagnostics.AddError("Error building multipart request", err.Error())
		return
	}

	// Add payload as a file part
	if len(payload.Payload) > 0 {
		part, err := writer.CreateFormFile("payload", "declaration.json")
		if err != nil {
			resp.Diagnostics.AddError("Error building multipart request", err.Error())
			return
		}
		if _, err := part.Write(payload.Payload); err != nil {
			resp.Diagnostics.AddError("Error building multipart request", err.Error())
			return
		}
	}

	// Add optional fields
	if payload.UserScope != nil {
		if err := writer.WriteField("user_scope", fmt.Sprintf("%t", *payload.UserScope)); err != nil {
			resp.Diagnostics.AddError("Error building multipart request", err.Error())
			return
		}
	}

	if payload.AttributeSupport != nil {
		if err := writer.WriteField("attribute_support", fmt.Sprintf("%t", *payload.AttributeSupport)); err != nil {
			resp.Diagnostics.AddError("Error building multipart request", err.Error())
			return
		}
	}

	if payload.EscapeAttributes != nil {
		if err := writer.WriteField("escape_attributes", fmt.Sprintf("%t", *payload.EscapeAttributes)); err != nil {
			resp.Diagnostics.AddError("Error building multipart request", err.Error())
			return
		}
	}

	if payload.ActivationPredicate != "" {
		if err := writer.WriteField("activation_predicate", payload.ActivationPredicate); err != nil {
			resp.Diagnostics.AddError("Error building multipart request", err.Error())
			return
		}
	}

	if payload.ReinstallAfterOsUpdate != nil {
		if err := writer.WriteField("reinstall_after_os_update", fmt.Sprintf("%t", *payload.ReinstallAfterOsUpdate)); err != nil {
			resp.Diagnostics.AddError("Error building multipart request", err.Error())
			return
		}
	}

	if err := writer.Close(); err != nil {
		resp.Diagnostics.AddError("Error building multipart request", err.Error())
		return
	}

	url := fmt.Sprintf("https://%s/api/v1/custom_declarations", r.client.HostName)
	httpReq, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating SimpleMDM custom declaration request", err.Error())
		return
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	// API may return 200 or 201 on success
	responseBody, err := r.client.RequestResponse200(httpReq)
	if err != nil {
		// Try 201 if 200 failed
		if strings.Contains(err.Error(), "200") {
			httpReq2, _ := http.NewRequest(http.MethodPost, url, &body)
			httpReq2.Header.Set("Content-Type", writer.FormDataContentType())
			responseBody, err = r.client.RequestResponse201(httpReq2)
		}
		if err != nil {
			resp.Diagnostics.AddError("Error creating SimpleMDM custom declaration", err.Error())
			return
		}
	}

	var declaration customDeclarationResponse
	if err := json.Unmarshal(responseBody, &declaration); err != nil {
		resp.Diagnostics.AddError("Error parsing SimpleMDM custom declaration response", err.Error())
		return
	}

	if len(declaration.Data.Attributes.Payload) == 0 {
		raw, err := downloadCustomDeclarationPayload(ctx, r.client, declaration.Data.ID)
		if err != nil {
			resp.Diagnostics.AddError("Error downloading SimpleMDM custom declaration payload", err.Error())
			return
		}

		declaration.Data.Attributes.Payload = raw
	}

	if diags := plan.refreshFromResponse(ctx, &declaration); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *customDeclarationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state customDeclarationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("https://%s/api/v1/custom_declarations/%s", r.client.HostName, state.ID.ValueString())
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error creating SimpleMDM custom declaration request", err.Error())
		return
	}

	responseBody, err := r.client.RequestResponse200(httpReq)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading SimpleMDM custom declaration", err.Error())
		return
	}

	var declaration customDeclarationResponse
	if err := json.Unmarshal(responseBody, &declaration); err != nil {
		resp.Diagnostics.AddError("Error parsing SimpleMDM custom declaration response", err.Error())
		return
	}

	if len(declaration.Data.Attributes.Payload) == 0 {
		raw, err := downloadCustomDeclarationPayload(ctx, r.client, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error downloading SimpleMDM custom declaration payload", err.Error())
			return
		}

		declaration.Data.Attributes.Payload = raw
	}

	if diags := state.refreshFromResponse(ctx, &declaration); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *customDeclarationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan customDeclarationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diags := buildCustomDeclarationPayload(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build multipart form data
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add fields to update
	if err := writer.WriteField("name", payload.Name); err != nil {
		resp.Diagnostics.AddError("Error building multipart request", err.Error())
		return
	}

	if err := writer.WriteField("declaration_type", payload.DeclarationType); err != nil {
		resp.Diagnostics.AddError("Error building multipart request", err.Error())
		return
	}

	// Add payload as a file part if present
	if len(payload.Payload) > 0 {
		part, err := writer.CreateFormFile("payload", "declaration.json")
		if err != nil {
			resp.Diagnostics.AddError("Error building multipart request", err.Error())
			return
		}
		if _, err := part.Write(payload.Payload); err != nil {
			resp.Diagnostics.AddError("Error building multipart request", err.Error())
			return
		}
	}

	// Add optional fields
	if payload.UserScope != nil {
		if err := writer.WriteField("user_scope", fmt.Sprintf("%t", *payload.UserScope)); err != nil {
			resp.Diagnostics.AddError("Error building multipart request", err.Error())
			return
		}
	}

	if payload.AttributeSupport != nil {
		if err := writer.WriteField("attribute_support", fmt.Sprintf("%t", *payload.AttributeSupport)); err != nil {
			resp.Diagnostics.AddError("Error building multipart request", err.Error())
			return
		}
	}

	if payload.EscapeAttributes != nil {
		if err := writer.WriteField("escape_attributes", fmt.Sprintf("%t", *payload.EscapeAttributes)); err != nil {
			resp.Diagnostics.AddError("Error building multipart request", err.Error())
			return
		}
	}

	if payload.ActivationPredicate != "" {
		if err := writer.WriteField("activation_predicate", payload.ActivationPredicate); err != nil {
			resp.Diagnostics.AddError("Error building multipart request", err.Error())
			return
		}
	}

	if payload.ReinstallAfterOsUpdate != nil {
		if err := writer.WriteField("reinstall_after_os_update", fmt.Sprintf("%t", *payload.ReinstallAfterOsUpdate)); err != nil {
			resp.Diagnostics.AddError("Error building multipart request", err.Error())
			return
		}
	}

	if err := writer.Close(); err != nil {
		resp.Diagnostics.AddError("Error building multipart request", err.Error())
		return
	}

	url := fmt.Sprintf("https://%s/api/v1/custom_declarations/%s", r.client.HostName, plan.ID.ValueString())
	httpReq, err := http.NewRequest(http.MethodPatch, url, &body)
	if err != nil {
		resp.Diagnostics.AddError("Error creating SimpleMDM custom declaration request", err.Error())
		return
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	responseBody, err := r.client.RequestResponse200(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating SimpleMDM custom declaration", err.Error())
		return
	}

	var declaration customDeclarationResponse
	if err := json.Unmarshal(responseBody, &declaration); err != nil {
		resp.Diagnostics.AddError("Error parsing SimpleMDM custom declaration response", err.Error())
		return
	}

	if len(declaration.Data.Attributes.Payload) == 0 {
		raw, err := downloadCustomDeclarationPayload(ctx, r.client, plan.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error downloading SimpleMDM custom declaration payload", err.Error())
			return
		}

		declaration.Data.Attributes.Payload = raw
	}

	if diags := plan.refreshFromResponse(ctx, &declaration); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *customDeclarationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state customDeclarationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("https://%s/api/v1/custom_declarations/%s", r.client.HostName, state.ID.ValueString())
	httpReq, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error creating SimpleMDM custom declaration request", err.Error())
		return
	}

	_, err = r.client.RequestResponse204(httpReq)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return
		}

		resp.Diagnostics.AddError("Error deleting SimpleMDM custom declaration", err.Error())
	}
}

func (r *customDeclarationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func buildCustomDeclarationPayload(ctx context.Context, model *customDeclarationResourceModel) (*customDeclarationPayload, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Validate required fields (BUG-CD-007)
	if model.Name.IsNull() || model.Name.ValueString() == "" {
		diags.AddError("Missing required field", "name is required")
	}
	if model.DeclarationType.IsNull() || model.DeclarationType.ValueString() == "" {
		diags.AddError("Missing required field", "declaration_type is required")
	}
	if model.Payload.IsNull() || model.Payload.ValueString() == "" {
		diags.AddError("Missing required field", "payload is required")
	}

	if diags.HasError() {
		return nil, diags
	}

	normalizedPayload, err := normalizeJSON(model.Payload.ValueString(), "payload", model.ID.ValueString())
	if err != nil {
		diags.AddError("Invalid JSON payload", fmt.Sprintf("Unable to parse declaration payload: %s", err))
		return nil, diags
	}

	payload := &customDeclarationPayload{
		Name:            model.Name.ValueString(),
		DeclarationType: model.DeclarationType.ValueString(),
		Payload:         []byte(normalizedPayload),
	}

	if !model.UserScope.IsNull() {
		userScope := model.UserScope.ValueBool()
		payload.UserScope = &userScope
	}

	if !model.AttributeSupport.IsNull() {
		attributeSupport := model.AttributeSupport.ValueBool()
		payload.AttributeSupport = &attributeSupport
	}

	if !model.EscapeAttributes.IsNull() {
		escapeAttributes := model.EscapeAttributes.ValueBool()
		payload.EscapeAttributes = &escapeAttributes
	}

	if !model.ActivationPredicate.IsNull() && model.ActivationPredicate.ValueString() != "" {
		payload.ActivationPredicate = model.ActivationPredicate.ValueString()
	}

	if !model.ReinstallAfterOsUpdate.IsNull() {
		reinstall := model.ReinstallAfterOsUpdate.ValueBool()
		payload.ReinstallAfterOsUpdate = &reinstall
	}

	return payload, diags
}

func (m *customDeclarationResourceModel) refreshFromResponse(ctx context.Context, response *customDeclarationResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	attributes := response.Data.Attributes

	m.ID = types.StringValue(response.Data.ID)
	m.Name = types.StringValue(attributes.Name)
	m.DeclarationType = types.StringValue(attributes.DeclarationType)

	if attributes.UserScope != nil {
		m.UserScope = types.BoolValue(*attributes.UserScope)
	} else {
		m.UserScope = types.BoolNull()
	}

	if attributes.AttributeSupport != nil {
		m.AttributeSupport = types.BoolValue(*attributes.AttributeSupport)
	} else {
		m.AttributeSupport = types.BoolNull()
	}

	if attributes.EscapeAttributes != nil {
		m.EscapeAttributes = types.BoolValue(*attributes.EscapeAttributes)
	} else {
		m.EscapeAttributes = types.BoolNull()
	}

	if attributes.ActivationPredicate != "" {
		m.ActivationPredicate = types.StringValue(attributes.ActivationPredicate)
	} else {
		m.ActivationPredicate = types.StringNull()
	}

	if attributes.ReinstallAfterOsUpdate != nil {
		m.ReinstallAfterOsUpdate = types.BoolValue(*attributes.ReinstallAfterOsUpdate)
	} else {
		m.ReinstallAfterOsUpdate = types.BoolNull()
	}

	if attributes.ProfileIdentifier != "" {
		m.ProfileIdentifier = types.StringValue(attributes.ProfileIdentifier)
	} else {
		m.ProfileIdentifier = types.StringNull()
	}

	if attributes.GroupCount != nil {
		m.GroupCount = types.Int64Value(*attributes.GroupCount)
	} else {
		m.GroupCount = types.Int64Null()
	}

	if attributes.DeviceCount != nil {
		m.DeviceCount = types.Int64Value(*attributes.DeviceCount)
	} else {
		m.DeviceCount = types.Int64Null()
	}

	if len(attributes.Payload) > 0 {
		normalized, err := normalizeJSON(string(attributes.Payload), "payload", m.ID.ValueString())
		if err != nil {
			diags.AddError("Invalid JSON payload", fmt.Sprintf("Unable to normalize declaration payload: %s", err))
			return diags
		}

		m.Payload = types.StringValue(normalized)
	} else {
		m.Payload = types.StringNull()
	}

	return diags
}

func downloadCustomDeclarationPayload(ctx context.Context, client *simplemdm.Client, declarationID string) (json.RawMessage, error) {
	url := fmt.Sprintf("https://%s/api/v1/custom_declarations/%s/download", client.HostName, declarationID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := client.RequestResponse200(httpReq)
	if err != nil {
		return nil, err
	}

	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return nil, nil
	}

	return json.RawMessage(trimmed), nil
}

func normalizeJSON(input string, fieldName string, declarationID string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		contextInfo := ""
		if declarationID != "" {
			contextInfo = fmt.Sprintf(" in declaration %s", declarationID)
		}
		return "", fmt.Errorf("field '%s'%s: expected JSON object or array", fieldName, contextInfo)
	}

	decoder := json.NewDecoder(strings.NewReader(trimmed))
	decoder.UseNumber()

	var value any
	if err := decoder.Decode(&value); err != nil {
		contextInfo := ""
		if declarationID != "" {
			contextInfo = fmt.Sprintf(" in declaration %s", declarationID)
		}
		return "", fmt.Errorf("field '%s'%s: %w", fieldName, contextInfo, err)
	}

	normalized, err := json.Marshal(value)
	if err != nil {
		contextInfo := ""
		if declarationID != "" {
			contextInfo = fmt.Sprintf(" in declaration %s", declarationID)
		}
		return "", fmt.Errorf("field '%s'%s: %w", fieldName, contextInfo, err)
	}

	return string(normalized), nil
}
