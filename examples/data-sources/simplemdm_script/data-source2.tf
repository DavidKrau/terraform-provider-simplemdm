# Advanced Example - Reference script in script job
data "simplemdm_script" "maintenance" {
  id = "123456"
}

resource "simplemdm_scriptjob" "run_maintenance" {
  script_id  = data.simplemdm_script.maintenance.id
  device_ids = ["1001", "1002", "1003"]
}

output "script_details" {
  description = "Details about the maintenance script"
  value = {
    id              = data.simplemdm_script.maintenance.id
    name            = data.simplemdm_script.maintenance.name
    variablesupport = data.simplemdm_script.maintenance.variablesupport
  }
}