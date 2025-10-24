# SimpleMDM API Coverage Review

This document captures the results of comparing the Terraform provider resources/data sources with the documented SimpleMDM API endpoints.

## Source documentation

The SimpleMDM API documentation is published at https://api.simplemdm.com/v1. The current environment does not allow downloading the OpenAPI definition directly (HTTPS CONNECT is blocked and returns `403 Forbidden`). The catalog of endpoints below therefore reflects the endpoints already tracked by the provider in `internal/apicatalog/catalog.go`.

## Endpoint coverage summary

| API Section | Endpoint | Terraform Resource | Terraform Data Source | Coverage Notes |
|-------------|----------|--------------------|-----------------------|----------------|
| Apps | `/api/v1/apps` | `simplemdm_app` | `simplemdm_app` | Covered |
| Assignment Groups | `/api/v1/assignment_groups` | `simplemdm_assignmentgroup` | `simplemdm_assignmentgroup` | Covered |
| Custom Attributes | `/api/v1/attributes` | `simplemdm_attribute` | `simplemdm_attribute` | Covered |
| Custom Profiles | `/api/v1/custom_profiles` | `simplemdm_customprofile` | `simplemdm_customprofile` | Covered |
| Devices | `/api/v1/devices` | `simplemdm_device` | `simplemdm_device` | Covered |
| Device Groups | `/api/v1/device_groups` | `simplemdm_devicegroup` | `simplemdm_devicegroup` | Covered |
| Profiles | `/api/v1/profiles` | _Not implemented_ | `simplemdm_profile` | ❗ Resource missing |
| Scripts | `/api/v1/scripts` | `simplemdm_script` | `simplemdm_script` | Covered |
| Script Jobs | `/api/v1/script_jobs` | `simplemdm_scriptjob` | `simplemdm_scriptjob` | Covered |

## TODOs

- [ ] Implement a `simplemdm_profile` resource to cover the `/api/v1/profiles` endpoint.
- [ ] Re-run the coverage review after the SimpleMDM OpenAPI specification can be downloaded from https://api.simplemdm.com/v1/openapi.json.

