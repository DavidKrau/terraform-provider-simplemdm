package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &scriptResource{}
	_ resource.ResourceWithConfigure   = &scriptResource{}
	_ resource.ResourceWithImportState = &scriptResource{}
)

// scriptResourceModel maps the resource schema data.
type scriptResourceModel struct {
	Name            types.String `tfsdk:"name"`
	ScriptFile      types.String `tfsdk:"scriptfile"`
	ID              types.String `tfsdk:"id"`
	VariableSupport types.Bool   `tfsdk:"variablesupport"`
	CreatedBy       types.String `tfsdk:"created_by"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
}

// scriptResource is a helper function to simplify the provider implementation.
func ScriptResource() resource.Resource {
	return &scriptResource{}
}

// scriptResource is the resource implementation.
type scriptResource struct {
	client *simplemdm.Client
}

// Configure adds the provider configured client to the resource.
func (r *scriptResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*simplemdm.Client)
}

// Metadata returns the resource type name.
func (r *scriptResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_script"
}

// Schema defines the schema for the resource.
func (r *scriptResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Script resource can be used to manage Scripts.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Required. A name for the Script. Example: \"My First Script managed by terraform\"",
			},
			"scriptfile": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Required. Can be directly string or you can use function 'file' or 'templatefile' to load string from file. Emaple: scriptfile = file(\"./scripts/script.sh\") or scriptfile = <<-EOT\n #!/bin/bash\n echo \"Hello!!\"\n EOT ",
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "ID of a Script in SimpleMDM",
			},
			"variablesupport": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Optional. A boolean true or false. Whether or not to enable variable support in this script. Defaults to false",
			},
			"created_by": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "User that created the Script in SimpleMDM",
			},
			"created_at": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Date when script was created in SimpleMDM",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Date when script was updated in SimpleMDM",
			},
		},
	}
}

func (r *scriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource
func (r *scriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan scriptResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	script, err := r.client.ScriptCreate(plan.Name.ValueString(), plan.VariableSupport.ValueBool(), plan.ScriptFile.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating script",
			"Could not create script, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(strconv.Itoa(script.Data.ID))
	plan.CreatedAt = types.StringValue(script.Data.Attributes.CreatedAt)
	plan.UpdatedAt = types.StringValue(script.Data.Attributes.UpdatedAt)
	plan.CreatedBy = types.StringValue(script.Data.Attributes.CreateBy)
	plan.ScriptFile = types.StringValue(script.Data.Attributes.Content)
	plan.VariableSupport = types.BoolValue(script.Data.Attributes.VariableSupport)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *scriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state scriptResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get script values from SimpleMDM
	script, err := r.client.ScriptGet(state.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM Script",
			"Could not read SimpleMDM Script "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.Name = types.StringValue(script.Data.Attributes.Name)
	state.ScriptFile = types.StringValue(script.Data.Attributes.Content)
	state.CreatedAt = types.StringValue(script.Data.Attributes.CreatedAt)
	state.UpdatedAt = types.StringValue(script.Data.Attributes.UpdatedAt)
	state.VariableSupport = types.BoolValue(script.Data.Attributes.VariableSupport)
	state.CreatedBy = types.StringValue(script.Data.Attributes.CreateBy)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *scriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Retrieve values from plan
	var plan scriptResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	script, err := r.client.ScriptUpdate(plan.Name.ValueString(), plan.VariableSupport.ValueBool(), plan.ScriptFile.ValueString(), plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating script",
			"Could not update script, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ScriptFile = types.StringValue(script.Data.Attributes.Content)
	plan.UpdatedAt = types.StringValue(script.Data.Attributes.UpdatedAt)
	plan.CreatedAt = types.StringValue(script.Data.Attributes.CreatedAt)
	plan.CreatedBy = types.StringValue(script.Data.Attributes.CreateBy)
	plan.VariableSupport = types.BoolValue(script.Data.Attributes.VariableSupport)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *scriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state scriptResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing script
	err := r.client.ScriptDelete(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SimpleMDM script",
			"Could not delete script, unexpected error: "+err.Error(),
		)
		return
	}
}
