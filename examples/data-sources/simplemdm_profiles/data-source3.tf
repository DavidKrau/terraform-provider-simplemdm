# Example with sorting
data "simplemdm_profiles" "sorted" {
  direction = "desc"
}

output "first_profile" {
  value = length(data.simplemdm_profiles.sorted.profiles) > 0 ? {
    id   = data.simplemdm_profiles.sorted.profiles[0].id
    name = data.simplemdm_profiles.sorted.profiles[0].name
    type = data.simplemdm_profiles.sorted.profiles[0].type
  } : null
}