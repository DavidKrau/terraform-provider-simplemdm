package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomProfileResource(t *testing.T) {
	testAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name= "testprofile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			userscope = true
			attributesupport = true
			escapeattributes = true
			reinstallafterosupdate = true
			
		  }
`,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "name", "testprofile"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "mobileconfig", "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\n<dict>\n    <key>PayloadContent</key>\n    <array>\n        <dict>\n            <key>stickyKey</key>\n            <true/>\n            <key>PayloadIdentifier</key>\n            <string>com.example.myaccessibilitypayload</string>\n            <key>PayloadType</key>\n            <string>com.apple.universalaccess</string>\n            <key>PayloadUUID</key>\n            <string>bff2939d-cb4c-4f6d-8521-e26bc7c03e96</string>\n            <key>PayloadVersion</key>\n            <integer>1</integer>\n            <key>mouseDriverCursorSize</key>\n            <integer>3</integer>\n        </dict>\n    </array>\n    <key>PayloadDisplayName</key>\n    <string>Accessibility</string>\n    <key>PayloadIdentifier</key>\n    <string>com.example.myprofile</string>\n    <key>PayloadType</key>\n    <string>Configuration</string>\n    <key>PayloadUUID</key>\n    <string>e7b55cc7-0d94-4045-8868-dcc1b1c58159</string>\n    <key>PayloadVersion</key>\n    <integer>1</integer>\n</dict>\n</plist>"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "userscope", "true"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "attributesupport", "true"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "escapeattributes", "true"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "reinstallafterosupdate", "true"),
					resource.TestCheckResourceAttrSet("simplemdm_customprofile.test", "profileidentifier"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "groupcount", "0"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "devicecount", "0"),
					resource.TestCheckResourceAttrSet("simplemdm_customprofile.test", "profilesha"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("simplemdm_customprofile.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "simplemdm_customprofile.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The filesha and  mobileconfig attributes does not exist in SimpleMDM
				// API, therefore there is no value for it during import.
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "simplemdm_customprofile" "test" {
					name= "testprofile2"
					mobileconfig = file("./testfiles/testprofile2.mobileconfig")
					userscope = false
					attributesupport = false
					escapeattributes = false
					reinstallafterosupdate = false
					
				  }
`,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify attributes
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "name", "testprofile2"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "mobileconfig", "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\n<dict>\n    <key>PayloadContent</key>\n    <array>\n        <dict>\n            <key>stickyKey</key>\n            <true/>\n            <key>PayloadIdentifier</key>\n            <string>com.example.myaccessibilitypayload</string>\n            <key>PayloadType</key>\n            <string>com.apple.universalaccess</string>\n            <key>PayloadUUID</key>\n            <string>bff2939d-cb4c-4f6d-8521-e26bc7c03e96</string>\n            <key>PayloadVersion</key>\n            <integer>1</integer>\n            <key>mouseDriverCursorSize</key>\n            <integer>10</integer>\n        </dict>\n    </array>\n    <key>PayloadDisplayName</key>\n    <string>Accessibility</string>\n    <key>PayloadIdentifier</key>\n    <string>com.example.myprofile</string>\n    <key>PayloadType</key>\n    <string>Configuration</string>\n    <key>PayloadUUID</key>\n    <string>e7b55cc7-0d94-4045-8868-dcc1b1c58159</string>\n    <key>PayloadVersion</key>\n    <integer>1</integer>\n</dict>\n</plist>"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "userscope", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "attributesupport", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "escapeattributes", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "reinstallafterosupdate", "false"),
					resource.TestCheckResourceAttrSet("simplemdm_customprofile.test", "profileidentifier"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "groupcount", "0"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "devicecount", "0"),
					resource.TestCheckResourceAttrSet("simplemdm_customprofile.test", "profilesha"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
