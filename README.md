# Terraform Provider for SimpleMDM

The Terraform SimpleMDM provider is a plugin for Terraform that allows for the full lifecycle management of SimpleMDM resources.

Provider was written by David Kraushuber from FREENOW and was open sourced to the comunity.

## Using the Provider

To use a released provider in your Terraform environment,
run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the
provider. To specify a particular provider version when installing released providers, see
the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions)

## Examples

All the resources and data sources has [one or more examples](./examples) to give you an idea of how to use this
provider to build your own SimpleMDM infrastructure.Provider's official documentation is located in the
[official terraform registry](https://registry.terraform.io/providers/DavidKrau/simplemdm/latest/docs), or [here](./docs/) in form of raw markdown files.

## Acceptance testing

Acceptance tests run automatically whenever `SIMPLEMDM_APIKEY` is defined. The traditional
`TF_ACC=1` flag is still honored for compatibility, but it is no longer required in CI. The suite also
expects a handful of fixture identifiers that point at objects in a SimpleMDM test account. Tests that
rely on these variables automatically skip when values are missing, so day-to-day development can rely
on unit coverage while CI can opt into the full suite simply by exporting the necessary values.

### GitLab CI secrets

Configure the following as masked GitLab CI variables so that pipeline jobs can authenticate against
SimpleMDM:

* `SIMPLEMDM_APIKEY` – API key used by the provider during acceptance tests.
* `SIMPLEMDM_HOST` – optional override when testing against a dedicated SimpleMDM host.

### GitLab CI variables for fixture objects

Set these variables to the identifiers of fixtures that live in your SimpleMDM test tenant. They enable
tests that read existing objects:

* `SIMPLEMDM_APP_ID`
* `SIMPLEMDM_ASSIGNMENT_GROUP_ID`
* `SIMPLEMDM_ASSIGNMENT_GROUP_APP_ID`
* `SIMPLEMDM_ASSIGNMENT_GROUP_GROUP_ID`
* `SIMPLEMDM_ASSIGNMENT_GROUP_PROFILE_ID`
* `SIMPLEMDM_ASSIGNMENT_GROUP_DEVICE_ID`
* `SIMPLEMDM_ASSIGNMENT_GROUP_UPDATED_APP_ID`
* `SIMPLEMDM_ASSIGNMENT_GROUP_UPDATED_DEVICE_ID`
* `SIMPLEMDM_ATTRIBUTE_NAME`
* `SIMPLEMDM_CUSTOM_PROFILE_ID`
* `SIMPLEMDM_DEVICE_GROUP_ID`
* `SIMPLEMDM_DEVICE_ID`
* `SIMPLEMDM_PROFILE_ID`
* `SIMPLEMDM_SCRIPT_ID`
* `SIMPLEMDM_SCRIPT_JOB_ID`
* `SIMPLEMDM_SCRIPT_JOB_SCRIPT_ID`
* `SIMPLEMDM_SCRIPT_JOB_GROUP_ID`
* `SIMPLEMDM_SCRIPT_JOB_DEVICE_ID`

The device group resource tests also expect additional fixture metadata so that profile and attribute
assignments can be exercised end-to-end:

* `SIMPLEMDM_DEVICE_GROUP_NAME`
* `SIMPLEMDM_DEVICE_GROUP_ATTRIBUTE_KEY`
* `SIMPLEMDM_DEVICE_GROUP_ATTRIBUTE_VALUE`
* `SIMPLEMDM_DEVICE_GROUP_ATTRIBUTE_UPDATED_VALUE`
* `SIMPLEMDM_DEVICE_GROUP_PROFILE_ID`
* `SIMPLEMDM_DEVICE_GROUP_PROFILE_UPDATED_ID`
* `SIMPLEMDM_DEVICE_GROUP_CUSTOM_PROFILE_ID`
* `SIMPLEMDM_DEVICE_GROUP_CUSTOM_PROFILE_UPDATED_ID`

The acceptance suite now runs automatically when a `SIMPLEMDM_APIKEY` is available. GitHub Actions
pipelines only need to expose the API key (and any of the optional fixture variables above) to execute the
tests—no additional feature flags are required. Locally, tests that require extra fixtures will continue to
skip until the corresponding environment variables are provided.

## Known Issues

- Device name updates require a workaround via direct PATCH request for API compatibility
- Custom Profiles and Profiles for Assignment Groups and Devices cannot be updated due to API limitations. They are compared only between the plan and state from the previous apply. Changes made directly in the SimpleMDM web interface will not be detected in subsequent terraform apply operations.
