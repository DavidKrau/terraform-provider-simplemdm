resource "simplemdm_customprofile" "myprofile" {
  name                      = "My First profiles"
  mobileconfig              = file("./profiles/profile.mobileconfig")
  user_scope                = true
  attribute_support         = true
  escape_attributes         = true
  reinstall_after_os_update = false
}