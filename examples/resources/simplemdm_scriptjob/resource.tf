resource "simplemdm_scriptjob" "test" {
  script_id              = "35"
  device_ids             = ["1", "2", "3"]
  assignment_group_ids   = ["6", "7"]
  custom_attribute       = "greeting_attribute"
  custom_attribute_regex = "\\n"
}