data "simplemdm_enrollments" "all" {
}

output "enrollment_count" {
  value = length(data.simplemdm_enrollments.all.enrollments)
}