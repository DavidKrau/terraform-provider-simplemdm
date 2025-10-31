package provider

import (
	"fmt"
	"os"
	"testing"

	simplemdm "github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

// testAccGetEnv fetches an environment variable and returns its value
// or an empty string if not defined. Does not skip the test.
func testAccGetEnv(t *testing.T, name string) string {
	return os.Getenv(name)
}

// getTestClient returns a SimpleMDM client configured from environment variables
// for use in test CheckDestroy functions
func getTestClient() (*simplemdm.Client, error) {
	apiKey := os.Getenv("SIMPLEMDM_APIKEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SIMPLEMDM_APIKEY environment variable not set")
	}

	host := os.Getenv("SIMPLEMDM_HOST")
	if host == "" {
		host = "a.simplemdm.com"
	}

	return simplemdm.NewClient(host, apiKey), nil
}

// testAccCheckDestroy is a helper to verify resource destruction
// It takes a resource type name and a function to check if the resource exists
func testAccCheckResourceDestroyed(resourceType string, checkExists func(*simplemdm.Client, string) error) func(*terraform.State) error {
	return func(s *terraform.State) error {
		client, err := getTestClient()
		if err != nil {
			return fmt.Errorf("failed to create test client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			err := checkExists(client, rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("%s %s still exists after destroy", resourceType, rs.Primary.ID)
			}

			if !isNotFoundError(err) {
				return fmt.Errorf("unexpected error checking %s %s: %w", resourceType, rs.Primary.ID, err)
			}
		}

		return nil
	}
}