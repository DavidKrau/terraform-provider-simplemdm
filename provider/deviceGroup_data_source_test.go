package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceGroupDataSource(t *testing.T) {
	testAccPreCheck(t)

	groupID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(`data "simplemdm_devicegroup" "test" {id ="%s"}`, groupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify returned values
					resource.TestCheckResourceAttr("data.simplemdm_devicegroup.test", "id", groupID),
					resource.TestCheckResourceAttrSet("data.simplemdm_devicegroup.test", "name"),
				),
			},
		},
	})
}
