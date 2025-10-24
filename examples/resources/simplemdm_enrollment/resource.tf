resource "simplemdm_enrollment" "account_driven" {
  device_group_id    = "1234"
  user_enrollment    = false
  welcome_screen     = true
  authentication     = false
  invitation_contact = "user@example.com"
}
