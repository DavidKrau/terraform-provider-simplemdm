# Advanced Example - Audit device profiles
data "simplemdm_device_profiles" "device_profiles" {
  device_id = "123456"
}

# Filter user-scoped profiles
locals {
  user_scoped_profiles = [
    for profile in data.simplemdm_device_profiles.device_profiles.profiles :
    profile if profile.user_scope
  ]
}

output "profile_audit" {
  description = "Profile audit information for the device"
  value = {
    total_profiles       = length(data.simplemdm_device_profiles.device_profiles.profiles)
    user_scoped_profiles = length(local.user_scoped_profiles)
    profiles = [
      for profile in data.simplemdm_device_profiles.device_profiles.profiles : {
        id                 = profile.id
        name               = profile.name
        profile_identifier = profile.profile_identifier
        user_scope         = profile.user_scope
        attribute_support  = profile.attribute_support
      }
    ]
  }
}