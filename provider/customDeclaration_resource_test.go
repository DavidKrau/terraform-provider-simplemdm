package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	simplemdm "github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCheckCustomDeclarationDestroy(s *terraform.State) error {
	return testAccCheckResourceDestroyed("simplemdm_customdeclaration", func(client *simplemdm.Client, id string) error {
		url := fmt.Sprintf("https://%s/api/v1/custom_declarations/%s", client.HostName, id)
		httpReq, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		responseBody, err := client.RequestResponse200(httpReq)
		if err != nil {
			return err
		}

		var declaration customDeclarationResponse
		return json.Unmarshal(responseBody, &declaration)
	})(s)
}

func TestAccCustomDeclarationResource(t *testing.T) {
	testAccPreCheck(t)

	t.Skip("Custom declaration creation requires specific Apple declaration types and payloads. " +
		"This test needs to be updated with valid declaration payload for your SimpleMDM instance. " +
		"See Apple's Declarative Device Management documentation for valid declaration types and structures.")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCustomDeclarationDestroy,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				                            resource "simplemdm_customdeclaration" "test" {
				                                    name             = "Terraform Custom Declaration"
				                                    identifier       = "com.example.terraform"
				                                    declaration_type = "com.apple.configuration.management.test"
				                                    user_scope       = false
				                                    attribute_support = true
				                                    escape_attributes = true
				                                    activation_predicate = "TRUEPREDICATE"
				                                    platforms        = ["macos"]
				                                    data             = jsonencode({
				                                            Type       = "com.example.test"
				                                            Identifier = "com.example.terraform.payload"
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
				                                            Type       = "com.example.updated"
				                                            Identifier = "com.example.terraform.updated.payload"
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
