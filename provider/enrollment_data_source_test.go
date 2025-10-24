package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnrollmentDataSource(t *testing.T) {
	testAccPreCheck(t)

	enrollmentID := testAccRequireEnv(t, "SIMPLEMDM_ENROLLMENT_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`data "simplemdm_enrollment" "test" { id = "%s" }`, enrollmentID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.simplemdm_enrollment.test", "id", enrollmentID),
					resource.TestCheckResourceAttrSet("data.simplemdm_enrollment.test", "user_enrollment"),
					resource.TestCheckResourceAttrSet("data.simplemdm_enrollment.test", "welcome_screen"),
					resource.TestCheckResourceAttrSet("data.simplemdm_enrollment.test", "authentication"),
				),
			},
		},
	})
}
