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

## Know issues

-app assignemnt is not exactly working as expected because of missing data from API (changes requested already), currently there will be always diff in app assignement.