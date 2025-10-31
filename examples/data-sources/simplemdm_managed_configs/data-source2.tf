# Advanced Example - Analyze all managed configs for an app
data "simplemdm_app" "enterprise_app" {
  id = "123456"
}

data "simplemdm_managed_configs" "app_configs" {
  app_id = data.simplemdm_app.enterprise_app.id
}

# Extract specific configurations by key
locals {
  environment_config = [
    for config in data.simplemdm_managed_configs.app_configs.managed_configs :
    config if config.key == "environment"
  ]

  boolean_configs = [
    for config in data.simplemdm_managed_configs.app_configs.managed_configs :
    config if config.value_type == "boolean"
  ]
}

output "all_configs" {
  description = "All managed configurations for the app"
  value = {
    total_count = length(data.simplemdm_managed_configs.app_configs.managed_configs)
    configs     = data.simplemdm_managed_configs.app_configs.managed_configs
  }
}

output "environment_setting" {
  description = "Environment configuration value"
  value       = length(local.environment_config) > 0 ? local.environment_config[0].value : "not set"
}

output "boolean_settings" {
  description = "All boolean configuration settings"
  value = {
    count   = length(local.boolean_configs)
    configs = local.boolean_configs
  }
}