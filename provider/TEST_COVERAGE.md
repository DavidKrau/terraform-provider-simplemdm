# Test Coverage and Fixture Requirements

## Overview

This document outlines which tests are fully dynamic (require only `SIMPLEMDM_APIKEY`) and which tests require additional fixtures due to API or operational limitations.

## GitHub Secrets Setup (2025-10-31)

**📖 Quick Start**: See [`GITHUB_SECRETS_GUIDE.md`](../GITHUB_SECRETS_GUIDE.md) for setting up missing test fixtures

Four GitHub secrets are configured in the workflow but need values added:
- `SIMPLEMDM_DEVICE_ID` - Enables 5 device-related tests
- `SIMPLEMDM_DEVICE_GROUP_CLONE_SOURCE_ID` - Enables device group cloning tests
- `SIMPLEMDM_SCRIPT_JOB_ID` - Enables script job data source tests
- `SIMPLEMDM_CUSTOM_DECLARATION_DEVICE_ID` - Enables DDM assignment tests

**Impact**: Adding these secrets will increase acceptance test coverage from ~70% to ~80% (+8 tests)

**Auto-Discovery**: Run `./scripts/discover-test-fixtures.sh` to automatically find fixture IDs from your SimpleMDM account

See also: [`TESTING_SETUP.md`](../TESTING_SETUP.md) for detailed setup documentation

---

## Recent Coverage Improvement Effort (2025-10-31)

An attempt was made to increase unit test coverage from 2.4% to 60%+. See [`TEST_COVERAGE_IMPROVEMENT_REPORT.md`](./TEST_COVERAGE_IMPROVEMENT_REPORT.md) for:
- Detailed analysis of current architecture
- Why 60% coverage is challenging without refactoring
- Three approaches to achieve 60%+ coverage
- Recommended path forward

**Current Unit Test Coverage**: 2.7%
- Added comprehensive data transformation tests
- Identified architectural limitations
- Documented pragmatic testing strategy

**Key Finding**: For Terraform providers, acceptance test coverage is more valuable than unit test coverage percentage.

---

## ✅ Fully Dynamic Tests (No Fixtures Required)

These tests can run with only `SIMPLEMDM_APIKEY` and `TF_ACC=1` environment variables. They create all necessary resources dynamically during test execution.

### Resources
- **App Resource** - Creates apps dynamically
- **Assignment Group Resource** - Creates assignment groups dynamically
- **Attribute Resource** - Creates custom attributes dynamically
- **Custom Profile Resource** - Creates custom profiles from test files
- **Device Group Resource** - Creates device groups dynamically (NOTE: Requires fixture profile IDs for full testing)
- **Enrollment Resource** - Creates enrollments and device groups dynamically
- **Managed Config Resource** - Creates managed configs dynamically
- **Script Resource** - Creates scripts from test files dynamically

### Data Sources
- **App Data Source** - Uses dynamically created app
- **Assignment Group Data Source** - Uses dynamically created assignment group
- **Attribute Data Source** - Uses dynamically created attribute
- **Custom Profile Data Source** - ✅ **NEWLY DYNAMIC** - Uses dynamically created custom profile
- **Device Group Data Source** - ✅ **NEWLY DYNAMIC** - Uses dynamically created device group
- **Enrollment Data Source** - Uses dynamically created enrollment
- **Managed Config Data Source** - Uses dynamically created managed config
- **Script Data Source** - Uses dynamically created script
- **Script Job Data Source** - Uses dynamically created script job
- **Script Job Resource** - ✅ **NEWLY DYNAMIC** - Creates script and device group dynamically

### Running Dynamic Tests

```bash
export SIMPLEMDM_APIKEY="your-api-key"
export TF_ACC="1"

# Run all dynamic tests
go test -v ./provider/ -timeout 30m

# Run specific dynamic test
go test -v ./provider/ -run TestAccCustomProfileDataSource
```

---

## 🔒 Tests Requiring Fixtures

These tests require external fixtures because of API or operational limitations.

### Profile Resources (API Limitation)

**Why Fixtures Required:** Profiles can only be created through the SimpleMDM web UI. The API only supports reading and updating existing profiles.

