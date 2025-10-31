package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type ResourceDefinition struct {
	TypeName      string
	Factory       func() resource.Resource
	DocsPath      string
	ExampleDirs   []string
	TestFiles     []string
	APIEndpoints  []string
	TestsOptional bool
}

type DataSourceDefinition struct {
	TypeName     string
	Factory      func() datasource.DataSource
	DocsPath     string
	ExampleDirs  []string
	TestFiles    []string
	APIEndpoints []string
}

var resourceDefinitions = []ResourceDefinition{
	{
		TypeName:     "simplemdm_app",
		Factory:      AppResource,
		DocsPath:     "docs/resources/app.md",
		ExampleDirs:  []string{"examples/resources/simplemdm_app"},
		TestFiles:    []string{"provider/app_resource_test.go"},
		APIEndpoints: []string{"/api/v1/apps"},
	},
	{
		TypeName:     "simplemdm_attribute",
		Factory:      AttributeResource,
		DocsPath:     "docs/resources/attribute.md",
		ExampleDirs:  []string{"examples/resources/simplemdm_attribute"},
		TestFiles:    []string{"provider/attribute_resource_test.go"},
		APIEndpoints: []string{"/api/v1/custom_attributes"},
	},
	{
		TypeName:     "simplemdm_assignmentgroup",
		Factory:      AssignmentGroupResource,
		DocsPath:     "docs/resources/assignmentgroup.md",
		ExampleDirs:  []string{"examples/resources/simplemdm_assignmentgroup"},
		TestFiles:    []string{"provider/assignmentGroup_resource_test.go"},
		APIEndpoints: []string{"/api/v1/assignment_groups"},
	},
	{
		TypeName:     "simplemdm_customprofile",
		Factory:      CustomProfileResource,
		DocsPath:     "docs/resources/customprofile.md",
		ExampleDirs:  []string{"examples/resources/simplemdm_customprofile"},
		TestFiles:    []string{"provider/customProfile_resource_test.go"},
		APIEndpoints: []string{"/api/v1/custom_profiles"},
	},
	{
		TypeName:     "simplemdm_profile",
		Factory:      ProfileResource,
		DocsPath:     "docs/resources/profile.md",
		ExampleDirs:  []string{"examples/resources/simplemdm_profile"},
		TestFiles:    []string{"provider/profile_resource_test.go"},
		APIEndpoints: []string{"/api/v1/profiles"},
	},
	{
		TypeName:     "simplemdm_customdeclaration",
		Factory:      CustomDeclarationResource,
		DocsPath:     "docs/resources/customdeclaration.md",
		ExampleDirs:  []string{"examples/resources/simplemdm_customdeclaration"},
		TestFiles:    []string{"provider/customDeclaration_resource_test.go"},
		APIEndpoints: []string{"/api/v1/custom_declarations"},
	},
	{
		TypeName:      "simplemdm_customdeclaration_device_assignment",
		Factory:       CustomDeclarationDeviceAssignmentResource,
		DocsPath:      "docs/resources/customdeclaration_device_assignment.md",
		ExampleDirs:   []string{"examples/resources/simplemdm_customdeclaration_device_assignment"},
		TestFiles:     []string{"provider/customDeclaration_device_assignment_resource_test.go"},
		APIEndpoints:  []string{"/api/v1/custom_declarations/{custom_declaration_id}/devices/{device_id}"},
		TestsOptional: true,
	},
	{
		TypeName:     "simplemdm_device",
		Factory:      DeviceResource,
		DocsPath:     "docs/resources/device.md",
		ExampleDirs:  []string{"examples/resources/simplemdm_device"},
		TestFiles:    []string{"provider/device_resource_test.go"},
		APIEndpoints: []string{"/api/v1/devices"},
	},
	{
		TypeName:    "simplemdm_device_command",
		Factory:     DeviceCommandResource,
		DocsPath:    "docs/resources/device_command.md",
		ExampleDirs: []string{"examples/resources/simplemdm_device_command"},
		TestFiles:   []string{"provider/device_command_resource_test.go"},
		APIEndpoints: []string{
			"/api/v1/devices/{DEVICE_ID}/push_apps",
			"/api/v1/devices/{DEVICE_ID}/refresh",
			"/api/v1/devices/{DEVICE_ID}/restart",
			"/api/v1/devices/{DEVICE_ID}/shutdown",
			"/api/v1/devices/{DEVICE_ID}/lock",
			"/api/v1/devices/{DEVICE_ID}/clear_passcode",
			"/api/v1/devices/{DEVICE_ID}/clear_firmware_password",
			"/api/v1/devices/{DEVICE_ID}/rotate_firmware_password",
			"/api/v1/devices/{DEVICE_ID}/clear_recovery_lock_password",
			"/api/v1/devices/{DEVICE_ID}/clear_restrictions_password",
			"/api/v1/devices/{DEVICE_ID}/rotate_recovery_lock_password",
			"/api/v1/devices/{DEVICE_ID}/rotate_filevault_key",
			"/api/v1/devices/{DEVICE_ID}/set_admin_password",
			"/api/v1/devices/{DEVICE_ID}/rotate_admin_password",
			"/api/v1/devices/{DEVICE_ID}/wipe",
			"/api/v1/devices/{DEVICE_ID}/update_os",
			"/api/v1/devices/{DEVICE_ID}/remote_desktop",
			"/api/v1/devices/{DEVICE_ID}/bluetooth",
			"/api/v1/devices/{DEVICE_ID}/set_time_zone",
			"/api/v1/devices/{DEVICE_ID}/unenroll",
			"/api/v1/devices/{DEVICE_ID}/users/{USER_ID}",
		},
		TestsOptional: true,
	},
	{
		TypeName:     "simplemdm_devicegroup",
		Factory:      DeviceGroupResource,
		DocsPath:     "docs/resources/devicegroup.md",
		ExampleDirs:  []string{"examples/resources/simplemdm_devicegroup"},
		TestFiles:    []string{"provider/deviceGroup_resource_test.go"},
		APIEndpoints: []string{"/api/v1/device_groups"},
	},
	{
		TypeName:    "simplemdm_enrollment",
		Factory:     EnrollmentResource,
		DocsPath:    "docs/resources/enrollment.md",
		ExampleDirs: []string{"examples/resources/simplemdm_enrollment"},
		TestFiles: []string{
			"provider/enrollment_resource_test.go",
		},
		APIEndpoints: []string{
			"/api/v1/enrollments",
			"/api/v1/enrollments/{enrollment_id}/invitations",
		},
	},
	{
		TypeName:     "simplemdm_script",
		Factory:      ScriptResource,
		DocsPath:     "docs/resources/script.md",
		ExampleDirs:  []string{"examples/resources/simplemdm_script"},
		TestFiles:    []string{"provider/script_resource_test.go"},
		APIEndpoints: []string{"/api/v1/scripts"},
	},
	{
		TypeName:     "simplemdm_scriptjob",
		Factory:      ScriptJobResource,
		DocsPath:     "docs/resources/scriptjob.md",
		ExampleDirs:  []string{"examples/resources/simplemdm_scriptjob"},
		TestFiles:    []string{"provider/scriptJob_resource_test.go"},
		APIEndpoints: []string{"/api/v1/script_jobs"},
	},
	{
		TypeName:     "simplemdm_managed_config",
		Factory:      ManagedConfigResource,
		DocsPath:     "docs/resources/managed_config.md",
		ExampleDirs:  []string{"examples/resources/simplemdm_managed_config"},
		TestFiles:    []string{"provider/managedConfig_resource_test.go"},
		APIEndpoints: []string{"/api/v1/apps/{APP_ID}/managed_configs", "/api/v1/apps/{APP_ID}/managed_configs/push"},
	},
}

