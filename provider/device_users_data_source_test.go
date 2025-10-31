package provider

import "testing"

// TestAccDeviceUsersDataSource requires an actual enrolled device because:
// 1. Devices cannot be created via the SimpleMDM API
// 2. Device must have associated user accounts from enrollment
// 3. User data is populated during device enrollment and usage
//
// To run this test, set SIMPLEMDM_DEVICE_ID to an enrolled device with user accounts.
func TestAccDeviceUsersDataSource(t *testing.T) {
	testAccPreCheck(t)
	t.Skip("Requires enrolled device with user accounts. Set SIMPLEMDM_DEVICE_ID to enable.")
}
