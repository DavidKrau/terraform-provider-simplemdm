data "simplemdm_enrollment" "by_id" {
  id = "1234"
}

output "enrollment_url" {
  value = data.simplemdm_enrollment.by_id.url
}
