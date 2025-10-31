
terraform {
  required_providers {
    simplemdm = {
      source = "github.com/DavidKrau/simplemdm"
    }
  }
}

provider "simplemdm" {
  host   = "a.simplemdm.com"
  apikey = "yourapikey"
}