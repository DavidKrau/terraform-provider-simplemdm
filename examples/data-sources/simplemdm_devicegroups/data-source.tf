data "simplemdm_devicegroups" "all" {
}

output "device_group_count" {
  value = length(data.simplemdm_devicegroups.all.device_groups)
}