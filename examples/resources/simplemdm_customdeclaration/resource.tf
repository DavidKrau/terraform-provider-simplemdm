resource "simplemdm_customdeclaration" "test" {
  name                 = "testdeclaration"
  declaration          = jsonencode(jsondecode(file("./testfiles/testdeclaration.json")))
  userscope            = true
  attributesupport     = true
  escapeattributes     = true
  declaration_type     = "com.apple.configuration.safari.bookmarks"
  activation_predicate = ""
}