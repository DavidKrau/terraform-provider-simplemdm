data "simplemdm_devices" "all" {
  search = "MacBook"
}

output "device_ids" {
  value = [for device in data.simplemdm_devices.all.devices : device.id]
}

# Search by serial number
data "simplemdm_devices" "by_serial" {
  search = "DNFJE9DNG5MG"
}

# Search by device name
data "simplemdm_devices" "by_name" {
  search = "Mike's iPhone"
}

# Search by UDID
data "simplemdm_devices" "by_udid" {
  search = "00008030-001234567890401E"
}
