# Advanced Example - Reference custom profile in assignment
data "simplemdm_customprofile" "security_config" {
  id = "123456"
}

resource "simplemdm_assignmentgroup" "secure_devices" {
  name     = "High Security Devices"
  profiles = [data.simplemdm_customprofile.security_config.id]
}

output "custom_profile_details" {
  description = "Details about the security configuration profile"
  value = {
    id                     = data.simplemdm_customprofile.security_config.id
    name                   = data.simplemdm_customprofile.security_config.name
    userscope              = data.simplemdm_customprofile.security_config.userscope
    attributesupport       = data.simplemdm_customprofile.security_config.attributesupport
    escapeattributes       = data.simplemdm_customprofile.security_config.escapeattributes
    reinstallafterosupdate = data.simplemdm_customprofile.security_config.reinstallafterosupdate
  }
}