package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomDeclarationDeviceAssignmentResource(t *testing.T) {
	testAccPreCheck(t)
	_ = testAccRequireEnv(t, "SIMPLEMDM_RUN_CUSTOM_DECLARATION_TESTS")
	deviceID := testAccRequireEnv(t, "SIMPLEMDM_CUSTOM_DECLARATION_DEVICE_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(providerConfig+`
                                resource "simplemdm_customdeclaration" "test" {
                                        name             = "Terraform Custom Declaration Assignment"
                                        identifier       = "com.example.terraform.assignment"
                                        declaration_type = "com.apple.configuration.management.assignment"
                                        platforms        = ["macos"]
                                        data             = jsonencode({
                                                declaration_identifier = "com.example.terraform.assignment"
                                                declaration_type       = "com.apple.configuration.management.assignment"
                                                payload = {
                                                        type       = "com.example"
                                                        identifier = "com.example.payload.assignment"
                                                }
                                        })
                                }

                                resource "simplemdm_customdeclaration_device_assignment" "test" {
                                        custom_declaration_id = simplemdm_customdeclaration.test.id
                                        device_id             = "%s"
                                }
                                `, deviceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("simplemdm_customdeclaration_device_assignment.test", "id"),
					resource.TestCheckResourceAttrPair("simplemdm_customdeclaration_device_assignment.test", "custom_declaration_id", "simplemdm_customdeclaration.test", "id"),
				),
			},
			{
				ResourceName:      "simplemdm_customdeclaration_device_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
