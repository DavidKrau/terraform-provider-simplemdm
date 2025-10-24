package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccManagedConfigDataSource_basic(t *testing.T) {
	testAccPreCheck(t)
	_ = testAccRequireEnv(t, "SIMPLEMDM_RUN_MANAGED_CONFIG_TESTS")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                                resource "simplemdm_app" "testapp" {
                                        app_store_id = "586447913"
                                }

                                resource "simplemdm_managed_config" "config" {
                                        app_id     = simplemdm_app.testapp.id
                                        key        = "environment"
                                        value      = "production"
                                        value_type = "string"
                                }

                                data "simplemdm_managed_config" "config" {
                                        id = simplemdm_managed_config.config.id
                                }
                                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.simplemdm_managed_config.config", "id", "simplemdm_managed_config.config", "id"),
					resource.TestCheckResourceAttrPair("data.simplemdm_managed_config.config", "app_id", "simplemdm_app.testapp", "id"),
					resource.TestCheckResourceAttr("data.simplemdm_managed_config.config", "key", "environment"),
					resource.TestCheckResourceAttr("data.simplemdm_managed_config.config", "value", "production"),
					resource.TestCheckResourceAttr("data.simplemdm_managed_config.config", "value_type", "string"),
				),
			},
		},
	})
}
