package provider

import "testing"

func TestAccDeviceProfilesDataSource(t *testing.T) {
	testAccPreCheck(t)
	_ = testAccRequireEnv(t, "SIMPLEMDM_RUN_DEVICE_PROFILES_DATA_SOURCE_TESTS")
	t.Skip("Acceptance test requires explicit fixtures and is skipped by default")
}
