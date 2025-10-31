package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccManagedConfigsDataSource_basic(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                                resource "simplemdm_app" "testapp" {
                                        app_store_id = "586447913"
                                }

                                resource "simplemdm_managed_config" "config1" {
                                        app_id     = simplemdm_app.testapp.id
                                        key        = "environment"
                                        value      = "production"
                                        value_type = "string"
                                }

                                resource "simplemdm_managed_config" "config2" {
                                        app_id     = simplemdm_app.testapp.id
                                        key        = "debug_mode"
                                        value      = "false"
                                        value_type = "boolean"
                                }

                                data "simplemdm_managed_configs" "all" {
                                        app_id = simplemdm_app.testapp.id
                                        depends_on = [
                                                simplemdm_managed_config.config1,
                                                simplemdm_managed_config.config2,
                                        ]
                                }
                                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.simplemdm_managed_configs.all", "app_id", "simplemdm_app.testapp", "id"),
					resource.TestCheckResourceAttr("data.simplemdm_managed_configs.all", "managed_configs.#", "2"),
				),
			},
		},
	})
}
