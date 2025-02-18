resource "simplemdm_app" "app" {
  bundle_id = "com.myCompany.MyApp1"
  deploy_to = "all" // Default to "none" if not added but possible values are "outdated" and "all"
}