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
		Name:           "Managed App Configs",
		Endpoint:       "/api/v1/apps/{APP_ID}/managed_configs",
		ResourceType:   "simplemdm_managed_config",
		DataSourceType: "simplemdm_managed_config",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Managed-App-Configs",
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
		Name:           "Devices Collection",
		Endpoint:       "/api/v1/devices",
		DataSourceType: "simplemdm_devices",
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
	{
		Name:           "Device Profiles",
		Endpoint:       "/api/v1/devices/{DEVICE_ID}/profiles",
		DataSourceType: "simplemdm_device_profiles",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:           "Device Installed Apps",
		Endpoint:       "/api/v1/devices/{DEVICE_ID}/installed_apps",
		DataSourceType: "simplemdm_device_installed_apps",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:           "Device Users",
		Endpoint:       "/api/v1/devices/{DEVICE_ID}/users",
		DataSourceType: "simplemdm_device_users",
		DocsURL:        "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Push Assigned Apps",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/push_apps",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Refresh",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/refresh",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Restart",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/restart",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Shutdown",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/shutdown",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Lock",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/lock",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Clear Passcode",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/clear_passcode",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Clear Firmware Password",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/clear_firmware_password",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Rotate Firmware Password",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/rotate_firmware_password",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Clear Recovery Lock Password",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/clear_recovery_lock_password",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Clear Restrictions Password",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/clear_restrictions_password",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Rotate Recovery Lock Password",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/rotate_recovery_lock_password",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Rotate FileVault Recovery Key",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/rotate_filevault_key",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Set Admin Password",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/set_admin_password",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Rotate Admin Password",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/rotate_admin_password",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Wipe",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/wipe",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Update OS",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/update_os",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Enable Remote Desktop",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/remote_desktop (POST)",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Disable Remote Desktop",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/remote_desktop (DELETE)",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Enable Bluetooth",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/bluetooth (POST)",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Disable Bluetooth",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/bluetooth (DELETE)",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Set Time Zone",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/set_time_zone",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Unenroll",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/unenroll",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
	{
		Name:         "Device Command - Delete User",
		Endpoint:     "/api/v1/devices/{DEVICE_ID}/users/{USER_ID}",
		ResourceType: "simplemdm_device_command",
		DocsURL:      "https://api.simplemdm.com/v1/#tag/Devices",
	},
}
