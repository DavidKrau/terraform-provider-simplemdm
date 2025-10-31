# Advanced Example - Analyze installed apps on a device
data "simplemdm_device_installed_apps" "device_apps" {
  device_id = "123456"
}

# Process the installed apps data
locals {
  managed_apps = [
    for app in data.simplemdm_device_installed_apps.device_apps.installed_apps :
    jsondecode(app.attributes_json) if lookup(jsondecode(app.attributes_json), "is_managed", false)
  ]
}

output "app_inventory" {
  description = "Comprehensive app inventory for the device"
  value = {
    total_apps   = length(data.simplemdm_device_installed_apps.device_apps.installed_apps)
    managed_apps = length(local.managed_apps)
    app_details = [
      for app in data.simplemdm_device_installed_apps.device_apps.installed_apps :
      jsondecode(app.attributes_json)
    ]
  }
}