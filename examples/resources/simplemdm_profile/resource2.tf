# Advanced Example - Using profile data source with assignment group
data "simplemdm_profile" "wifi_config" {
  id = "123456"
}

resource "simplemdm_assignmentgroup" "office_devices" {
  name     = "Office Devices"
  profiles = [data.simplemdm_profile.wifi_config.id]
}

output "profile_info" {
  description = "Details about the WiFi profile"
  value = {
    id   = data.simplemdm_profile.wifi_config.id
    name = data.simplemdm_profile.wifi_config.name
  }
}