package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomProfileDataSource(t *testing.T) {
	testAccPreCheck(t)

	profileID := testAccRequireEnv(t, "SIMPLEMDM_CUSTOM_PROFILE_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`data "simplemdm_customprofile" "test" {id ="%s"}`, profileID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify returned values
					resource.TestCheckResourceAttr("data.simplemdm_customprofile.test", "id", profileID),
					resource.TestCheckResourceAttrSet("data.simplemdm_customprofile.test", "name"),
				),
			},
		},
	})
}
