---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "simplemdm Provider"
subcategory: ""
description: |-
  SimpleMDM terraform provider developed by FreeNow.
---

# simplemdm Provider

SimpleMDM terraform provider developed by FreeNow.

## Example Usage

```terraform
provider "simplemdm" {
  host   = "a.simplemdm.com"
  apikey = "yourapikey"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `apikey` (String) API key for you instance, can be set as environment variable SIMPLEMDM_APIKEY
- `host` (String) API host for you instance, can be set as environment variable SIMPLEMDM_HOST, if not set it will default to a.simplemdm.com