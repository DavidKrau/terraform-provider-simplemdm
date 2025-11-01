# Terraform Provider for SimpleMDM

This repository contains the Terraform provider that manages resources in
[SimpleMDM](https://simplemdm.com). The provider is implemented with the
[Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)
and uses the official
[`simplemdm-go-client`](https://github.com/DavidKrau/simplemdm-go-client) to talk to
SimpleMDM's REST API.

## Installation

Add the provider to your Terraform configuration. Published releases are available on
[the Terraform Registry](https://registry.terraform.io/providers/DavidKrau/simplemdm):

```terraform
terraform {
  required_providers {
    simplemdm = {
      source  = "DavidKrau/simplemdm"
      version = "~> 0.1"
    }
  }
}

provider "simplemdm" {
  # Optional. Defaults to https://a.simplemdm.com
  host   = "a.simplemdm.com"

  # Required unless provided via the SIMPLEMDM_APIKEY environment variable.
  apikey = var.simplemdm_api_key
}
```

The provider accepts two configuration attributes:

| Attribute | Environment variable | Notes |
|-----------|----------------------|-------|
| `apikey`  | `SIMPLEMDM_APIKEY`   | Required. API key for your tenant. |
| `host`    | `SIMPLEMDM_HOST`     | Optional. Override the API hostname (defaults to `a.simplemdm.com`). |

## Documentation and examples

* Generated documentation for every resource and data source lives in [`docs/`](./docs/).
* Copyable end-to-end examples for the provider, resources, and data sources are in [`examples/`](./examples/).

Regenerate the documentation whenever schemas change:

```bash
go generate ./...
```

## Project layout

| Path | Description |
|------|-------------|
| [`provider/`](./provider/) | Provider, resource, data source, and acceptance test implementations. |
| [`internal/`](./internal/) | Helper packages that support the provider. |
| [`docs/`](./docs/) | Terraform Plugin Docs output used by the Terraform Registry. |
| [`examples/`](./examples/) | Working configuration snippets for the provider. |
| [`scripts/`](./scripts/) | Utility scripts, including fixture discovery helpers for tests. |

## Development workflow

This project uses Go 1.24 (see [`go.mod`](./go.mod)). Typical development steps:

```bash
# Install dependencies and verify the provider builds
go mod download
go build ./...

# Run linting (same tool as CI)
golangci-lint run

# Run unit tests
go test ./...
```

### Acceptance tests

Acceptance tests are located under [`provider/`](./provider/) and can be run with:

```bash
TF_ACC=1 SIMPLEMDM_APIKEY="your-api-key" go test -v -cover ./provider/
```

The suite skips tests automatically when required fixtures are missing, allowing day-to-day
development to rely on dynamic coverage while CI can opt into additional cases by
setting the appropriate environment variables. GitHub Actions runs the same command in
[`.github/workflows/test.yml`](.github/workflows/test.yml).

#### Fixture environment variables

The following optional variables unlock additional tests. Values should reference existing
objects in a SimpleMDM test tenant:

| Variable | Used by | Purpose |
|----------|---------|---------|
| `SIMPLEMDM_APP_ID` | App data source, assignment group resource | ID of an app available to your tenant. |
| `SIMPLEMDM_ASSIGNMENT_GROUP_ID` | Assignment group data source and resource | Fixture assignment group (device groups are deprecated). |
| `SIMPLEMDM_ATTRIBUTE_NAME` | Attribute data source | Name of an existing custom attribute. |
| `SIMPLEMDM_CUSTOM_DECLARATION_DEVICE_ID` | Custom declaration device assignment resource | Device capable of receiving DDM declarations. |
| `SIMPLEMDM_DEVICE_GROUP_ID` | Device group data source, enrollment and script job resources | Existing device group when cloning or referencing real groups. |
| `SIMPLEMDM_DEVICE_GROUP_CLONE_SOURCE_ID` | Device group resource | Source device group for clone operations. |
| `SIMPLEMDM_DEVICE_GROUP_NAME` | Device group resource | Name reused across updates during acceptance tests. |
| `SIMPLEMDM_DEVICE_GROUP_ATTRIBUTE_KEY` | Device group resource | Attribute key validated during updates. |
| `SIMPLEMDM_DEVICE_GROUP_ATTRIBUTE_VALUE` | Device group resource | Initial attribute value. |
| `SIMPLEMDM_DEVICE_GROUP_ATTRIBUTE_UPDATED_VALUE` | Device group resource | Updated attribute value. |
| `SIMPLEMDM_DEVICE_GROUP_PROFILE_ID` | Device and device group resources | Profile assigned during tests. |
| `SIMPLEMDM_DEVICE_GROUP_PROFILE_UPDATED_ID` | Device and device group resources | Updated profile reference. |
| `SIMPLEMDM_DEVICE_GROUP_CUSTOM_PROFILE_ID` | Device and device group resources | Custom profile assigned during tests. |
| `SIMPLEMDM_DEVICE_GROUP_CUSTOM_PROFILE_UPDATED_ID` | Device and device group resources | Updated custom profile reference. |
| `SIMPLEMDM_DEVICE_ID` | Device data source, device command resource, script job data source | ID of an enrolled device. Required for device-centric commands. |
| `SIMPLEMDM_ENROLLMENT_CONTACT` | Enrollment resource | Contact email or phone used to create an enrollment invitation. |
| `SIMPLEMDM_ENROLLMENT_CONTACT_UPDATE` | Enrollment resource | Optional updated contact value to exercise update paths. |
| `SIMPLEMDM_ENROLLMENT_ID` | Enrollment data source | ID of an existing enrollment. |
| `SIMPLEMDM_PROFILE_ID` | Profile data source, assignment group resource | ID of a profile created via the SimpleMDM UI. |
| `SIMPLEMDM_SCRIPT_ID` | Script data source | ID of an existing script. |
| `SIMPLEMDM_SCRIPT_JOB_ID` | Script job data source | ID of an existing script job. |

Use [`scripts/discover-test-fixtures.sh`](./scripts/discover-test-fixtures.sh) to collect common fixture
IDs automatically from your tenant and output `gh secret set` commands that match the CI workflow.

## Known issues

* Device groups are deprecated in SimpleMDM. The legacy `simplemdm_devicegroup` resource and data source remain for backward compatibility, but new deployments should favor `simplemdm_assignmentgroup`.
* Device name updates require a manual PATCH request outside of Terraform.
* Profiles and custom profiles applied to assignment groups or devices cannot be updated via API; Terraform compares the desired configuration against the previous state only.
