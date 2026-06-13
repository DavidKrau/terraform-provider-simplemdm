package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomProfileResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "testprofile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			userscope = true
			attributesupport = true
			escapeattributes = true
			reinstallafterosupdate = true
		}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "name", "testprofile"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "mobileconfig", "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\n<dict>\n    <key>PayloadContent</key>\n    <array>\n        <dict>\n            <key>stickyKey</key>\n            <true/>\n            <key>PayloadIdentifier</key>\n            <string>com.example.myaccessibilitypayload</string>\n            <key>PayloadType</key>\n            <string>com.apple.universalaccess</string>\n            <key>PayloadUUID</key>\n            <string>bff2939d-cb4c-4f6d-8521-e26bc7c03e96</string>\n            <key>PayloadVersion</key>\n            <integer>1</integer>\n            <key>mouseDriverCursorSize</key>\n            <integer>3</integer>\n        </dict>\n    </array>\n    <key>PayloadDisplayName</key>\n    <string>Accessibility</string>\n    <key>PayloadIdentifier</key>\n    <string>com.example.myprofile</string>\n    <key>PayloadType</key>\n    <string>Configuration</string>\n    <key>PayloadUUID</key>\n    <string>e7b55cc7-0d94-4045-8868-dcc1b1c58159</string>\n    <key>PayloadVersion</key>\n    <integer>1</integer>\n</dict>\n</plist>"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "userscope", "true"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "attributesupport", "true"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "escapeattributes", "true"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "reinstallafterosupdate", "true"),
					// Verify defaults for new fields
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "declarative", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "auto_renew_scep_based_certificates", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "allowed_macos_architecture", "any"),
					// Verify default platforms (all platforms)
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "allowed_platforms.#", "5"),
					resource.TestCheckTypeSetElemAttr("simplemdm_customprofile.test", "allowed_platforms.*", "macos"),
					resource.TestCheckTypeSetElemAttr("simplemdm_customprofile.test", "allowed_platforms.*", "ios"),
					resource.TestCheckTypeSetElemAttr("simplemdm_customprofile.test", "allowed_platforms.*", "ipados"),
					resource.TestCheckTypeSetElemAttr("simplemdm_customprofile.test", "allowed_platforms.*", "tvos"),
					resource.TestCheckTypeSetElemAttr("simplemdm_customprofile.test", "allowed_platforms.*", "visionos"),
					resource.TestCheckResourceAttrSet("simplemdm_customprofile.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "simplemdm_customprofile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "simplemdm_customprofile" "test" {
					name = "testprofile2"
					mobileconfig = file("./testfiles/testprofile2.mobileconfig")
					userscope = false
					attributesupport = false
					escapeattributes = false
					reinstallafterosupdate = false
				}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "name", "testprofile2"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "mobileconfig", "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\n<dict>\n    <key>PayloadContent</key>\n    <array>\n        <dict>\n            <key>stickyKey</key>\n            <true/>\n            <key>PayloadIdentifier</key>\n            <string>com.example.myaccessibilitypayload</string>\n            <key>PayloadType</key>\n            <string>com.apple.universalaccess</string>\n            <key>PayloadUUID</key>\n            <string>bff2939d-cb4c-4f6d-8521-e26bc7c03e96</string>\n            <key>PayloadVersion</key>\n            <integer>1</integer>\n            <key>mouseDriverCursorSize</key>\n            <integer>10</integer>\n        </dict>\n    </array>\n    <key>PayloadDisplayName</key>\n    <string>Accessibility</string>\n    <key>PayloadIdentifier</key>\n    <string>com.example.myprofile</string>\n    <key>PayloadType</key>\n    <string>Configuration</string>\n    <key>PayloadUUID</key>\n    <string>e7b55cc7-0d94-4045-8868-dcc1b1c58159</string>\n    <key>PayloadVersion</key>\n    <integer>1</integer>\n</dict>\n</plist>"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "userscope", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "attributesupport", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "escapeattributes", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "reinstallafterosupdate", "false"),
				),
			},
		},
	})
}

func TestAccCustomProfileResource_allFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "full-featured-profile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			userscope = true
			attributesupport = true
			escapeattributes = false
			reinstallafterosupdate = true
			declarative = true
			auto_renew_scep_based_certificates = false
			allowed_platforms = ["macos", "ios"]
			minimum_macos_version = "14.0"
			maximum_macos_version = "15.0"
			allowed_macos_architecture = "arm"
		}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "name", "full-featured-profile"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "mobileconfig", "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\n<dict>\n    <key>PayloadContent</key>\n    <array>\n        <dict>\n            <key>stickyKey</key>\n            <true/>\n            <key>PayloadIdentifier</key>\n            <string>com.example.myaccessibilitypayload</string>\n            <key>PayloadType</key>\n            <string>com.apple.universalaccess</string>\n            <key>PayloadUUID</key>\n            <string>bff2939d-cb4c-4f6d-8521-e26bc7c03e96</string>\n            <key>PayloadVersion</key>\n            <integer>1</integer>\n            <key>mouseDriverCursorSize</key>\n            <integer>3</integer>\n        </dict>\n    </array>\n    <key>PayloadDisplayName</key>\n    <string>Accessibility</string>\n    <key>PayloadIdentifier</key>\n    <string>com.example.myprofile</string>\n    <key>PayloadType</key>\n    <string>Configuration</string>\n    <key>PayloadUUID</key>\n    <string>e7b55cc7-0d94-4045-8868-dcc1b1c58159</string>\n    <key>PayloadVersion</key>\n    <integer>1</integer>\n</dict>\n</plist>"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "declarative", "true"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "auto_renew_scep_based_certificates", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "allowed_platforms.#", "2"),
					resource.TestCheckTypeSetElemAttr("simplemdm_customprofile.test", "allowed_platforms.*", "macos"),
					resource.TestCheckTypeSetElemAttr("simplemdm_customprofile.test", "allowed_platforms.*", "ios"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "minimum_macos_version", "14.0"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "maximum_macos_version", "15.0"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "allowed_macos_architecture", "arm"),
					resource.TestCheckResourceAttrSet("simplemdm_customprofile.test", "id"),
				),
			},
		},
	})
}

