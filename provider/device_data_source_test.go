package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDeviceDataSource requires an actual enrolled device because devices
// cannot be created via the SimpleMDM API - they must be enrolled through
// the normal device enrollment process.
//
// To run this test, set SIMPLEMDM_DEVICE_ID to an enrolled device's ID.
func TestAccDeviceDataSource(t *testing.T) {
	testAccPreCheck(t)

	deviceID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`data "simplemdm_device" "test" {id ="%s"}`, deviceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify returned values
					resource.TestCheckResourceAttr("data.simplemdm_device.test", "id", deviceID),
					resource.TestCheckResourceAttrSet("data.simplemdm_device.test", "name"),
				),
			},
		},
	})
}
