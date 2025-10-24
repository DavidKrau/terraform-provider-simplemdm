package provider

import (
	"os"
	"testing"
)

// testAccPreCheck ensures acceptance tests run only when TF_ACC is enabled
// and the required authentication information is present.
func testAccPreCheck(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests require TF_ACC to be set")
	}

	if os.Getenv("SIMPLEMDM_APIKEY") == "" {
		t.Skip("Acceptance tests require SIMPLEMDM_APIKEY to be set")
	}
}

// testAccRequireEnv fetches an environment variable or skips the current
// test if the variable is not defined. This allows contributors to opt-in to
// tests that require additional fixture data without breaking CI runs.
func testAccRequireEnv(t *testing.T, name string) string {
	value := os.Getenv(name)
	if value == "" {
		t.Skipf("Acceptance test requires %s to be set", name)
	}

	return value
}
