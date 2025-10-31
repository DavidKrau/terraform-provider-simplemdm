package provider

import (
	"context"
	"testing"

	simplemdm "github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCheckAppDestroy(s *terraform.State) error {
	return testAccCheckResourceDestroyed("simplemdm_app", func(client *simplemdm.Client, id string) error {
		_, err := fetchApp(context.Background(), client, id)
		return err
	})(s)
}

func TestAccAppResourceWithAppStoreIdAttr(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAppDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_app" "testapp" {
					app_store_id = "284882215"
					deploy_to    = "none"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_app.testapp", "app_store_id", "284882215"),
					resource.TestCheckResourceAttr("simplemdm_app.testapp", "deploy_to", "none"),
					resource.TestCheckResourceAttrSet("simplemdm_app.testapp", "name"),
					resource.TestCheckResourceAttrSet("simplemdm_app.testapp", "app_type"),
					resource.TestCheckResourceAttrSet("simplemdm_app.testapp", "platform_support"),
					resource.TestCheckResourceAttrSet("simplemdm_app.testapp", "processing_status"),

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
			//Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAppResourceWithBundleIdAttr(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAppDestroy,
		Steps: []resource.TestStep{
			// Update without deploy_to in tf code but use bundle_id instead of app_store_id
			{
				Config: providerConfig + `
				resource "simplemdm_app" "testapp" {
					bundle_id     = "com.microsoft.Office.Excel"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_app.testapp", "bundle_id", "com.microsoft.Office.Excel"),
					resource.TestCheckResourceAttr("simplemdm_app.testapp", "deploy_to", "none"),
					resource.TestCheckResourceAttrSet("simplemdm_app.testapp", "app_type"),
					resource.TestCheckResourceAttrSet("simplemdm_app.testapp", "platform_support"),
					resource.TestCheckResourceAttrSet("simplemdm_app.testapp", "processing_status"),

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
			//Delete testing automatically occurs in TestCase
		},
	})
}
