package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// boolPointerFromType converts a types.Bool to *bool, returning nil for null/unknown values
func boolPointerFromType(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	v := value.ValueBool()
	return &v
}

// stringPointerFromType converts a types.String to *string, returning nil for null/unknown values
func stringPointerFromType(value types.String) *string {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	v := value.ValueString()
	return &v
}

// int64PointerFromType converts a types.Int64 to *int64, returning nil for null/unknown values
func int64PointerFromType(value types.Int64) *int64 {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	v := value.ValueInt64()
	return &v
}

// boolValueOrDefault returns the bool value if it exists, otherwise returns the fallback
func boolValueOrDefault(value types.Bool, fallback bool) bool {
	if value.IsNull() || value.IsUnknown() {
		return fallback
	}

	return value.ValueBool()
}

// stringValueOrDefault returns the string value if it exists, otherwise returns the fallback
func stringValueOrDefault(value types.String, fallback string) string {
	if value.IsNull() || value.IsUnknown() {
		return fallback
	}

	return value.ValueString()
}

// int64ValueOrDefault returns the int64 value if it exists, otherwise returns the fallback
func int64ValueOrDefault(value types.Int64, fallback int64) int64 {
	if value.IsNull() || value.IsUnknown() {
		return fallback
	}

	return value.ValueInt64()
}
