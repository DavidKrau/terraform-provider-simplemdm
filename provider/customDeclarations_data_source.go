package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &customDeclarationsDataSource{}
	_ datasource.DataSourceWithConfigure = &customDeclarationsDataSource{}
)

type customDeclarationsDataSource struct {
	client *simplemdm.Client
}

type customDeclarationsDataSourceModel struct {
	Search             types.String                                   `tfsdk:"search"`
	CustomDeclarations []customDeclarationsDataSourceDeclarationModel `tfsdk:"custom_declarations"`
}

type customDeclarationsDataSourceDeclarationModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	DeclarationType        types.String `tfsdk:"declaration_type"`
	UserScope              types.Bool   `tfsdk:"user_scope"`
	AttributeSupport       types.Bool   `tfsdk:"attribute_support"`
	EscapeAttributes       types.Bool   `tfsdk:"escape_attributes"`
	ActivationPredicate    types.String `tfsdk:"activation_predicate"`
	ReinstallAfterOsUpdate types.Bool   `tfsdk:"reinstall_after_os_update"`
	ProfileIdentifier      types.String `tfsdk:"profile_identifier"`
	GroupCount             types.Int64  `tfsdk:"group_count"`
	DeviceCount            types.Int64  `tfsdk:"device_count"`
}

func CustomDeclarationsDataSource() datasource.DataSource {
	return &customDeclarationsDataSource{}
}

func (d *customDeclarationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customdeclarations"
}

func (d *customDeclarationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the collection of custom declarations from your SimpleMDM account.",
		Attributes: map[string]schema.Attribute{
			"search": schema.StringAttribute{
				Optional:    true,
				Description: "Limit response to declarations with matching name.",
			},
		},
		Blocks: map[string]schema.Block{
			"custom_declarations": schema.ListNestedBlock{
				Description: "Collection of custom declaration records returned by the API.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Custom declaration identifier.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "A name for the custom declaration.",
						},
						"declaration_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of declaration being defined.",
						},
						"user_scope": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the declaration is scoped to users or devices.",
						},
						"attribute_support": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether variable expansion is enabled.",
						},
						"escape_attributes": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether custom variables are escaped.",
						},
						"activation_predicate": schema.StringAttribute{
							Computed:    true,
							Description: "Predicate controlling when the declaration activates.",
						},
						"reinstall_after_os_update": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether to reinstall the declaration after macOS updates.",
						},
						"profile_identifier": schema.StringAttribute{
							Computed:    true,
							Description: "Identifier assigned by SimpleMDM.",
						},
						"group_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of device groups assigned to the declaration.",
						},
						"device_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of devices assigned to the declaration.",
						},
					},
				},
			},
		},
	}
}

func (d *customDeclarationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config customDeclarationsDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	search := ""
	if !config.Search.IsNull() {
		search = config.Search.ValueString()
	}

	declarations, err := fetchAllCustomDeclarations(ctx, d.client, search)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list SimpleMDM custom declarations",
			err.Error(),
		)
		return
	}

	entries := make([]customDeclarationsDataSourceDeclarationModel, 0, len(declarations))
	for _, decl := range declarations {
		entry := customDeclarationsDataSourceDeclarationModel{
			ID:                     types.StringValue(decl.ID),
			Name:                   types.StringValue(decl.Attributes.Name),
			DeclarationType:        types.StringValue(decl.Attributes.DeclarationType),
			ActivationPredicate:    stringValueOrNull(decl.Attributes.ActivationPredicate),
			ProfileIdentifier:      stringValueOrNull(decl.Attributes.ProfileIdentifier),
		}

		if decl.Attributes.UserScope != nil {
			entry.UserScope = types.BoolValue(*decl.Attributes.UserScope)
		} else {
			entry.UserScope = types.BoolNull()
		}

		if decl.Attributes.AttributeSupport != nil {
			entry.AttributeSupport = types.BoolValue(*decl.Attributes.AttributeSupport)
		} else {
			entry.AttributeSupport = types.BoolNull()
		}

		if decl.Attributes.EscapeAttributes != nil {
			entry.EscapeAttributes = types.BoolValue(*decl.Attributes.EscapeAttributes)
		} else {
			entry.EscapeAttributes = types.BoolNull()
		}

		if decl.Attributes.ReinstallAfterOsUpdate != nil {
			entry.ReinstallAfterOsUpdate = types.BoolValue(*decl.Attributes.ReinstallAfterOsUpdate)
		} else {
			entry.ReinstallAfterOsUpdate = types.BoolNull()
		}

		if decl.Attributes.GroupCount != nil {
			entry.GroupCount = types.Int64Value(*decl.Attributes.GroupCount)
		} else {
			entry.GroupCount = types.Int64Null()
		}

		if decl.Attributes.DeviceCount != nil {
			entry.DeviceCount = types.Int64Value(*decl.Attributes.DeviceCount)
		} else {
			entry.DeviceCount = types.Int64Null()
		}

		entries = append(entries, entry)
	}

	config.CustomDeclarations = entries

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func (d *customDeclarationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*simplemdm.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// fetchAllCustomDeclarations retrieves all custom declarations with pagination support
func fetchAllCustomDeclarations(ctx context.Context, client *simplemdm.Client, search string) ([]customDeclarationDataList, error) {
	var allDeclarations []customDeclarationDataList
	startingAfter := ""
	limit := 100

	for {
		url := fmt.Sprintf("https://%s/api/v1/custom_declarations?limit=%d", client.HostName, limit)
		if search != "" {
			url += fmt.Sprintf("&search=%s", strings.ReplaceAll(search, " ", "+"))
		}
		if startingAfter != "" {
			url += fmt.Sprintf("&starting_after=%s", startingAfter)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		body, err := client.RequestResponse200(req)
		if err != nil {
			return nil, err
		}

		var response customDeclarationsAPIResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}

		allDeclarations = append(allDeclarations, response.Data...)

		if !response.HasMore {
			break
		}

		if len(response.Data) > 0 {
			startingAfter = response.Data[len(response.Data)-1].ID
		} else {
			break
		}
	}

	return allDeclarations, nil
}

// customDeclarationsAPIResponse represents the paginated API response
type customDeclarationsAPIResponse struct {
	Data    []customDeclarationDataList `json:"data"`
	HasMore bool                        `json:"has_more"`
}

// customDeclarationDataList represents a single declaration in the list response
type customDeclarationDataList struct {
	ID         string                      `json:"id"`
	Type       string                      `json:"type"`
	Attributes customDeclarationAttributes `json:"attributes"`
}
