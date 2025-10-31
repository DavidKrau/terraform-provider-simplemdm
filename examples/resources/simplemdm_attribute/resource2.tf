# Advanced Example - Attribute with device assignment
resource "simplemdm_attribute" "department" {
  name          = "department"
  default_value = "general"
}

resource "simplemdm_device" "sales_laptop" {
  name = "Sales Team Laptop"
  attributes = {
    (simplemdm_attribute.department.name) = "sales"
  }
}

output "attribute_info" {
  description = "Details about the department attribute"
  value = {
    id            = simplemdm_attribute.department.id
    name          = simplemdm_attribute.department.name
    default_value = simplemdm_attribute.department.default_value
  }
}