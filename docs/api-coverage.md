# SimpleMDM API Coverage Review

This document captures the results of comparing the Terraform provider resources/data sources with the documented SimpleMDM API endpoints.

## Source documentation

The SimpleMDM API documentation is published at https://api.simplemdm.com/v1. The reference HTML was retrieved with `curl https://api.simplemdm.com/v1` on 2025-10-24 and refreshed on 2025-02-14 while preparing the coverage summary below.

## Endpoint coverage summary

| API Section | Representative Endpoint(s) | Terraform Resource | Terraform Data Source | Coverage Notes |
|-------------|----------------------------|--------------------|-----------------------|----------------|
| Account | `/api/v1/account` | - | - | Not covered |
| Apps | `/api/v1/apps` | `simplemdm_app` | `simplemdm_app` | Covered (resource supports App Store, bundle identifier, and binary uploads; state reflects deploy status, processing status, install channels, timestamps, and store metadata) |
| Assignment Groups | `/api/v1/assignment_groups` | `simplemdm_assignmentgroup` | `simplemdm_assignmentgroup` | Covered |
| Custom Configuration Profiles | `/api/v1/custom_configuration_profiles` | `simplemdm_customprofile` | `simplemdm_customprofile` | Covered (resource & data source expose profile metadata, download content/SHA, and assignments are handled through device and device group resources) |
| Custom Attributes | `/api/v1/custom_attributes` | `simplemdm_attribute` | `simplemdm_attribute` | Covered (resource and data source manage attribute definitions; device, assignment group, and device group resources integrate with the corresponding `custom_attribute_values` endpoints) |
| Custom Declarations | `/api/v1/custom_declarations` | `simplemdm_customdeclaration` | `simplemdm_customdeclaration` | Covered |
| DEP Servers | `/api/v1/dep_servers` | - | - | Not covered |
| Devices | `/api/v1/devices` | `simplemdm_device` | `simplemdm_device` | Covered |
| Device Groups (deprecated) | `/api/v1/device_groups` | `simplemdm_devicegroup` | `simplemdm_devicegroup` | Covered (resource provisions legacy groups via create or clone, manages name, attributes, configuration profile assignments, and calls the delete endpoint; data source reads existing groups) |
| Enrollments | `/api/v1/enrollments` | - | - | Not covered |
| Installed Apps | `/api/v1/installed_apps` | - | - | Not covered |
| Logs | `/api/v1/logs` | - | - | Not covered |
| Lost Mode | `/api/v1/devices/{DEVICE_ID}/lost_mode` | - | - | Not covered |
| Managed App Configs | `/api/v1/apps/{APP_ID}/managed_configs` | - | - | Not covered |
| Profiles | `/api/v1/profiles` | `simplemdm_profile` | `simplemdm_profile` | Covered |
| Push Certificate | `/api/v1/push_certificate` | - | - | Not covered |
| Scripts | `/api/v1/scripts` | `simplemdm_script` | `simplemdm_script` | Covered |
| Script Jobs | `/api/v1/script_jobs` | `simplemdm_scriptjob` | `simplemdm_scriptjob` | Covered |
| Webhooks | (Event delivery) | - | - | Not covered |

## Observations

- The provider implements Terraform resources and data sources for all configuration collections backed by the `/api/v1` endpoints enumerated in `internal/apicatalog/catalog.go`. These include apps, assignment groups, custom attributes, custom configuration profiles, custom declarations, devices, device groups (including create/clone/update/delete support), profiles, scripts, and script jobs.
- Additional API sections such as account management, DEP servers, enrollments, installed apps, logging, device actions (lost mode, push certificate lifecycle, and managed app configs), and webhook subscriptions currently have no Terraform coverage.
