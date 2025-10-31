# Advanced Example - Device with multiple profiles and attributes
resource "simplemdm_attribute" "location" {
  name          = "office_location"
  default_value = "headquarters"
}

resource "simplemdm_devicegroup" "corporate_devices" {
  name = "Corporate Devices"
}

resource "simplemdm_customprofile" "vpn_profile" {
  name         = "Corporate VPN"
  mobileconfig = file("${path.module}/profiles/vpn.mobileconfig")
}

resource "simplemdm_device" "executive_laptop" {
  name       = "CEO MacBook Pro"
  devicename = "ceo-mbp-2024"

  # Assign to device group
  devicegroup = simplemdm_devicegroup.corporate_devices.id

  # Apply custom profiles
  customprofiles = [
    simplemdm_customprofile.vpn_profile.id,
  ]

  # Set custom attributes
  attributes = {
    (simplemdm_attribute.location.name) = "executive_floor"
    "asset_tag"                         = "EXEC-001"
    "owner"                             = "executive_team"
  }
}

output "device_details" {
  description = "Details of the configured device"
  value = {
    id          = simplemdm_device.executive_laptop.id
    name        = simplemdm_device.executive_laptop.name
    devicename  = simplemdm_device.executive_laptop.devicename
    devicegroup = simplemdm_device.executive_laptop.devicegroup
  }
}