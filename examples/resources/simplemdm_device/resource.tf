resource "simplemdm_device" "firstdevice" {
  // Attribute name (required)
  name        = "mydevice"
  devicename  = "OSmydevice"
  devicegroup = 123456
  profiles    = [456123]
  attributes = {
    "myattribute" = "testvalue"
  }
}