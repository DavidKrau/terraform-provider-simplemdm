package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceResource(t *testing.T) {
	testAccPreCheck(t)
	_ = testAccRequireEnv(t, "SIMPLEMDM_RUN_DEVICE_RESOURCE_TESTS")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
		resource "simplemdm_device" "test" {
			name= "Created test device"
			devicename  = "Created test device"
			devicegroup = 140188
  			profiles = [172801]
  			customprofiles = [172804]
		}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_device.test", "name", "Created test device"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "devicename", "Created test device"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "devicegroup", "140188"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "profiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "profiles.0", "172801"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "customprofiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "customprofiles.0", "172804"),
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
				ImportStateVerifyIgnore: []string{"profiles", "customprofiles"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_device" "test" {
					name= "Created test device changed"
					devicename  = "Created test device changed"
					devicegroup = 140189
					profiles = [172802]
					customprofiles = [172805]
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_device.test", "name", "Created test device changed"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "devicename", "Created test device changed"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "devicegroup", "140189"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "profiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "profiles.0", "172802"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "customprofiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "customprofiles.0", "172805"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_device.test", "id"),
					resource.TestCheckResourceAttrSet("simplemdm_device.test", "enrollmenturl"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
