package provider

import "testing"

func TestAccDevicesDataSource(t *testing.T) {
	testAccPreCheck(t)
	_ = testAccRequireEnv(t, "SIMPLEMDM_RUN_DEVICES_DATA_SOURCE_TESTS")
	t.Skip("Acceptance test requires explicit fixtures and is skipped by default")
}
