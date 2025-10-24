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

Acceptance tests are guarded behind the standard `TF_ACC=1` flag. They also require a valid
`SIMPLEMDM_APIKEY` and, for data sources that target pre-existing objects, additional identifiers via
environment variables:

* `SIMPLEMDM_ASSIGNMENT_GROUP_ID`
* `SIMPLEMDM_SCRIPT_JOB_ID`

Tests that rely on these variables will automatically skip when a value is not provided, which keeps the
default developer experience lightweight while still enabling full end-to-end coverage when fixtures are
available.

## Know issues

- API current doesnt support Create and Delete for Device Groups
- API currently doesnt support update of "name" attribute for Device Groups
- Custom Profiles/Profiles for Assignment group and Devices can no be updated because of API limitation (they are compared only between plan and state from previous plan), aka adding profile via web will not be considered in next apply.
