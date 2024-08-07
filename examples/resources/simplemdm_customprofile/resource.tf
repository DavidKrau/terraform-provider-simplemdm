resource "simplemdm_customprofile" "myprofile" {
  name                   = "My First profiles"
  mobileconfig           = "./profiles/profile.mobileconfig"
  filesha                = filesha256("./profiles/profile.mobileconfig")
  userscope              = true
  attributesupport       = true
  escapeattributes       = true
  reinstallafterosupdate = false
}