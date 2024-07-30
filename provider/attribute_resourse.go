package provider

import (
	"context"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &attributeResource{}
	_ resource.ResourceWithConfigure   = &attributeResource{}
	_ resource.ResourceWithImportState = &attributeResource{}
)

// attributeResourceModel maps the resource schema data.
type attributeResourceModel struct {
	DefaultValue types.String `tfsdk:"default_value"`
	Name         types.String `tfsdk:"name"`
	ID           types.String `tfsdk:"id"`
}

// AttributeResource is a helper function to simplify the provider implementation.
func AttributeResource() resource.Resource {
	return &attributeResource{}
}

// AttributeResource is the resource implementation.
type attributeResource struct {
	client *simplemdm.Client
}

// Configure adds the provider configured client to the resource.
func (r *attributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*simplemdm.Client)
}

// Metadata returns the resource type name.
func (r *attributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_attribute"
}

// Schema defines the schema for the resource.
func (r *attributeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Attribute resourse can be used to manage SimpleMDM Attribute. Can be used together with Device(s) or Device Group(s) to set values or in lifecycle management.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Required. The name of the Custom Attribute. This name will be used when referencing the Custom Attribute throughout the provider. Alphanumeric characters and underscores only. Case insensitive. Changing name after plan apply will result in replacement(Destroy and Create of new)",
			},
			"default_value": schema.StringAttribute{
				Optional:    true,
				Description: "Optional. The value that will be used if the Attribute value is not provided on Group or Device level.",
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "ID of a Attribute in SimpleMDM",
			},
		},
	}
}

func (r *attributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// Create a new resource
func (r *attributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan attributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	_, err := r.client.AttributeCreate(plan.Name.ValueString(), plan.DefaultValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating attribute",
			"Could not create attribute, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	//mataDataLink := fmt.Sprintf("%s/%s/%s", r.client.HostName, "private", secret.MetadataKey)
	//plan.MetaDataLink = types.StringValue(mataDataLink)
	//plan.SecretValue = types.StringValue(secret.Value)

	//secretLink := fmt.Sprintf("%s/%s/%s", r.client.HostName, "secret", secret.SecretKey)
	//plan.SecretLink = types.StringValue(secretLink)

	plan.ID = plan.Name

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *attributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state attributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed attribute value from SimpleMDM
	attribute, err := r.client.AttributeGet(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM Attribute",
			"Could not read SimpleMDM Attribute ID "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.Name = types.StringValue(attribute.Data.Attributes.Name)
	if attribute.Data.Attributes.DefaultValue != "" {
		state.DefaultValue = types.StringValue(attribute.Data.Attributes.DefaultValue)
	}
	state.ID = types.StringValue(attribute.Data.Attributes.Name)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *attributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Retrieve values from plan
	var plan attributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	err := r.client.AttributeUpdate(plan.Name.ValueString(), plan.DefaultValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating attribute",
			"Could not create attribute, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *attributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state attributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.AttributeDelete(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SimpleMDM attribute",
			"Could not attribute, unexpected error: "+err.Error(),
		)
		return
	}
}
