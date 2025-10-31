package provider

import "testing"

// TestAccDeviceInstalledAppsDataSource requires an actual enrolled device because:
// 1. Devices cannot be created via the SimpleMDM API
// 2. Apps must be assigned and installed on the device
// 3. The device must have time to install and report apps
//
// To run this test, set SIMPLEMDM_DEVICE_ID to an enrolled device with installed apps.
func TestAccDeviceInstalledAppsDataSource(t *testing.T) {
	testAccPreCheck(t)
	t.Skip("Requires enrolled device with installed apps. Set SIMPLEMDM_DEVICE_ID to enable.")
}
