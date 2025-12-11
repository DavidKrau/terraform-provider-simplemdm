# Fetch script jobs with pagination
data "simplemdm_scriptjobs" "limited" {
  limit          = 50
  starting_after = 1000
}

output "limited_job_count" {
  value = length(data.simplemdm_scriptjobs.limited.script_jobs)
}

# Access new fields
output "first_job_details" {
  value = length(data.simplemdm_scriptjobs.limited.script_jobs) > 0 ? {
    id               = data.simplemdm_scriptjobs.limited.script_jobs[0].id
    job_name         = data.simplemdm_scriptjobs.limited.script_jobs[0].job_name
    script_name      = data.simplemdm_scriptjobs.limited.script_jobs[0].script_name
    status           = data.simplemdm_scriptjobs.limited.script_jobs[0].status
    pending_count    = data.simplemdm_scriptjobs.limited.script_jobs[0].pending_count
    success_count    = data.simplemdm_scriptjobs.limited.script_jobs[0].success_count
    errored_count    = data.simplemdm_scriptjobs.limited.script_jobs[0].errored_count
  } : null
}