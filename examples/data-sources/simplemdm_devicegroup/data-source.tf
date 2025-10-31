# ⚠️ DEPRECATED: Device Groups have been superseded by Assignment Groups in SimpleMDM
# Please use simplemdm_assignmentgroup data source instead
# This example is maintained for backward compatibility only

data "simplemdm_devicegroup" "devicegroup" {
  id = "123456"
}

# Recommended: Use Assignment Groups instead
# data "simplemdm_assignmentgroup" "devicegroup" {
#   id = "123456"
# }