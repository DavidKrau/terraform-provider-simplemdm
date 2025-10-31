# Advanced Example - Restart device command
resource "simplemdm_device_command" "restart_device" {
  device_id = "123456"
  command   = "restart"
}

output "restart_command_id" {
  description = "ID of the restart command"
  value       = simplemdm_device_command.restart_device.id
}