package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccScriptJobResource(t *testing.T) {

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

		resource "simplemdm_device" "test" {
			name= "Created test device"
			devicename  = "Created test device"
			devicegroup = 152757
  			profiles = []
  			customprofiles = []
		}
		
		resource "simplemdm_scriptjob" "test_job" {
			script_id              = simplemdm_script.test.id
			device_ids             = []
			group_ids              = [152757]
			assignment_group_ids   = []
		}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check the script job attributes
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "id"),
					resource.TestCheckResourceAttrPair(
						"simplemdm_scriptjob.test_job", "script_id",
						"simplemdm_script.test", "id",
					),
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
		resource "simplemdm_scriptjob" "test_job" {
			script_id              = simplemdm_script.test.id
			device_ids             = [1903896]
			group_ids              = []
			assignment_group_ids   = []
			custom_attribute       = "updated_attribute"
			custom_attribute_regex = "\\r"
		}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check the updated script job attributes
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "id"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "custom_attribute", "updated_attribute"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "custom_attribute_regex", "\\r"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "script_id", "simplemdm_script.test.id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
