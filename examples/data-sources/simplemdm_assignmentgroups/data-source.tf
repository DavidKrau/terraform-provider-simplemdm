# List all assignment groups in your SimpleMDM account
data "simplemdm_assignmentgroups" "all" {}

# Output the first assignment group's name
output "first_group_name" {
  value = length(data.simplemdm_assignmentgroups.all.assignment_groups) > 0 ? data.simplemdm_assignmentgroups.all.assignment_groups[0].name : "No groups found"
}

# Output all assignment group IDs
output "group_ids" {
  value = [for group in data.simplemdm_assignmentgroups.all.assignment_groups : group.id]
}