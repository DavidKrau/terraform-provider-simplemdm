package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type customDeclarationDataSource struct {
	client *simplemdm.Client
}

type customDeclarationDataSourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Identifier          types.String `tfsdk:"identifier"`
	DeclarationType     types.String `tfsdk:"declaration_type"`
	Topic               types.String `tfsdk:"topic"`
	Transport           types.String `tfsdk:"transport"`
	Description         types.String `tfsdk:"description"`
	Platforms           types.Set    `tfsdk:"platforms"`
	Data                types.String `tfsdk:"data"`
	Active              types.Bool   `tfsdk:"active"`
	Priority            types.Int64  `tfsdk:"priority"`
	Payload             types.String `tfsdk:"payload"`
	UserScope           types.Bool   `tfsdk:"user_scope"`
	AttributeSupport    types.Bool   `tfsdk:"attribute_support"`
	EscapeAttributes    types.Bool   `tfsdk:"escape_attributes"`
	ActivationPredicate types.String `tfsdk:"activation_predicate"`
	ProfileIdentifier   types.String `tfsdk:"profile_identifier"`
	GroupCount          types.Int64  `tfsdk:"group_count"`
	DeviceCount         types.Int64  `tfsdk:"device_count"`
	CreatedAt           types.String `tfsdk:"created_at"`
	UpdatedAt           types.String `tfsdk:"updated_at"`
}

var _ datasource.DataSource = &customDeclarationDataSource{}
var _ datasource.DataSourceWithConfigure = &customDeclarationDataSource{}

func CustomDeclarationDataSource() datasource.DataSource {
	return &customDeclarationDataSource{}
}

func (d *customDeclarationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customdeclaration"
}

func (d *customDeclarationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Custom Declaration data source retrieves Declarative Device Management custom declarations from SimpleMDM.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the declaration to retrieve.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Human readable name for the declaration.",
			},
			"identifier": schema.StringAttribute{
				Computed:    true,
				Description: "Unique declaration identifier.",
			},
			"declaration_type": schema.StringAttribute{
				Computed:    true,
				Description: "Declaration type reported to Apple devices.",
			},
			"topic": schema.StringAttribute{
				Computed:    true,
				Description: "Topic used for declarative management payloads.",
			},
			"transport": schema.StringAttribute{
				Computed:    true,
				Description: "Transport mechanism for the declaration.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Description of the declaration.",
			},
			"platforms": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of platforms that receive the declaration.",
			},
			"data": schema.StringAttribute{
				Computed:    true,
				Description: "JSON payload of the declaration data.",
			},
			"payload": schema.StringAttribute{
				Computed:    true,
				Description: "Alias that mirrors the JSON payload returned by the SimpleMDM download endpoint.",
			},
			"active": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the declaration is active.",
			},
			"priority": schema.Int64Attribute{
				Computed:    true,
				Description: "Priority value used for ordering declarations.",
			},
			"user_scope": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the declaration is scoped to users (true) or devices (false).",
			},
			"attribute_support": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether variable expansion is enabled for the declaration payload.",
			},
			"escape_attributes": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether custom variable values are escaped before being delivered.",
			},
			"activation_predicate": schema.StringAttribute{
				Computed:    true,
				Description: "Predicate that controls when the declaration activates on a device.",
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

func (d *customDeclarationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*simplemdm.Client)
}

func (d *customDeclarationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state customDeclarationDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("https://%s/api/v1/custom_declarations/%s", d.client.HostName, state.ID.ValueString())
	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error creating SimpleMDM custom declaration request", err.Error())
		return
	}

	responseBody, err := d.client.RequestResponse200(httpReq)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.Diagnostics.AddError("Custom declaration not found", err.Error())
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

	if len(declaration.Data.Attributes.Data) == 0 && len(declaration.Data.Attributes.Payload) == 0 {
		raw, err := downloadCustomDeclarationPayload(ctx, d.client, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error downloading SimpleMDM custom declaration payload", err.Error())
			return
		}

		declaration.Data.Attributes.Data = raw
	}

	var model customDeclarationResourceModel
	if diags := model.refreshFromResponse(ctx, &declaration); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Copy data from resource model into data source state.
	state.ID = model.ID
	state.Name = model.Name
	state.Identifier = model.Identifier
	state.DeclarationType = model.DeclarationType
	state.Topic = model.Topic
	state.Transport = model.Transport
	state.Description = model.Description
	state.Platforms = model.Platforms
	state.Data = model.Data
	state.Active = model.Active
	state.Priority = model.Priority
	state.Payload = model.Payload
	state.UserScope = model.UserScope
	state.AttributeSupport = model.AttributeSupport
	state.EscapeAttributes = model.EscapeAttributes
	state.ActivationPredicate = model.ActivationPredicate
	state.ProfileIdentifier = model.ProfileIdentifier
	state.GroupCount = model.GroupCount
	state.DeviceCount = model.DeviceCount
	state.CreatedAt = model.CreatedAt
	state.UpdatedAt = model.UpdatedAt

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
