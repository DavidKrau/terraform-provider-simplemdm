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

Acceptance tests are guarded behind the standard `TF_ACC=1` flag. They require a valid
`SIMPLEMDM_APIKEY` and a handful of fixture identifiers that point at objects in a SimpleMDM test
account. Tests that rely on these variables automatically skip when values are missing, so day-to-day
development can rely on unit coverage while CI can opt into the full suite.

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
* `SIMPLEMDM_ATTRIBUTE_NAME`
* `SIMPLEMDM_CUSTOM_PROFILE_ID`
* `SIMPLEMDM_DEVICE_GROUP_ID`
* `SIMPLEMDM_DEVICE_ID`
* `SIMPLEMDM_PROFILE_ID`
* `SIMPLEMDM_SCRIPT_ID`
* `SIMPLEMDM_SCRIPT_JOB_ID`

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

### Opt-in flags for mutating tests

Resource acceptance tests that modify or create data are disabled by default. Set the following variables
to `1` (or any non-empty value) in GitLab CI when you want to exercise them:

* `SIMPLEMDM_RUN_APP_RESOURCE_TESTS`
* `SIMPLEMDM_RUN_ASSIGNMENT_GROUP_TESTS`
* `SIMPLEMDM_RUN_CUSTOM_DECLARATION_TESTS`
* `SIMPLEMDM_RUN_DEVICE_RESOURCE_TESTS`
* `SIMPLEMDM_RUN_DEVICE_GROUP_RESOURCE_TESTS`
* `SIMPLEMDM_RUN_PROFILE_RESOURCE_TESTS`
* `SIMPLEMDM_RUN_SCRIPT_JOB_TESTS`

With all of the above defined (alongside `TF_ACC=1`), GitLab CI pipelines can execute the provider's full
acceptance test coverage without manual intervention.

## Know issues

- API current doesnt support Create and Delete for Device Groups
- API currently doesnt support update of "name" attribute for Device Groups
- Custom Profiles/Profiles for Assignment group and Devices can no be updated because of API limitation (they are compared only between plan and state from previous plan), aka adding profile via web will not be considered in next apply.
