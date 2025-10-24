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
					resource.TestCheckResourceAttrSet("data.simplemdm_customprofile.test", "mobileconfig"),
					resource.TestCheckResourceAttrSet("data.simplemdm_customprofile.test", "profileidentifier"),
					resource.TestCheckResourceAttrSet("data.simplemdm_customprofile.test", "profilesha"),
					resource.TestCheckResourceAttrSet("data.simplemdm_customprofile.test", "userscope"),
					resource.TestCheckResourceAttrSet("data.simplemdm_customprofile.test", "attributesupport"),
					resource.TestCheckResourceAttrSet("data.simplemdm_customprofile.test", "escapeattributes"),
					resource.TestCheckResourceAttrSet("data.simplemdm_customprofile.test", "reinstallafterosupdate"),
					resource.TestCheckResourceAttrSet("data.simplemdm_customprofile.test", "groupcount"),
					resource.TestCheckResourceAttrSet("data.simplemdm_customprofile.test", "devicecount"),
				),
			},
		},
	})
}