#### Tests Affected:
- `TestAccProfileResource_ReadOnly`
- `TestAccProfileResource_NonExistent`
- `TestAccProfileDataSource`

#### Required Environment Variables:
```bash
export SIMPLEMDM_PROFILE_ID="212749"  # Replace with actual profile ID
```

#### Running Profile Tests:
```bash
export SIMPLEMDM_APIKEY="your-api-key"
export TF_ACC="1"
export SIMPLEMDM_PROFILE_ID="your-profile-id"

go test -v ./provider/ -run "TestAccProfile"
```

---

### Device Resources (Operational Limitation)

**Why Fixtures Required:** Devices cannot be created via API. They must be physically enrolled through Apple's Device Enrollment Program (DEP) or manual enrollment.

#### Tests Affected:
- `TestAccDeviceDataSource` - Requires enrolled device
- `TestAccDeviceCommandResource_*` - Requires enrolled device for command execution
- `TestAccDeviceInstalledAppsDataSource` - Requires enrolled device with installed apps (currently skipped)
- `TestAccDeviceProfilesDataSource` - Requires enrolled device with profiles (currently skipped)
- `TestAccDeviceUsersDataSource` - Requires enrolled device with user accounts (currently skipped)

#### Required Environment Variables:
```bash
export SIMPLEMDM_DEVICE_ID="123456"  # Replace with actual enrolled device ID
```

#### Running Device Tests:
```bash
export SIMPLEMDM_APIKEY="your-api-key"
export TF_ACC="1"
export SIMPLEMDM_DEVICE_ID="your-device-id"

# Run device data source test
go test -v ./provider/ -run "TestAccDeviceDataSource"

# Run device command tests (safe commands)
go test -v ./provider/ -run "TestAccDeviceCommandResource_PushApps"
go test -v ./provider/ -run "TestAccDeviceCommandResource_Refresh"

# WARNING: This will lock the device!
go test -v ./provider/ -run "TestAccDeviceCommandResource_Lock"
```

---

## 📊 Fixture Requirements Summary

### Minimal Test Setup (Most Tests)
```bash
export SIMPLEMDM_APIKEY="your-api-key"
export TF_ACC="1"
```
This runs all dynamic tests including:
- All resource CRUD operations that support API creation
- All data sources that use dynamically created resources
- Script jobs (now fully dynamic)
- Custom profile data source (now fully dynamic)
- Device group data source (now fully dynamic)

### Full Test Coverage (All Tests)
```bash
export SIMPLEMDM_APIKEY="your-api-key"
export TF_ACC="1"
export SIMPLEMDM_PROFILE_ID="your-profile-id"        # For profile tests
export SIMPLEMDM_DEVICE_ID="your-device-id"          # For device tests
```

---

## 🔍 Test Categories

### 1. Device Command Tests

**Type:** Fixture-dependent (requires enrolled device)

**Commands Tested:**
- ✅ `push_assigned_apps` - Safe, commonly used
- ✅ `refresh` - Safe, triggers device check-in
- ✅ `lock` - **WARNING:** Will lock device screen
- ✅ Invalid command - Error handling

**Additional Supported Commands (Not Currently Tested):**
- `restart`, `shutdown`, `clear_passcode`
- `rotate_firmware_password`, `rotate_recovery_lock_password`
- `rotate_filevault_recovery_key`, `rotate_admin_password`
- `wipe`, `update_os`, `unenroll`

