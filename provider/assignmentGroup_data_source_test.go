package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAssignmentGroupDataSource(t *testing.T) {
	testAccPreCheck(t)

	assignmentGroupID := testAccRequireEnv(t, "SIMPLEMDM_ASSIGNMENT_GROUP_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
                data "simplemdm_assignmentgroup" "test" {
                  id = "%s"
                }
                `, assignmentGroupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.simplemdm_assignmentgroup.test", "id", assignmentGroupID),
					resource.TestCheckResourceAttrSet("data.simplemdm_assignmentgroup.test", "created_at"),
					resource.TestCheckResourceAttrSet("data.simplemdm_assignmentgroup.test", "updated_at"),
				),
			},
		},
	})
}
