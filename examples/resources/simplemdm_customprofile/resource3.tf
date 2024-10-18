resource "simplemdm_customprofile" "myprofile" {
  name                   = "My First profiles"
  mobileconfig           = <<-EOT
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>PayloadIdentifier</key>
	.....redacted....
 </dict>
 </plist>
EOT
  userscope              = true
  attributesupport       = true
  escapeattributes       = true
  reinstallafterosupdate = false
}