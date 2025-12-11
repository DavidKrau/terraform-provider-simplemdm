# Advanced Example - Analyze installed apps on a device
data "simplemdm_device_installed_apps" "device_apps" {
  device_id = "123456"
}

# Filter managed apps
locals {
  managed_apps = [
    for app in data.simplemdm_device_installed_apps.device_apps.installed_apps :
    app if app.managed
  ]
}

output "app_inventory" {
  description = "Comprehensive app inventory for the device"
  value = {
    total_apps   = length(data.simplemdm_device_installed_apps.device_apps.installed_apps)
    managed_apps = length(local.managed_apps)
    app_details = [
      for app in data.simplemdm_device_installed_apps.device_apps.installed_apps :
      {
        name            = app.name
        identifier      = app.identifier
        version         = app.version
        managed         = app.managed
        bundle_size     = app.bundle_size
        discovered_at   = app.discovered_at
      }
    ]
  }
}