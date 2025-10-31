# List all apps in your SimpleMDM account
data "simplemdm_apps" "all" {}

# Output the first app's name
output "first_app_name" {
  value = length(data.simplemdm_apps.all.apps) > 0 ? data.simplemdm_apps.all.apps[0].name : "No apps found"
}

# Output all app IDs
output "app_ids" {
  value = [for app in data.simplemdm_apps.all.apps : app.id]
}