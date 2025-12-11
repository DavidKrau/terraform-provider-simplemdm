# Advanced Example - Use predefined profile in assignment group
data "simplemdm_profile" "company_wifi" {
  id = "123456"
}

resource "simplemdm_assignmentgroup" "all_devices" {
  name     = "Company Devices"
  profiles = [data.simplemdm_profile.company_wifi.id]
}

output "profile_info" {
  description = "WiFi profile information"
  value = {
    id   = data.simplemdm_profile.company_wifi.id
    name = data.simplemdm_profile.company_wifi.name
    type = data.simplemdm_profile.company_wifi.type
  }
}