package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccScriptJobResource(t *testing.T) {
	testAccPreCheck(t)
	_ = testAccRequireEnv(t, "SIMPLEMDM_RUN_SCRIPT_JOB_TESTS")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `		
		resource "simplemdm_scriptjob" "test_job" {
			script_id              = 5727
			device_ids             = []
			group_ids              = [140188]
			assignment_group_ids   = []
		}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check the script job attributes
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "id"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "script_id", "5727"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "group_ids.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "group_ids.0", "140188"),
				),
			},
			// // ImportState testing
			// {
			// 	ResourceName:      "simplemdm_scriptjob.test_job",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	// The filesha and  scriptfile attributes does not exist in SimpleMDM
			// 	// API, therefore there is no value for it during import.
			// },
			// Update and Read testing
			{
				Config: providerConfig + `
		resource "simplemdm_scriptjob" "test_job" {
			script_id              = 5727
			device_ids             = [1905524]
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
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "script_id", "5727"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
