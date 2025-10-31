# Advanced Example - Clear passcode command
resource "simplemdm_device_command" "clear_passcode" {
  device_id = "123456"
  command   = "clear_passcode"
}

# Advanced Example - Update inventory command
resource "simplemdm_device_command" "update_inventory" {
  device_id = "789012"
  command   = "update_inventory"
}

output "command_ids" {
  description = "IDs of executed commands"
  value = {
    clear_passcode   = simplemdm_device_command.clear_passcode.id
    update_inventory = simplemdm_device_command.update_inventory.id
  }
}