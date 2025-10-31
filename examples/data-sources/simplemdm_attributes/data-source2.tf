# List all custom attributes and filter by default value
data "simplemdm_attributes" "all" {}

# Find attributes with default values set
output "attributes_with_defaults" {
  value = [for attr in data.simplemdm_attributes.all.attributes : attr.name if attr.default_value != ""]
}

# Find attributes without default values
output "attributes_without_defaults" {
  value = [for attr in data.simplemdm_attributes.all.attributes : attr.name if attr.default_value == ""]
}

# Create a map of attribute names to default values
output "attribute_defaults_map" {
  value = {
    for attr in data.simplemdm_attributes.all.attributes : attr.name => attr.default_value
  }
}