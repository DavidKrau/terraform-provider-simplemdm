package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAttributeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_attribute" "testattribute" {
					name= "newAttribute"
					default_value= "test value for test attribute"
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_attribute.testattribute", "name", "newAttribute"),
					resource.TestCheckResourceAttr("simplemdm_attribute.testattribute", "default_value", "test value for test attribute"),
					resource.TestCheckResourceAttr("simplemdm_attribute.testattribute", "id", "newAttribute"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "simplemdm_attribute.testattribute",
				ImportState:       true,
				ImportStateVerify: true,
				//ImportStateVerifyIgnore: []string{"filesha", "mobileconfig"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_attribute" "testattribute" {
					name= "newAttribute2"
					default_value= ""
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_attribute.testattribute", "name", "newAttribute2"),
					resource.TestCheckResourceAttr("simplemdm_attribute.testattribute", "default_value", ""),
					resource.TestCheckResourceAttr("simplemdm_attribute.testattribute", "id", "newAttribute2"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
