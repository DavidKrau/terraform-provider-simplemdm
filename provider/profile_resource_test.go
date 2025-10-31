package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProfileResource(t *testing.T) {
	testAccPreCheck(t)

	profileID := os.Getenv("SIMPLEMDM_PROFILE_ID")
	if profileID == "" {
		t.Skip("SIMPLEMDM_PROFILE_ID environment variable not set - skipping test that requires existing profile fixture")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "simplemdm_profile" "test" {
  id = "` + profileID + `"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_profile.test", "id", profileID),
					resource.TestCheckResourceAttrSet("simplemdm_profile.test", "name"),
				),
			},
		},
	})
}
