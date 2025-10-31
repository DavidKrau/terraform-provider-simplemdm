package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceGroupDataSource(t *testing.T) {
	testAccPreCheck(t)

	// Device groups cannot be created via API - they must exist
	// This test requires SIMPLEMDM_DEVICE_GROUP_ID or will skip
	deviceGroupID := testAccGetEnv(t, "SIMPLEMDM_DEVICE_GROUP_ID")
	
	if deviceGroupID == "" {
		t.Skip("SIMPLEMDM_DEVICE_GROUP_ID not set - skipping test as device groups cannot be created via API")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read existing device group with data source
			{
				Config: providerConfig + `
					data "simplemdm_devicegroup" "test" {
						id = "` + deviceGroupID + `"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.simplemdm_devicegroup.test", "id", deviceGroupID),
					resource.TestCheckResourceAttrSet("data.simplemdm_devicegroup.test", "name"),
				),
			},
		},
	})
}
