package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceGroupResource(t *testing.T) {
	testAccPreCheck(t)
	_ = testAccRequireEnv(t, "SIMPLEMDM_RUN_DEVICE_GROUP_RESOURCE_TESTS")

	groupID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_ID")
	name := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_NAME")
	attributeKey := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_ATTRIBUTE_KEY")
	attributeValue := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_ATTRIBUTE_VALUE")
	attributeUpdatedValue := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_ATTRIBUTE_UPDATED_VALUE")
	profileID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_PROFILE_ID")
	profileUpdatedID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_PROFILE_UPDATED_ID")
	customProfileID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_CUSTOM_PROFILE_ID")
	customProfileUpdatedID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_CUSTOM_PROFILE_UPDATED_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            providerConfig + testAccDeviceGroupResourceConfig(name, attributeKey, attributeValue, profileID, customProfileID),
				ResourceName:      "simplemdm_devicegroup.test",
				ImportState:       true,
				ImportStateId:     groupID,
				ImportStateVerify: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "name", name),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", fmt.Sprintf("attributes.%s", attributeKey), attributeValue),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "profiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "profiles.0", profileID),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "customprofiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "customprofiles.0", customProfileID),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "id", groupID),
				),
			},
			{
				Config: providerConfig + testAccDeviceGroupResourceConfig(name, attributeKey, attributeUpdatedValue, profileUpdatedID, customProfileUpdatedID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "name", name),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", fmt.Sprintf("attributes.%s", attributeKey), attributeUpdatedValue),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "profiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "profiles.0", profileUpdatedID),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "customprofiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "customprofiles.0", customProfileUpdatedID),
					resource.TestCheckResourceAttrSet("simplemdm_devicegroup.test", "id"),
				),
			},
			{
				Config: providerConfig + testAccDeviceGroupResourceConfig(name, attributeKey, attributeValue, profileID, customProfileID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "name", name),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", fmt.Sprintf("attributes.%s", attributeKey), attributeValue),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "profiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "profiles.0", profileID),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "customprofiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_devicegroup.test", "customprofiles.0", customProfileID),
					resource.TestCheckResourceAttrSet("simplemdm_devicegroup.test", "id"),
				),
			},
		},
	})
}

func testAccDeviceGroupResourceConfig(name, attributeKey, attributeValue, profileID, customProfileID string) string {
	return fmt.Sprintf(`
resource "simplemdm_devicegroup" "test" {
  name = %q

  attributes = {
    %q = %q
  }

  profiles = [%s]

  customprofiles = [%s]
}
`, name, attributeKey, attributeValue, profileID, customProfileID)
}
