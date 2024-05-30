resource "simplemdm_script" "test" {
  name            = "This is test script"
  scriptfile      = "./testfiles/testscript.sh"
  filesha         = filesha256("./testfiles/testscript.sh")
  variablesupport = true
}