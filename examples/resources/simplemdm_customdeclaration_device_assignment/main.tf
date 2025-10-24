resource "simplemdm_customdeclaration" "example" {
  name             = "Example Custom Declaration"
  identifier       = "com.example.customdeclaration.assignment"
  declaration_type = "com.apple.configuration.management.assignment"
  platforms        = ["macos"]
  data = jsonencode({
    declaration_identifier = "com.example.customdeclaration.assignment"
    declaration_type       = "com.apple.configuration.management.assignment"
    payload = {
      type       = "com.example"
      identifier = "com.example.payload"
    }
  })
}

resource "simplemdm_customdeclaration_device_assignment" "example" {
  custom_declaration_id = simplemdm_customdeclaration.example.id
  device_id             = var.device_id
}

variable "device_id" {
  description = "Identifier of the device that should receive the custom declaration."
  type        = string
}
