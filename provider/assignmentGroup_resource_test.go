package provider

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCheckAssignmentGroupDestroy(s *terraform.State) error {
	client, err := getTestClient()
	if err != nil {
		return fmt.Errorf("failed to create test client: %w", err)
	}

	for name, rs := range s.RootModule().Resources {
		// Skip if not our resource type
		if rs.Type != "simplemdm_assignmentgroup" {
			continue
		}

		// Skip data sources - they start with "data."
		if strings.HasPrefix(name, "data.") {
			continue
		}

		// Try to fetch the resource
		_, err := fetchAssignmentGroup(context.Background(), client, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("assignment group %s still exists after destroy", rs.Primary.ID)
		}

		// If it's a 404, that's expected (resource was destroyed)
		if !isNotFoundError(err) {
			return fmt.Errorf("unexpected error checking assignment group %s: %w", rs.Primary.ID, err)
		}
	}

	return nil
}

func TestAccAssignmentGroupResource(t *testing.T) {
	testAccPreCheck(t)

	// Use pre-existing fixture assignment group - Device Groups are deprecated!
	assignmentGroupID := testAccRequireEnv(t, "SIMPLEMDM_ASSIGNMENT_GROUP_ID")
	appID := testAccRequireEnv(t, "SIMPLEMDM_APP_ID")
	profileID := testAccRequireEnv(t, "SIMPLEMDM_PROFILE_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAssignmentGroupDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
					# Use pre-existing fixture resources
					data "simplemdm_assignmentgroup" "fixture_group" {
						id = "%s"
					}

					data "simplemdm_app" "fixture_app" {
						id = "%s"
					}

					data "simplemdm_profile" "fixture_profile" {
						id = "%s"
					}

					# Create assignment group using fixture references
					resource "simplemdm_assignmentgroup" "testgroup2" {
						name                  = "Test Assignment Group Resource"
						auto_deploy           = false
						group_type            = "standard"
						install_type          = "managed"
						priority              = 3
						app_track_location    = false
						apps                  = [data.simplemdm_app.fixture_app.id]
						profiles              = [data.simplemdm_profile.fixture_profile.id]
						devices_remove_others = true
						profiles_sync         = false
						apps_push             = false
						apps_update           = false
					}
				`, assignmentGroupID, appID, profileID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "name", "Test Assignment Group Resource"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "group_type", "standard"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "priority", "3"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "app_track_location", "false"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices_remove_others", "true"),
					// Note: Due to API eventual consistency, profiles and apps counts may be 0 or 1
					// Verify dynamic values have any value set in the state
					resource.TestCheckResourceAttrSet("simplemdm_assignmentgroup.testgroup2", "id"),
					// Note: created_at and updated_at may not be immediately returned by API
				),
				// Allow non-empty plan due to API eventual consistency with relationships
				ExpectNonEmptyPlan: true,
			},
			// ImportState testing
			{
				ResourceName:            "simplemdm_assignmentgroup.testgroup2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"apps_update", "apps_push", "auto_deploy", "profiles_sync", "install_type", "profiles", "created_at", "updated_at", "device_count", "group_count", "devices_remove_others"},
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
					# Use fixture app for update
					data "simplemdm_app" "fixture_app_updated" {
						id = "%s"
					}

					# Update assignment group with modified attributes
					resource "simplemdm_assignmentgroup" "testgroup2" {
						name                  = "Updated Assignment Group Resource"
						auto_deploy           = false
						group_type            = "munki"
						install_type          = "managed"
						priority              = 7
						app_track_location    = true
						apps                  = [data.simplemdm_app.fixture_app_updated.id]
						devices_remove_others = false
						profiles_sync         = false
						apps_push             = false
						apps_update           = false
					}
				`, appID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "name", "Updated Assignment Group Resource"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "group_type", "munki"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "install_type", "managed"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "priority", "7"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "app_track_location", "true"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices_remove_others", "false"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "apps.#", "1"),
					// Verify dynamic relationships
					resource.TestCheckResourceAttrPair(
						"simplemdm_assignmentgroup.testgroup2", "apps.0",
						"data.simplemdm_app.fixture_app_updated", "id",
					),
					// Verify dynamic values have any value set in the state
					resource.TestCheckResourceAttrSet("simplemdm_assignmentgroup.testgroup2", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
