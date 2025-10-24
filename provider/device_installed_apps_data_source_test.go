package provider

import "testing"

func TestAccDeviceInstalledAppsDataSource(t *testing.T) {
	testAccPreCheck(t)
	_ = testAccRequireEnv(t, "SIMPLEMDM_RUN_DEVICE_INSTALLED_APPS_DATA_SOURCE_TESTS")
	t.Skip("Acceptance test requires explicit fixtures and is skipped by default")
}
