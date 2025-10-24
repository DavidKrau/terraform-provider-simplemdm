package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAssignmentGroupResource(t *testing.T) {
	testAccPreCheck(t)

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
					apps= [577575]
					groups = [140188]
					profiles = [172801]
					devices = [1601809]
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
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "profiles.0", "172801"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices.0", "1601809"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "apps.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "apps.0", "577575"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "groups.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "groups.0", "140188"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_assignmentgroup.testgroup2", "id"),
					resource.TestCheckResourceAttrSet("simplemdm_assignmentgroup.testgroup2", "created_at"),
					resource.TestCheckResourceAttrSet("simplemdm_assignmentgroup.testgroup2", "updated_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "simplemdm_assignmentgroup.testgroup2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"apps_update", "apps_push", "auto_deploy", "profiles_sync", "install_type", "profiles", "created_at", "updated_at", "device_count", "group_count"},
			},
			//Update and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_assignmentgroup" "testgroup2" {
					name= "renamed assignemnt group"
					auto_deploy = false
					group_type   = "munki"
					install_type = "managed"
					apps= [553192]
					devices = [1601810]
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
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "devices.0", "1601810"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "apps.#", "1"),
					resource.TestCheckResourceAttr("simplemdm_assignmentgroup.testgroup2", "apps.0", "553192"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_assignmentgroup.testgroup2", "id"),
				),
			},
			//Delete testing automatically occurs in TestCase
		},
	})
}
