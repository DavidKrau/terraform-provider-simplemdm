package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCheckCustomDeclarationDeviceAssignmentDestroy(s *terraform.State) error {
	client, err := getTestClient()
	if err != nil {
		return fmt.Errorf("failed to create test client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "simplemdm_customdeclaration_device_assignment" {
			continue
		}

		customDeclarationID := rs.Primary.Attributes["custom_declaration_id"]
		deviceID := rs.Primary.Attributes["device_id"]

		// Check if the device still has the custom declaration assigned
		url := fmt.Sprintf("https://%s/api/v1/devices/%s", client.HostName, deviceID)
		httpReq, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		body, err := client.RequestResponse200(httpReq)
		if err != nil {
			// If device doesn't exist, the assignment is definitely destroyed
			if isNotFoundError(err) {
				continue
			}
			return fmt.Errorf("unexpected error checking device %s: %w", deviceID, err)
		}

		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			return fmt.Errorf("failed to parse device response: %w", err)
		}

		assigned, err := deviceHasCustomDeclarationAssignment(body, customDeclarationID)
		if err != nil {
			return fmt.Errorf("error checking assignment: %w", err)
		}

		if assigned {
			return fmt.Errorf("custom declaration assignment %s still exists after destroy", rs.Primary.ID)
		}
	}

	return nil
}

func TestAccCustomDeclarationDeviceAssignmentResource(t *testing.T) {
	testAccPreCheck(t)
	deviceID := testAccRequireEnv(t, "SIMPLEMDM_CUSTOM_DECLARATION_DEVICE_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCustomDeclarationDeviceAssignmentDestroy,
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
