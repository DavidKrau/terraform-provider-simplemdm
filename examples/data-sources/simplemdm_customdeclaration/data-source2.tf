# Advanced Example - Reference custom declaration for device assignment
data "simplemdm_customdeclaration" "company_policy" {
  id = "123456"
}

resource "simplemdm_customdeclaration_device_assignment" "apply_policy" {
  custom_declaration_id = data.simplemdm_customdeclaration.company_policy.id
  device_id             = "789012"
}

output "declaration_details" {
  description = "Custom declaration configuration details"
  value = {
    id                = data.simplemdm_customdeclaration.company_policy.id
    name              = data.simplemdm_customdeclaration.company_policy.name
    identifier        = data.simplemdm_customdeclaration.company_policy.identifier
    declaration_type  = data.simplemdm_customdeclaration.company_policy.declaration_type
    platforms         = data.simplemdm_customdeclaration.company_policy.platforms
    user_scope        = data.simplemdm_customdeclaration.company_policy.user_scope
    attribute_support = data.simplemdm_customdeclaration.company_policy.attribute_support
  }
}