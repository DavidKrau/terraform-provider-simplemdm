# Advanced Example - Audit device profiles
data "simplemdm_device_profiles" "device_profiles" {
  device_id = "123456"
}

# Process profile data
locals {
  profile_ids = [
    for profile in data.simplemdm_device_profiles.device_profiles.profiles :
    profile.id
  ]
}

output "profile_audit" {
  description = "Profile audit information for the device"
  value = {
    total_profiles = length(data.simplemdm_device_profiles.device_profiles.profiles)
    profile_ids    = local.profile_ids
    profiles = [
      for profile in data.simplemdm_device_profiles.device_profiles.profiles : {
        id   = profile.id
        name = profile.name
      }
    ]
  }
}