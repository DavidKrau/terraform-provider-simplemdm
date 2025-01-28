package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAppResourceWithAppStoreIdAttr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_app" "testapp" {
					app_store_id = "1477376905"
					deploy_to    = "outdated"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_app.testapp", "app_store_id", "1477376905"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_app.testapp", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "simplemdm_app.testapp",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deploy_to"},
			},
			//Update and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_app" "testapp" {
					app_store_id = "586447913"
					deploy_to		 = "all"
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_app.testapp", "app_store_id", "586447913"),
					resource.TestCheckResourceAttr("simplemdm_app.testapp", "deploy_to", "all"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_app.testapp", "id"),
				),
			},
			//Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAppResourceWithBundleIdAttr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Update without deploy_to in tf code but use bundle_id insted of app_store_id
			{
				Config: providerConfig + `
				resource "simplemdm_app" "testapp" {
					bundle_id     = "com.microsoft.Office.Excel"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_app.testapp", "bundle_id", "com.microsoft.Office.Excel"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_app.testapp", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "simplemdm_app.testapp",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deploy_to"},
			},
			//Update and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_app" "testapp" {
					app_store_id = "586447913"
					deploy_to		 = "all"
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_app.testapp", "app_store_id", "586447913"),
					resource.TestCheckResourceAttr("simplemdm_app.testapp", "deploy_to", "all"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_app.testapp", "id"),
				),
			},
			//Delete testing automatically occurs in TestCase
		},
	})
}
