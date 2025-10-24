resource "simplemdm_app" "marketing" {
  app_store_id = "1090161858"
  name         = "Marketing App"
  deploy_to    = "outdated" // Defaults to "none". Valid values: "none", "outdated", "all".
}

# Computed attributes such as `installation_channels`, `status`, or `version`
# can be referenced once the app has been created. This output illustrates how
# to inspect the deployment channels that SimpleMDM reports for the app.
output "marketing_installation_channels" {
  description = "Deployment channels supported by the Marketing App."
  value       = simplemdm_app.marketing.installation_channels
}
