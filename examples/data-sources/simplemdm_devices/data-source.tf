data "simplemdm_devices" "all" {
  search = "MacBook"
}

output "device_ids" {
  value = [for device in data.simplemdm_devices.all.devices : device.id]
}
