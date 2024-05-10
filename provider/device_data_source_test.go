package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "simplemdm_device" "test" {id ="1579954"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify returned values
					resource.TestCheckResourceAttr("data.simplemdm_device.test", "name", "David Kraushuber MacBook Pro"),
				),
			},
		},
	})
}
