package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type customDeclarationResource struct {
	client *simplemdm.Client
}

type customDeclarationResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Identifier      types.String `tfsdk:"identifier"`
	DeclarationType types.String `tfsdk:"declaration_type"`
	Topic           types.String `tfsdk:"topic"`
	Transport       types.String `tfsdk:"transport"`
	Description     types.String `tfsdk:"description"`
	Platforms       types.Set    `tfsdk:"platforms"`
	Data            types.String `tfsdk:"data"`
	Active          types.Bool   `tfsdk:"active"`
	Priority        types.Int64  `tfsdk:"priority"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
}

type customDeclarationAttributes struct {
	Name            string          `json:"name"`
	Identifier      string          `json:"identifier"`
	DeclarationType string          `json:"declaration_type"`
	Topic           string          `json:"topic"`
	Transport       string          `json:"transport"`
	Description     string          `json:"description"`
	Platforms       []string        `json:"platforms"`
	Data            json.RawMessage `json:"data"`
	Active          *bool           `json:"active"`
	Priority        *int64          `json:"priority"`
	CreatedAt       string          `json:"created_at"`
	UpdatedAt       string          `json:"updated_at"`
}

type customDeclarationResponse struct {
	Data struct {
		ID         string                      `json:"id"`
		Attributes customDeclarationAttributes `json:"attributes"`
	} `json:"data"`
}

type customDeclarationRequest struct {
	Data struct {
		Type       string                   `json:"type"`
		Attributes customDeclarationPayload `json:"attributes"`
	} `json:"data"`
}

type customDeclarationPayload struct {
	Name            string          `json:"name"`
	Identifier      string          `json:"identifier"`
	DeclarationType string          `json:"declaration_type"`
	Topic           *string         `json:"topic,omitempty"`
	Transport       *string         `json:"transport,omitempty"`
	Description     *string         `json:"description,omitempty"`
	Platforms       []string        `json:"platforms"`
	Data            json.RawMessage `json:"data"`
	Active          *bool           `json:"active,omitempty"`
	Priority        *int64          `json:"priority,omitempty"`
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
				Description: "Human readable name for the declaration.",
			},
			"identifier": schema.StringAttribute{
				Required:    true,
				Description: "Unique declaration identifier. Changing forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"declaration_type": schema.StringAttribute{
				Required:    true,
				Description: "Declaration type reported to Apple devices.",
			},
			"topic": schema.StringAttribute{
				Optional:    true,
				Description: "Optional topic used for declarative management payloads.",
			},
			"transport": schema.StringAttribute{
				Optional:    true,
				Description: "Optional transport mechanism for the declaration.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Optional description of the declaration.",
			},
			"platforms": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of platforms that should receive the declaration.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"data": schema.StringAttribute{
				Required:    true,
				Description: "JSON payload of the declaration data.",
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the declaration is active.",
				Default:     booldefault.StaticBool(true),
			},
			"priority": schema.Int64Attribute{
				Optional:    true,
				Description: "Optional priority value used for ordering declarations.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the declaration was created in SimpleMDM.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the declaration was last updated in SimpleMDM.",
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

	requestBody := customDeclarationRequest{}
	requestBody.Data.Type = "custom_declaration"
	requestBody.Data.Attributes = *payload

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		resp.Diagnostics.AddError("Error marshalling SimpleMDM custom declaration payload", err.Error())
		return
	}

	url := fmt.Sprintf("https://%s/api/v1/custom_declarations", r.client.HostName)
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		resp.Diagnostics.AddError("Error creating SimpleMDM custom declaration request", err.Error())
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")

	responseBody, err := r.client.RequestResponse201(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating SimpleMDM custom declaration", err.Error())
		return
	}

	var declaration customDeclarationResponse
	if err := json.Unmarshal(responseBody, &declaration); err != nil {
		resp.Diagnostics.AddError("Error parsing SimpleMDM custom declaration response", err.Error())
		return
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

	requestBody := customDeclarationRequest{}
	requestBody.Data.Type = "custom_declaration"
	requestBody.Data.Attributes = *payload

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		resp.Diagnostics.AddError("Error marshalling SimpleMDM custom declaration payload", err.Error())
		return
	}

	url := fmt.Sprintf("https://%s/api/v1/custom_declarations/%s", r.client.HostName, plan.ID.ValueString())
	httpReq, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(bodyBytes))
	if err != nil {
		resp.Diagnostics.AddError("Error creating SimpleMDM custom declaration request", err.Error())
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")

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

	platforms := make([]string, 0)
	if !model.Platforms.IsNull() && !model.Platforms.IsUnknown() {
		var platformValues []string
		d := model.Platforms.ElementsAs(ctx, &platformValues, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		sort.Strings(platformValues)
		platforms = append(platforms, platformValues...)
	}

	normalizedData, err := normalizeJSON(model.Data.ValueString())
	if err != nil {
		diags.AddError("Invalid JSON data", fmt.Sprintf("Unable to parse declaration data: %s", err))
		return nil, diags
	}

	payload := &customDeclarationPayload{
		Name:            model.Name.ValueString(),
		Identifier:      model.Identifier.ValueString(),
		DeclarationType: model.DeclarationType.ValueString(),
		Platforms:       platforms,
		Data:            json.RawMessage(normalizedData),
	}

	if !model.Topic.IsNull() {
		topic := model.Topic.ValueString()
		payload.Topic = &topic
	}

	if !model.Transport.IsNull() {
		transport := model.Transport.ValueString()
		payload.Transport = &transport
	}

	if !model.Description.IsNull() {
		description := model.Description.ValueString()
		payload.Description = &description
	}

	if !model.Active.IsNull() {
		active := model.Active.ValueBool()
		payload.Active = &active
	}

	if !model.Priority.IsNull() {
		priority := model.Priority.ValueInt64()
		payload.Priority = &priority
	}

	return payload, diags
}

func (m *customDeclarationResourceModel) refreshFromResponse(ctx context.Context, response *customDeclarationResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	attributes := response.Data.Attributes

	m.ID = types.StringValue(response.Data.ID)
	m.Name = types.StringValue(attributes.Name)
	m.Identifier = types.StringValue(attributes.Identifier)
	m.DeclarationType = types.StringValue(attributes.DeclarationType)

	if attributes.Topic != "" {
		m.Topic = types.StringValue(attributes.Topic)
	} else {
		m.Topic = types.StringNull()
	}

	if attributes.Transport != "" {
		m.Transport = types.StringValue(attributes.Transport)
	} else {
		m.Transport = types.StringNull()
	}

	if attributes.Description != "" {
		m.Description = types.StringValue(attributes.Description)
	} else {
		m.Description = types.StringNull()
	}

	if attributes.Active != nil {
		m.Active = types.BoolValue(*attributes.Active)
	} else {
		m.Active = types.BoolNull()
	}

	if attributes.Priority != nil {
		m.Priority = types.Int64Value(*attributes.Priority)
	} else {
		m.Priority = types.Int64Null()
	}

	if len(attributes.Platforms) > 0 {
		sort.Strings(attributes.Platforms)
		platforms, d := types.SetValueFrom(ctx, types.StringType, attributes.Platforms)
		diags.Append(d...)
		m.Platforms = platforms
	} else {
		m.Platforms = types.SetNull(types.StringType)
	}

	if len(attributes.Data) > 0 {
		normalized, err := normalizeJSON(string(attributes.Data))
		if err != nil {
			diags.AddError("Invalid JSON data", fmt.Sprintf("Unable to normalize declaration data: %s", err))
			return diags
		}

		m.Data = types.StringValue(normalized)
	} else {
		m.Data = types.StringNull()
	}

	if attributes.CreatedAt != "" {
		m.CreatedAt = types.StringValue(attributes.CreatedAt)
	} else {
		m.CreatedAt = types.StringNull()
	}

	if attributes.UpdatedAt != "" {
		m.UpdatedAt = types.StringValue(attributes.UpdatedAt)
	} else {
		m.UpdatedAt = types.StringNull()
	}

	return diags
}

func normalizeJSON(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", fmt.Errorf("expected JSON object or array")
	}

	decoder := json.NewDecoder(strings.NewReader(trimmed))
	decoder.UseNumber()

	var value any
	if err := decoder.Decode(&value); err != nil {
		return "", err
	}

	normalized, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(normalized), nil
}
