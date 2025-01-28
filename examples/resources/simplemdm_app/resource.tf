resource "simplemdm_app" "app" {
  app_store_id  = "1090161858"
  deploy_to     = "outdated" // Default to "none" if not added but possible values are "outdated" and "all"
}