**Expected Behavior:**
- **Create**: Executes command, returns HTTP 202 status
- **Read**: Returns stored state (commands can't be read from API)
- **Update**: Not supported (returns error)
- **Delete**: No-op (commands can't be undone)
- **Import**: Supported (imports stored state)

### 2. Profile Tests

**Type:** Fixture-dependent (profiles created in UI only)

**Test Cases:**
- ✅ Read existing profile
- ✅ Verify all attributes
- ✅ ImportState
- ✅ Error handling (non-existent profile)

**Expected Behavior:**
- **Create**: Reads existing profile and imports into state
- **Read**: Refreshes profile data from SimpleMDM
- **Update**: Refreshes profile data (no actual updates)
- **Delete**: Removes from state only (doesn't delete from SimpleMDM)
- **Import**: Supports importing existing profiles by ID

### 3. Custom Profile Tests

**Type:** Fully dynamic ✅

**Test Cases:**
- ✅ Create custom profile from mobileconfig file
- ✅ Update custom profile attributes
- ✅ Data source reads dynamically created profile

---

## 🎯 Recent Improvements

### Made Dynamic (No Longer Require Fixtures)

1. **Custom Profile Data Source** (`customProfile_data_source_test.go`)
   - Previously required: `SIMPLEMDM_CUSTOM_PROFILE_ID`
   - Now: Creates custom profile dynamically, then reads with data source
   - Pattern: Resource creation → Data source read with reference

2. **Device Group Data Source** (`deviceGroup_data_source_test.go`)
   - Previously required: `SIMPLEMDM_DEVICE_GROUP_ID`
   - Now: Creates device group dynamically, then reads with data source
   - Pattern: Resource creation → Data source read with reference

3. **Script Job Resource** (`scriptJob_resource_test.go`)
   - Previously required: `SIMPLEMDM_DEVICE_GROUP_ID` and `SIMPLEMDM_SCRIPT_ID`
   - Now: Creates both script and device group dynamically
   - Pattern: Create dependencies → Create script job → Test operations

### Enhanced Documentation

All fixture-dependent tests now include clear comments explaining:
- Why fixtures are required
- What environment variables are needed
- How to run the tests
- Any safety warnings (e.g., device lock commands)

---

## 🚀 CI/CD Configuration

### GitHub Actions Example

```yaml
name: Acceptance Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run Dynamic Tests
        env:
          SIMPLEMDM_APIKEY: ${{ secrets.SIMPLEMDM_APIKEY }}
          TF_ACC: "1"
        run: go test -v ./provider/ -timeout 30m
      
      # Optional: Run fixture-dependent tests if secrets are configured
      - name: Run Profile Tests
        if: ${{ secrets.SIMPLEMDM_PROFILE_ID != '' }}
        env:
          SIMPLEMDM_APIKEY: ${{ secrets.SIMPLEMDM_APIKEY }}
          SIMPLEMDM_PROFILE_ID: ${{ secrets.SIMPLEMDM_PROFILE_ID }}
          TF_ACC: "1"
        run: go test -v ./provider/ -run "TestAccProfile" -timeout 10m
      
      - name: Run Device Tests
        if: ${{ secrets.SIMPLEMDM_DEVICE_ID != '' }}
        env:
          SIMPLEMDM_APIKEY: ${{ secrets.SIMPLEMDM_APIKEY }}
          SIMPLEMDM_DEVICE_ID: ${{ secrets.SIMPLEMDM_DEVICE_ID }}
          TF_ACC: "1"
        run: go test -v ./provider/ -run "TestAccDevice" -timeout 10m
```

---

## 📝 Developer Notes

### Adding New Tests

When adding new tests, prefer dynamic resource creation over fixtures:

**✅ Good (Dynamic):**
```go
Config: providerConfig + `
    resource "simplemdm_script" "test" {
        name = "Test Script"
        scriptfile = file("./testfiles/testscript.sh")
    }
    
    resource "simplemdm_scriptjob" "test" {
        script_id = simplemdm_script.test.id
        device_ids = []
        group_ids = []
    }
`
```

**❌ Avoid (Fixtures):**
```go
scriptID := testAccRequireEnv(t, "SIMPLEMDM_SCRIPT_ID")
Config: providerConfig + fmt.Sprintf(`
    resource "simplemdm_scriptjob" "test" {
        script_id = "%s"
        device_ids = []
        group_ids = []
    }
`, scriptID)
```

### When Fixtures Are Acceptable

Use fixtures only when:
1. Resources cannot be created via API (e.g., profiles, devices)
2. Resources require significant setup time (enrolled devices with apps/profiles)
3. Resources require external configuration (Apple DDM declarations)

Document clearly why fixtures are required and how to obtain them.

---

## 📚 Resources

- [SimpleMDM API Documentation](https://simplemdm.com/docs/api/)
- [Terraform Plugin Testing](https://developer.hashicorp.com/terraform/plugin/testing)
- [Test File Locations](./provider/)