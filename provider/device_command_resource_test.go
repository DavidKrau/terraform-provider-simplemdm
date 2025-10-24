package provider

import "testing"

func TestAccDeviceCommandResource(t *testing.T) {
	testAccPreCheck(t)
	_ = testAccRequireEnv(t, "SIMPLEMDM_RUN_DEVICE_COMMAND_TESTS")
	t.Skip("Acceptance test requires explicit fixtures and is skipped by default")
}
