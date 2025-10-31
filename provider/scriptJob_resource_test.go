package provider

import (
	"context"
	"testing"

	simplemdm "github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCheckScriptJobDestroy(s *terraform.State) error {
	return testAccCheckResourceDestroyed("simplemdm_scriptjob", func(client *simplemdm.Client, id string) error {
		_, err := fetchScriptJobDetails(context.Background(), client, id)
		return err
	})(s)
}

func TestAccScriptJobResource(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckScriptJobDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
					# Create prerequisite resources
					resource "simplemdm_script" "test_script" {
						name            = "Test Script Job Script"
						scriptfile      = file("./testfiles/testscript.sh")
						variablesupport = true
					}

					resource "simplemdm_devicegroup" "test_group" {
						name = "Test Script Job Group"
					}

					# Create script job using dynamic references
					resource "simplemdm_scriptjob" "test_job" {
						script_id            = simplemdm_script.test_script.id
						device_ids           = []
						group_ids            = [simplemdm_devicegroup.test_group.id]
						assignment_group_ids = []

						depends_on = [
							simplemdm_script.test_script,
							simplemdm_devicegroup.test_group
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check the script job attributes
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "id"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "group_ids.#", "1"),
					// Verify dynamic relationships
					resource.TestCheckResourceAttrPair(
						"simplemdm_scriptjob.test_job", "script_id",
						"simplemdm_script.test_script", "id",
					),
					resource.TestCheckResourceAttrPair(
						"simplemdm_scriptjob.test_job", "group_ids.0",
						"simplemdm_devicegroup.test_group", "id",
					),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "job_identifier"),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "status"),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "pending_count"),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "created_at"),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "variable_support"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "simplemdm_scriptjob.test_job",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					# Create prerequisite resources for update
					resource "simplemdm_script" "test_script" {
						name            = "Test Script Job Script"
						scriptfile      = file("./testfiles/testscript.sh")
						variablesupport = true
					}

					resource "simplemdm_device" "test_device" {
						name       = "Test Script Job Device"
						devicename = "Test Script Job Device"
					}

					# Update script job with different targets
					resource "simplemdm_scriptjob" "test_job" {
						script_id              = simplemdm_script.test_script.id
						device_ids             = [simplemdm_device.test_device.id]
						group_ids              = []
						assignment_group_ids   = []
						custom_attribute       = "updated_attribute"
						custom_attribute_regex = "\\r"

						depends_on = [
							simplemdm_script.test_script,
							simplemdm_device.test_device
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check the updated script job attributes
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "id"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "custom_attribute", "updated_attribute"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "custom_attribute_regex", "\\r"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "device_ids.#", "1"),
					// Verify dynamic relationships
					resource.TestCheckResourceAttrPair(
						"simplemdm_scriptjob.test_job", "script_id",
						"simplemdm_script.test_script", "id",
					),
					resource.TestCheckResourceAttrPair(
						"simplemdm_scriptjob.test_job", "device_ids.0",
						"simplemdm_device.test_device", "id",
					),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "status"),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "job_identifier"),
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "success_count"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
