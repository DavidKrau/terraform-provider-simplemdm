package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccManagedConfigResource_basic(t *testing.T) {
	testAccPreCheck(t)

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
                                        key        = "serverURL"
                                        value      = "https://example.com"
                                        value_type = "string"
                                }
                                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_managed_config.config", "key", "serverURL"),
					resource.TestCheckResourceAttr("simplemdm_managed_config.config", "value", "https://example.com"),
					resource.TestCheckResourceAttr("simplemdm_managed_config.config", "value_type", "string"),
					resource.TestCheckResourceAttrSet("simplemdm_managed_config.config", "id"),
					resource.TestCheckResourceAttrPair("simplemdm_managed_config.config", "app_id", "simplemdm_app.testapp", "id"),
				),
			},
			{
				ResourceName:      "simplemdm_managed_config.config",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerConfig + `
                                resource "simplemdm_app" "testapp" {
                                        app_store_id = "586447913"
                                }

                                resource "simplemdm_managed_config" "config" {
                                        app_id     = simplemdm_app.testapp.id
                                        key        = "serverURL"
                                        value      = "https://terraform.example.com"
                                        value_type = "string"
                                }
                                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_managed_config.config", "value", "https://terraform.example.com"),
					resource.TestCheckResourceAttrSet("simplemdm_managed_config.config", "id"),
				),
			},
		},
	})
}
