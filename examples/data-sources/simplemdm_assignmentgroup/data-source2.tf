# Advanced Example - Using assignment group data in configuration
data "simplemdm_assignmentgroup" "engineering" {
  id = "123456"
}

# Use the data source to reference group in device assignment
resource "simplemdm_device" "new_engineer_laptop" {
  name       = "Engineer Laptop"
  devicename = "eng-laptop-001"

  # Reference devices from the assignment group would typically be done
  # via the assignment group's device management, but this shows data source usage
}

output "group_details" {
  description = "Details about the engineering assignment group"
  value = {
    id           = data.simplemdm_assignmentgroup.engineering.id
    name         = data.simplemdm_assignmentgroup.engineering.name
    device_count = data.simplemdm_assignmentgroup.engineering.device_count
    auto_deploy  = data.simplemdm_assignmentgroup.engineering.auto_deploy
    created_at   = data.simplemdm_assignmentgroup.engineering.created_at
  }
}