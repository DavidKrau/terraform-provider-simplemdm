resource "simplemdm_app" "app" {
  bundle_id = "com.myCompany.MyApp1"
  deploy_to = "all" // Defaults to "none". Valid values: "none", "outdated", "all".
}
