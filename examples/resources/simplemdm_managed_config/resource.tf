provider "simplemdm" {
  # The provider reads SIMPLEMDM_APIKEY from the environment.
}

resource "simplemdm_app" "example" {
  app_store_id = "586447913"
}

resource "simplemdm_managed_config" "server" {
  app_id     = simplemdm_app.example.id
  key        = "serverURL"
  value      = "https://api.example.com"
  value_type = "string"
}
