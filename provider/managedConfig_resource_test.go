package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCheckManagedConfigDestroy(s *terraform.State) error {
	client, err := getTestClient()
	if err != nil {
		return fmt.Errorf("failed to create test client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "simplemdm_managed_config" {
			continue
		}

		appID, configID, err := parseManagedConfigID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to parse managed config ID %s: %w", rs.Primary.ID, err)
		}

		_, err = fetchManagedConfig(context.Background(), client, appID, configID)
		if err == nil {
			return fmt.Errorf("managed config %s still exists after destroy", rs.Primary.ID)
		}

		if err != errManagedConfigNotFound && !isNotFoundError(err) {
			return fmt.Errorf("unexpected error checking managed config %s: %w", rs.Primary.ID, err)
		}
	}

	return nil
}

func TestAccManagedConfigResource_basic(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckManagedConfigDestroy,
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
