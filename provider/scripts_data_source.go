package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &scriptsDataSource{}
	_ datasource.DataSourceWithConfigure = &scriptsDataSource{}
)

type scriptsDataSource struct {
	client *simplemdm.Client
}

type scriptsDataSourceModel struct {
	Scripts []scriptsDataSourceScriptModel `tfsdk:"scripts"`
}

type scriptsDataSourceScriptModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	VariableSupport types.Bool   `tfsdk:"variable_support"`
	CreatedBy       types.String `tfsdk:"created_by"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
}

func ScriptsDataSource() datasource.DataSource {
	return &scriptsDataSource{}
}

func (d *scriptsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scripts"
}

func (d *scriptsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the collection of scripts from your SimpleMDM account.",
		Blocks: map[string]schema.Block{
			"scripts": schema.ListNestedBlock{
				Description: "Collection of script records returned by the API.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Script identifier.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the script.",
						},
						"variable_support": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether variable support is enabled for this script.",
						},
						"created_by": schema.StringAttribute{
							Computed:    true,
							Description: "User that created the script.",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp when the script was created.",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp when the script was last updated.",
						},
					},
				},
			},
		},
	}
}

func (d *scriptsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config scriptsDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	scripts, err := fetchAllScripts(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list SimpleMDM scripts",
			err.Error(),
		)
		return
	}

	entries := make([]scriptsDataSourceScriptModel, 0, len(scripts))
	for _, script := range scripts {
		entry := scriptsDataSourceScriptModel{
			ID:              types.StringValue(strconv.Itoa(script.ID)),
			Name:            types.StringValue(script.Attributes.Name),
			VariableSupport: types.BoolValue(script.Attributes.VariableSupport),
			CreatedBy:       types.StringValue(script.Attributes.CreateBy),
			CreatedAt:       types.StringValue(script.Attributes.CreatedAt),
			UpdatedAt:       types.StringValue(script.Attributes.UpdatedAt),
		}

		entries = append(entries, entry)
	}

	config.Scripts = entries

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func (d *scriptsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// fetchAllScripts retrieves all scripts with pagination support
func fetchAllScripts(ctx context.Context, client *simplemdm.Client) ([]scriptDataList, error) {
	var allScripts []scriptDataList
	startingAfter := 0
	limit := 100

	for {
		url := fmt.Sprintf("https://%s/api/v1/scripts?limit=%d", client.HostName, limit)
		if startingAfter > 0 {
			url += fmt.Sprintf("&starting_after=%d", startingAfter)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		body, err := client.RequestResponse200(req)
		if err != nil {
			return nil, err
		}

		var response scriptsAPIResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}

		allScripts = append(allScripts, response.Data...)

		if !response.HasMore {
			break
		}

		if len(response.Data) > 0 {
			startingAfter = response.Data[len(response.Data)-1].ID
		} else {
			break
		}
	}

	return allScripts, nil
}

// scriptsAPIResponse represents the paginated API response
type scriptsAPIResponse struct {
	Data    []scriptDataList `json:"data"`
	HasMore bool             `json:"has_more"`
}

// scriptDataList represents a single script in the list response
type scriptDataList struct {
	ID         int                  `json:"id"`
	Type       string               `json:"type"`
	Attributes scriptListAttributes `json:"attributes"`
}

type scriptListAttributes struct {
	Name            string `json:"name"`
	VariableSupport bool   `json:"variable_support"`
	CreateBy        string `json:"created_by"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}
