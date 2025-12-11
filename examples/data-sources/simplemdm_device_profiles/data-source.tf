data "simplemdm_device_profiles" "profiles" {
  device_id = "123456"
}

output "direct_profile_names" {
  value = [for profile in data.simplemdm_device_profiles.profiles.profiles : profile.name]
}
