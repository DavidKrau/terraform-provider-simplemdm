resource "simplemdm_devicegroup" "testgroup" {
  name       = "group2"
  clone_from = "123456"
  attributes = {
    "myattribute" = "attributevalue"
  }
  profiles       = [123456, 654321]
  customprofiles = [456123]
}