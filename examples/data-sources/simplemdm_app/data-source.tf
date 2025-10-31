# Query account-specific app only (default behavior)
data "simplemdm_app" "myapp" {
  id = "123456"
}

# Query including shared catalog apps
data "simplemdm_app" "shared_app" {
  id             = "789012"
  include_shared = true
}

output "app_store_identifier" {
  description = "Apple App Store ID associated with the app."
  value       = data.simplemdm_app.myapp.app_store_id
}

output "app_installation_channels" {
  description = "Deployment channels supported by the app."
  value       = data.simplemdm_app.myapp.installation_channels
}
