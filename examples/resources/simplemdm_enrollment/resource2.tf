# Advanced Example - User enrollment with authentication
resource "simplemdm_devicegroup" "byod_devices" {
  name = "BYOD Devices"
}

resource "simplemdm_enrollment" "user_enrollment" {
  device_group_id = simplemdm_devicegroup.byod_devices.id
  user_enrollment = true
  welcome_screen  = true
  authentication  = true

  # Optional: Send invitation email
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