# List all apps including shared catalog apps
data "simplemdm_apps" "all_with_shared" {
  include_shared = true
}

# Filter apps by type
output "app_store_apps" {
  value = [for app in data.simplemdm_apps.all_with_shared.apps : app.name if app.app_type == "app store"]
}

# Filter apps by platform
output "ios_apps" {
  value = [for app in data.simplemdm_apps.all_with_shared.apps : app.name if app.platform_support == "iOS"]
}