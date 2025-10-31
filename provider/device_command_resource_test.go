package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDeviceCommandResource_PushApps tests the push_apps command.
// This is a safe, commonly-used command that pushes assigned apps to a device.
//
// Requirements:
//   - SIMPLEMDM_DEVICE_ID environment variable must be set with an enrolled device ID
//   - The device must be enrolled and active in SimpleMDM
//
// Note: Device commands are "one-shot" resources - they execute once during creation
// and cannot be read back from the API. The resource only maintains local state.
func TestAccDeviceCommandResource_PushApps(t *testing.T) {
	testAccPreCheck(t)
	deviceID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceCommandResourceConfig(deviceID, "push_assigned_apps"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the command was set correctly
					resource.TestCheckResourceAttr("simplemdm_device_command.test", "command", "push_assigned_apps"),
					resource.TestCheckResourceAttr("simplemdm_device_command.test", "device_id", deviceID),
					
					// Verify the command executed successfully (202 Accepted)
					resource.TestCheckResourceAttr("simplemdm_device_command.test", "status_code", "202"),
					
					// Verify ID was generated
					resource.TestCheckResourceAttrSet("simplemdm_device_command.test", "id"),
				),
			},
		},
	})
}

// TestAccDeviceCommandResource_Refresh tests the refresh command.
// This is a safe command that triggers a device to check in with SimpleMDM.
//
// Requirements:
//   - SIMPLEMDM_DEVICE_ID environment variable must be set with an enrolled device ID
//   - The device must be enrolled and active in SimpleMDM
func TestAccDeviceCommandResource_Refresh(t *testing.T) {
	testAccPreCheck(t)
	deviceID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceCommandResourceConfig(deviceID, "refresh"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_device_command.test", "command", "refresh"),
					resource.TestCheckResourceAttr("simplemdm_device_command.test", "device_id", deviceID),
					resource.TestCheckResourceAttr("simplemdm_device_command.test", "status_code", "202"),
					resource.TestCheckResourceAttrSet("simplemdm_device_command.test", "id"),
				),
			},
		},
	})
}

// TestAccDeviceCommandResource_Lock tests the lock command.
// This command locks the device screen with an optional message.
//
// Requirements:
//   - SIMPLEMDM_DEVICE_ID environment variable must be set with an enrolled device ID
//   - The device must be enrolled and active in SimpleMDM
//   - WARNING: This will lock the device and may require physical access to unlock
func TestAccDeviceCommandResource_Lock(t *testing.T) {
	testAccPreCheck(t)
	deviceID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceCommandResourceConfigWithParams(
					deviceID,
					"lock",
					map[string]string{
						"message": "Test lock from Terraform provider",
					},
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_device_command.test", "command", "lock"),
					resource.TestCheckResourceAttr("simplemdm_device_command.test", "device_id", deviceID),
					resource.TestCheckResourceAttr("simplemdm_device_command.test", "status_code", "202"),
					resource.TestCheckResourceAttrSet("simplemdm_device_command.test", "id"),
				),
			},
		},
	})
}

// TestAccDeviceCommandResource_InvalidCommand tests that an unsupported command
// results in an appropriate error.
func TestAccDeviceCommandResource_InvalidCommand(t *testing.T) {
	testAccPreCheck(t)
	deviceID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDeviceCommandResourceConfig(deviceID, "invalid_command"),
				ExpectError: regexp.MustCompile("Unsupported device command|not currently supported"),
			},
		},
	})
}

// testAccDeviceCommandResourceConfig returns a test configuration for a device command
// without parameters.
func testAccDeviceCommandResourceConfig(deviceID, command string) string {
	return providerConfig + fmt.Sprintf(`
resource "simplemdm_device_command" "test" {
  device_id = "%s"
  command   = "%s"
}
`, deviceID, command)
}

// testAccDeviceCommandResourceConfigWithParams returns a test configuration for a device
// command with parameters.
func testAccDeviceCommandResourceConfigWithParams(deviceID, command string, params map[string]string) string {
	config := providerConfig + fmt.Sprintf(`
resource "simplemdm_device_command" "test" {
  device_id = "%s"
  command   = "%s"
  parameters = {
`, deviceID, command)

        for key, value := range params {
                config += fmt.Sprintf("    %s = %q\n", key, value)
        }

	config += `  }
}
`
	return config
}
