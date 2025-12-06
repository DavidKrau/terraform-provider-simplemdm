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
    id                        = data.simplemdm_customprofile.security_config.id
    name                      = data.simplemdm_customprofile.security_config.name
    user_scope                = data.simplemdm_customprofile.security_config.user_scope
    attribute_support         = data.simplemdm_customprofile.security_config.attribute_support
    escape_attributes         = data.simplemdm_customprofile.security_config.escape_attributes
    reinstall_after_os_update = data.simplemdm_customprofile.security_config.reinstall_after_os_update
  }
}