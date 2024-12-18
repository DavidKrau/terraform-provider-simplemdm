# simplemdm_scriptjob (Resource)

The `simplemdm_scriptjob` resource allows you to manage script jobs in SimpleMDM. It is used to assign scripts to devices, groups, or assignment groups.

## Example Usage

```terraform
resource "simplemdm_scriptjob" "example" {
  script_id             = "12345"
  device_ids            = ["device1", "device2"]
  group_ids             = ["group1", "group2"]
  assignment_group_ids  = ["assignment_group1"]
  custom_attribute      = "custom_attribute_example"
  custom_attribute_regex = "\\n"
}
```

## Schema

### Required

- `script_id` (String) Required. The ID of the script to be run on the devices.
- `device_ids` (Set of String) Required. A list of device IDs to run the script on. At least one of `device_ids`, `group_ids`, or `assignment_group_ids` must be provided.
- `group_ids` (Set of String) Required. A list of group IDs to run the script on. At least one of `device_ids`, `group_ids`, or `assignment_group_ids` must be provided.
- `assignment_group_ids` (Set of String) Required. A list of assignment group IDs to run the script on. At least one of `device_ids`, `group_ids`, or `assignment_group_ids` must be provided.


The following must be declared, they can be empty but at least one must contain something to work: `device_ids`, `group_ids` and `assignment_group_ids`

### Optional

- `custom_attribute` (String) Optional. If provided, the output from the script will be stored in this custom attribute on each device.
- `custom_attribute_regex` (String) Optional. A regex pattern used to sanitize the output from the script before storing it in the custom attribute.

### Read-Only

- `id` (String) Read-Only. The ID of the script job in SimpleMDM. Automatically generated upon creation.

## Import

This resource supports importing using the script job ID. For example:

```shell
terraform import simplemdm_scriptjob.example <script_job_id>
```
