# Example with search filter
data "simplemdm_profiles" "wifi_profiles" {
  search = "wifi"
}

output "wifi_profiles" {
  value = [for p in data.simplemdm_profiles.wifi_profiles.profiles : {
    id   = p.id
    name = p.name
    type = p.type
  }]
}