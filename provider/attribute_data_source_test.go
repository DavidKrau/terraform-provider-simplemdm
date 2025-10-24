package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAttributeDataSource(t *testing.T) {
	testAccPreCheck(t)

	attributeName := testAccRequireEnv(t, "SIMPLEMDM_ATTRIBUTE_NAME")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`data "simplemdm_attribute" "test" {name ="%s"}`, attributeName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify returned values
					resource.TestCheckResourceAttr("data.simplemdm_attribute.test", "name", attributeName),
					resource.TestCheckResourceAttrSet("data.simplemdm_attribute.test", "default_value"),
				),
			},
		},
	})
}
