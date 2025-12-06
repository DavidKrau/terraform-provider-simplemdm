# Example 1: Basic Assignment Group
resource "simplemdm_assignmentgroup" "basic" {
  name        = "My Assignment Group"
  auto_deploy = true
  priority    = 10

  apps     = [123456]
  profiles = [123456, 987654]
  devices  = [135431, 987654]
}

# Example 2: Assignment Group with Command Triggers
resource "simplemdm_assignmentgroup" "with_commands" {
  name        = "Production Apps Group"
  auto_deploy = true

  apps     = [123456]
  profiles = [789012]

  # Trigger commands after group changes
  # apps_update: Only installs apps with available updates
  apps_update = false

  # apps_push: Reinstalls all assigned apps regardless of version
  apps_push = false

  # profiles_sync: Syncs all profiles to devices
  # ⚠️ RATE LIMITED: 1 request per 30 seconds
  # Wait 30 seconds between applies when using this
  profiles_sync = false
}

# Example 3: Legacy Munki Group (Deprecated)
resource "simplemdm_assignmentgroup" "legacy_munki" {
  name = "Legacy Munki Group"

  # ⚠️ DEPRECATED: group_type may be ignored for New Groups Experience
  # Valid values: "standard" or "munki", defaults to standard
  # Changing this will destroy and recreate the group
  group_type = "munki"

  # ⚠️ DEPRECATED: install_type should be set per-app instead
  # Valid values: "managed", "self_serve", "managed_updates", "default_installs"
  # Only applies to munki-type groups
  install_type = "managed"

  apps = [123456]
}

# Example 4: Using Deprecated Device Groups
resource "simplemdm_assignmentgroup" "with_device_groups" {
  name = "Group with Device Groups"

  # ⚠️ DEPRECATED: Device groups use a deprecated API
  # Only works with legacy_device_group_id from migrated groups
  # For New Groups Experience, assign devices directly instead
  groups = [135431, 654321]
}

# USAGE GUIDELINES:
#
# App Commands:
# - apps_push: Pushes ALL assigned apps to devices, regardless of current version
# - apps_update: Only pushes apps where newer version is available
# - Use apps_push after adding new apps or for reinstallation
# - Use apps_update for routine update deployments
#
# Profile Sync Rate Limiting:
# The profiles_sync command is rate limited to 1 request per 30 seconds by SimpleMDM.
# If you see rate limit errors:
# 1. Wait 30 seconds between applies when using profiles_sync
# 2. Consider using separate terraform apply commands for profile changes
# 3. Set profiles_sync to true only when needed, not on every apply
#
# Device Group Deprecation:
# The groups attribute uses a deprecated API. For new implementations:
# - Assign devices directly using the devices attribute
# - Use the New Groups Experience in SimpleMDM
# - Legacy device groups only work with legacy_device_group_id
#
# Priority:
# - Valid range: 0-999
# - Lower numbers are evaluated first
# - Use to control order when devices are in multiple groups
