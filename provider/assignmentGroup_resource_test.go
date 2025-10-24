package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAssignmentGroupResource(t *testing.T) {
	testAccPreCheck(t)

	appID := testAccRequireEnv(t, "SIMPLEMDM_ASSIGNMENT_GROUP_APP_ID")
	groupID := testAccRequireEnv(t, "SIMPLEMDM_ASSIGNMENT_GROUP_GROUP_ID")
	profileID := testAccRequireEnv(t, "SIMPLEMDM_ASSIGNMENT_GROUP_PROFILE_ID")
	deviceID := testAccRequireEnv(t, "SIMPLEMDM_ASSIGNMENT_GROUP_DEVICE_ID")
	updatedAppID := testAccRequireEnv(t, "SIMPLEMDM_ASSIGNMENT_GROUP_UPDATED_APP_ID")
	updatedDeviceID := testAccRequireEnv(t, "SIMPLEMDM_ASSIGNMENT_GROUP_UPDATED_DEVICE_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(providerConfig+`
                                resource "simplemdm_assignmentgroup" "testgroup2" {
                                        name= "This assignment group"
                                        auto_deploy = false
                                        group_type   = "standard"
                                        install_type = "managed"
                                        priority     = 3
                                        app_track_location = false
                                        apps= [%s]
                                        groups = [%s]
                                        profiles = [%s]
                                        devices = [%s]
                                        devices_remove_others = true
                                        profiles_sync = false
                                        apps_push = false
                                        apps_update = false
                                  }
`, appID, groupID, profileID, deviceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "name", "This assignment group"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "group_type", "standard"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "priority", "3"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "app_track_location", "false"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices_remove_others", "true"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "profiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "profiles.0", profileID),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices.0", deviceID),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "apps.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "apps.0", appID),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "groups.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "groups.0", groupID),
					// Verify dynamic values have any value set in the state.
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
			//Update and Read testing
			{
				Config: fmt.Sprintf(providerConfig+`
                                resource "simplemdm_assignmentgroup" "testgroup2" {
                                        name= "renamed assignemnt group"
                                        auto_deploy = false
                                        group_type   = "munki"
                                        install_type = "managed"
                                        priority     = 7
                                        app_track_location = true
                                        apps= [%s]
                                        devices = [%s]
                                        devices_remove_others = false
                                        profiles_sync = false
                                        apps_push = false
                                        apps_update = false
                                  }
                        `, updatedAppID, updatedDeviceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "name", "renamed assignemnt group"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "group_type", "munki"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "install_type", "managed"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "priority", "7"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "app_track_location", "true"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices_remove_others", "false"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices.0", updatedDeviceID),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "apps.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "apps.0", updatedAppID),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_assignmentgroup.testgroup2", "id"),
				),
			},
			//Delete testing automatically occurs in TestCase
		},
	})
}
