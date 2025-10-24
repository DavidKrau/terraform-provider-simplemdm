data "simplemdm_device_users" "users" {
  device_id = "123456"
}

output "secure_token_users" {
  value = [
    for user in data.simplemdm_device_users.users.users :
    jsondecode(user.attributes_json).username
  ]
}
