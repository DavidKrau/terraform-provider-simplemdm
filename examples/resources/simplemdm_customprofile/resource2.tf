resource "simplemdm_customprofile" "myprofile" {
  name                   = "My First profiles"
  mobileconfig           = templatefile("./profiles/profile.mobileconfig", { foo = "bar" })
  userscope              = true
  attributesupport       = true
  escapeattributes       = true
  reinstallafterosupdate = false
}