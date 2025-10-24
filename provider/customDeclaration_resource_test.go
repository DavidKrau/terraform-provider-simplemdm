package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomDeclarationResource(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                                resource "simplemdm_customdeclaration" "test" {
                                        name             = "Terraform Custom Declaration"
                                        identifier       = "com.example.terraform"
                                        declaration_type = "com.apple.configuration.management.test"
                                        topic            = "com.example.topic"
                                        user_scope       = false
                                        attribute_support = true
                                        escape_attributes = true
                                        activation_predicate = "TRUEPREDICATE"
                                        platforms        = ["macos"]
                                        data             = jsonencode({
                                                declaration_identifier = "com.example.terraform"
                                                declaration_type       = "com.apple.configuration.management.test"
                                                payload = {
                                                        type       = "com.example"
                                                        identifier = "com.example.payload"
                                                }
                                        })
                                }
                                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("simplemdm_customdeclaration.test", "id"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "name", "Terraform Custom Declaration"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "identifier", "com.example.terraform"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "declaration_type", "com.apple.configuration.management.test"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "user_scope", "false"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "attribute_support", "true"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "escape_attributes", "true"),
					resource.TestCheckResourceAttrPair("simplemdm_customdeclaration.test", "payload", "simplemdm_customdeclaration.test", "data"),
				),
			},
			{
				ResourceName:      "simplemdm_customdeclaration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: providerConfig + `
                                resource "simplemdm_customdeclaration" "test" {
                                        name             = "Terraform Custom Declaration Updated"
                                        identifier       = "com.example.terraform"
                                        declaration_type = "com.apple.configuration.management.updated"
                                        description      = "Updated description"
                                        user_scope       = true
                                        attribute_support = false
                                        escape_attributes = false
                                        activation_predicate = "FALSEPREDICATE"
                                        platforms        = ["macos", "ios"]
                                        active           = false
                                        data             = jsonencode({
                                                declaration_identifier = "com.example.terraform"
                                                declaration_type       = "com.apple.configuration.management.updated"
                                                payload = {
                                                        type       = "com.example"
                                                        identifier = "com.example.payload"
                                                }
                                        })
                                }
                                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "name", "Terraform Custom Declaration Updated"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "declaration_type", "com.apple.configuration.management.updated"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "active", "false"),
					resource.TestCheckResourceAttr("simplemdm_customdeclaration.test", "user_scope", "true"),
				),
			},
		},
	})
}
