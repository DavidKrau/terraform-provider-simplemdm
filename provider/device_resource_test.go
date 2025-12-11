package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
		resource "simplemdm_device" "test" {
			name= "Created test device"
			devicename  = "Created test device"
			devicegroups = [1978695,2170591]
  			profiles = [172801]
		}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_device.test", "name", "Created test device"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "devicename", "Created test device"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "devicegroups.#", "2"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "devicegroups.1", "2170591"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "devicegroups.0", "1978695"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "profiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "profiles.0", "172801"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_device.test", "id"),
					resource.TestCheckResourceAttrSet("simplemdm_device.test", "enrollmenturl"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "simplemdm_device.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The profiles and  customprofiles attributes does not exist in SimpleMDM
				// API, therefore there is no value for it during import.
				ImportStateVerifyIgnore: []string{"profiles"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_device" "test" {
					name= "Created test device changed"
					devicename  = "Created test device changed"
					devicegroups = [1538158]
					//profiles = [176844]
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_device.test", "name", "Created test device changed"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "devicename", "Created test device changed"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "devicegroups.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "devicegroups.0", "1538158"),
					//resource.TestCheckResourceAttr("simplemdm_device.test", "profiles.#", "1"),
					//resource.TestCheckResourceAttr("simplemdm_device.test", "profiles.0", "176844"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_device.test", "id"),
					resource.TestCheckResourceAttrSet("simplemdm_device.test", "enrollmenturl"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
