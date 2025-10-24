resource "simplemdm_customdeclaration" "example" {
  name             = "Terraform Custom Declaration"
  identifier       = "com.example.terraform"
  declaration_type = "com.apple.configuration.management"
  platforms        = ["macos"]
  data = jsonencode({
    declaration_identifier = "com.example.terraform"
    declaration_type       = "com.apple.configuration.management"
    payload = {
      type       = "com.example.payload"
      identifier = "com.example.payload"
    }
  })
}
