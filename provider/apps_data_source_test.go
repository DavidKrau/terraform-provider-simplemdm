package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAppsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.simplemdm_apps.test", "apps.#"),
				),
			},
		},
	})
}

func TestAccAppsDataSourceWithShared(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppsDataSourceConfigWithShared,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.simplemdm_apps.test", "apps.#"),
				),
			},
		},
	})
}

const testAccAppsDataSourceConfig = `
data "simplemdm_apps" "test" {}
`

const testAccAppsDataSourceConfigWithShared = `
data "simplemdm_apps" "test" {
  include_shared = true
}
`
