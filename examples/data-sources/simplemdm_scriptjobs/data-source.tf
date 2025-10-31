data "simplemdm_scriptjobs" "all" {
}

output "script_job_count" {
  value = length(data.simplemdm_scriptjobs.all.script_jobs)
}