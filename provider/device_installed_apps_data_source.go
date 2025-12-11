package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/DavidKrau/terraform-provider-simplemdm/internal/simplemdmext"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &deviceInstalledAppsDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceInstalledAppsDataSource{}
)

type deviceInstalledAppsDataSource struct {
	client *simplemdm.Client
}

type deviceInstalledAppsDataSourceModel struct {
	DeviceID      types.String                          `tfsdk:"device_id"`
	InstalledApps []deviceInstalledAppDataSourceModel   `tfsdk:"installed_apps"`
}

type deviceInstalledAppDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Type          types.String `tfsdk:"type"`
	Name          types.String `tfsdk:"name"`
	Identifier    types.String `tfsdk:"identifier"`
	Version       types.String `tfsdk:"version"`
	ShortVersion  types.String `tfsdk:"short_version"`
	BundleSize    types.Int64  `tfsdk:"bundle_size"`
	DynamicSize   types.Int64  `tfsdk:"dynamic_size"`
	Managed       types.Bool   `tfsdk:"managed"`
	DiscoveredAt  types.String `tfsdk:"discovered_at"`
	LastSeenAt    types.String `tfsdk:"last_seen_at"`
}

func DeviceInstalledAppsDataSource() datasource.DataSource {
	return &deviceInstalledAppsDataSource{}
}

func (d *deviceInstalledAppsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_installed_apps"
}

func (d *deviceInstalledAppsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves installed applications for a device.",
		Attributes: map[string]schema.Attribute{
			"device_id": schema.StringAttribute{
				Required:    true,
				Description: "Identifier of the device.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\d+$`),
						"device_id must be a numeric string",
					),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"installed_apps": schema.ListNestedBlock{
				Description: "Collection of installed applications on the device.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Installed app identifier.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Installed app resource type.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Application name.",
						},
						"identifier": schema.StringAttribute{
							Computed:    true,
							Description: "Application bundle identifier.",
						},
						"version": schema.StringAttribute{
							Computed:    true,
							Description: "Application version.",
						},
						"short_version": schema.StringAttribute{
							Computed:    true,
							Description: "Application short version string.",
						},
						"bundle_size": schema.Int64Attribute{
							Computed:    true,
							Description: "Size of the application bundle in bytes.",
						},
						"dynamic_size": schema.Int64Attribute{
							Computed:    true,
							Description: "Dynamic size of the application in bytes.",
						},
						"managed": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the application is managed by SimpleMDM.",
						},
						"discovered_at": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp when the application was first discovered.",
						},
						"last_seen_at": schema.StringAttribute{
							Computed:    true,
							Description: "Timestamp when the application was last seen.",
						},
					},
				},
			},
		},
	}
}

func (d *deviceInstalledAppsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deviceInstalledAppsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apps, err := simplemdmext.ListDeviceInstalledApps(ctx, d.client, state.DeviceID.ValueString())
	if err != nil {
		if isNotFoundError(err) {
			resp.Diagnostics.AddError(
				"Device not found",
				fmt.Sprintf("The device with ID %s was not found. It may have been deleted.", state.DeviceID.ValueString()),
			)
		} else {
			resp.Diagnostics.AddError(
				"Unable to list installed apps",
				fmt.Sprintf("Failed to retrieve installed apps for device %s: %s", state.DeviceID.ValueString(), err.Error()),
			)
		}
		return
	}

	items := make([]deviceInstalledAppDataSourceModel, 0, len(apps.Data))
	for _, item := range apps.Data {
		app := deviceInstalledAppDataSourceModel{
			ID:   types.StringValue(item.ID.String()),
			Type: types.StringValue(item.Type),
		}

		// Parse attributes from the map
		if name, ok := item.Attributes["name"].(string); ok {
			app.Name = types.StringValue(name)
		} else {
			app.Name = types.StringNull()
		}

		if identifier, ok := item.Attributes["identifier"].(string); ok {
			app.Identifier = types.StringValue(identifier)
		} else {
			app.Identifier = types.StringNull()
		}

		if version, ok := item.Attributes["version"].(string); ok {
			app.Version = types.StringValue(version)
		} else {
			app.Version = types.StringNull()
		}

		if shortVersion, ok := item.Attributes["short_version"].(string); ok {
			app.ShortVersion = types.StringValue(shortVersion)
		} else {
			app.ShortVersion = types.StringNull()
		}

		if bundleSize, ok := item.Attributes["bundle_size"].(float64); ok {
			app.BundleSize = types.Int64Value(int64(bundleSize))
		} else {
			app.BundleSize = types.Int64Null()
		}

		if dynamicSize, ok := item.Attributes["dynamic_size"].(float64); ok {
			app.DynamicSize = types.Int64Value(int64(dynamicSize))
		} else {
			app.DynamicSize = types.Int64Null()
		}

		if managed, ok := item.Attributes["managed"].(bool); ok {
			app.Managed = types.BoolValue(managed)
		} else {
			app.Managed = types.BoolNull()
		}

		if discoveredAt, ok := item.Attributes["discovered_at"].(string); ok {
			app.DiscoveredAt = types.StringValue(discoveredAt)
		} else {
			app.DiscoveredAt = types.StringNull()
		}

		if lastSeenAt, ok := item.Attributes["last_seen_at"].(string); ok {
			app.LastSeenAt = types.StringValue(lastSeenAt)
		} else {
			app.LastSeenAt = types.StringNull()
		}

		items = append(items, app)
	}

	state.InstalledApps = items

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *deviceInstalledAppsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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