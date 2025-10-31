# Advanced Example - Using attribute data source with conditional logic
data "simplemdm_attribute" "department" {
  name = "department"
}

# Use attribute in device configuration
resource "simplemdm_device" "conditional_device" {
  name = "Department Device"

  attributes = {
    (data.simplemdm_attribute.department.name) = "engineering"
  }
}

output "attribute_metadata" {
  description = "Metadata about the department attribute"
  value = {
    id            = data.simplemdm_attribute.department.id
    name          = data.simplemdm_attribute.department.name
    default_value = data.simplemdm_attribute.department.default_value
  }
}