resource "simplemdm_customprofile" "myprofile" {
  name                   = "My First profiles"
  mobileconfig           = file("./profiles/profile.mobileconfig")
  userscope              = true
  attributesupport       = true
  escapeattributes       = true
  reinstallafterosupdate = false
}