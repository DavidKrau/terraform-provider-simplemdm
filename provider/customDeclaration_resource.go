package provider

import (
	"context"
	"encoding/json"
	"strconv"

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
	_ resource.Resource                = &customDeclarationResource{}
	_ resource.ResourceWithConfigure   = &customDeclarationResource{}
	_ resource.ResourceWithImportState = &customDeclarationResource{}
)

// declarationResourceModel maps the resource schema data.
type customDeclarationResourceModel struct {
	Name                  types.String `tfsdk:"name"`
	Declaration           types.String `tfsdk:"declaration"`
	DeclarationType       types.String `tfsdk:"declaration_type"`
	UserScope             types.Bool   `tfsdk:"userscope"`
	AttributeSupport      types.Bool   `tfsdk:"attributesupport"`
	EscapeAttributes      types.Bool   `tfsdk:"escapeattributes"`
	ActivatetionPredicate types.String `tfsdk:"activation_predicate"`
	ID                    types.String `tfsdk:"id"`
}

// DeclarationResource is a helper function to simplify the provider implementation.
func CustomDeclarationResource() resource.Resource {
	return &customDeclarationResource{}
}

// declarationResource is the resource implementation.
type customDeclarationResource struct {
	client *simplemdm.Client
}

// Configure adds the provider configured client to the resource.
func (r *customDeclarationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*simplemdm.Client)
}

// Metadata returns the resource type name.
func (r *customDeclarationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_customdeclaration"
}

// Schema defines the schema for the resource.
func (r *customDeclarationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Custom Declaration resource can be used to manage Custom Declaration. Can be used together with Device(s) and Group(s) and set addition details regarding Custom Declaration.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Required. A name for the declaration. Example: \"My First declaration by terraform\"",
			},
			"declaration": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Required. Can be string or you can use function 'file' or 'templatefile' to load string from file (see examples folder). Example: declaration = file(\"./declarations/declaration.json\") or declaration = <<-EOT DECLARATION STRING EOT",
			},
			"declaration_type": schema.StringAttribute{
				Required:    true,
				Optional:    false,
				Description: "Required. The type of declaration being defined",
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "ID of a Custom Declaration in SimpleMDM",
			},
			"userscope": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(true),
				Computed:    true,
				Description: "Optional. A boolean true or false. If false, deploy as a device declaration instead of a user declaration for macOS devices. Defaults to true.",
			},
			"attributesupport": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Optional. A boolean true or false. When enabled, SimpleMDM will process variables in the uploaded declaration. Defaults to false",
			},
			"escapeattributes": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Optional. A boolean true or false. When enabled, SimpleMDM escape the values of the custom variables in the uploaded declaration. Defaults to false",
			},
			"activation_predicate": schema.StringAttribute{
				Optional:    true,
				Description: "Optional. A predicate format string as Apple's Predicate Programming describes. The activation only installs when the predicate evaluates to true or if it is left blank.",
			},
		},
	}
}

func (r *customDeclarationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource
func (r *customDeclarationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	//Retrieve values from plan
	var plan customDeclarationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	declaration, err := r.client.CustomDeclarationCreate(plan.Name.ValueString(), plan.DeclarationType.ValueString(), plan.Declaration.ValueString(), plan.UserScope.ValueBool(), plan.AttributeSupport.ValueBool(), plan.EscapeAttributes.ValueBool(), plan.ActivatetionPredicate.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating declaration",
			"Could not create declration, unexpected error: "+err.Error(),
		)
		return
	}

	// h := sha256.New()
	// h.Write([]byte(plan.MobileConfig.ValueString()))
	// sha256_hash := hex.EncodeToString(h.Sum(nil))[0:32]

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(strconv.Itoa(declaration.Data.ID))
	//plan.FileSHA = types.StringValue(sha256_hash)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *customDeclarationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state customDeclarationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//retrive one needs to be implemented first
	// Get refreshed declaration values from SimpleMDM
	//https://a.simplemdm.com/api/v1/profiles/211589
	declaration, err := r.client.ProfileGet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM custom declaration",
			"Could not read custom declaration ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	declarationStruct, err := r.client.CustomDeclarationDownload(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM custom declaration",
			"Could not read custom profles ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	rawPayload := declarationStruct.Payload

	var payloadMap map[string]interface{}
	err = json.Unmarshal(rawPayload, &payloadMap)
	if err != nil {
		// Handle error if the payload JSON is invalid
		resp.Diagnostics.AddError("Error Unmarshaling Payload", err.Error())
		return
	}

	var activationPredicate string
	if val, ok := payloadMap["activation_predicate"].(string); ok {
		activationPredicate = val
		delete(payloadMap, "activation_predicate") // Remove it after extraction
	}

	delete(payloadMap, "declaration_name")

	cleanedPayloadBytes, err := json.Marshal(payloadMap)
	//cleanedPayloadBytes, err := json.MarshalIndent(payloadMap, "", "  ")
	if err != nil {
		// Handle error if re-marshaling fails
		resp.Diagnostics.AddError("Error Marshaling Cleaned Payload", err.Error())
		return
	}

	finalPayloadString := string(cleanedPayloadBytes)

	//get it from call line 170
	state.DeclarationType = types.StringValue(declarationStruct.Type)
	state.ActivatetionPredicate = types.StringValue(activationPredicate)
	state.Declaration = types.StringValue(finalPayloadString)

	state.Name = types.StringValue(declaration.Data.Attributes.Name)
	state.UserScope = types.BoolValue(declaration.Data.Attributes.UserScope)
	state.AttributeSupport = types.BoolValue(declaration.Data.Attributes.AttributeSupport)
	state.EscapeAttributes = types.BoolValue(declaration.Data.Attributes.EscapeAttributes)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *customDeclarationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//Retrieve values from plan
	var plan customDeclarationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	_, err := r.client.CustomDeclarationUpdate(plan.Name.ValueString(), plan.DeclarationType.ValueString(), plan.Declaration.ValueString(), plan.UserScope.ValueBool(), plan.AttributeSupport.ValueBool(), plan.EscapeAttributes.ValueBool(), plan.ID.ValueString(), plan.ActivatetionPredicate.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating declration",
			"Could not update declration, unexpected error: "+err.Error(),
		)
		return
	}

	// h := sha256.New()
	// h.Write([]byte(plan.MobileConfig.ValueString()))
	// sha256_hash := hex.EncodeToString(h.Sum(nil))[0:32]

	// plan.FileSHA = types.StringValue(sha256_hash)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *customDeclarationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state customDeclarationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing custom declaration
	err := r.client.CustomDeclarationDelete(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SimpleMDM custom declaration",
			"Could not delete custom declaration, unexpected error: "+err.Error(),
		)
		return
	}
}
