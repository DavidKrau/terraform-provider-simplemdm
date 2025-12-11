package provider

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	Content         types.String `tfsdk:"content"`
	ID              types.String `tfsdk:"id"`
	VariableSupport types.Bool   `tfsdk:"variable_support"`
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
			"content": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Required. The script content. Must begin with a valid shebang (e.g., #!/bin/sh). Can be loaded from a file using the file() or templatefile() functions. Example: content = file(\"./scripts/script.sh\")",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^#!`),
						"script content must begin with a shebang (e.g., #!/bin/sh)",
					),
				},
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "ID of a Script in SimpleMDM",
			},
			"variable_support": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Optional. Whether to enable variable support in this script. The provider converts boolean values to the API's expected format. Defaults to false.",
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
	script, err := r.client.ScriptCreate(plan.Name.ValueString(), plan.VariableSupport.ValueBool(), plan.Content.ValueString())
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
	plan.Content = types.StringValue(script.Data.Attributes.Content)
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
	state.Content = types.StringValue(script.Data.Attributes.Content)
	state.CreatedAt = types.StringValue(script.Data.Attributes.CreatedAt)
	state.UpdatedAt = types.StringValue(script.Data.Attributes.UpdatedAt)
	state.VariableSupport = types.BoolValue(script.Data.Attributes.VariableSupport)

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
	script, err := r.client.ScriptUpdate(plan.Name.ValueString(), plan.VariableSupport.ValueBool(), plan.Content.ValueString(), plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating script",
			"Could not update script, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Content = types.StringValue(script.Data.Attributes.Content)
	plan.UpdatedAt = types.StringValue(script.Data.Attributes.UpdatedAt)
	plan.CreatedAt = types.StringValue(script.Data.Attributes.CreatedAt)
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
