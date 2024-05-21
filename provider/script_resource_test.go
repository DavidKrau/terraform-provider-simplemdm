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
			scriptfile = "./testfiles/testscript.sh"
			filesha =    "${filesha256("./testfiles/testscript.sh")}"
			variablesupport = true
		  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_script.test", "name", "This is test script"),
					resource.TestCheckResourceAttr("simplemdm_script.test", "scriptfile", "./testfiles/testscript.sh"),
					resource.TestCheckResourceAttr("simplemdm_script.test", "filesha", "b6fcca636ed070775f43ef09bcf04458086aeeb364fc427d7b203821fa7e4727"),
					resource.TestCheckResourceAttr("simplemdm_script.test", "content", "#!/bin/bash\necho \"Hello!\""),
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
				ImportStateVerifyIgnore: []string{"filesha", "scriptfile"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_script" "test" {
					name= "This is test script 2"
					scriptfile = "./testfiles/testscript2.sh"
					filesha =    "${filesha256("./testfiles/testscript2.sh")}"
					variablesupport = false				
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_script.test", "name", "This is test script 2"),
					resource.TestCheckResourceAttr("simplemdm_script.test", "scriptfile", "./testfiles/testscript2.sh"),
					resource.TestCheckResourceAttr("simplemdm_script.test", "filesha", "a2ef0628495f6f5a40e0faf0792afe52670461fce6f87124e96cb5f68066214e"),
					resource.TestCheckResourceAttr("simplemdm_script.test", "variablesupport", "false"),
					resource.TestCheckResourceAttr("simplemdm_script.test", "content", "#!/bin/bash\necho \"Hello again!\""),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_script.test", "id"),
					resource.TestCheckResourceAttrSet("simplemdm_script.test", "created_at"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
