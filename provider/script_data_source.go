package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &scriptDataSource{}
	_ datasource.DataSourceWithConfigure = &scriptDataSource{}
)

// scriptDataSourceModel maps the data source schema data.
type scriptDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	ScriptFile      types.String `tfsdk:"scriptfile"`
	VariableSupport types.Bool   `tfsdk:"variablesupport"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
	CreatedBy       types.String `tfsdk:"created_by"`
}

// scriptDataSource is a helper function to simplify the provider implementation.
func ScriptDataSource() datasource.DataSource {
	return &scriptDataSource{}
}

// scriptDataSource is the data source implementation.
type scriptDataSource struct {
	client *simplemdm.Client
}

// Metadata returns the data source type name.
func (d *scriptDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_script"
}

// Schema defines the schema for the data source.
func (d *scriptDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Script data source can be used together with Script Jobs(s).",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the Script.",
			},
			"scriptfile": schema.StringAttribute{
				Computed:    true,
				Description: "The script content.",
			},
			"variablesupport": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether variable support is enabled for this script.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the Script was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the Script was last updated.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "User that created the Script in SimpleMDM.",
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Script.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *scriptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state scriptDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	script, err := d.client.ScriptGet(state.ID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"SimpleMDM script not found",
				fmt.Sprintf("The script with ID %s was not found. It may have been deleted.", state.ID.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to Read SimpleMDM Script",
				err.Error(),
			)
		}
		return
	}

	// Map response body to model
	state.Name = types.StringValue(script.Data.Attributes.Name)
	state.ID = types.StringValue(strconv.Itoa(script.Data.ID))
	state.ScriptFile = types.StringValue(script.Data.Attributes.Content)
	state.VariableSupport = types.BoolValue(script.Data.Attributes.VariableSupport)
	state.CreatedAt = types.StringValue(script.Data.Attributes.CreatedAt)
	state.UpdatedAt = types.StringValue(script.Data.Attributes.UpdatedAt)
	state.CreatedBy = types.StringValue(script.Data.Attributes.CreateBy)

	// Set state

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *scriptDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
