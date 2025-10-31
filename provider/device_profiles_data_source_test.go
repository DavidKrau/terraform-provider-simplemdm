package provider

import "testing"

// TestAccDeviceProfilesDataSource requires an actual enrolled device because:
// 1. Devices cannot be created via the SimpleMDM API
// 2. Profiles must be assigned and installed on the device
// 3. The device must have time to install and report profiles
//
// To run this test, set SIMPLEMDM_DEVICE_ID to an enrolled device with assigned profiles.
func TestAccDeviceProfilesDataSource(t *testing.T) {
	testAccPreCheck(t)
	t.Skip("Requires enrolled device with installed profiles. Set SIMPLEMDM_DEVICE_ID to enable.")
}
