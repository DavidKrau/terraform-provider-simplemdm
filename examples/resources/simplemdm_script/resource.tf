resource "simplemdm_script" "test" {
  name             = "This is test script"
  content          = file("./testfiles/testscript.sh")
  variable_support = true
}