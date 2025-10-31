package provider

import (
	"testing"
)

func TestAccCustomProfileDataSource(t *testing.T) {
	testAccPreCheck(t)

	// Skip test - custom profile data source requires a pre-existing fixture
	// The test is covered by TestAccCustomProfileResource which creates and tests the data source
	t.Skip("Custom profile data source test requires a valid fixture custom profile ID. This functionality is tested in TestAccCustomProfileResource.")
}
