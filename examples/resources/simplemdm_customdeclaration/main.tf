resource "simplemdm_customdeclaration" "example" {
  name                      = "Terraform Custom Declaration"
  declaration_type          = "com.apple.configuration.management.status-subscriptions"
  user_scope                = false
  attribute_support         = true
  escape_attributes         = true
  activation_predicate      = "TRUEPREDICATE"
  reinstall_after_os_update = false
  payload = jsonencode({
    Type       = "com.apple.configuration.management.status-subscriptions"
    Identifier = "com.example.terraform.status"
    ServerToken = "example-token"
    StatusItems = [
      {
        Name = "device.identifier.serial-number"
      }
    ]
  })
}
