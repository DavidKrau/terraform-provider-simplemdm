package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProfileResource(t *testing.T) {
	testAccPreCheck(t)
	_ = testAccRequireEnv(t, "SIMPLEMDM_RUN_PROFILE_RESOURCE_TESTS")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "simplemdm_profile" "test" {
  id = "172801"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_profile.test", "id", "172801"),
					resource.TestCheckResourceAttr("simplemdm_profile.test", "name", "Restriction test profiles"),
				),
			},
		},
	})
}
