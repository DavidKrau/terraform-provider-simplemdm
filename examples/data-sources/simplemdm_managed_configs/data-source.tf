provider "simplemdm" {}

# List all managed configs for an app
resource "simplemdm_app" "example" {
  app_store_id = "586447913"
}

resource "simplemdm_managed_config" "config1" {
  app_id     = simplemdm_app.example.id
  key        = "environment"
  value      = "production"
  value_type = "string"
}

resource "simplemdm_managed_config" "config2" {
  app_id     = simplemdm_app.example.id
  key        = "debug_mode"
  value      = "false"
  value_type = "boolean"
}

data "simplemdm_managed_configs" "all" {
  app_id = simplemdm_app.example.id
  depends_on = [
    simplemdm_managed_config.config1,
    simplemdm_managed_config.config2,
  ]
}

output "config_count" {
  description = "Number of managed configurations"
  value       = length(data.simplemdm_managed_configs.all.managed_configs)
}