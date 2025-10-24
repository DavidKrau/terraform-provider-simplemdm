package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDeviceGroupResource(t *testing.T) {
	testAccPreCheck(t)
	_ = testAccRequireEnv(t, "SIMPLEMDM_RUN_DEVICE_GROUP_RESOURCE_TESTS")

	cloneSourceID := testAccRequireEnv(t, "SIMPLEMDM_DEVICE_GROUP_CLONE_SOURCE_ID")
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
				Config: providerConfig + testAccDeviceGroupResourceConfig(name, attributeKey, attributeValue, profileID, customProfileID, cloneSourceID),
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
			{
				ResourceName:      "simplemdm_devicegroup.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["simplemdm_devicegroup.test"]
					if !ok {
						return "", fmt.Errorf("resource not found: simplemdm_devicegroup.test")
					}

					return rs.Primary.ID, nil
				},
			},
			{
				Config: providerConfig + testAccDeviceGroupResourceConfig(name, attributeKey, attributeUpdatedValue, profileUpdatedID, customProfileUpdatedID, cloneSourceID),
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
				Config: providerConfig + testAccDeviceGroupResourceConfig(name, attributeKey, attributeValue, profileID, customProfileID, cloneSourceID),
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

func testAccDeviceGroupResourceConfig(name, attributeKey, attributeValue, profileID, customProfileID, cloneFrom string) string {
	return fmt.Sprintf(`
resource "simplemdm_devicegroup" "test" {
  name = %q
  clone_from = %q

  attributes = {
    %q = %q
  }

  profiles = [%s]

  customprofiles = [%s]
}
`, name, cloneFrom, attributeKey, attributeValue, profileID, customProfileID)
}
