resource "simplemdm_devicegroup" "testgroup" {
  name = "group2"
  attributes = {
    "myattribute" = "attributevalue"
  }
  profiles       = [123456, 654321]
  customprofiles = [456123]
}