resource "simplemdm_profile" "myprofile" {
  //Custom Profile name (required)
  name = "My First profiles"
  //path to the file (required)
  mobileconfig = "./profiles/profile.mobileconfig"
  //function for count SHA256 (required)
  filesha = filesha256("./profiles/profile.mobileconfig")
  //Scope of the profile true/false, default to false
  userscope = true
  //Enable attribute support true/false, default to false
  attributesupport = true
  //Escape attributes in the profile true/false, default to false
  escapeattributes = true
  //Reisntall profile after OS update true/false, default to false
  reinstallafterosupdate = false
}