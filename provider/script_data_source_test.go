package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccScriptDataSource(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "simplemdm_script" "test" {id ="5727"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name of the app
					resource.TestCheckResourceAttr("data.simplemdm_script.test", "name", "Test script"),
				),
			},
		},
	})
}
