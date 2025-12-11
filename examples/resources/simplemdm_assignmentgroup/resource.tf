resource "simplemdm_assignmentgroup" "myfirstgroup" {
  // Group name required
  name = "My group name"
  //auto deploy true or false, default is true
  auto_deploy = true
  //group type "standard" or "munki", defaults to standard. If this parameter is changed it will destroy/create whole group
  profiles      = [123456, 987654]
  devices       = [135431, 987654]
  profiles_sync = false
  apps_push     = false
  apps_update   = false
  attributes = {
    "testAttribute" = "attributevalue"
  }
  apps = [{ app_id = 553192, deployment_type = "munki", install_type = "managed" }]
}
