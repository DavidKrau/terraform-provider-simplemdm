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
				Description: "A name for the custom declaration.",
			},
			"declaration_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of declaration being defined.",
			},
			"payload": schema.StringAttribute{
				Computed:    true,
				Description: "The JSON payload for the declaration.",
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
			"reinstall_after_os_update": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether to reinstall the declaration after macOS updates.",
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

	if len(declaration.Data.Attributes.Payload) == 0 {
		raw, err := downloadCustomDeclarationPayload(ctx, d.client, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error downloading SimpleMDM custom declaration payload", err.Error())
			return
		}

		declaration.Data.Attributes.Payload = raw
	}

	var model customDeclarationResourceModel
	if diags := model.refreshFromResponse(ctx, &declaration); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Copy data from resource model into data source state.
	state.ID = model.ID
	state.Name = model.Name
	state.DeclarationType = model.DeclarationType
	state.Payload = model.Payload
	state.UserScope = model.UserScope
	state.AttributeSupport = model.AttributeSupport
	state.EscapeAttributes = model.EscapeAttributes
	state.ActivationPredicate = model.ActivationPredicate
	state.ReinstallAfterOsUpdate = model.ReinstallAfterOsUpdate
	state.ProfileIdentifier = model.ProfileIdentifier
	state.GroupCount = model.GroupCount
	state.DeviceCount = model.DeviceCount

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
