# Advanced Example - Query script job status
data "simplemdm_scriptjob" "completed_job" {
  id = "789012"
}

output "job_execution_details" {
  description = "Details about the completed script job"
  value = {
    id                     = data.simplemdm_scriptjob.completed_job.id
    script_id              = data.simplemdm_scriptjob.completed_job.script_id
    device_ids             = data.simplemdm_scriptjob.completed_job.device_ids
    group_ids              = data.simplemdm_scriptjob.completed_job.group_ids
    assignment_group_ids   = data.simplemdm_scriptjob.completed_job.assignment_group_ids
    custom_attribute       = data.simplemdm_scriptjob.completed_job.custom_attribute
    custom_attribute_regex = data.simplemdm_scriptjob.completed_job.custom_attribute_regex
  }
}