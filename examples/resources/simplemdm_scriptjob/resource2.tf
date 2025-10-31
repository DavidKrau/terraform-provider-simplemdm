# Advanced Example - Script job targeting specific devices with custom attribute filtering
resource "simplemdm_script" "update_software" {
  name            = "Software Update Script"
  scriptfile      = file("${path.module}/scripts/update.sh")
  variablesupport = true
}

resource "simplemdm_assignmentgroup" "production_servers" {
  name = "Production Servers"
}

resource "simplemdm_scriptjob" "scheduled_update" {
  script_id = simplemdm_script.update_software.id

  # Target specific assignment groups
  assignment_group_ids = [
    simplemdm_assignmentgroup.production_servers.id,
  ]

  # Use custom attribute to pass data to the script
  custom_attribute = "update_window"

  # Filter devices based on attribute value using regex
  custom_attribute_regex = "^(night|weekend)$"
}

output "job_status" {
  description = "Script job execution details"
  value = {
    id        = simplemdm_scriptjob.scheduled_update.id
    script_id = simplemdm_scriptjob.scheduled_update.script_id
  }
}