package provider

import (
	"context"
	"testing"

	simplemdm "github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCheckAssignmentGroupDestroy(s *terraform.State) error {
	return testAccCheckResourceDestroyed("simplemdm_assignmentgroup", func(client *simplemdm.Client, id string) error {
		_, err := fetchAssignmentGroup(context.Background(), client, id)
		return err
	})(s)
}

func TestAccAssignmentGroupResource(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAssignmentGroupDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
					# Create prerequisite resources
					resource "simplemdm_app" "test_app" {
						app_store_id = "586447913"
					}

					resource "simplemdm_devicegroup" "test_group" {
						name = "Test Assignment Group Device Group"
					}

					resource "simplemdm_customprofile" "test_profile" {
						name         = "Test Assignment Profile"
						mobileconfig = file("./testfiles/testprofile.mobileconfig")
						userscope    = false
					}

					resource "simplemdm_device" "test_device" {
						name       = "Test Assignment Device"
						devicename = "Test Assignment Device"
					}

					# Create assignment group using dynamic references
					resource "simplemdm_assignmentgroup" "testgroup2" {
						name                  = "This assignment group"
						auto_deploy           = false
						group_type            = "standard"
						install_type          = "managed"
						priority              = 3
						app_track_location    = false
						apps                  = [simplemdm_app.test_app.id]
						groups                = [simplemdm_devicegroup.test_group.id]
						profiles              = [simplemdm_customprofile.test_profile.id]
						devices               = [simplemdm_device.test_device.id]
						devices_remove_others = true
						profiles_sync         = false
						apps_push             = false
						apps_update           = false

						depends_on = [
							simplemdm_app.test_app,
							simplemdm_devicegroup.test_group,
							simplemdm_customprofile.test_profile,
							simplemdm_device.test_device
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "name", "This assignment group"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "group_type", "standard"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "priority", "3"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "app_track_location", "false"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices_remove_others", "true"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "profiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "apps.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "groups.#", "1"),
					// Verify dynamic relationships
					resource.TestCheckResourceAttrPair(
						"simplemdm_assignmentgroup.testgroup2", "apps.0",
						"simplemdm_app.test_app", "id",
					),
					resource.TestCheckResourceAttrPair(
						"simplemdm_assignmentgroup.testgroup2", "groups.0",
						"simplemdm_devicegroup.test_group", "id",
					),
					resource.TestCheckResourceAttrPair(
						"simplemdm_assignmentgroup.testgroup2", "profiles.0",
						"simplemdm_customprofile.test_profile", "id",
					),
					resource.TestCheckResourceAttrPair(
						"simplemdm_assignmentgroup.testgroup2", "devices.0",
						"simplemdm_device.test_device", "id",
					),
					// Verify dynamic values have any value set in the state
					resource.TestCheckResourceAttrSet("simplemdm_assignmentgroup.testgroup2", "id"),
					resource.TestCheckResourceAttrSet("simplemdm_assignmentgroup.testgroup2", "created_at"),
					resource.TestCheckResourceAttrSet("simplemdm_assignmentgroup.testgroup2", "updated_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "simplemdm_assignmentgroup.testgroup2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"apps_update", "apps_push", "auto_deploy", "profiles_sync", "install_type", "profiles", "created_at", "updated_at", "device_count", "group_count"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					# Create updated prerequisite resources
					resource "simplemdm_app" "updated_app" {
						app_store_id = "1477376905"
					}

					resource "simplemdm_device" "updated_device" {
						name       = "Updated Assignment Device"
						devicename = "Updated Assignment Device"
					}

					# Update assignment group with new references
					resource "simplemdm_assignmentgroup" "testgroup2" {
						name                  = "renamed assignemnt group"
						auto_deploy           = false
						group_type            = "munki"
						install_type          = "managed"
						priority              = 7
						app_track_location    = true
						apps                  = [simplemdm_app.updated_app.id]
						devices               = [simplemdm_device.updated_device.id]
						devices_remove_others = false
						profiles_sync         = false
						apps_push             = false
						apps_update           = false

						depends_on = [
							simplemdm_app.updated_app,
							simplemdm_device.updated_device
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "name", "renamed assignemnt group"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "group_type", "munki"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "install_type", "managed"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "priority", "7"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "app_track_location", "true"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices_remove_others", "false"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "apps.#", "1"),
					// Verify dynamic relationships
					resource.TestCheckResourceAttrPair(
						"simplemdm_assignmentgroup.testgroup2", "apps.0",
						"simplemdm_app.updated_app", "id",
					),
					resource.TestCheckResourceAttrPair(
						"simplemdm_assignmentgroup.testgroup2", "devices.0",
						"simplemdm_device.updated_device", "id",
					),
					// Verify dynamic values have any value set in the state
					resource.TestCheckResourceAttrSet("simplemdm_assignmentgroup.testgroup2", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