var dataSourceDefinitions = []DataSourceDefinition{
	{
		TypeName:     "simplemdm_app",
		Factory:      AppDataSource,
		DocsPath:     "docs/data-sources/app.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_app"},
		TestFiles:    []string{"provider/app_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/apps"},
	},
	{
		TypeName:     "simplemdm_attribute",
		Factory:      AttributeDataSource,
		DocsPath:     "docs/data-sources/attribute.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_attribute"},
		TestFiles:    []string{"provider/attribute_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/custom_attributes"},
	},
	{
		TypeName:     "simplemdm_assignmentgroup",
		Factory:      AssignmentGroupDataSource,
		DocsPath:     "docs/data-sources/assignmentgroup.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_assignmentgroup"},
		TestFiles:    []string{"provider/assignmentGroup_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/assignment_groups"},
	},
	{
		TypeName:     "simplemdm_customprofile",
		Factory:      CustomProfileDataSource,
		DocsPath:     "docs/data-sources/customprofile.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_customprofile"},
		TestFiles:    []string{"provider/customProfile_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/custom_profiles"},
	},
	{
		TypeName:     "simplemdm_customdeclaration",
		Factory:      CustomDeclarationDataSource,
		DocsPath:     "docs/data-sources/customdeclaration.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_customdeclaration"},
		TestFiles:    []string{"provider/customDeclaration_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/custom_declarations"},
	},
	{
		TypeName:     "simplemdm_device",
		Factory:      DeviceDataSource,
		DocsPath:     "docs/data-sources/device.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_device"},
		TestFiles:    []string{"provider/device_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/devices"},
	},
	{
		TypeName:     "simplemdm_devices",
		Factory:      DevicesDataSource,
		DocsPath:     "docs/data-sources/devices.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_devices"},
		TestFiles:    []string{"provider/devices_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/devices"},
	},
	{
		TypeName:     "simplemdm_device_profiles",
		Factory:      DeviceProfilesDataSource,
		DocsPath:     "docs/data-sources/device_profiles.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_device_profiles"},
		TestFiles:    []string{"provider/device_profiles_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/devices/{DEVICE_ID}/profiles"},
	},
	{
		TypeName:     "simplemdm_device_installed_apps",
		Factory:      DeviceInstalledAppsDataSource,
		DocsPath:     "docs/data-sources/device_installed_apps.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_device_installed_apps"},
		TestFiles:    []string{"provider/device_installed_apps_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/devices/{DEVICE_ID}/installed_apps"},
	},
	{
		TypeName:     "simplemdm_device_users",
		Factory:      DeviceUsersDataSource,
		DocsPath:     "docs/data-sources/device_users.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_device_users"},
		TestFiles:    []string{"provider/device_users_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/devices/{DEVICE_ID}/users"},
	},
	{
		TypeName:     "simplemdm_devicegroup",
		Factory:      DeviceGroupDataSource,
		DocsPath:     "docs/data-sources/devicegroup.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_devicegroup"},
		TestFiles:    []string{"provider/deviceGroup_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/device_groups"},
	},
	{
		TypeName:     "simplemdm_enrollment",
		Factory:      EnrollmentDataSource,
		DocsPath:     "docs/data-sources/enrollment.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_enrollment"},
		TestFiles:    []string{"provider/enrollment_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/enrollments"},
	},
	{
		TypeName:     "simplemdm_profile",
		Factory:      ProfileDataSource,
		DocsPath:     "docs/data-sources/profile.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_profile"},
		TestFiles:    []string{"provider/profile_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/profiles"},
	},
	{
		TypeName:     "simplemdm_script",
		Factory:      ScriptDataSource,
		DocsPath:     "docs/data-sources/script.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_script"},
		TestFiles:    []string{"provider/script_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/scripts"},
	},
	{
		TypeName:     "simplemdm_scriptjob",
		Factory:      ScriptJobDataSource,
		DocsPath:     "docs/data-sources/scriptjob.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_scriptjob"},
		TestFiles:    []string{"provider/scriptJob_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/script_jobs"},
	},
	{
		TypeName:     "simplemdm_managed_config",
		Factory:      ManagedConfigDataSource,
		DocsPath:     "docs/data-sources/managed_config.md",
		ExampleDirs:  []string{"examples/data-sources/simplemdm_managed_config"},
		TestFiles:    []string{"provider/managedConfig_data_source_test.go"},
		APIEndpoints: []string{"/api/v1/apps/{APP_ID}/managed_configs"},
	},
}

