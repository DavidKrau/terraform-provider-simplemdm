package apicatalog

// EndpointCoverage describes how a Terraform construct maps to a SimpleMDM API
// endpoint. The catalog is used in tests to guarantee that every documented API
// endpoint has a matching Terraform resource and/or data source implementation.
type EndpointCoverage struct {
	Name           string
	Endpoint       string
	ResourceType   string
	DataSourceType string
	DocsURL        string
}

var Catalog = []EndpointCoverage{
	{
		Name:           "Apps",
		Endpoint:       "/api/v1/apps",
		ResourceType:   "simplemdm_app",
		DataSourceType: "simplemdm_app",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Apps",
	},
	{
		Name:           "Assignment Groups",
		Endpoint:       "/api/v1/assignment_groups",
		ResourceType:   "simplemdm_assignmentgroup",
		DataSourceType: "simplemdm_assignmentgroup",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Assignment-Groups",
	},
	{
		Name:           "Attributes",
		Endpoint:       "/api/v1/attributes",
		ResourceType:   "simplemdm_attribute",
		DataSourceType: "simplemdm_attribute",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Custom-Attributes",
	},
	{
		Name:           "Custom Profiles",
		Endpoint:       "/api/v1/custom_profiles",
		ResourceType:   "simplemdm_customprofile",
		DataSourceType: "simplemdm_customprofile",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Custom-Profiles",
	},
	{
		Name:           "Custom Declarations",
		Endpoint:       "/api/v1/custom_declarations",
		ResourceType:   "simplemdm_customdeclaration",
		DataSourceType: "simplemdm_customdeclaration",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Custom-Declarations",
	},
	{
		Name:           "Devices",
		Endpoint:       "/api/v1/devices",
		ResourceType:   "simplemdm_device",
		DataSourceType: "simplemdm_device",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:           "Device Groups",
		Endpoint:       "/api/v1/device_groups",
		ResourceType:   "simplemdm_devicegroup",
		DataSourceType: "simplemdm_devicegroup",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Device-Groups",
	},
	{
		Name:           "Profiles",
		Endpoint:       "/api/v1/profiles",
		ResourceType:   "",
		DataSourceType: "simplemdm_profile",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Profiles",
		// TODO: Add a simplemdm_profile resource to cover profile CRUD operations.
	},
	{
		Name:           "Scripts",
		Endpoint:       "/api/v1/scripts",
		ResourceType:   "simplemdm_script",
		DataSourceType: "simplemdm_script",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Scripts",
	},
	{
		Name:           "Script Jobs",
		Endpoint:       "/api/v1/script_jobs",
		ResourceType:   "simplemdm_scriptjob",
		DataSourceType: "simplemdm_scriptjob",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Script-Jobs",
	},
}
