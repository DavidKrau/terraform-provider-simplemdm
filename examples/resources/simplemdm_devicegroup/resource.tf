# ⚠️ DEPRECATED: Device Groups have been superseded by Assignment Groups in SimpleMDM
# Please use simplemdm_assignmentgroup instead
# This example is maintained for backward compatibility only

resource "simplemdm_devicegroup" "testgroup" {
  name       = "group2"
  clone_from = "123456"
  attributes = {
    "myattribute" = "attributevalue"
  }
  profiles       = [123456, 654321]
  customprofiles = [456123]
}

# Recommended: Use Assignment Groups instead
# resource "simplemdm_assignmentgroup" "testgroup" {
#   name     = "group2"
#   profiles = [123456, 654321]
#   # Assignment groups provide additional features like auto_deploy, priority, etc.
# }