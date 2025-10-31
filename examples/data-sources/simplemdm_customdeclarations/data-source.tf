data "simplemdm_customdeclarations" "all" {
}

output "custom_declaration_count" {
  value = length(data.simplemdm_customdeclarations.all.custom_declarations)
}