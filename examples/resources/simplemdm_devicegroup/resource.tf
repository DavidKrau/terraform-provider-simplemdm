# ⚠️ DEPRECATED AND READ-ONLY: Device Groups cannot be created
# This resource is READ-ONLY and can only import existing legacy device groups
# Device groups CANNOT be created via Terraform (API does not support it)

# IMPORTANT: This example will FAIL because device groups cannot be created
# To use this resource, you must import an existing legacy device group:
# terraform import simplemdm_devicegroup.legacy_group 123456

# Example showing attempted creation (THIS WILL FAIL):
# resource "simplemdm_devicegroup" "testgroup" {
#   name       = "group2"
#   clone_from = "123456"  # Clone is also not supported - no create API
#   attributes = {
#     "myattribute" = "attributevalue"
#   }
#   profiles       = [123456, 654321]
#   customprofiles = [456123]
# }
# Error: Device Group Creation Not Supported
# Device Groups are deprecated and cannot be created via the API.

# Recommended: Use Assignment Groups for all new groups
resource "simplemdm_assignmentgroup" "testgroup" {
  name        = "group2"
  auto_deploy = true
  profiles    = [123456, 654321]
  # Assignment groups provide full CRUD functionality
}