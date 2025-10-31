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

var (
	_ datasource.DataSource              = &scriptJobsDataSource{}
	_ datasource.DataSourceWithConfigure = &scriptJobsDataSource{}
)

type scriptJobsDataSource struct {
	client *simplemdm.Client
}

type scriptJobsDataSourceModel struct {
	ScriptJobs []scriptJobsDataSourceJobModel `tfsdk:"script_jobs"`
}

type scriptJobsDataSourceJobModel struct {
	ID                  types.String `tfsdk:"id"`
	ScriptID            types.String `tfsdk:"script_id"`
	AssignmentGroupID   types.String `tfsdk:"assignment_group_id"`
	AssignmentGroupName types.String `tfsdk:"assignment_group_name"`
	Status              types.String `tfsdk:"status"`
	CreatedAt           types.String `tfsdk:"created_at"`
	UpdatedAt           types.String `tfsdk:"updated_at"`
}

func ScriptJobsDataSource() datasource.DataSource {
	return &scriptJobsDataSource{}
}

func (d *scriptJobsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scriptjobs"
}

func (d *scriptJobsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the collection of script jobs from your SimpleMDM account.",
		Blocks: map[string]schema.Block{
			"script_jobs": schema.ListNestedBlock{
				Description: "Collection of script job records returned by the API.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Script job identifier.",
						},
						"script_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the script associated with this job.",
						},
						"assignment_group_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the assignment group this job targets.",
						},
						"assignment_group_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the assignment group.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Current status of the script job.",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp when the job was created.",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp when the job was last updated.",
						},
					},
				},
			},
		},
	}
}

func (d *scriptJobsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config scriptJobsDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	scriptJobs, err := fetchAllScriptJobs(ctx, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list SimpleMDM script jobs",
			err.Error(),
		)
		return
	}

	entries := make([]scriptJobsDataSourceJobModel, 0, len(scriptJobs))
	for _, job := range scriptJobs {
		flat := flattenScriptJob(&job)

		entry := scriptJobsDataSourceJobModel{
			ID:                  types.StringValue(strconv.Itoa(flat.ID)),
			Status:              types.StringValue(flat.Status),
			CreatedAt:           types.StringValue(flat.CreatedAt),
			UpdatedAt:           types.StringValue(flat.UpdatedAt),
			AssignmentGroupName: types.StringValue(flat.AssignmentGroupName),
		}

		if flat.ScriptID != nil {
			entry.ScriptID = types.StringValue(strconv.Itoa(*flat.ScriptID))
		} else {
			entry.ScriptID = types.StringNull()
		}

		if flat.AssignmentGroupID != nil {
			entry.AssignmentGroupID = types.StringValue(strconv.Itoa(*flat.AssignmentGroupID))
		} else {
			entry.AssignmentGroupID = types.StringNull()
		}

		entries = append(entries, entry)
	}

	config.ScriptJobs = entries

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

func (d *scriptJobsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// fetchAllScriptJobs retrieves all script jobs with pagination support
func fetchAllScriptJobs(ctx context.Context, client *simplemdm.Client) ([]scriptJobResponse, error) {
	return listScriptJobs(ctx, client, 0)
}
