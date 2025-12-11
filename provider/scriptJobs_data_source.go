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
	_ datasource.DataSource              = &scriptJobsDataSource{}
	_ datasource.DataSourceWithConfigure = &scriptJobsDataSource{}
)

type scriptJobsDataSource struct {
	client *simplemdm.Client
}

type scriptJobsDataSourceModel struct {
	Limit          types.Int64                    `tfsdk:"limit"`
	StartingAfter  types.Int64                    `tfsdk:"starting_after"`
	ScriptJobs     []scriptJobsDataSourceJobModel `tfsdk:"script_jobs"`
}

type scriptJobsDataSourceJobModel struct {
	ID                   types.String `tfsdk:"id"`
	ScriptName           types.String `tfsdk:"script_name"`
	JobName              types.String `tfsdk:"job_name"`
	JobIdentifier        types.String `tfsdk:"job_identifier"`
	Content              types.String `tfsdk:"content"`
	VariableSupport      types.Bool   `tfsdk:"variable_support"`
	Status               types.String `tfsdk:"status"`
	PendingCount         types.Int64  `tfsdk:"pending_count"`
	SuccessCount         types.Int64  `tfsdk:"success_count"`
	ErroredCount         types.Int64  `tfsdk:"errored_count"`
	CustomAttribute      types.String `tfsdk:"custom_attribute"`
	CustomAttributeRegex types.String `tfsdk:"custom_attribute_regex"`
	CreatedBy            types.String `tfsdk:"created_by"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
	ScriptID             types.String `tfsdk:"script_id"`
	AssignmentGroupID    types.String `tfsdk:"assignment_group_id"`
	Devices              types.List   `tfsdk:"devices"`
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
		Attributes: map[string]schema.Attribute{
			"limit": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum number of script jobs to return (1-100). Defaults to 100. When not specified, all jobs will be fetched with automatic pagination.",
			},
			"starting_after": schema.Int64Attribute{
				Optional:    true,
				Description: "Return script jobs with IDs after this value for pagination.",
			},
		},
		Blocks: map[string]schema.Block{
			"script_jobs": schema.ListNestedBlock{
				Description: "Collection of script job records returned by the API.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Script job identifier.",
						},
						"script_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the script that was executed.",
						},
						"job_name": schema.StringAttribute{
							Computed:    true,
							Description: "Human friendly name of the job.",
						},
						"job_identifier": schema.StringAttribute{
							Computed:    true,
							Description: "Short identifier string for the job (different from the numeric ID).",
						},
						"content": schema.StringAttribute{
							Computed:    true,
							Description: "Script contents that were executed by the job.",
						},
						"variable_support": schema.BoolAttribute{
							Computed:    true,
							Description: "Indicates whether the script supports variables.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Current execution status of the job (pending, completed-with-errors, completed, cancelled).",
						},
						"pending_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of devices that have not yet reported a result.",
						},
						"success_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of devices that completed successfully.",
						},
						"errored_count": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of devices that failed to execute the script.",
						},
						"custom_attribute": schema.StringAttribute{
							Computed:    true,
							Description: "Custom attribute ID that stores the job output, when configured.",
						},
						"custom_attribute_regex": schema.StringAttribute{
							Computed:    true,
							Description: "Regular expression used to sanitize the custom attribute output.",
						},
						"created_by": schema.StringAttribute{
							Computed:    true,
							Description: "User or API key that created the job.",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp when the job was created.",
						},
						"updated_at": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp when the job was last updated.",
						},
						"script_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the script associated with this job.",
						},
						"assignment_group_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the assignment group this job targets.",
						},
						"devices": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Execution results for each targeted device.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:    true,
										Description: "Device identifier.",
									},
									"status": schema.StringAttribute{
										Computed:    true,
										Description: "Execution status reported for the device.",
									},
									"status_code": schema.StringAttribute{
										Computed:    true,
										Description: "Optional status code returned by the device.",
									},
									"response": schema.StringAttribute{
										Computed:    true,
										Description: "Output returned by the device, when available.",
									},
								},
							},
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

	// Determine starting point for pagination
	startingAfter := 0
	if !config.StartingAfter.IsNull() && !config.StartingAfter.IsUnknown() {
		startingAfter = int(config.StartingAfter.ValueInt64())
	}

	// Fetch script jobs (with automatic pagination if limit not specified)
	var scriptJobs []scriptJobResponse
	var err error
	
	if !config.Limit.IsNull() && !config.Limit.IsUnknown() {
		// User specified a limit, fetch only one page
		scriptJobs, err = listScriptJobsWithLimit(ctx, d.client, startingAfter, int(config.Limit.ValueInt64()))
	} else {
		// No limit specified, fetch all pages automatically
		scriptJobs, err = fetchAllScriptJobs(ctx, d.client)
	}
	
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list SimpleMDM script jobs",
			err.Error(),
		)
		return
	}

	entries := make([]scriptJobsDataSourceJobModel, 0, len(scriptJobs))
	for _, job := range scriptJobs {
		// Fetch full details for each job to get all fields
		details, err := fetchScriptJobDetails(ctx, d.client, strconv.Itoa(job.Data.ID))
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to fetch script job details",
				fmt.Sprintf("Failed to fetch details for script job %d: %s", job.Data.ID, err.Error()),
			)
			return
		}

		entry := scriptJobsDataSourceJobModel{
			ID:                   types.StringValue(details.ID),
			ScriptName:           stringValueOrNull(details.ScriptName),
			JobName:              stringValueOrNull(details.JobName),
			JobIdentifier:        stringValueOrNull(details.JobIdentifier),
			Content:              stringValueOrNull(details.Content),
			VariableSupport:      types.BoolValue(details.VariableSupport),
			Status:               stringValueOrNull(details.Status),
			PendingCount:         types.Int64Value(details.PendingCount),
			SuccessCount:         types.Int64Value(details.SuccessCount),
			ErroredCount:         types.Int64Value(details.ErroredCount),
			CreatedBy:            stringValueOrNull(details.CreatedBy),
			CreatedAt:            stringValueOrNull(details.CreatedAt),
			UpdatedAt:            stringValueOrNull(details.UpdatedAt),
		}

		if details.CustomAttribute != "" {
			entry.CustomAttribute = types.StringValue(details.CustomAttribute)
		} else {
			entry.CustomAttribute = types.StringNull()
		}

		if details.CustomAttributeRegex != "" {
			entry.CustomAttributeRegex = types.StringValue(details.CustomAttributeRegex)
		} else {
			entry.CustomAttributeRegex = types.StringNull()
		}

		// Get script ID and assignment group ID from the job response
		flat := flattenScriptJob(&job)
		if flat.ScriptID != nil {
			entry.ScriptID = types.StringValue(*flat.ScriptID)
		} else {
			entry.ScriptID = types.StringNull()
		}

		if flat.AssignmentGroupID != nil {
			entry.AssignmentGroupID = types.StringValue(*flat.AssignmentGroupID)
		} else {
			entry.AssignmentGroupID = types.StringNull()
		}

		// Add device information
		devices, diags := scriptJobDevicesListValue(ctx, details.Devices)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		entry.Devices = devices

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

// fetchAllScriptJobs retrieves all script jobs with automatic pagination
func fetchAllScriptJobs(ctx context.Context, client *simplemdm.Client) ([]scriptJobResponse, error) {
	return listScriptJobs(ctx, client, 0)
}

// listScriptJobsWithLimit retrieves a single page of script jobs with specified limit
func listScriptJobsWithLimit(ctx context.Context, client *simplemdm.Client, startingAfter int, limit int) ([]scriptJobResponse, error) {
	if limit < 1 || limit > 100 {
		limit = 100
	}
	
	url := fmt.Sprintf("https://%s/api/v1/script_jobs?limit=%d", client.HostName, limit)
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

	var response scriptJobsListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	var jobs []scriptJobResponse
	for _, data := range response.Data {
		jobs = append(jobs, scriptJobResponse{Data: data})
	}

	return jobs, nil
}
