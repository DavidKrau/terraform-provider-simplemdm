# List all custom attributes in your SimpleMDM account
data "simplemdm_attributes" "all" {}

# Output the first attribute's name
output "first_attribute_name" {
  value = length(data.simplemdm_attributes.all.attributes) > 0 ? data.simplemdm_attributes.all.attributes[0].name : "No attributes found"
}

# Output all attribute names
output "attribute_names" {
  value = [for attr in data.simplemdm_attributes.all.attributes : attr.name]
}