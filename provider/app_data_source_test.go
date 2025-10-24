package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAppDataSource(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "simplemdm_app" "test" {id ="577575"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name of the app
					resource.TestCheckResourceAttr("data.simplemdm_app.test", "name", "SimpleMDM"),
					resource.TestCheckResourceAttrSet("data.simplemdm_app.test", "app_store_id"),
					resource.TestCheckResourceAttrSet("data.simplemdm_app.test", "bundle_id"),
					resource.TestCheckResourceAttrSet("data.simplemdm_app.test", "deploy_to"),
					resource.TestCheckResourceAttrSet("data.simplemdm_app.test", "status"),
				),
			},
		},
	})
}
