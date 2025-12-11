# Using legacy device groups (deprecated)
resource "simplemdm_enrollment" "legacy" {
  device_group_id    = "1234"
  user_enrollment    = false
  welcome_screen     = true
  authentication     = false
  invitation_contact = "user@example.com"
}

# Using modern assignment groups (recommended)
resource "simplemdm_enrollment" "modern" {
  assignment_group_id = "5678"
  user_enrollment     = false
  welcome_screen      = true
  authentication      = false
  invitation_contact  = "user@example.com"
}
