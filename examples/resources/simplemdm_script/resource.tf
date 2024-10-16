resource "simplemdm_script" "test" {
  name            = "This is test script"
  scriptfile      = file("./testfiles/testscript.sh")
  variablesupport = true
}