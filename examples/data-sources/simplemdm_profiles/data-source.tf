data "simplemdm_profiles" "all" {
}

output "profile_count" {
  value = length(data.simplemdm_profiles.all.profiles)
}