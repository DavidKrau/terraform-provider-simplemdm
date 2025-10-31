package provider

import (
	"context"
	"os"
	"testing"

	simplemdm "github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCheckEnrollmentDestroy(s *terraform.State) error {
	return testAccCheckResourceDestroyed("simplemdm_enrollment", func(client *simplemdm.Client, id string) error {
		_, err := fetchEnrollment(context.Background(), client, id)
		return err
	})(s)
}

func TestAccEnrollmentResource(t *testing.T) {
	testAccPreCheck(t)

	// Enrollments require actual device groups which cannot be created via API
	// Skip this test if no device group ID is available
	deviceGroupID := testAccGetEnv(t, "SIMPLEMDM_DEVICE_GROUP_ID")
	if deviceGroupID == "" {
		t.Skip("SIMPLEMDM_DEVICE_GROUP_ID not set - skipping test as enrollments require actual device groups which cannot be created via API")
	}

	// Get the required contact email (still needed as it's user-specific)
	invitationContact := testAccRequireEnv(t, "SIMPLEMDM_ENROLLMENT_CONTACT")

	steps := []resource.TestStep{
		{
			Config: providerConfig + `
				# Use existing device group (cannot be created via API)
				resource "simplemdm_enrollment" "test" {
					device_group_id    = "` + deviceGroupID + `"
					user_enrollment    = false
					welcome_screen     = true
					authentication     = false
					invitation_contact = "` + invitationContact + `"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("simplemdm_enrollment.test", "user_enrollment", "false"),
				resource.TestCheckResourceAttr("simplemdm_enrollment.test", "device_group_id", deviceGroupID),
				resource.TestCheckResourceAttr("simplemdm_enrollment.test", "welcome_screen", "true"),
				resource.TestCheckResourceAttr("simplemdm_enrollment.test", "authentication", "false"),
				resource.TestCheckResourceAttr("simplemdm_enrollment.test", "invitation_contact", invitationContact),
				resource.TestCheckResourceAttrSet("simplemdm_enrollment.test", "id"),
			),
		},
		{
			ResourceName:      "simplemdm_enrollment.test",
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"invitation_contact",
			},
		},
	}

	// Optional: test updating the invitation contact if provided
	if updatedContact := os.Getenv("SIMPLEMDM_ENROLLMENT_CONTACT_UPDATE"); updatedContact != "" {
		steps = append(steps, resource.TestStep{
			Config: providerConfig + `
				# Update enrollment with new contact
				resource "simplemdm_enrollment" "test" {
					device_group_id    = "` + deviceGroupID + `"
					user_enrollment    = false
					welcome_screen     = true
					authentication     = false
					invitation_contact = "` + updatedContact + `"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("simplemdm_enrollment.test", "device_group_id", deviceGroupID),
				resource.TestCheckResourceAttr("simplemdm_enrollment.test", "invitation_contact", updatedContact),
			),
		})
	}

	// Test removing the invitation contact
	steps = append(steps, resource.TestStep{
		Config: providerConfig + `
			# Remove invitation contact
			resource "simplemdm_enrollment" "test" {
				device_group_id = "` + deviceGroupID + `"
				user_enrollment = false
				welcome_screen  = true
				authentication  = false
			}
		`,
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckNoResourceAttr("simplemdm_enrollment.test", "invitation_contact"),
		),
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEnrollmentDestroy,
		Steps:                    steps,
	})
}
