resource "simplemdm_app" "app" {
  app_store_id = "1090161858"
  name         = "Marketing App"
  deploy_to    = "outdated" // Defaults to "none". Valid values: "none", "outdated", "all".
}
