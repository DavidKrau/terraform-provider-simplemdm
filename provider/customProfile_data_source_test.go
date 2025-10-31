package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomProfileDataSource(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create the profile first
			{
				Config: providerConfig + `
					resource "simplemdm_customprofile" "fixture" {
						name                   = "Test Data Source Profile"
						mobileconfig           = file("./testfiles/testprofile.mobileconfig")
						userscope              = true
						attributesupport       = true
						escapeattributes       = true
						reinstallafterosupdate = true
					}
				`,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("simplemdm_customprofile.fixture", "id"),
				),
			},
			// Then test reading it - giving API time to stabilize
			{
				Config: providerConfig + `
					resource "simplemdm_customprofile" "fixture" {
						name                   = "Test Data Source Profile"
						mobileconfig           = file("./testfiles/testprofile.mobileconfig")
						userscope              = true
						attributesupport       = true
						escapeattributes       = true
						reinstallafterosupdate = true
					}

					data "simplemdm_customprofile" "test" {
						id = simplemdm_customprofile.fixture.id
					}
				`,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.simplemdm_customprofile.test", "id",
						"simplemdm_customprofile.fixture", "id",
					),
					resource.TestCheckResourceAttrSet("data.simplemdm_customprofile.test", "name"),
					resource.TestCheckResourceAttrSet("data.simplemdm_customprofile.test", "mobileconfig"),
					resource.TestCheckResourceAttrSet("data.simplemdm_customprofile.test", "profileidentifier"),
				),
			},
		},
	})
}
