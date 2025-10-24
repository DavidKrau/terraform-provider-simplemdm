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
		Name:           "Custom Attributes",
		Endpoint:       "/api/v1/custom_attributes",
		ResourceType:   "simplemdm_attribute",
		DataSourceType: "simplemdm_attribute",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Custom-Attributes",
	},
	{
		Name:           "Custom Configuration Profiles",
		Endpoint:       "/api/v1/custom_configuration_profiles",
		ResourceType:   "simplemdm_customprofile",
		DataSourceType: "simplemdm_customprofile",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Custom-Configuration-Profiles",
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
		Name:           "Enrollments",
		Endpoint:       "/api/v1/enrollments",
		ResourceType:   "simplemdm_enrollment",
		DataSourceType: "simplemdm_enrollment",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Enrollments",
	},
	{
		Name:           "Profiles",
		Endpoint:       "/api/v1/profiles",
		ResourceType:   "simplemdm_profile",
		DataSourceType: "simplemdm_profile",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Profiles",
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
