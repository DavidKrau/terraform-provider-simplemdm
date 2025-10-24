# Upload a custom enterprise or macOS package app by providing a binary file.
# The provider will post the binary to SimpleMDM and keep the metadata in sync.
resource "simplemdm_app" "enterprise" {
  name        = "Internal Tools"
  binary_file = "${path.module}/files/internal-tools.pkg"
}

output "enterprise_processing_status" {
  description = "Processing state for the uploaded enterprise app binary."
  value       = simplemdm_app.enterprise.processing_status
}
