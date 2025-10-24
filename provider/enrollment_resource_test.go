package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnrollmentResource(t *testing.T) {
	testAccPreCheck(t)

	deviceGroupID := testAccRequireEnv(t, "SIMPLEMDM_ENROLLMENT_DEVICE_GROUP_ID")
	invitationContact := testAccRequireEnv(t, "SIMPLEMDM_ENROLLMENT_CONTACT")

	steps := []resource.TestStep{
		{
			Config: providerConfig + fmt.Sprintf(`
resource "simplemdm_enrollment" "test" {
  device_group_id     = "%s"
  user_enrollment     = false
  welcome_screen      = true
  authentication      = false
  invitation_contact  = "%s"
}
`, deviceGroupID, invitationContact),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("simplemdm_enrollment.test", "device_group_id", deviceGroupID),
				resource.TestCheckResourceAttr("simplemdm_enrollment.test", "user_enrollment", "false"),
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

	if updatedContact := os.Getenv("SIMPLEMDM_ENROLLMENT_CONTACT_UPDATE"); updatedContact != "" {
		steps = append(steps, resource.TestStep{
			Config: providerConfig + fmt.Sprintf(`
resource "simplemdm_enrollment" "test" {
  device_group_id     = "%s"
  user_enrollment     = false
  welcome_screen      = true
  authentication      = false
  invitation_contact  = "%s"
}
`, deviceGroupID, updatedContact),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("simplemdm_enrollment.test", "device_group_id", deviceGroupID),
				resource.TestCheckResourceAttr("simplemdm_enrollment.test", "invitation_contact", updatedContact),
			),
		})
	}

	steps = append(steps, resource.TestStep{
		Config: providerConfig + fmt.Sprintf(`
resource "simplemdm_enrollment" "test" {
  device_group_id     = "%s"
  user_enrollment     = false
  welcome_screen      = true
  authentication      = false
}
`, deviceGroupID),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckNoResourceAttr("simplemdm_enrollment.test", "invitation_contact"),
		),
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}
