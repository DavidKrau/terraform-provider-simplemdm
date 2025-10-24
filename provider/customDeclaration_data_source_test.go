package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomDeclarationDataSource(t *testing.T) {
	testAccPreCheck(t)
	_ = testAccRequireEnv(t, "SIMPLEMDM_RUN_CUSTOM_DECLARATION_TESTS")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                                resource "simplemdm_customdeclaration" "test" {
                                        name             = "Terraform Custom Declaration"
                                        identifier       = "com.example.terraform.ds"
                                        declaration_type = "com.apple.configuration.management.test"
                                        user_scope       = false
                                        attribute_support = true
                                        escape_attributes = true
                                        platforms        = ["macos"]
                                        data             = jsonencode({
                                                declaration_identifier = "com.example.terraform.ds"
                                                declaration_type       = "com.apple.configuration.management.test"
                                                payload = {
                                                        type       = "com.example"
                                                        identifier = "com.example.payload"
                                                }
                                        })
                                }

                                data "simplemdm_customdeclaration" "test" {
                                        id = simplemdm_customdeclaration.test.id
                                }
                                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.simplemdm_customdeclaration.test", "name", "simplemdm_customdeclaration.test", "name"),
					resource.TestCheckResourceAttrPair("data.simplemdm_customdeclaration.test", "identifier", "simplemdm_customdeclaration.test", "identifier"),
					resource.TestCheckResourceAttrPair("data.simplemdm_customdeclaration.test", "data", "simplemdm_customdeclaration.test", "data"),
					resource.TestCheckResourceAttrPair("data.simplemdm_customdeclaration.test", "payload", "simplemdm_customdeclaration.test", "data"),
					resource.TestCheckResourceAttrPair("data.simplemdm_customdeclaration.test", "user_scope", "simplemdm_customdeclaration.test", "user_scope"),
					resource.TestCheckResourceAttrPair("data.simplemdm_customdeclaration.test", "attribute_support", "simplemdm_customdeclaration.test", "attribute_support"),
					resource.TestCheckResourceAttrPair("data.simplemdm_customdeclaration.test", "escape_attributes", "simplemdm_customdeclaration.test", "escape_attributes"),
				),
			},
		},
	})
}
