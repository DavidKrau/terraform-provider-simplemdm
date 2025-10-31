# Test Coverage for Profile and Device Command Resources

## Profile Resource Tests

### Overview
The profile resource is READ-ONLY because profiles can only be created and managed through the SimpleMDM web UI. The Terraform resource is for state management only.

### Test Cases

1. **TestAccProfileResource_ReadOnly**
   - Reads an existing profile using fixture ID
   - Verifies all profile attributes are populated correctly
   - Tests ImportState functionality
   - Validates that profiles can be managed in Terraform state

2. **TestAccProfileResource_NonExistent**
   - Tests error handling for non-existent profile IDs
   - Verifies appropriate error messages

### Running Profile Tests

```bash
export SIMPLEMDM_APIKEY="your-api-key"
export TF_ACC="1"
export SIMPLEMDM_PROFILE_ID="212749"  # Replace with actual profile ID

go test -v ./provider/ -run "TestAccProfileResource"
```

### Expected Behavior

- **Create**: Reads existing profile and imports into state
- **Read**: Refreshes profile data from SimpleMDM
- **Update**: Refreshes profile data (no actual updates)
- **Delete**: Removes from state only (doesn't delete from SimpleMDM)
- **Import**: Supports importing existing profiles by ID

---

## Device Command Resource Tests

### Overview
Device commands are "one-shot" resources that execute once during creation. They cannot be read back from the API, and the resource only maintains local state after execution.

### Test Cases

1. **TestAccDeviceCommandResource_PushApps**
   - Tests the `push_assigned_apps` command
   - Safe, commonly-used command
   - Verifies successful execution (HTTP 202)

2. **TestAccDeviceCommandResource_Refresh**
   - Tests the `refresh` command
   - Safe command that triggers device check-in
   - Verifies successful execution (HTTP 202)

3. **TestAccDeviceCommandResource_Lock**
   - Tests the `lock` command with parameters
   - Demonstrates parameter passing
   - **WARNING**: Will lock the device screen

4. **TestAccDeviceCommandResource_InvalidCommand**
   - Tests error handling for unsupported commands
   - Verifies appropriate error messages

### Running Device Command Tests

```bash
export SIMPLEMDM_APIKEY="your-api-key"
export TF_ACC="1"
export SIMPLEMDM_DEVICE_ID="123456"  # Replace with actual enrolled device ID

# Run all device command tests
go test -v ./provider/ -run "TestAccDeviceCommandResource"

# Run specific test (e.g., only safe commands)
go test -v ./provider/ -run "TestAccDeviceCommandResource_PushApps"
go test -v ./provider/ -run "TestAccDeviceCommandResource_Refresh"
```

### Expected Behavior

- **Create**: Executes command against device, returns status code and response
- **Read**: Returns stored state (commands can't be read back from API)
- **Update**: Not supported (returns error)
- **Delete**: No-op (commands can't be undone)
- **Import**: Supported (imports stored state)

### Supported Commands

The following commands are tested:
- `push_assigned_apps` - Push assigned apps to device
- `refresh` - Trigger device to check in
- `lock` - Lock device screen (with optional message)

Additional supported but untested commands include:
- `restart`, `shutdown`, `clear_passcode`
- `rotate_firmware_password`, `rotate_recovery_lock_password`
- `rotate_filevault_recovery_key`, `rotate_admin_password`
- `wipe`, `update_os`, `unenroll`
- And more (see device_command_resource.go)

---

## Test Requirements

### Environment Variables

Both test suites require:
- `SIMPLEMDM_APIKEY` - Your SimpleMDM API key
- `TF_ACC=1` - Enable acceptance tests

Profile tests require:
- `SIMPLEMDM_PROFILE_ID` - ID of an existing profile in SimpleMDM

Device command tests require:
- `SIMPLEMDM_DEVICE_ID` - ID of an enrolled, active device

### Test Execution Notes

1. **Fixtures**: Both resources require external fixtures (profiles and devices) that must already exist in SimpleMDM
2. **Safety**: Device command tests execute real commands - be cautious with destructive commands
3. **Skip Logic**: Tests automatically skip if required environment variables are not set
4. **CI/CD**: Configure environment variables in your CI pipeline to enable these tests

### Example: Running All Tests

```bash
# Set up environment
export SIMPLEMDM_APIKEY="tKgXtjvJkExAVtIjOCBzg5IPvZafweqIKKmiRtFcZEjCv2iJTdk3gn6Uy8Gmb9j7"
export TF_ACC="1"
export SIMPLEMDM_PROFILE_ID="212749"
export SIMPLEMDM_DEVICE_ID="123456"  # Replace with actual device

# Run tests
go test -v ./provider/ -run "TestAccProfileResource|TestAccDeviceCommandResource" -timeout 30m
```

---

## Test Coverage Summary

### Profile Resource
- ✅ Read existing profile
- ✅ Verify all attributes
- ✅ ImportState
- ✅ Error handling (non-existent profile)
- ⏭️ Create (N/A - read-only resource)
- ⏭️ Update (N/A - read-only resource)
- ⏭️ Delete (removes from state only)

### Device Command Resource
- ✅ Execute push_apps command
- ✅ Execute refresh command
- ✅ Execute command with parameters (lock)
- ✅ Error handling (invalid command)
- ✅ Verify status codes
- ⏭️ Update (not supported by design)
- ⏭️ Delete (no-op by design)

Both resources now have appropriate test coverage considering their unique characteristics as read-only and one-shot resources respectively.