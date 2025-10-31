# Advanced Example - Query device and use its properties
data "simplemdm_device" "primary_server" {
  id = "138262"
}

# Use device information in outputs or conditional logic
output "device_info" {
  description = "Comprehensive device information"
  value = {
    id              = data.simplemdm_device.primary_server.id
    name            = data.simplemdm_device.primary_server.name
    serial_number   = data.simplemdm_device.primary_server.serial_number
    model           = data.simplemdm_device.primary_server.model
    os_version      = data.simplemdm_device.primary_server.os_version
    device_group_id = data.simplemdm_device.primary_server.device_group_id
    last_seen_at    = data.simplemdm_device.primary_server.last_seen_at
    status          = data.simplemdm_device.primary_server.status
  }
}