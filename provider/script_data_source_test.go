package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccScriptDataSource(t *testing.T) {
	testAccPreCheck(t)

	scriptID := testAccRequireEnv(t, "SIMPLEMDM_SCRIPT_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`data "simplemdm_script" "test" {id ="%s"}`, scriptID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("data.simplemdm_script.test", "id", scriptID),
					resource.TestCheckResourceAttrSet("data.simplemdm_script.test", "name"),
					resource.TestCheckResourceAttrSet("data.simplemdm_script.test", "content"),
					resource.TestCheckResourceAttrSet("data.simplemdm_script.test", "variable_support"),
					resource.TestCheckResourceAttrSet("data.simplemdm_script.test", "created_at"),
					resource.TestCheckResourceAttrSet("data.simplemdm_script.test", "updated_at"),
				),
			},
		},
	})
}
