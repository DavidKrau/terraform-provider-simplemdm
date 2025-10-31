# Advanced Example - Query managed app configuration
data "simplemdm_app" "enterprise_app" {
  id = "123456"
}

data "simplemdm_managed_config" "app_config" {
  id = "789012"
}

output "config_details" {
  description = "Managed configuration details"
  value = {
    id         = data.simplemdm_managed_config.app_config.id
    app_id     = data.simplemdm_managed_config.app_config.app_id
    key        = data.simplemdm_managed_config.app_config.key
    value      = data.simplemdm_managed_config.app_config.value
    value_type = data.simplemdm_managed_config.app_config.value_type
  }
}