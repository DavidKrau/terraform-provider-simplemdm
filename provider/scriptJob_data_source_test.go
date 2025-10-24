package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccScriptJobDataSource(t *testing.T) {
	testAccPreCheck(t)

	scriptJobID := testAccRequireEnv(t, "SIMPLEMDM_SCRIPT_JOB_ID")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
                data "simplemdm_scriptjob" "test" {
                  id = "%s"
                }
                `, scriptJobID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.simplemdm_scriptjob.test", "id", scriptJobID),
				),
			},
		},
	})
}
