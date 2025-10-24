package provider

import "testing"

func TestAccDeviceInstalledAppsDataSource(t *testing.T) {
	testAccPreCheck(t)
	t.Skip("Acceptance test requires explicit fixtures and is skipped by default")
}
