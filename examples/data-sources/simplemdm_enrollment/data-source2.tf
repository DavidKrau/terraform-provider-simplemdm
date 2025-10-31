# Advanced Example - Query enrollment and use URL in automation
data "simplemdm_enrollment" "onboarding" {
  id = "1234"
}

# Share enrollment URL via output
output "onboarding_details" {
  description = "Enrollment details for new device onboarding"
  value = {
    id                 = data.simplemdm_enrollment.onboarding.id
    url                = data.simplemdm_enrollment.onboarding.url
    device_group_id    = data.simplemdm_enrollment.onboarding.device_group_id
    user_enrollment    = data.simplemdm_enrollment.onboarding.user_enrollment
    welcome_screen     = data.simplemdm_enrollment.onboarding.welcome_screen
    authentication     = data.simplemdm_enrollment.onboarding.authentication
    invitation_contact = data.simplemdm_enrollment.onboarding.invitation_contact
  }
  sensitive = true
}