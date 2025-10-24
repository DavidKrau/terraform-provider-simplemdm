package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceResource(t *testing.T) {
	testAccPreCheck(t)

	deviceGroupID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_ID")
	profileID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_PROFILE_ID")
	profileUpdatedID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_PROFILE_UPDATED_ID")
	customProfileID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_CUSTOM_PROFILE_ID")
	customProfileUpdatedID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_CUSTOM_PROFILE_UPDATED_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(providerConfig+`
                resource "simplemdm_device" "test" {
                        name          = "Created test device"
                        devicename    = "Created test device"
                        devicegroup   = %s
                        profiles      = [%s]
                        customprofiles = [%s]
                }
`, deviceGroupID, profileID, customProfileID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_device.test", "name", "Created test device"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "devicename", "Created test device"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "devicegroup", deviceGroupID),
					resource.TestCheckResourceAttr("simplemdm_device.test", "profiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "profiles.0", profileID),
					resource.TestCheckResourceAttr("simplemdm_device.test", "customprofiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "customprofiles.0", customProfileID),
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
				Config: fmt.Sprintf(providerConfig+`
                                resource "simplemdm_device" "test" {
                                        name           = "Created test device changed"
                                        devicename     = "Created test device changed"
                                        devicegroup    = %s
                                        profiles       = [%s]
                                        customprofiles = [%s]
                                  }
`, deviceGroupID, profileUpdatedID, customProfileUpdatedID),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_device.test", "name", "Created test device changed"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "devicename", "Created test device changed"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "devicegroup", deviceGroupID),
					resource.TestCheckResourceAttr("simplemdm_device.test", "profiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "profiles.0", profileUpdatedID),
					resource.TestCheckResourceAttr("simplemdm_device.test", "customprofiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_device.test", "customprofiles.0", customProfileUpdatedID),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_device.test", "id"),
					resource.TestCheckResourceAttrSet("simplemdm_device.test", "enrollmenturl"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
