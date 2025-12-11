# Advanced Example - Query device users
data "simplemdm_device_users" "device_users" {
  device_id = "123456"  # Must be a macOS device
}

# Filter logged in users
locals {
  logged_in_users = [
    for user in data.simplemdm_device_users.device_users.users :
    user if user.logged_in
  ]
}

output "user_info" {
  description = "User information for the device"
  value = {
    total_users     = length(data.simplemdm_device_users.device_users.users)
    logged_in_users = length(local.logged_in_users)
    users = [
      for user in data.simplemdm_device_users.device_users.users : {
        username       = user.username
        full_name      = user.full_name
        secure_token   = user.secure_token
        logged_in      = user.logged_in
        mobile_account = user.mobile_account
      }
    ]
  }
}