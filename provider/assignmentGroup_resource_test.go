package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAssignmentGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_assignmentgroup" "testgroup2" {
					name= "This assignment group"
					auto_deploy = false
					group_type   = "standard"
					install_type = "managed"
					apps= [550788]
					groups = [140129]
					profiles = [172500]
					devices = [1597950]
					profiles_sync = false
					apps_push = false
					apps_update = false
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "name", "This assignment group"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "group_type", "standard"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "profiles.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "profiles.0", "172500"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices.0", "1597950"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "apps.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "apps.0", "550788"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "groups.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "groups.0", "140129"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_assignmentgroup.testgroup2", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "simplemdm_assignmentgroup.testgroup2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"apps_update", "apps_push", "auto_deploy", "profiles_sync", "install_type", "profiles"},
			},
			//Update and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_assignmentgroup" "testgroup2" {
					name= "renamed assignemnt group"
					auto_deploy = false
					group_type   = "munki"
					install_type = "managed"
					apps= [550788]
					devices = [1597950]
					profiles_sync = false
					apps_push = false
					apps_update = false
				  }
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "name", "renamed assignemnt group"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "group_type", "munki"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "install_type", "managed"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices.0", "1597950"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "apps.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "apps.0", "550788"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_assignmentgroup.testgroup2", "id"),
				),
			},
			//Delete testing automatically occurs in TestCase
		},
	})
}
