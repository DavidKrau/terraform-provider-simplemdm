# ⚠️ DEPRECATED: Use simplemdm_assignmentgroups instead
# This only returns legacy device groups

data "simplemdm_devicegroups" "all_legacy" {
}

output "legacy_group_count" {
  description = "Number of legacy device groups (deprecated)"
  value       = length(data.simplemdm_devicegroups.all_legacy.device_groups)
}

# Recommended: Use Assignment Groups instead
data "simplemdm_assignmentgroups" "all_modern" {
}

output "modern_group_count" {
  description = "Number of assignment groups"
  value       = length(data.simplemdm_assignmentgroups.all_modern.assignment_groups)
}