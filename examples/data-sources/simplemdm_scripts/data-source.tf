data "simplemdm_scripts" "all" {
}

output "script_count" {
  value = length(data.simplemdm_scripts.all.scripts)
}