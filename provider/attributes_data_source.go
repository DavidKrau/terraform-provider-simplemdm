package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &attributesDataSource{}
	_ datasource.DataSourceWithConfigure = &attributesDataSource{}
)

type attributesDataSource struct {
	client *simplemdm.Client
}

type attributesDataSourceModel struct {
	Attributes []attributesDataSourceAttributeModel `tfsdk:"attributes"`
}

type attributesDataSourceAttributeModel struct {
	Name         types.String `tfsdk:"name"`
	DefaultValue types.String `tfsdk:"default_value"`
}

func AttributesDataSource() datasource.DataSource {
	return &attributesDataSource{}
}

func (d *attributesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_attributes"
}

func (d *attributesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the collection of custom attributes from your SimpleMDM account.",
		Blocks: map[string]schema.Block{
			"attributes": schema.ListNestedBlock{
				Description: "Collection of custom attribute records returned by the API.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the custom attribute.",
						},
						"default_value": schema.StringAttribute{
							Computed:    true,
							Description: "Default (global) value of the custom attribute.",
						},
					},
				},
			},
		},
	}
}

func (d *attributesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config attributesDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	attributes, err := fetchAllAttributes(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list SimpleMDM custom attributes",
			err.Error(),
		)
		return
	}

	entries := make([]attributesDataSourceAttributeModel, 0, len(attributes))
	for _, attr := range attributes {
		entry := attributesDataSourceAttributeModel{
			Name:         types.StringValue(attr.Attributes.Name),
			DefaultValue: types.StringValue(attr.Attributes.DefaultValue),
		}

		entries = append(entries, entry)
	}

	config.Attributes = entries

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// fetchAllAttributes retrieves all custom attributes using the API
func fetchAllAttributes(ctx context.Context, client *simplemdm.Client) ([]attributeData, error) {
	url := fmt.Sprintf("https://%s/api/v1/custom_attributes", client.HostName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := client.RequestResponse200(req)
	if err != nil {
		return nil, err
	}

	var response attributesAPIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// attributesAPIResponse represents the API response for attributes list
type attributesAPIResponse struct {
	Data []attributeData `json:"data"`
}

// attributeData represents a single attribute in the list response
type attributeData struct {
	Type       string                  `json:"type"`
	Attributes attributeDataAttributes `json:"attributes"`
}

type attributeDataAttributes struct {
	Name         string `json:"name"`
	DefaultValue string `json:"default_value"`
}

func (d *attributesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
