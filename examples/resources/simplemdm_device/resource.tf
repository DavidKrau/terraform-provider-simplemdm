resource "simplemdm_device" "firstdevice" {
  // Device name (required)
  name = "mydevice"
  
  // Optional: Set the device hostname (requires supervision, async operation)
  devicename = "OSmydevice"
  
  // Optional: Assign to device group (deprecated - use static_group_ids instead)
  devicegroup = 123456
  
  // Optional: Apply profiles
  profiles       = [456123]
  customprofiles = [456123]
  
  // Optional: Set custom attributes
  attributes = {
    "myattribute" = "testvalue"
  }
}

# Example using static_group_ids (recommended)
resource "simplemdm_device" "modern_device" {
  name = "modern-device"
  
  # Assign to multiple static groups
  static_group_ids = ["123", "456"]
  
  attributes = {
    "location" = "office"
  }
}