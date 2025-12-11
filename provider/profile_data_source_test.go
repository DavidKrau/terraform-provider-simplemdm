package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccProfileDataSource requires an existing profile because profiles are read-only
// in the SimpleMDM API. Profiles can only be created through the SimpleMDM web UI or
// using the custom_configuration_profiles API endpoint.
//
// To run this test, set SIMPLEMDM_PROFILE_ID to an existing profile's ID from your SimpleMDM account.
func TestAccProfileDataSource(t *testing.T) {
	testAccPreCheck(t)

	profileID := testAccRequireEnv(t, "SIMPLEMDM_PROFILE_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`data "simplemdm_profile" "test" {id ="%s"}`, profileID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify returned values
					resource.TestCheckResourceAttr("data.simplemdm_profile.test", "id", profileID),
					resource.TestCheckResourceAttrSet("data.simplemdm_profile.test", "name"),
					resource.TestCheckResourceAttrSet("data.simplemdm_profile.test", "type"),
					resource.TestCheckResourceAttrSet("data.simplemdm_profile.test", "profile_identifier"),
					resource.TestCheckResourceAttrSet("data.simplemdm_profile.test", "user_scope"),
					resource.TestCheckResourceAttrSet("data.simplemdm_profile.test", "reinstall_after_os_update"),
					resource.TestCheckResourceAttrSet("data.simplemdm_profile.test", "group_count"),
					resource.TestCheckResourceAttrSet("data.simplemdm_profile.test", "device_count"),
				),
			},
		},
	})
}