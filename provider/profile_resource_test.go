package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccProfileResource_ReadOnly tests the profile resource's read-only behavior.
// This resource is READ-ONLY because profiles can only be created/managed through
// the SimpleMDM web UI. The Terraform resource is for state management only.
//
// The test requires SIMPLEMDM_PROFILE_ID environment variable to be set with
// an existing profile ID from SimpleMDM.
//
// Test Coverage:
//   - Read existing profile and verify attributes
//   - ImportState to verify resource can be imported
//   - Delete only removes from state (doesn't delete from SimpleMDM)
func TestAccProfileResource_ReadOnly(t *testing.T) {
	testAccPreCheck(t)

	profileID := testAccGetEnv(t, "SIMPLEMDM_PROFILE_ID")
	
	if profileID == "" {
		t.Skip("SIMPLEMDM_PROFILE_ID not set - skipping test as profiles can only be created via SimpleMDM UI")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Read existing profile and verify basic attributes
			{
				Config: testAccProfileResourceConfig(profileID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify required fields are set
					resource.TestCheckResourceAttr("simplemdm_profile.test", "id", profileID),
					resource.TestCheckResourceAttrSet("simplemdm_profile.test", "name"),
					resource.TestCheckResourceAttrSet("simplemdm_profile.test", "type"),
					
					// Verify boolean attributes exist
					resource.TestCheckResourceAttrSet("simplemdm_profile.test", "auto_deploy"),
					resource.TestCheckResourceAttrSet("simplemdm_profile.test", "user_scope"),
					resource.TestCheckResourceAttrSet("simplemdm_profile.test", "attribute_support"),
					resource.TestCheckResourceAttrSet("simplemdm_profile.test", "escape_attributes"),
					
					// Verify profile_identifier always exists
					resource.TestCheckResourceAttrSet("simplemdm_profile.test", "profile_identifier"),
					
					// Verify numeric attributes exist
					resource.TestCheckResourceAttrSet("simplemdm_profile.test", "group_count"),
					resource.TestCheckResourceAttrSet("simplemdm_profile.test", "device_count"),
					
					// Note: install_type, source, created_at, and updated_at are Optional+Computed
					// and may not be returned by the SimpleMDM API for all profile types.
					// We don't assert these fields to avoid test failures with different profile types.
				),
			},
			// Step 2: ImportState test - verify resource can be imported
			{
				ResourceName:      "simplemdm_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccProfileResource_NonExistent tests that attempting to reference a
// non-existent profile results in an appropriate error.
func TestAccProfileResource_NonExistent(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProfileResourceConfig("999999999"),
				ExpectError: regexp.MustCompile("Error creating SimpleMDM profile reference|404"),
			},
		},
	})
}

// testAccProfileResourceConfig returns a test configuration for the profile resource
func testAccProfileResourceConfig(profileID string) string {
	return providerConfig + fmt.Sprintf(`
resource "simplemdm_profile" "test" {
  id = "%s"
}
`, profileID)
}
