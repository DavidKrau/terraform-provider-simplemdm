# Advanced Example - Multiple managed configurations for an app
resource "simplemdm_app" "enterprise_app" {
  bundle_id = "com.example.enterprise"
  name      = "Enterprise App"
}

# String configuration
resource "simplemdm_managed_config" "api_endpoint" {
  app_id     = simplemdm_app.enterprise_app.id
  key        = "apiEndpoint"
  value      = "https://api.company.com/v2"
  value_type = "string"
}

# Boolean configuration
resource "simplemdm_managed_config" "enable_logging" {
  app_id     = simplemdm_app.enterprise_app.id
  key        = "enableLogging"
  value      = "true"
  value_type = "boolean"
}

# Integer configuration
resource "simplemdm_managed_config" "refresh_interval" {
  app_id     = simplemdm_app.enterprise_app.id
  key        = "refreshIntervalSeconds"
  value      = "300"
  value_type = "integer"
}

output "app_configs" {
  description = "Managed configuration IDs for the enterprise app"
  value = {
    api_endpoint     = simplemdm_managed_config.api_endpoint.id
    enable_logging   = simplemdm_managed_config.enable_logging.id
    refresh_interval = simplemdm_managed_config.refresh_interval.id
  }
}