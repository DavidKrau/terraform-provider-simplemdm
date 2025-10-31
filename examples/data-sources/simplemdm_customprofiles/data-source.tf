data "simplemdm_customprofiles" "all" {
}

output "custom_profile_count" {
  value = length(data.simplemdm_customprofiles.all.custom_profiles)
}

output "custom_profile_names" {
  value = [for profile in data.simplemdm_customprofiles.all.custom_profiles : profile.name]
}