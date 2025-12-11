# Advanced Example - Reference legacy device group (deprecated)
# ⚠️ DEPRECATED: This only works with legacy device group IDs
# Use simplemdm_assignmentgroup for modern groups

data "simplemdm_devicegroup" "legacy_group" {
  id = "123456"  # Must be a legacy device group ID
}

output "device_group_info" {
  description = "Legacy device group information (deprecated)"
  value = {
    id   = data.simplemdm_devicegroup.legacy_group.id
    name = data.simplemdm_devicegroup.legacy_group.name
  }
}

# Recommended: Use Assignment Groups for all new implementations
data "simplemdm_assignmentgroup" "modern_group" {
  id = "123456"
}

output "modern_group_info" {
  value = {
    id          = data.simplemdm_assignmentgroup.modern_group.id
    name        = data.simplemdm_assignmentgroup.modern_group.name
    auto_deploy = data.simplemdm_assignmentgroup.modern_group.auto_deploy
  }
}