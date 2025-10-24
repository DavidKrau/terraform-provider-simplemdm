# Deploy an existing Volume Purchase Program (VPP) app by referencing its
# bundle identifier. SimpleMDM resolves the bundle ID to the latest catalog
# metadata and makes it available to devices.
resource "simplemdm_app" "bundle_id_example" {
  bundle_id = "com.myCompany.MyApp1"
  deploy_to = "all" // Defaults to "none". Valid values: "none", "outdated", "all".
}

output "bundle_id_app_status" {
  description = "Deployment status reported by SimpleMDM for the bundle ID app."
  value       = simplemdm_app.bundle_id_example.status
}
