# Advanced Example - Query device users
data "simplemdm_device_users" "device_users" {
  device_id = "123456"
}

output "user_info" {
  description = "User information for the device"
  value = {
    total_users = length(data.simplemdm_device_users.device_users.users)
    users = [
      for user in data.simplemdm_device_users.device_users.users : {
        id        = user.id
        full_name = user.full_name
        user_name = user.user_name
        email     = user.email
      }
    ]
  }
}