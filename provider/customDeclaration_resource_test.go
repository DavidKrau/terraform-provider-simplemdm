package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomDeclarationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
		resource "simplemdm_customdeclaration" "test" {
			name= "testdeclaration"
			declaration = jsonencode(jsondecode(file("./testfiles/testdeclaration.json")))
			userscope = true
			attributesupport = true
			escapeattributes = true
			declaration_type = "com.apple.configuration.safari.bookmarks"
			activation_predicate = ""			
		  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "name", "testdeclaration"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "declaration", "{\"ManagedBookmarks\":[{\"Bookmarks\":[{\"Title\":\"Public Site\",\"URL\":\"https://www.example.com\"}],\"GroupIdentifier\":\"Group1\",\"Title\":\"Company Bookmarks\"}]}"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "userscope", "true"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "attributesupport", "true"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "escapeattributes", "true"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "declaration_type", "com.apple.configuration.safari.bookmarks"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "activation_predicate", ""),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_customdeclaration.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "simplemdm_customdeclaration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			//Update and Read testing
			{
				Config: providerConfig + `
							resource "simplemdm_customdeclaration" "test" {
								name= "testdeclaration2"
								declaration = jsonencode(jsondecode(file("./testfiles/testdeclaration2.json")))
								userscope = false
								attributesupport = false
								escapeattributes = false
								declaration_type = "com.apple.configuration.safari.bookmarks2"
								activation_predicate = "1234"			
		  					}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "name", "testdeclaration2"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "declaration", "{\"Calculator\":{\"ProgrammerMode\":{\"Enabled\":false},\"ScientificMode\":{\"Enabled\":false}}}"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "userscope", "false"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "attributesupport", "false"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "escapeattributes", "false"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "declaration_type", "com.apple.configuration.safari.bookmarks2"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "activation_predicate", "1234"),
				),
			},
			//Delete testing automatically occurs in TestCase
		},
	})
}
