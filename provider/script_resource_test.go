package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccScriptResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
		resource "simplemdm_script" "test" {
			name= "This is test script"
			scriptfile = file("./testfiles/testscript.sh")
			variablesupport = true
		  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_script.test", "name", "This is test script"),
					resource.TestCheckResourceAttr("simplemdm_script.test", "scriptfile", "#!/bin/bash\necho \"Hello!\""),
					resource.TestCheckResourceAttr("simplemdm_script.test", "variablesupport", "true"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_script.test", "id"),
					resource.TestCheckResourceAttrSet("simplemdm_script.test", "created_at"),
					resource.TestCheckResourceAttrSet("simplemdm_script.test", "updated_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "simplemdm_script.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The filesha and  scriptfile attributes does not exist in SimpleMDM
				// API, therefore there is no value for it during import.
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_script" "test" {
					name= "This is test script 2"
					scriptfile = file("./testfiles/testscript2.sh")
					variablesupport = false				
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_script.test", "name", "This is test script 2"),
					resource.TestCheckResourceAttr("simplemdm_script.test", "scriptfile", "#!/bin/bash\necho \"Hello again!\""),
					resource.TestCheckResourceAttr("simplemdm_script.test", "variablesupport", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_script.test", "id"),
					resource.TestCheckResourceAttrSet("simplemdm_script.test", "created_at"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
