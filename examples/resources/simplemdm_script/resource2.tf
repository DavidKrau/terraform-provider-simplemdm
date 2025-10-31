# Advanced Example - Script with variable support for templating
resource "simplemdm_script" "device_setup" {
  name            = "Initial Device Setup"
  scriptfile      = file("${path.module}/scripts/setup.sh")
  variablesupport = true
}

# Example script that uses variables:
# #!/bin/bash
# echo "Setting up device for: $DEVICE_NAME"
# echo "Location: $OFFICE_LOCATION"
# # Script continues with device-specific setup...

output "script_id" {
  description = "ID of the setup script for use in script jobs"
  value       = simplemdm_script.device_setup.id
}