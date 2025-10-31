package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAssignmentGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAssignmentGroupsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.simplemdm_assignmentgroups.test", "assignment_groups.#"),
				),
			},
		},
	})
}

const testAccAssignmentGroupsDataSourceConfig = `
data "simplemdm_assignmentgroups" "test" {}
`
