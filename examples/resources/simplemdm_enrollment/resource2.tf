# Advanced Example - User enrollment with authentication using assignment groups
resource "simplemdm_assignmentgroup" "byod_devices" {
  name            = "BYOD Devices"
  auto_deploy     = true
  device_families = ["iPhone", "iPad"]
}

resource "simplemdm_enrollment" "user_enrollment" {
  assignment_group_id = simplemdm_assignmentgroup.byod_devices.id
  user_enrollment     = true
  welcome_screen      = true
  authentication      = true

  # Optional: Send invitation email (only works for one-time enrollments)
  invitation_contact = "byod-users@example.com"
}

output "enrollment_url" {
  description = "URL for users to enroll their devices"
  value       = simplemdm_enrollment.user_enrollment.url
  sensitive   = true
}

output "enrollment_id" {
  description = "Enrollment configuration ID"
  value       = simplemdm_enrollment.user_enrollment.id
}