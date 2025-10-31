# Advanced Example - Reference device group (deprecated - use assignment groups)
# ⚠️ DEPRECATED: Device Groups have been superseded by Assignment Groups in SimpleMDM
# This example is maintained for backward compatibility only

data "simplemdm_devicegroup" "legacy_group" {
  id = "123456"
}

output "device_group_info" {
  description = "Device group information (deprecated)"
  value = {
    id   = data.simplemdm_devicegroup.legacy_group.id
    name = data.simplemdm_devicegroup.legacy_group.name
  }
}

# Recommended: Use Assignment Groups instead
# data "simplemdm_assignmentgroup" "modern_group" {
#   id = "123456"
# }