package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomDeclarationDataSource(t *testing.T) {
	testAccPreCheck(t)

	t.Skip("Custom declaration data source test requires a valid custom declaration resource. " +
		"This test is skipped until the custom declaration resource test is fixed.")

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
				                                            Type       = "com.example.test"
				                                            Identifier = "com.example.terraform.ds.payload"
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
