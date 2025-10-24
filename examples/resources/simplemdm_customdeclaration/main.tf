resource "simplemdm_customdeclaration" "example" {
  name                 = "Terraform Custom Declaration"
  identifier           = "com.example.terraform"
  declaration_type     = "com.apple.configuration.management"
  topic                = "com.example.topic"
  description          = "Example declaration managed by Terraform"
  user_scope           = false
  attribute_support    = true
  escape_attributes    = true
  activation_predicate = "TRUEPREDICATE"
  platforms            = ["macos"]
  data = jsonencode({
    declaration_identifier = "com.example.terraform"
    declaration_type       = "com.apple.configuration.management"
    payload = {
      type       = "com.example.payload"
      identifier = "com.example.payload"
    }
  })
}
