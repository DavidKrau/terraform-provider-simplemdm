package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAppDataSource(t *testing.T) {
	testAccPreCheck(t)

	appID := testAccRequireEnv(t, "SIMPLEMDM_APP_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`data "simplemdm_app" "test" {id ="%s"}`, appID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.simplemdm_app.test", "id", appID),
					resource.TestCheckResourceAttrSet("data.simplemdm_app.test", "name"),
					resource.TestCheckResourceAttrSet("data.simplemdm_app.test", "app_store_id"),
					resource.TestCheckResourceAttrSet("data.simplemdm_app.test", "bundle_id"),
					resource.TestCheckResourceAttrSet("data.simplemdm_app.test", "deploy_to"),
					resource.TestCheckResourceAttrSet("data.simplemdm_app.test", "status"),
					resource.TestCheckResourceAttrSet("data.simplemdm_app.test", "app_type"),
					resource.TestCheckResourceAttrSet("data.simplemdm_app.test", "version"),
					resource.TestCheckResourceAttrSet("data.simplemdm_app.test", "platform_support"),
					resource.TestCheckResourceAttrSet("data.simplemdm_app.test", "processing_status"),
					resource.TestCheckResourceAttr("data.simplemdm_app.test", "installation_channels.#", "1"),
				),
			},
		},
	})
}
