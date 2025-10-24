resource "simplemdm_device_command" "lock_device" {
  device_id = "123456"
  command   = "lock"

  parameters = {
    message      = "Device locked by Terraform"
    phone_number = "+15555551234"
  }
}
