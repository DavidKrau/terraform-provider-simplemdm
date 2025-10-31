# Advanced Example - Filter and process multiple devices
data "simplemdm_devices" "all_macbooks" {
  search = "MacBook"
}

# Extract specific device IDs
locals {
  macbook_ids = [for device in data.simplemdm_devices.all_macbooks.devices : device.id]
}

# Use in assignment group
resource "simplemdm_assignmentgroup" "macbook_group" {
  name    = "All MacBooks"
  devices = local.macbook_ids
}

output "macbook_summary" {
  description = "Summary of all MacBooks in the organization"
  value = {
    total_count = length(data.simplemdm_devices.all_macbooks.devices)
    device_ids  = local.macbook_ids
    devices = [
      for device in data.simplemdm_devices.all_macbooks.devices : {
        id     = device.id
        name   = device.name
        model  = device.model
        status = device.status
      }
    ]
  }
}