func TestAccCustomProfileResource_platformVariations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test with macOS only
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "macos-only-profile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			allowed_platforms = ["macos"]
		}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "allowed_platforms.#", "1"),
					resource.TestCheckTypeSetElemAttr("simplemdm_customprofile.test", "allowed_platforms.*", "macos"),
				),
			},
			// Update to iOS and iPadOS
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "ios-ipados-profile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			allowed_platforms = ["ios", "ipados"]
		}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "name", "ios-ipados-profile"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "allowed_platforms.#", "2"),
					resource.TestCheckTypeSetElemAttr("simplemdm_customprofile.test", "allowed_platforms.*", "ios"),
					resource.TestCheckTypeSetElemAttr("simplemdm_customprofile.test", "allowed_platforms.*", "ipados"),
				),
			},
			// Update to all platforms including visionOS
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "all-platforms-profile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			allowed_platforms = ["macos", "ios", "ipados", "tvos", "visionos"]
		}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "allowed_platforms.#", "5"),
					resource.TestCheckTypeSetElemAttr("simplemdm_customprofile.test", "allowed_platforms.*", "macos"),
					resource.TestCheckTypeSetElemAttr("simplemdm_customprofile.test", "allowed_platforms.*", "ios"),
					resource.TestCheckTypeSetElemAttr("simplemdm_customprofile.test", "allowed_platforms.*", "ipados"),
					resource.TestCheckTypeSetElemAttr("simplemdm_customprofile.test", "allowed_platforms.*", "tvos"),
					resource.TestCheckTypeSetElemAttr("simplemdm_customprofile.test", "allowed_platforms.*", "visionos"),
				),
			},
		},
	})
}

func TestAccCustomProfileResource_architectureVariations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test with Intel only
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "intel-only-profile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			allowed_macos_architecture = "x86"
		}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "allowed_macos_architecture", "x86"),
				),
			},
			// Update to Apple Silicon
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "arm-only-profile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			allowed_macos_architecture = "arm"
		}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "allowed_macos_architecture", "arm"),
				),
			},
			// Update to any architecture
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "any-arch-profile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			allowed_macos_architecture = "any"
		}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "allowed_macos_architecture", "any"),
				),
			},
		},
	})
}

func TestAccCustomProfileResource_versionConstraints(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test with minimum version only
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "min-version-profile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			minimum_macos_version = "14.0"
		}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "minimum_macos_version", "14.0"),
				),
			},
			// Update with both min and max versions
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "version-range-profile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			minimum_macos_version = "13.0"
			maximum_macos_version = "15.0"
		}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "minimum_macos_version", "13.0"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "maximum_macos_version", "15.0"),
				),
			},
		},
	})
}

func TestAccCustomProfileResource_declarativeManagement(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "declarative-profile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			declarative = true
		}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "declarative", "true"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "auto_renew_scep_based_certificates", "false"),
				),
			},
			// Update to disable declarative
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "declarative-profile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			declarative = false
			auto_renew_scep_based_certificates = true
		}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "declarative", "false"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "auto_renew_scep_based_certificates", "true"),
				),
			},
		},
	})
}

func TestAccCustomProfileResource_scepAutoRenew(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "scep-profile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			auto_renew_scep_based_certificates = true
		}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "auto_renew_scep_based_certificates", "true"),
					resource.TestCheckResourceAttr("simplemdm_customprofile.test", "declarative", "false"),
				),
			},
		},
	})
}

func TestAccCustomProfileResource_conflictingAttributes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "conflicting-profile"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			declarative = true
			auto_renew_scep_based_certificates = true
		}
`,
				ExpectError: regexp.MustCompile("Conflicting Attribute Configuration"),
			},
		},
	})
}

func TestAccCustomProfileResource_invalidPlatform(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "invalid-platform"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			allowed_platforms = ["macos", "windows"]
		}
`,
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
		},
	})
}

func TestAccCustomProfileResource_invalidArchitecture(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
		resource "simplemdm_customprofile" "test" {
			name = "invalid-arch"
			mobileconfig = file("./testfiles/testprofile.mobileconfig")
			allowed_macos_architecture = "powerpc"
		}
`,
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
		},
	})
}
