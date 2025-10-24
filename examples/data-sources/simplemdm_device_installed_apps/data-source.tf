data "simplemdm_device_installed_apps" "apps" {
  device_id = "123456"
}

output "managed_apps" {
  value = [
    for app in data.simplemdm_device_installed_apps.apps.installed_apps :
    jsondecode(app.attributes_json)
  ]
}
