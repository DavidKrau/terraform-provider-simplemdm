package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccScriptJobResource(t *testing.T) {
	testAccPreCheck(t)

	scriptID := testAccRequireEnv(t, "SIMPLEMDM_SCRIPT_JOB_SCRIPT_ID")
	groupID := testAccRequireEnv(t, "SIMPLEMDM_SCRIPT_JOB_GROUP_ID")
	deviceID := testAccRequireEnv(t, "SIMPLEMDM_SCRIPT_JOB_DEVICE_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(providerConfig+`
                resource "simplemdm_scriptjob" "test_job" {
                        script_id              = %s
                        device_ids             = []
                        group_ids              = [%s]
                        assignment_group_ids   = []
                }
                                `, scriptID, groupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check the script job attributes
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "id"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "script_id", scriptID),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "group_ids.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "group_ids.0", groupID),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "job_identifier"),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "status"),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "pending_count"),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "created_at"),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "variable_support"),
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
				Config: fmt.Sprintf(providerConfig+`
                resource "simplemdm_scriptjob" "test_job" {
                        script_id              = %s
                        device_ids             = [%s]
                        group_ids              = []
                        assignment_group_ids   = []
                        custom_attribute       = "updated_attribute"
                        custom_attribute_regex = "\\r"
                }
                                `, scriptID, deviceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check the updated script job attributes
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "id"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "custom_attribute", "updated_attribute"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "custom_attribute_regex", "\\r"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "script_id", scriptID),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "status"),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "job_identifier"),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "success_count"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
