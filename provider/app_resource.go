package provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &appResource{}
	_ resource.ResourceWithConfigure   = &appResource{}
	_ resource.ResourceWithImportState = &appResource{}
)

// appResourceModel maps the resource schema data.
type appResourceModel struct {
	Name       types.String `tfsdk:"name"`
	ID         types.String `tfsdk:"id"`
	AppStoreId types.String `tfsdk:"app_store_id"`
	BundleId   types.String `tfsdk:"bundle_id"`
	DeployTo   types.String `tfsdk:"deploy_to"`
}

func AppResource() resource.Resource {
	return &appResource{}
}

// appResource is the resource implementation.
type appResource struct {
	client *simplemdm.Client
}

// Configure adds the provider configured client to the resource.
func (r *appResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*simplemdm.Client)
}

// Metadata returns the resource type name.
func (r *appResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

// Schema defines the schema for the resource.
func (r *appResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "App resource can be used to manage Apps.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name that SimpleMDM will use to reference this app. If left blank, SimpleMDM will automatically set this to the app name specified by the binary.",
			},
			"app_store_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Required. The Apple App Store ID of the app to be added. Example: 1090161858.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("app_store_id"),
						path.MatchRoot("bundle_id"),
					),
				},
			},
			"bundle_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Required. The bundle identifier of the Apple App Store app to be added. Example: com.myCompany.MyApp1",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("app_store_id"),
						path.MatchRoot("bundle_id"),
					),
				},
			},
			"deploy_to": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Optional. Deploy the app to associated devices immediately after the app has been uploaded and processed. Possible values are none, outdated or all. Defaults to none.",
				Default:     stringdefault.StaticString("none"),
				Validators: []validator.String{
					stringvalidator.OneOf("outdated", "all"),
				},
			},
		},
	}
}

func (r *appResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create a new resource
func (r *appResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan appResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var appStoreId, bundleId, name string
	if !plan.AppStoreId.IsNull() {
		appStoreId = plan.AppStoreId.ValueString()
	}
	if !plan.BundleId.IsNull() {
		bundleId = plan.BundleId.ValueString()
	}
	if !plan.Name.IsNull() {
		name = plan.Name.ValueString()
	}

	// Définir la valeur par défaut pour deploy_to
	if plan.DeployTo.IsNull() {
		plan.DeployTo = types.StringValue("none")
	}

	// Generate API request body from plan
	app, err := r.client.AppCreate(
		appStoreId,
		bundleId,
		name,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating app",
			"Could not create app, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(app.Data.ID))

	// Mettre à jour les attributs computés
	if app.Data.Attributes.AppStoreId != 0 {
		plan.AppStoreId = types.StringValue(strconv.Itoa(app.Data.Attributes.AppStoreId))
	}
	if app.Data.Attributes.BundleId != "" {
		plan.BundleId = types.StringValue(app.Data.Attributes.BundleId)
	}
	if app.Data.Attributes.Name != "" {
		plan.Name = types.StringValue(app.Data.Attributes.Name)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *appResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state appResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing app
	err := r.client.AppDelete(state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting SimpleMDM app",
			"Could not delete app, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *appResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state appResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := r.client.AppGet(state.ID.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SimpleMDM App",
			"Could not read SimpleMDM App "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.ID = types.StringValue(strconv.Itoa(app.Data.ID))
	state.Name = types.StringValue(app.Data.Attributes.Name)

	// Les attributs computés sont mis à jour automatiquement
	if app.Data.Attributes.AppStoreId != 0 {
		state.AppStoreId = types.StringValue(strconv.Itoa(app.Data.Attributes.AppStoreId))
	}
	if app.Data.Attributes.BundleId != "" {
		state.BundleId = types.StringValue(app.Data.Attributes.BundleId)
	}

	if state.DeployTo.IsNull() {
		state.DeployTo = types.StringValue("none")
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *appResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state appResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deployTo := plan.DeployTo
	if deployTo.IsNull() {
		deployTo = types.StringValue("none")
	}

	_, err := r.client.AppUpdate(
		plan.ID.ValueString(),
		plan.Name.ValueString(),
		deployTo.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating app",
			"Failed to update app: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}
