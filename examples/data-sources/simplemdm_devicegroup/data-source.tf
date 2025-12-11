# ⚠️ DEPRECATED: Device Groups have been superseded by Assignment Groups in SimpleMDM
# This data source only works with legacy device group IDs
# Please use simplemdm_assignmentgroup data source instead
# This example is maintained for backward compatibility only

data "simplemdm_devicegroup" "legacy_devicegroup" {
  id = "123456"  # Must be a legacy device group ID from a migrated group
}

# Recommended: Use Assignment Groups instead
data "simplemdm_assignmentgroup" "modern_group" {
  id = "123456"
}