func ResourceFactories() []func() resource.Resource {
	factories := make([]func() resource.Resource, 0, len(resourceDefinitions))
	for _, definition := range resourceDefinitions {
		factories = append(factories, definition.Factory)
	}

	return factories
}

func DataSourceFactories() []func() datasource.DataSource {
	factories := make([]func() datasource.DataSource, 0, len(dataSourceDefinitions))
	for _, definition := range dataSourceDefinitions {
		factories = append(factories, definition.Factory)
	}

	return factories
}

func ResourceDefinitionMap() map[string]ResourceDefinition {
	result := make(map[string]ResourceDefinition, len(resourceDefinitions))
	for _, definition := range resourceDefinitions {
		result[definition.TypeName] = definition
	}

	return result
}

func DataSourceDefinitionMap() map[string]DataSourceDefinition {
	result := make(map[string]DataSourceDefinition, len(dataSourceDefinitions))
	for _, definition := range dataSourceDefinitions {
		result[definition.TypeName] = definition
	}

	return result
}

func ResourceDefinitions() []ResourceDefinition {
	return append([]ResourceDefinition(nil), resourceDefinitions...)
}

func DataSourceDefinitions() []DataSourceDefinition {
	return append([]DataSourceDefinition(nil), dataSourceDefinitions...)
}
