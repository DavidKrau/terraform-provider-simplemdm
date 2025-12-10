# Advanced Example - Enable Lost Mode
resource "simplemdm_device_command" "enable_lost_mode" {
  device_id = "123456"
  command   = "enable_lost_mode"

  parameters = {
    message      = "This device has been lost. Please contact IT."
    phone_number = "+15555551234"
    footnote     = "Reward if found"
  }
}

# Advanced Example - Play sound in lost mode
resource "simplemdm_device_command" "lost_mode_sound" {
  device_id = "123456"
  command   = "lost_mode_play_sound"

  # This command depends on lost mode being enabled first
  depends_on = [simplemdm_device_command.enable_lost_mode]
}

# Advanced Example - Update location in lost mode
resource "simplemdm_device_command" "lost_mode_location" {
  device_id = "123456"
  command   = "lost_mode_update_location"

  # This command depends on lost mode being enabled first
  depends_on = [simplemdm_device_command.enable_lost_mode]
}

# Advanced Example - Disable Lost Mode
resource "simplemdm_device_command" "disable_lost_mode" {
  device_id = "123456"
  command   = "disable_lost_mode"

  # Only disable after other lost mode operations complete
  depends_on = [
    simplemdm_device_command.lost_mode_sound,
    simplemdm_device_command.lost_mode_location
  ]
}

output "lost_mode_command_ids" {
  description = "IDs of lost mode commands"
  value = {
    enable  = simplemdm_device_command.enable_lost_mode.id
    sound   = simplemdm_device_command.lost_mode_sound.id
    location = simplemdm_device_command.lost_mode_location.id
    disable = simplemdm_device_command.disable_lost_mode.id
  }
}