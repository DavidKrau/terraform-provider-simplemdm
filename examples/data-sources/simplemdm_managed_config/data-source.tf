provider "simplemdm" {}

# Lookup an existing managed config after defining it in the same configuration.
resource "simplemdm_app" "example" {
  app_store_id = "586447913"
}

resource "simplemdm_managed_config" "server" {
  app_id     = simplemdm_app.example.id
  key        = "environment"
  value      = "staging"
  value_type = "string"
}

data "simplemdm_managed_config" "server" {
  id = simplemdm_managed_config.server.id
}
