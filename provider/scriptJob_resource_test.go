package provider

import (
	"context"
	"fmt"
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

	// Use pre-existing fixture device group for script jobs
	// Note: SimpleMDM API does not support assignment_group_ids for script jobs despite what the docs may suggest
	deviceGroupID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_ID")
	scriptID := testAccRequireEnv(t, "SIMPLEMDM_SCRIPT_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckScriptJobDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
					# Use pre-existing fixture resources
					data "simplemdm_script" "fixture_script" {
						id = "%s"
					}

					data "simplemdm_devicegroup" "fixture_group" {
						id = "%s"
					}

					# Create script job using fixture device group
					resource "simplemdm_scriptjob" "test_job" {
						script_id  = data.simplemdm_script.fixture_script.id
						device_ids = []
						group_ids  = [data.simplemdm_devicegroup.fixture_group.id]
					}
				`, scriptID, deviceGroupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check the script job attributes
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "id"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "group_ids.#", "1"),
					// Verify dynamic relationships
					resource.TestCheckResourceAttrPair(
						"simplemdm_scriptjob.test_job", "script_id",
						"data.simplemdm_script.fixture_script", "id",
					),
					resource.TestCheckResourceAttrPair(
						"simplemdm_scriptjob.test_job", "group_ids.0",
						"data.simplemdm_devicegroup.fixture_group", "id",
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
				Config: providerConfig + fmt.Sprintf(`
					# Use fixture script for update
					data "simplemdm_script" "fixture_script_updated" {
						id = "%s"
					}

					data "simplemdm_devicegroup" "fixture_group_updated" {
						id = "%s"
					}

					# Update script job with custom attributes
					resource "simplemdm_scriptjob" "test_job" {
						script_id              = data.simplemdm_script.fixture_script_updated.id
						device_ids             = []
						group_ids              = [data.simplemdm_devicegroup.fixture_group_updated.id]
						custom_attribute       = "SomeAttribute"
						custom_attribute_regex = ".*"
					}
				`, scriptID, deviceGroupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check the updated script job attributes
					resource.TestCheckResourceAttrSet("simplemdm_scriptjob.test_job", "id"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "custom_attribute", "SomeAttribute"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "custom_attribute_regex", ".*"),
					resource.TestCheckResourceAttr("simplemdm_scriptjob.test_job", "group_ids.#", "1"),
					// Verify dynamic relationships
					resource.TestCheckResourceAttrPair(
						"simplemdm_scriptjob.test_job", "script_id",
						"data.simplemdm_script.fixture_script_updated", "id",
					),
					resource.TestCheckResourceAttrPair(
						"simplemdm_scriptjob.test_job", "group_ids.0",
						"data.simplemdm_devicegroup.fixture_group_updated", "id",
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
