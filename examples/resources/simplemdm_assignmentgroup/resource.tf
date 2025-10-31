resource "simplemdm_assignmentgroup" "myfirstgroup" {
  // Group name required
  name = "My group name"

  // Auto deploy true or false, default is true
  auto_deploy = true

  // ⚠️ DEPRECATED: group_type is deprecated by SimpleMDM API
  // May be ignored for accounts using the New Groups Experience
  // Valid values: "standard" or "munki", defaults to standard
  // If this parameter is changed it will destroy/create whole group
  group_type = "standard"

  // ⚠️ DEPRECATED: install_type is deprecated by SimpleMDM API
  // SimpleMDM recommends setting install_type per-app instead of at group level
  // Valid values: "managed", "self_serve", "managed_updates", "default_installs"
  // Only applies to munki-type assignment groups
  install_type = "managed"

  priority           = 10
  app_track_location = true

  // Assignment relationships
  apps                  = [123456]
  profiles              = [123456, 987654]
  groups                = [135431, 654321]
  devices               = [135431, 987654]
  devices_remove_others = false

  // Post-operation commands
  profiles_sync = false
  apps_push     = false
  apps_update   = false
}
