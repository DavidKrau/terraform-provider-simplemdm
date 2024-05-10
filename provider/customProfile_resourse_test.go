package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomProfileResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name= "testprofile"
			mobileconfig = "./testfiles/testprofile.mobileconfig"
			filesha =    "${filesha256("./testfiles/testprofile.mobileconfig")}"
			userscope = true
			attributesupport = true
			escapeattributes = true
			reinstallafterosupdate = true
			
		  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "name", "testprofile"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "mobileconfig", "./testfiles/testprofile.mobileconfig"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "filesha", "46bda319b9fdc88856cf37fb7556b20990bb10538484af2eb7679f5e39a6ea51"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "userscope", "true"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "attributesupport", "true"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "escapeattributes", "true"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "reinstallafterosupdate", "true"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_customprofile.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "simplemdm_customprofile.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the HashiCups
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"filesha", "mobileconfig"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_customprofile" "test" {
					name= "testprofile2"
					mobileconfig = "./testfiles/testprofile2.mobileconfig"
					filesha = "${filesha256("./testfiles/testprofile2.mobileconfig")}"
					userscope = false
					attributesupport = false
					escapeattributes = false
					reinstallafterosupdate = false
					
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "name", "testprofile2"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "mobileconfig", "./testfiles/testprofile2.mobileconfig"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "filesha", "251b9ba70dd8d24b2ed6c4f785751c3cd4bd7c13170f008c2d7edb86bc1db989"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "userscope", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "attributesupport", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "escapeattributes", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "reinstallafterosupdate", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
