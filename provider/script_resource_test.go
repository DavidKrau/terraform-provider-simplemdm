package provider

import (
	"testing"

	simplemdm "github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCheckScriptDestroy(s *terraform.State) error {
	return testAccCheckResourceDestroyed("simplemdm_script", func(client *simplemdm.Client, id string) error {
		_, err := client.ScriptGet(id)
		return err
	})(s)
}

func TestAccScriptResource(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckScriptDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
		resource "simplemdm_script" "test" {
			name = "This is test script"
			content = file("./testfiles/testscript.sh")
			variable_support = true
			 }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_script.test", "name", "This is test script"),
					resource.TestCheckResourceAttr("simplemdm_script.test", "content", "#!/bin/bash\necho \"Hello!\""),
					resource.TestCheckResourceAttr("simplemdm_script.test", "variable_support", "true"),

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
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_script" "test" {
					name = "This is test script 2"
					content = file("./testfiles/testscript2.sh")
					variable_support = false
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_script.test", "name", "This is test script 2"),
					resource.TestCheckResourceAttr("simplemdm_script.test", "content", "#!/bin/bash\necho \"Hello again!\""),
					resource.TestCheckResourceAttr("simplemdm_script.test", "variable_support", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_script.test", "id"),
					resource.TestCheckResourceAttrSet("simplemdm_script.test", "created_at"),
					resource.TestCheckResourceAttrSet("simplemdm_script.test", "updated_at"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
