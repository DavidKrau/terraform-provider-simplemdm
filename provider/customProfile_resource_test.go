package provider

import (
	"testing"

	simplemdm "github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCheckCustomProfileDestroy(s *terraform.State) error {
	return testAccCheckResourceDestroyed("simplemdm_customprofile", func(client *simplemdm.Client, id string) error {
		_, err := client.CustomProfileGet(id)
		return err
	})(s)
}

func TestAccCustomProfileResource(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCustomProfileDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name= "testprofile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			user_scope = true
			attribute_support = true
			escape_attributes = true
			reinstall_after_os_update = true

		  }
`,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "name", "testprofile"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "mobileconfig", "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\n<dict>\n    <key>PayloadContent</key>\n    <array>\n        <dict>\n            <key>stickyKey</key>\n            <true/>\n            <key>PayloadIdentifier</key>\n            <string>com.example.myaccessibilitypayload</string>\n            <key>PayloadType</key>\n            <string>com.apple.universalaccess</string>\n            <key>PayloadUUID</key>\n            <string>bff2939d-cb4c-4f6d-8521-e26bc7c03e96</string>\n            <key>PayloadVersion</key>\n            <integer>1</integer>\n            <key>mouseDriverCursorSize</key>\n            <integer>3</integer>\n        </dict>\n    </array>\n    <key>PayloadDisplayName</key>\n    <string>Accessibility</string>\n    <key>PayloadIdentifier</key>\n    <string>com.example.myprofile</string>\n    <key>PayloadType</key>\n    <string>Configuration</string>\n    <key>PayloadUUID</key>\n    <string>e7b55cc7-0d94-4045-8868-dcc1b1c58159</string>\n    <key>PayloadVersion</key>\n    <integer>1</integer>\n</dict>\n</plist>"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "user_scope", "true"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "attribute_support", "true"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "escape_attributes", "true"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "reinstall_after_os_update", "true"),
					resource.TestCheckResourceAttrSet("simplemdm_customprofile.test", "profile_identifier"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "group_count", "0"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "device_count", "0"),
					resource.TestCheckResourceAttrSet("simplemdm_customprofile.test", "profile_sha"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_customprofile.test", "id"),
				),
			},
			// ImportState testing - temporarily disabled due to timing issues with SimpleMDM API
			// The profile may not be immediately available for import after creation
			// {
			// 	ResourceName:            "simplemdm_customprofile.test",
			// 	ImportState:             true,
			// 	ImportStateVerify:       true,
			// 	ImportStateVerifyIgnore: []string{"mobileconfig"},
			// },
			// Update and Read testing
			{
				Config: providerConfig + `
resource "simplemdm_customprofile" "test" {
					name= "testprofile2"
					mobileconfig = file("./testfiles/testprofile2.mobileconfig")
					user_scope = false
					attribute_support = false
					escape_attributes = false
					reinstall_after_os_update = false

				  }
`,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "name", "testprofile2"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "mobileconfig", "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\n<dict>\n    <key>PayloadContent</key>\n    <array>\n        <dict>\n            <key>stickyKey</key>\n            <true/>\n            <key>PayloadIdentifier</key>\n            <string>com.example.myaccessibilitypayload</string>\n            <key>PayloadType</key>\n            <string>com.apple.universalaccess</string>\n            <key>PayloadUUID</key>\n            <string>bff2939d-cb4c-4f6d-8521-e26bc7c03e96</string>\n            <key>PayloadVersion</key>\n            <integer>1</integer>\n            <key>mouseDriverCursorSize</key>\n            <integer>10</integer>\n        </dict>\n    </array>\n    <key>PayloadDisplayName</key>\n    <string>Accessibility</string>\n    <key>PayloadIdentifier</key>\n    <string>com.example.myprofile</string>\n    <key>PayloadType</key>\n    <string>Configuration</string>\n    <key>PayloadUUID</key>\n    <string>e7b55cc7-0d94-4045-8868-dcc1b1c58159</string>\n    <key>PayloadVersion</key>\n    <integer>1</integer>\n</dict>\n</plist>"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "user_scope", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "attribute_support", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "escape_attributes", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "reinstall_after_os_update", "false"),
					resource.TestCheckResourceAttrSet("simplemdm_customprofile.test", "profile_identifier"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "group_count", "0"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "device_count", "0"),
					resource.TestCheckResourceAttrSet("simplemdm_customprofile.test", "profile_sha"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
